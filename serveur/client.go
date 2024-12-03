package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

func processMessage(msg Message, id int) {
	switch msg.Type {
	case "end":
		endPartie(id)
	case "cursor_update":
		var payload map[string]int
		if decodePayload(msg.Payload, &payload) == nil {
			cursorUpdate(payload, id)
		}
	case "color":
		var payload ColorPayload
		if decodePayload(msg.Payload, &payload) == nil {
			colorSelection(payload, id)
		}
	case "move":
		var payload MovePayload
		if decodePayload(msg.Payload, &payload) == nil {
			move(payload, id)
		}
	case "ready":
		ready(id)
	case "disconnect":
		disconnectClient(clients[id], id)
	case "resetAll":
		resetAll()
	default:
		log.Printf("Type de message inconnu : %s\n", msg.Type)
	}
}

func cursorUpdate(payload map[string]int, id int) {
	if color, ok := payload["color"]; ok {
		// Notifier l'autre joueur de la position du curseur
		notifyOtherPlayers(id, Message{
			Type: "cursor_update",
			Payload: map[string]int{
				"color": color,
			},
		})
		log.Printf("Mise à jour du curseur du joueur %d : couleur %d\n", id, color)
	}
}

// Décodage générique de la charge utile
func decodePayload(payload interface{}, target interface{}) error {
	jsonData, err := json.Marshal(payload) // Re-marshal pour convertir en bytes
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}

func colorSelection(payload ColorPayload, id int) {
	color := payload.Color

	clientMux.Lock()
	playerColors[id] = color
	if firstPlayer == -1 {
		firstPlayer = id
	}
	clientMux.Unlock()

	// Créer le message structuré pour la notification
	message := Message{
		Type: "color",
		Payload: map[string]interface{}{
			"id":    id,
			"color": color,
		},
	}

	// Notifier les autres clients
	notifyOtherPlayers(id, message)

	// Vérifier si tous les joueurs ont choisi leurs couleurs
	if allPlayersSelectedColors() {
		notifyPlayers(Message{
			Type: "color_select_complete",
			Payload: map[string]interface{}{
				"firstPlayer": firstPlayer,
			},
		})
		log.Printf("Les deux joueurs ont choisi leurs couleurs. Le joueur %d commence la partie.", firstPlayer)
	}
}

func move(payload MovePayload, id int) {
	x, y := payload.X, payload.Y

	// Créer un message structuré pour la notification
	message := Message{
		Type: "move",
		Payload: map[string]int{
			"x": x,
			"y": y,
		},
	}

	// Notifier les autres joueurs avec un message JSON
	notifyOtherPlayers(id, message)

	log.Printf("Mouvement reçu de %d : (%d, %d)\n", id, x, y)
}

func ready(id int) {
	// Mettre à jour l'état du joueur
	clientMux.Lock()
	readyPlayers[id] = true
	clientMux.Unlock()

	// Vérifier si tous les joueurs sont prêts
	if allPlayersReady() {
		// Créer un message structuré pour notifier les clients
		message := Message{
			Type: "ready",
			Payload: map[string]string{
				"message": "Deux joueurs sont connectés. Vous pouvez commencer à jouer.",
			},
		}

		// Notifier tous les joueurs
		notifyPlayers(message)

		log.Println("Tous les joueurs sont prêts. Notification envoyée.")
	}
}

func allPlayersReady() bool {
	clientMux.Lock()
	defer clientMux.Unlock()

	if len(clients) < 2 {
		return false // Pas assez de joueurs
	}

	for _, ready := range readyPlayers {
		if !ready {
			return false
		}
	}
	return true
}

func allPlayersSelectedColors() bool {
	clientMux.Lock()
	defer clientMux.Unlock()

	if len(playerColors) < len(clients) {
		return false // Tous les joueurs n'ont pas encore sélectionné leur couleur
	}

	return true
}

func endPartie(id int) {
	log.Printf("Le joueur %d a terminé la partie.\n", id)

	// Mettre à jour l'état du joueur comme prêt à redémarrer
	clientMux.Lock()
	playersRestarted[id] = true
	clientMux.Unlock()

	// Informer les autres joueurs qu'un joueur est prêt pour un rematch
	notifyOtherPlayers(id, Message{
		Type: "rematch_waiting",
		Payload: map[string]interface{}{
			"message": fmt.Sprintf("Le joueur %d est en attente de rematch.", id),
		},
	})

	// Vérifier si tous les joueurs sont prêts à redémarrer
	if allPlayersRestarted() {
		// Notifier les joueurs que la partie peut redémarrer
		notifyPlayers(Message{
			Type: "restart_ok",
			Payload: map[string]string{
				"message": "Tous les joueurs sont prêts. La partie peut redémarrer.",
			},
		})

		// Réinitialiser l'état du serveur pour une nouvelle partie
		resetServerState()
		log.Println("Le serveur a été réinitialisé pour une nouvelle partie.")
	}
}

func resetAll() {
	log.Println("Commande 'resetAll' reçue. Réinitialisation du serveur...")
	clients = make(map[int]net.Conn)  // Liste des clients connectés
	readyPlayers = make(map[int]bool) // Suivi des joueurs prêts
	playerColors = make(map[int]int)
	firstPlayer = -1 // -1 indique qu'aucun joueur n'a encore été désigné
	playersRestarted = make(map[int]bool)
	log.Printf("All server have been reset.")
}

func allPlayersRestarted() bool {
	clientMux.Lock()
	defer clientMux.Unlock()

	if len(playersRestarted) < len(clients) {
		return false // Pas tous les joueurs n'ont indiqué qu'ils souhaitent redémarrer
	}

	for _, restarted := range playersRestarted {
		if !restarted {
			return false
		}
	}
	return true
}

func resetServerState() {
	clientMux.Lock()
	defer clientMux.Unlock()

	playersRestarted = make(map[int]bool)
	log.Println("Serveur pret pour une nouvelle partie")
}

func disconnectClient(conn net.Conn, id int) {
	clientMux.Lock()
	delete(clients, id)
	delete(readyPlayers, id)
	delete(playersRestarted, id)
	delete(playerColors, id)
	err := conn.Close()
	if err != nil {
		return
	}
	log.Printf("Client %d déconnecté\n", id)

	// Vérifiez si tous les joueurs sont déconnectés
	if len(clients) == 0 {
		log.Println("Tous les joueurs sont déconnectés. Réinitialisation du serveur...")
		resetAll()
		log.Println("En attente de connexions...")
	}
	clientMux.Unlock()
}
