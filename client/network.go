package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"net"
	"strings"
	"time"
)

// Connecter le joueur au serveur
func connectToServer(g *game) {
	// Vérifier si l'adresse contient déjà un port (si elle contient ":")
	if !strings.Contains(g.serverAddress, ":") {
		// Ajouter le port par défaut :8080
		g.serverAddress += ":8080"
	}

	conn, err := net.Dial("tcp", g.serverAddress)
	if err != nil {
		log.Println("Erreur de connexion au serveur :", err)
		g.errorConnection = "Erreur : Adresse incorrecte."
		g.gameState = inputServerState // Retour à l'état de saisie
		g.serverAddress = ""
		return
	}

	g.conn = conn
	g.connectionMessage = "Connecté au serveur. En attente d'autres joueurs..."
	g.nbJoueurConnecte++
	log.Println(g.connectionMessage)

	// Écouter les messages du serveur
	go listenToServer(conn, g)
}

func listenToServer(conn net.Conn, g *game) {
	reader := bufio.NewReader(conn)

	// Goroutine pour la surveillance de l'état et envoi automatique de messages
	go func() {
		messageSent := false
		for {
			if g.gameState == waitingColorSelect && !messageSent {
				err := sendJSONMessage(conn, "choix effectué", nil) // Envoyer un message initial si nécessaire
				if err != nil {
					log.Printf("Erreur lors de l'envoi automatique du message : %v\n", err)
				}
				messageSent = true
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Boucle principale pour lire et traiter les messages du serveur
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Erreur de lecture : %v\n", err)
			g.connectionMessage = "Erreur de communication. Entrez une nouvelle adresse."
			g.gameState = inputServerState
			return
		}

		message = strings.TrimSpace(message)

		// Désérialiser le message JSON
		var msg Message
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			log.Printf("Erreur de décodage JSON : %v\n", err)
			continue
		}

		// Gérer le message JSON
		handleServerMessage(msg, g)
	}
}

func handleServerMessage(msg Message, g *game) {
	switch msg.Type {
	case "id":
		// Récupérer l'ID du joueur
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if id, ok := payload["id"].(float64); ok {
				g.playerID = int(id)
				log.Printf("ID reçu : %d\n", g.playerID)
				sendJSONMessage(g.conn, "ready", nil) // Informer le serveur que le client est prêt
			}
		}
	case "color":
		// Récupérer la couleur de l'autre joueur
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if color, ok := payload["color"].(float64); ok {
				g.p2Color = int(color)
				log.Printf("Couleur de l'autre joueur reçue : %d\n", g.p2Color)
			}
		}
	case "move":
		// Mettre à jour la grille avec le mouvement de l'autre joueur
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if x, ok := payload["x"].(float64); ok {
				if y, ok := payload["y"].(float64); ok {
					log.Printf("Mouvement reçu : (%d, %d)\n", int(x), int(y))
					updated, _ := g.updateGrid(p2Token, int(x))
					if updated {
						finished, result := g.checkGameEnd(int(x), int(y))
						if finished {
							g.result = result
							if g.result == p1wins {
								g.turn = p2Turn
							} else if g.result == p2wins {
								g.turn = p1Turn
							} else {
								g.turn = g.firstPlayer
							}
							g.gameState = resultState
							g.restartOk = false
						} else {
							g.turn = p1Turn // C'est maintenant au tour du joueur 1
						}
					} else {
						log.Println("Erreur : Mise à jour de la grille échouée.")
					}
				}
			}
		}
	case "rematch_waiting":
		// Afficher le message d'attente pour le rematch
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if message, ok := payload["message"].(string); ok {
				g.messageWaitRematch = message
				log.Println(message)
			}
		}
	case "restart_ok":
		// Indiquer que le jeu peut redémarrer
		g.restartOk = true
		g.messageWaitRematch = ""
		log.Println("Le jeu peut redémarrer.")
	case "ready":
		// Indiquer que le serveur est prêt
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if message, ok := payload["message"].(string); ok {
				g.connectionMessage = message
				g.serverReady = true
				g.gameState = colorSelectState
				g.nbJoueurConnecte++
				log.Println(message)
			}
		}
	case "color_select_complete":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if starterID, ok := payload["firstPlayer"].(float64); ok {
				// Déterminer qui commence
				if int(starterID) == g.playerID {
					g.turn = p1Turn // Ce client commence
					g.firstPlayer = p1Turn
				} else {
					g.turn = p2Turn // L'autre joueur commence
					g.firstPlayer = p2Turn
				}

				// Mettre à jour l'état du jeu
				g.connectionMessage = "La partie commence. Préparez-vous !"
				g.gameReady = true
				g.gameState = playState

				log.Printf("Le joueur %d commence la partie. Votre tour : %v\n", int(starterID), g.turn == p1Turn)
			} else {
				log.Printf("Erreur : ID du premier joueur manquant dans le payload.")
			}
		} else {
			log.Printf("Erreur : Structure de payload invalide pour 'color_select_complete'.")
		}
	case "cursor_update":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if color, ok := payload["color"].(float64); ok {
				g.p2CursorColor = int(color)
			}
		}
	default:
		log.Printf("Type de message inconnu : %s\n", msg.Type)
	}
}

func sendCursorUpdateToServer(conn net.Conn, color int) {
	payload := map[string]int{"color": color}
	err := sendJSONMessage(conn, "cursor_update", payload)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de la position du curseur : %v\n", err)
	}
}

func (g *game) inputServerUpdate() bool {
	// Réinitialiser l'adresse si une erreur est survenue
	if g.connectionMessage == "Erreur : Adresse incorrecte. Entrez une nouvelle adresse." {
		g.serverAddress = ""
		g.connectionMessage = ""
	}

	// Ajouter les caractères saisis
	for _, char := range ebiten.AppendInputChars(nil) {
		g.serverAddress += string(char)
	}
	// Supprimer le dernier caractère (Backspace)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(g.serverAddress) > 0 {
		g.serverAddress = g.serverAddress[:len(g.serverAddress)-1]
	}
	// Valider l'adresse (Entrée)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.serverAddress == "" {
			g.serverAddress = "localhost:8080"
		}
		return true
	}
	return false
}

func sendColorToServer(conn net.Conn, color int) {
	payload := ColorPayload{Color: color}
	err := sendJSONMessage(conn, "color", payload)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de la couleur : %v\n", err)
	}
	log.Printf("Couleur envoyée au serveur : %d\n", color)
}

func sendMoveToServer(conn net.Conn, x int, y int) {
	payload := MovePayload{X: x, Y: y}
	err := sendJSONMessage(conn, "move", payload)
	if err != nil {
		log.Printf("Erreur lors de l'envoi du mouvement : %v\n", err)
	}
	log.Printf("Mouvement envoyé : (%d, %d)\n", x, y)
}

func endGame(g *game) {
	if g.conn != nil {
		err := sendJSONMessage(g.conn, "end", nil)
		if err != nil {
			log.Printf("Erreur lors de l'envoi de 'end' : %v\n", err)
		} else {
			log.Println("Message 'end' envoyé au serveur.")
		}
	}
}

func sendJSONMessage(conn net.Conn, messageType string, payload interface{}) error {
	message := Message{
		Type:    messageType,
		Payload: payload,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("erreur de sérialisation JSON : %w", err)
	}

	_, err = conn.Write(append(jsonData, '\n')) // Ajout de '\n' pour marquer la fin du message
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi du message JSON : %w", err)
	}

	return nil
}
