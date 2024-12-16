package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	clients               = make(map[int]net.Conn)                            // Table de hachage pour stocker les connexions actives des joueurs, associées à leur ID unique.
	readyPlayers          = make(map[int]bool)                                // Table de hachage pour indiquer si un joueur est prêt à jouer.
	clientMux             sync.Mutex                                          // Mutex utilisé pour synchroniser l'accès à des ressources partagées (comme les maps) entre plusieurs goroutines.
	playerColors                                   = make(map[int]int)        // Table de hachage pour associer chaque joueur (par son ID) à une couleur spécifique (représentée par un entier).
	firstPlayer           int                      = -1                       // ID du premier joueur à jouer. Initialisé à -1 pour indiquer qu'il n'a pas encore été défini.
	historiquePartie                               = make(map[int]Coordinate) // Historique des coordonnées des parties jouées, associant un entier (par exemple, un tour ou une action) à une structure `Coordinate`.
	turnPartie            int                      = 0                        // Numéro du tour actuel dans la partie (par exemple, 0 pour le début).
	restartReadyChannel                            = make(chan int, 2)        // Canal avec une capacité tamponnée de 2 pour signaler que deux joueurs sont prêts à redémarrer une partie.
	restartControlChannel                          = make(chan struct{})      // Canal non tamponné pour contrôler ou signaler un redémarrage ou une réinitialisation.
	playerSelections                               = make(map[int]string)     // Table de hachage pour stocker les sélections des joueurs (par exemple, une chaîne représentant leur choix).
)

// startServer démarre le serveur et gère les connexions des clients.
// Pour chaque nouvelle connexion acceptée, un ID unique est attribué au client,
// et la connexion est ajoutée à la table de hachage `clients`.
// La fonction lance ensuite une goroutine `handleClient` pour gérer la communication avec ce client.
func startServer(listener net.Listener) {
	clientID := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Erreur lors de l'acceptation d'une connexion :", err)
			continue
		}

		clientMux.Lock()
		clients[clientID] = conn
		clientMux.Unlock()

		log.Printf("Client %d connecté", clientID)

		go handleClient(conn, clientID)
		clientID++
	}
}

// handleClient gère la communication avec un client spécifique après sa connexion.
// Cette fonction envoie d'abord un message initial contenant l'ID du client,
// puis entre dans une boucle pour lire, désérialiser et traiter les messages envoyés par le client.
// En cas d'erreur ou de déconnexion, le client est déconnecté proprement.
func handleClient(conn net.Conn, id int) {
	defer disconnectClient(conn, id)

	reader := bufio.NewReader(conn)

	// Envoyer l'ID au client en utilisant JSON
	initialMessage := Message{
		Type:    "id",
		Payload: map[string]int{"id": id},
	}
	if err := sendJSONMessage(conn, initialMessage); err != nil {
		log.Printf("Erreur lors de l'envoi de l'ID au client %d : %v\n", id, err)
		return
	}

	// Boucle principale pour lire et traiter les messages
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Erreur de lecture du client %d : %v\n", id, err)
			return
		}

		message = strings.TrimSpace(message)

		// Désérialiser le message JSON
		var msg Message
		if err := json.Unmarshal([]byte(message), &msg); err != nil {
			log.Printf("Erreur de décodage JSON pour le client %d : %v\n", id, err)
			continue // Ignorer ce message et passer au suivant
		}

		// Traiter le message via processMessage
		processMessage(msg, id)
	}
}

// sendJSONMessage envoie un message structuré au format JSON à travers une connexion réseau (net.Conn).
func sendJSONMessage(conn net.Conn, msg Message) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation JSON : %w", err)
	}

	_, err = conn.Write(append(jsonData, '\n')) // Ajout de '\n' pour délimiter les messages
	return err
}

