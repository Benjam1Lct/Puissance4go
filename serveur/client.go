package main

import (
	"encoding/json"
	"log"
	"net"
	"strings"
)

// processMessage traite les messages reçus d'un client en fonction de leur type.
// Chaque type de message déclenche une action spécifique (par exemple, mise à jour de curseur, sélection de couleur, mouvement).
func processMessage(msg Message, id int) {
	switch msg.Type {
	case "restartReady":
		log.Printf("Joueur %d prêt à redémarrer.\n", id)
		restartReadyChannel <- id // Envoyer l'ID dans le channel
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
	case "token_update":
		log.Println("token client recu")
		var payload map[string]int
		if decodePayload(msg.Payload, &payload) == nil {
			sendPosition(payload, id)
		}
	case "require_history":
		sendHistory(id)
	case "chat":
		var payload ChatMessage
		if decodePayload(msg.Payload, &payload) == nil {
			broadcastChatMessage(id, payload.Text)
		}
	case "selected":
		var payload SelectedPayload
		if decodePayload(msg.Payload, &payload) == nil {
			handleSelection(payload, id)
		}
	default:
		log.Printf("Type de message inconnu : %s\n", msg.Type)
	}
}

// Envoie les messages du chat d'un client vers l'autre
func broadcastChatMessage(senderID int, text string) {
	message := Message{
		Type: "chat",
		Payload: map[string]interface{}{
			"id":   senderID,
			"text": text,
		},
	}

	notifyPlayers(message) // Envoyer à tous les joueurs
	log.Printf("Message de chat de %d : %s\n", senderID, text)
}

// Envoie l'historique des coups de la partie pour le replay du client
func sendHistory(id int) {
	clientMux.Lock()
	conn, ok := clients[id]
	clientMux.Unlock()

	if !ok {
		log.Printf("Client %d introuvable pour l'envoi de l'historique\n", id)
		return
	}

	// Envoyer l'historique au client demandeur
	historyMessage := Message{
		Type:    "sent_history",
		Payload: historiquePartie,
	}

	if err := sendJSONMessage(conn, historyMessage); err != nil {
		log.Printf("Erreur lors de l'envoi de l'historique au client %d : %v\n", id, err)
	} else {
		log.Printf("Historique envoyé au client %d\n", id)
	}
}

// Envoie la position du cursor du jouer au dessus de la grille pendant la partie
func sendPosition(payload map[string]int, id int) {
	if position, ok := payload["position"]; ok {
		// Notifier l'autre joueur de la position du curseur
		notifyOtherPlayers(id, Message{
			Type: "token_update",
			Payload: map[string]int{
				"position": position,
			},
		})
		log.Printf("Mise à jour du curseur du joueur %d : position %d\n", id, position)
	}
}

// Envoie la position du curseur d'un client sur la grille de couleur a l'autre joueur
func cursorUpdate(payload map[string]int, id int) {
	if color, ok := payload["color"]; ok {
		// Notifier l'autre joueur de la position du curseur
		notifyOtherPlayers(id, Message{
			Type: "cursor_update",
			Payload: map[string]int{
				"color": color,
			},
		})
	}
}

// decodePayload désérialise le payload générique en une structure cible spécifique.
// Cette fonction est utile lorsque payload est reçu sous forme d'interface{} et doit être converti
// en une structure typée spécifique target.
// Elle utilise un double passage (Marshal -> Unmarshal) pour garantir une conversion correcte.
func decodePayload(payload interface{}, target interface{}) error {
	jsonData, err := json.Marshal(payload) // Re-marshal pour convertir en bytes
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, target)
}

// colorSelection gère la sélection de couleur par un joueur.
// Elle met à jour la couleur choisie par le joueur, notifie les autres joueurs,
// et vérifie si tous les joueurs ont terminé leur sélection.
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

