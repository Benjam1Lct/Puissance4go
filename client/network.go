package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/text"
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
				err := sendJSONMessage(g.conn, "ready", nil)
				if err != nil {
					return
				} // Informer le serveur que le client est prêt
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
						finished, result, posWinnerCheck := g.checkGameEnd(int(x), int(y))
						if posWinnerCheck != nil {
							g.posWinner = posWinnerCheck
							log.Println("posWinnerCheck", g.posWinner)
						}
						if finished {
							g.result = result
							if g.result == p1wins {
								g.nbPartieWin++
								g.turn = p2Turn
							} else if g.result == p2wins {
								g.nbPartieAdversaireWin++
								g.turn = p1Turn
							} else {
								g.nbPartieWin++
								g.nbPartieAdversaireWin++
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
	case "chat":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if text, ok := payload["text"].(string); ok {
				if id, ok := payload["id"].(float64); ok {
					if id == float64(g.playerID) {
						g.addChatMessage(text, "You:")
					} else {
						g.addChatMessage(text, "Other:")
						g.chatNewMessage = true
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
				if g.nbJoueurConnecte < 2 {
					g.nbJoueurConnecte++
				}
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
				g.gameState = shifumiState

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
	case "token_update":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if position, ok := payload["position"].(float64); ok {
				g.adversaryTokenPosition = int(position)
				log.Printf("Position adversaire: %d\n", g.adversaryTokenPosition)
			}
		}
	case "sent_history":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			// Re-sérialiser le payload en JSON
			jsonData, err := json.Marshal(payload)
			if err != nil {
				log.Printf("Erreur lors de la re-sérialisation de l'historique : %v\n", err)
				return
			}

			// Désérialiser dans une map[string]Coordinate
			tempHistory := make(map[string]Coordinate)
			if err := json.Unmarshal(jsonData, &tempHistory); err != nil {
				log.Printf("Erreur lors du décodage de l'historique : %v\n", err)
				return
			}

			// Convertir les clés string en int et stocker dans votre variable globale
			for keyStr, coord := range tempHistory {
				var keyInt int
				_, err := fmt.Sscanf(keyStr, "%d", &keyInt)
				if err != nil {
					log.Printf("Erreur de conversion de clé '%s' en int : %v\n", keyStr, err)
					continue
				}
				history[keyInt] = coord
			}

			log.Println("Historique mis à jour côté client :", history)
		} else {
			log.Println("Payload invalide pour 'sent_history'")
		}
	case "other_disconnected":
		log.Println("L'autre joueur s'est deconnecté")
		g.disconnectClient()
	case "selected":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if selected, ok := payload["selected"].(string); ok {
				g.selected = selected
				log.Printf("Sélection reçue : %s\n", selected)
			}
		}
	case "shifumi_result":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if result, ok := payload["result"].(string); ok {
				// Stocker le résultat
				g.shifumiResult = result
				g.showShifumiResult = true
				g.shifumiResultTimer = 120 // environ 2 secondes à 60 FPS

				// Gérer les sélections des joueurs
				if p1, ok := payload["player1"].(map[string]interface{}); ok {
					if p1Selection, ok := p1["selection"].(string); ok {
						if g.playerID == 1 {
							g.selected = p1Selection
						} else {
							g.adversaryChoice = p1Selection
						}
					}
				}
				if p2, ok := payload["player2"].(map[string]interface{}); ok {
					if p2Selection, ok := p2["selection"].(string); ok {
						if g.playerID == 2 {
							g.selected = p2Selection
						} else {
							g.adversaryChoice = p2Selection
						}
					}
				}

				// Déterminer le résultat local
				g.determineShifumiWinner()

				if result == "draw" {
					// En cas d'égalité, rester dans l'état shifumi pour rejouer
					g.gameState = shifumiState
					g.selected = "" // Réinitialiser la sélection pour le prochain tour
					g.adversaryChoice = ""
				}
				log.Printf("Résultat Shifumi - Joueur: %s, Adversaire: %s, Résultat: %s\n",
					g.selected, g.adversaryChoice, g.shifumiResult)
			}
		}
	case "shifumi_complete":
		if payload, ok := msg.Payload.(map[string]interface{}); ok {
			if winnerID, ok := payload["winner"].(float64); ok {
				if int(winnerID) == g.playerID {
					g.turn = p1Turn
					g.firstPlayer = p1Turn
				} else {
					g.turn = p2Turn
					g.firstPlayer = p2Turn
				}
				// Passer à l'état de sélection des couleurs

				g.gameReady = true
				g.gameState = playState
				g.connectionMessage = "Shifumi terminé ! Sélectionnez votre couleur."
				log.Printf("Le joueur %d a gagné le shifumi et commence la partie\n", int(winnerID))
			}
		}
	default:
		log.Printf("Type de message inconnu : %s\n", msg.Type)

	}
}

