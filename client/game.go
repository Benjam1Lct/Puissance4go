package main

import (
	"log"
	"net"
)

// Structure de données pour représenter l'état courant du jeu.
type game struct {
	playerID               int
	gameState              int
	stateFrame             int
	grid                   [globalNumTilesX][globalNumTilesY]int
	p1Color                int
	p1ColorValidate        int
	p2Color                int
	p2CursorColor          int
	turn                   int
	firstPlayer            int
	tokenPosition          int
	result                 int
	serverAddress          string
	serverReady            bool
	connectionMessage      string
	errorConnection        string
	gameReady              bool
	conn                   net.Conn
	errorMessage           string
	restartOk              bool
	nbJoueurConnecte       int
	messageWaitRematch     string
	stateFrameIntro        int
	mouseReleased          bool
	debugMode              bool
	adversaryTokenPosition int
	isMuted                bool
	nbPartieWin            int
	nbPartieAdversaireWin  int
	nbBackground           int
	nbBackgroundTheme      int
	posWinner              [][2]int
	lastUpdateTime         int
	blinking               bool
	chatMessages           []string // Historique des messages de chat
	chatInput              string
	chatIsFocus            bool
	chatNewMessage         bool
	isReset                bool
	selected               string
	adversaryChoice        string    // Choix de l'adversaire dans le shifumi
	shifumiResult         string    // Résultat du shifumi (Gagné/Perdu/Égalité)
	showShifumiResult     bool      // Indique si on doit afficher le résultat
	shifumiResultTimer    int       // Timer pour l'affichage du résultat
}

// Constantes pour représenter la séquence de jeu actuelle (écran titre,
// écran de sélection des couleurs, jeu, écran de résultats).
const (
	introStateLogo int = iota
	introStateTexte
	titleState
	themeState
	colorSelectState
	playState
	resultState
	replayState
	inputServerState
	waitingState
	waitingColorSelect
	shifumiState
)

// Constantes pour représenter les pions dans la grille de puissance 4
// (absence de pion, pion du joueur 1, pion du joueur 2).
const (
	noToken int = iota
	p1Token
	p2Token
)

// Constantes pour représenter le tour de jeu (joueur 1 ou joueur 2).
const (
	p1Turn int = iota
	p2Turn
)

// Constantes pour représenter le résultat d'une partie (égalité si
// la grille a été remplie sans qu'un joueur n'ait gagné, joueur 1
// gagnant ou joueur 2 gagnant).
const (
	equality int = iota
	p1wins
	p2wins
)

// Remise à 0 du jeu pour recommencer une partie. Le joueur qui a
// perdu la dernière partie commence.
func (g *game) reset() {
	// Informer le serveur que la partie est terminer et que l'on est pret a rejouer
	if g.conn != nil {
		err := sendJSONMessage(g.conn, "restartReady", nil)
		if err != nil {
			log.Printf("Erreur lors de l'envoi de 'end' : %v\n", err)
		} else {
			log.Println("Message 'end' envoyé au serveur.")
		}
	}

	// Réinitialiser la grille
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	} // Réinitialiser la grille
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	}

	// Réinitialiser les variables du jeu
	g.stateFrame = 0    // Réinitialiser le compteur d'états
	g.tokenPosition = 0 // Réinitialiser la position du jeton
	g.result = noToken  // Aucun gagnant pour la nouvelle partie
	g.adversaryTokenPosition = 0

	log.Printf("Grille réinitialisée. Votre tour : %v\n", g.turn == p1Turn)
}

func (g *game) resetGrid() {

	// Réinitialiser la grille
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	} // Réinitialiser la grille
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	}
}

// determineShifumiWinner détermine le résultat du shifumi
func (g *game) determineShifumiWinner() {
	if g.selected == g.adversaryChoice {
		g.shifumiResult = "Égalité"
		return
	}

	switch g.selected {
	case "Pierre":
		if g.adversaryChoice == "Ciseaux" {
			g.shifumiResult = "Gagné"
		} else {
			g.shifumiResult = "Perdu"
		}
	case "Papier":
		if g.adversaryChoice == "Pierre" {
			g.shifumiResult = "Gagné"
		} else {
			g.shifumiResult = "Perdu"
		}
	case "Ciseaux":
		if g.adversaryChoice == "Papier" {
			g.shifumiResult = "Gagné"
		} else {
			g.shifumiResult = "Perdu"
		}
	}
}