// notifyPlayers envoie un message structuré au format JSON à tous les clients connectés.
func notifyPlayers(message Message) {
	clientMux.Lock()
	defer clientMux.Unlock()

	// Convertir le message en JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Erreur lors de la sérialisation du message JSON : %v\n", err)
		return
	}

	// Envoyer le message JSON à tous les clients
	for id, conn := range clients {
		_, err := conn.Write(append(jsonMessage, '\n')) // Ajouter '\n' pour délimiter le message
		if err != nil {
			log.Printf("Erreur lors de l'envoi au client %d : %v\n", id, err)
		} else {
			log.Printf("Message envoyé au client %d : %s\n", id, string(jsonMessage))
		}
	}
}

// notifyOtherPlayers envoie un message structuré au format JSON à tous les clients connectés
// sauf au client spécifié par `senderID`.
func notifyOtherPlayers(senderID int, message interface{}) {
	clientMux.Lock()
	defer clientMux.Unlock()

	// Convertir le message en JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Erreur lors de la sérialisation du message JSON : %v\n", err)
		return
	}

	// Parcourir les clients et envoyer le message à tous sauf l'expéditeur
	for id, conn := range clients {
		if id != senderID {
			_, err := conn.Write(append(jsonMessage, '\n')) // Ajouter '\n' pour marquer la fin du message
			if err != nil {
				log.Printf("Erreur lors de l'envoi au client %d : %v\n", id, err)
			} else {
				log.Printf("Message envoyé au client %d : %s\n", id, string(jsonMessage))
			}
		}
	}
}

// waitForRestart gère l'état de préparation des joueurs pour redémarrer une partie.
// Cette fonction est exécutée dans une goroutine et écoute deux canaux :
// 1. restartReadyChannel pour savoir quels joueurs sont prêts à un rematch.
// 2. restartControlChannel pour redemarrer proprement la fonction.
func waitForRestart() {
	go func() {
		readyPlayersList := make(map[int]bool) // Suivi des joueurs prêts

		for {
			select {
			case id := <-restartReadyChannel:
				readyPlayersList[id] = true
				log.Printf("Joueur %d prêt pour un rematch.\n", id)

				notifyOtherPlayers(id, Message{
					Type: "rematch_waiting",
					Payload: map[string]string{
						"message": "L'autre joueur est en attente de rematch",
					},
				})

				// Vérifie si tous les joueurs sont prêts
				if len(readyPlayersList) == len(clients) && len(clients) == 2 {
					log.Println("Tous les joueurs sont prêts. Redémarrage de la partie.")

					// Notifier tous les joueurs que la partie peut redémarrer
					notifyPlayers(Message{
						Type: "restart_ok",
						Payload: map[string]string{
							"message": "Tous les joueurs sont prêts. La partie peut redémarrer.",
						},
					})

					// Réinitialise l'état pour une nouvelle partie
					resetServerState()
					readyPlayersList = make(map[int]bool) // Réinitialise pour la prochaine partie
				}

			case <-restartControlChannel:
				log.Println("Arrêt de la fonction waitForRestart.")
				return // Termine la goroutine
			}
		}
	}()
}

// stopWaitForRestart arrête proprement la goroutine créée par waitForRestart
// Elle ferme le canal restartControlChannel pour signaler l'arrêt
func stopWaitForRestart() {
	// Signale l'arrêt de la goroutine
	close(restartControlChannel)

	// Crée un nouveau canal pour redémarrer la goroutine plus tard
	restartControlChannel = make(chan struct{})
}

// restartWaitForRestart redémarre la logique gérée par waitForRestart, elle permet de gerer les deconnexion des clients
// Elle arrête d'abord la goroutine existante, réinitialise les canaux utilisés,
// puis relance la fonction waitForRestart dans une nouvelle goroutine.
func restartWaitForRestart() {
	log.Println("Relance de la fonction waitForRestart...")
	stopWaitForRestart() // Arrête la fonction existante
	close(restartReadyChannel)
	restartReadyChannel = make(chan int, 2)
	waitForRestart() // Relance une nouvelle goroutine
}