func (g *game) disconnectClient() {
	err := g.conn.Close()
	if err != nil {
		return
	}

	g.isReset = true

	//init
	g.chatNewMessage = false
	g.stateFrame = 0
	g.restartOk = true
	g.p2Color = -1
	g.p1ColorValidate = -1
	g.playerID = -1
	g.mouseReleased = true

	//reset all
	g.playerID = 0
	g.resetGrid()
	g.p1Color = 0
	g.p2CursorColor = 0
	g.turn = noToken
	g.firstPlayer = noToken
	g.tokenPosition = 0
	g.result = 0
	g.serverAddress = ""
	g.serverReady = false
	g.connectionMessage = ""
	g.errorConnection = ""
	g.gameReady = false
	g.errorMessage = ""
	g.nbJoueurConnecte = 0
	g.messageWaitRematch = ""
	g.mouseReleased = true
	g.adversaryTokenPosition = g.tokenPosition
	g.nbPartieWin = 0
	g.nbPartieAdversaireWin = 0
	g.chatMessages = []string{} // Historique des messages de chat
	g.chatInput = ""
	g.chatIsFocus = false
	g.chatNewMessage = false
}

func (g *game) addChatMessage(message string, id string) {
	// Récupérer l'horodatage actuel
	timestamp := time.Now().Format("2006/01/02 15:04:05")

	// Ajouter la date et l'heure au message
	messageWithTimestamp := fmt.Sprintf("%s %s %s", timestamp, id, message)

	// Découper le message pour qu'il respecte la largeur définie
	messageWithTimestamp = truncateMessage(messageWithTimestamp, messageWidth+350)

	// Ajouter le message avec l'horodatage à la liste
	g.chatMessages = append(g.chatMessages, messageWithTimestamp)

	// Limiter à 10 messages
	g.limitChatHeight()

}

func truncateMessage(message string, maxWidth int) string {
	words := strings.Split(message, " ")
	var lines []string
	currentLine := ""

	for _, word := range words {
		testLine := strings.TrimSpace(currentLine + " " + word)
		if text.BoundString(smallFont, testLine).Dx() > maxWidth {
			// Si currentLine est vide, le mot seul dépasse la largeur
			if currentLine == "" {
				lines = append(lines, word) // Ajouter directement le mot
			} else {
				lines = append(lines, strings.TrimSpace(currentLine))
				currentLine = word
			}
		} else {
			currentLine = testLine
		}
	}

	// Ajouter la dernière ligne uniquement si elle n'est pas vide
	if strings.TrimSpace(currentLine) != "" {
		lines = append(lines, strings.TrimSpace(currentLine))
	}

	return strings.Join(lines, "\n")
}

// Limiter la hauteur totale des messages dans la zone de texte
func (g *game) limitChatHeight() {
	maxHeight := 200 // Hauteur maximale autorisée pour la zone de texte
	totalHeight := 0

	// Calculer la hauteur totale des messages
	for _, message := range g.chatMessages {
		totalHeight += calculateMessageHeight(message, messageWidth+250)
	}

	// Supprimer les messages les plus anciens jusqu'à ce que la hauteur soit inférieure à maxHeight
	for totalHeight > maxHeight && len(g.chatMessages) > 0 {
		// Supprimer le message le plus ancien
		oldMessage := g.chatMessages[0]
		totalHeight -= calculateMessageHeight(oldMessage, messageWidth+250)
		g.chatMessages = g.chatMessages[1:]
	}
}

// Calculer la hauteur d'un message en fonction de son contenu et de la largeur définie
func calculateMessageHeight(message string, maxWidth int) int {
	lines := strings.Split(message, "\n")
	return len(lines) * 20 // Chaque ligne fait 20 pixels de hauteur (ou ajustez selon votre font)
}

func sendChatMessage(conn net.Conn, text string) {
	payload := ChatMessage{Text: text}
	err := sendJSONMessage(conn, "chat", payload)
	if err != nil {
		log.Printf("Erreur lors de l'envoi du message de chat : %v\n", err)
	}
}

func requestHistory(conn net.Conn) error {
	return sendJSONMessage(conn, "require_history", nil)
}

func sendCursorUpdateToServer(conn net.Conn, color int) {
	payload := map[string]int{"color": color}
	err := sendJSONMessage(conn, "cursor_update", payload)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de la position du curseur : %v\n", err)
	}
}

func sendTokenUpdateToServer(conn net.Conn, position int) {
	// Créer un payload contenant la position du curseur
	payload := map[string]int{
		"position": position,
	}

	// Envoyer un message JSON de type "cursor_update" avec la position
	err := sendJSONMessage(conn, "token_update", payload)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de la mise à jour du curseur : %v\n", err)
	}
	log.Println("token envoyé au serveur")
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

// SelectedPayload représente les données de sélection à envoyer au serveur
type SelectedPayload struct {
	Selected string `json:"selected"`
}

func sendSelectedToServer(conn net.Conn, selected string) {
	payload := SelectedPayload{Selected: selected}
	err := sendJSONMessage(conn, "selected", payload)
	if err != nil {
		log.Printf("Erreur lors de l'envoi de la sélection : %v\n", err)
	}
	log.Printf("Sélection envoyée : %s\n", selected)
}
