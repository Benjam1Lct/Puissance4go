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
	clients          = make(map[int]net.Conn)
	readyPlayers     = make(map[int]bool)
	clientMux        sync.Mutex
	playerColors         = make(map[int]int)
	firstPlayer      int = -1
	playersRestarted     = make(map[int]bool)
)

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

func sendJSONMessage(conn net.Conn, msg Message) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation JSON : %w", err)
	}

	_, err = conn.Write(append(jsonData, '\n')) // Ajout de '\n' pour délimiter les messages
	return err
}

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