// move gère le déplacement effectué par un joueur.
// La fonction enregistre le déplacement dans l'historique de la partie,
// incrémente le numéro de tour, et notifie les autres joueurs du mouvement.
func move(payload MovePayload, id int) {
	x, y := payload.X, payload.Y
	historiquePartie[turnPartie] = Coordinate{ID: id, X: x, Y: y}
	turnPartie++

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

// ready gère le signalement d'un joueur indiquant qu'il est prêt à jouer.
// Elle met à jour l'état de préparation du joueur dans readyPlayers,
// puis vérifie si tous les joueurs sont prêts pour démarrer la partie.
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

// Est appelé par ready pour verifier si tous les joueurs connectés sont prêts à jouer.
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

// Est appelé par colorSelection pour verifier si tous les joueurs connectés ont choisi leur couleur.
func allPlayersSelectedColors() bool {
	clientMux.Lock()
	defer clientMux.Unlock()

	if len(playerColors) < len(clients) {
		return false // Tous les joueurs n'ont pas encore sélectionné leur couleur
	}

	return true
}

// Appelé lorsque les clients sont déconnectés du serveur afin de remettre toutes les variables à zero.
func resetAll() {
	log.Println("Commande 'resetAll' reçue. Réinitialisation du serveur...")
	clients = make(map[int]net.Conn)  // Liste des clients connectés
	readyPlayers = make(map[int]bool) // Suivi des joueurs prêts
	playerColors = make(map[int]int)
	firstPlayer = -1 // -1 indique qu'aucun joueur n'a encore été désigné
	historiquePartie = make(map[int]Coordinate)
	playerSelections = make(map[int]string) // Table de hachage pour stocker les sélections des joueurs (par exemple, une chaîne représentant leur choix).

	log.Printf("All server have been reset.")
}

// Appelé à la fin d'une partie lors du rematch pour recommencer une nouvelle partie
func resetServerState() {
	clientMux.Lock()
	turnPartie = 0
	historiquePartie = make(map[int]Coordinate)
	defer clientMux.Unlock()

	log.Println("Serveur pret pour une nouvelle partie")
}

// Appelé lorsqu'un client se déconnecte du serveur, notifie l'autre client pour qu'il se déconnecte
// et reinitialise tout le serveur.
func disconnectClient(conn net.Conn, id int) {

	message := Message{
		Type:    "other_disconnected",
		Payload: nil,
	}

	// Notifier tous les joueurs
	notifyOtherPlayers(id, message)

	clientMux.Lock()
	delete(clients, id)
	delete(readyPlayers, id)
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
		restartWaitForRestart()
		log.Println("En attente de connexions...")
	}
	clientMux.Unlock()
}

// Stock dans la table de hachage le coup effectué par un client au pierre/feuille/ciseaux
func handleSelection(payload SelectedPayload, id int) {
	clientMux.Lock()
	playerSelections[id] = payload.Selected
	nbSelections := len(playerSelections)
	clientMux.Unlock()

	log.Printf("Joueur %d a choisi : %s\n", id, payload.Selected)

	// Vérifier si les deux joueurs ont fait leur sélection
	if nbSelections == 2 {
		determineWinner()
	}
}

// Appelé par handleSelection pour determiner le gagnant du jeu
// Determine par la suite le joueur qui commence la partie et le notifie au client
func determineWinner() {
	// Récupérer les sélections des deux joueurs
	var player1Selection, player2Selection string
	var player1ID, player2ID int

	clientMux.Lock()
	// S'assurer que nous avons bien deux sélections
	if len(playerSelections) != 2 {
		clientMux.Unlock()
		return
	}

	// Récupérer les sélections de manière ordonnée
	for id, selection := range playerSelections {
		if player1ID == 0 {
			player1ID = id
			player1Selection = selection
		} else {
			player2ID = id
			player2Selection = selection
		}
	}
	clientMux.Unlock()

	// Vérifier que les deux sélections sont valides
	if player1Selection == "" || player2Selection == "" {
		log.Printf("Erreur : sélections invalides (j1: %s, j2: %s)\n", player1Selection, player2Selection)
		return
	}

	// Déterminer le gagnant selon les règles du shifumi
	var winnerID int
	switch {
	case player1Selection == player2Selection:
		winnerID = -1 // Match nul
	case (strings.ToLower(player1Selection) == "pierre" && strings.ToLower(player2Selection) == "ciseaux") ||
		(strings.ToLower(player1Selection) == "ciseaux" && strings.ToLower(player2Selection) == "papier") ||
		(strings.ToLower(player1Selection) == "papier" && strings.ToLower(player2Selection) == "pierre"):
		winnerID = player1ID
	default:
		winnerID = player2ID
	}

	// Si match nul, on refait une partie
	if winnerID == -1 {
		result := Message{
			Type: "shifumi_result",
			Payload: map[string]interface{}{
				"result": "draw",
				"player1": map[string]interface{}{
					"id":        player1ID,
					"selection": player1Selection,
				},
				"player2": map[string]interface{}{
					"id":        player2ID,
					"selection": player2Selection,
				},
			},
		}
		notifyPlayers(result)

		// Réinitialiser explicitement les sélections pour le prochain tour
		clientMux.Lock()
		playerSelections = make(map[int]string)
		clientMux.Unlock()
	} else {
		// Le gagnant devient le premier joueur
		firstPlayer = winnerID

		// Notifier les joueurs du résultat et qui commence
		result := Message{
			Type: "shifumi_complete",
			Payload: map[string]interface{}{
				"winner":      winnerID,
				"firstPlayer": winnerID,
				"player1": map[string]interface{}{
					"id":        player1ID,
					"selection": player1Selection,
				},
				"player2": map[string]interface{}{
					"id":        player2ID,
					"selection": player2Selection,
				},
			},
		}
		notifyPlayers(result)
	}

	// Réinitialiser les sélections pour la prochaine partie
	clientMux.Lock()
	playerSelections = make(map[int]string)
	clientMux.Unlock()
}
