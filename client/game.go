package main

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
	"log"
	"net"
	"os"
)

// Structure de données pour représenter l'état courant du jeu.
type game struct {
	playerID           int
	gameState          int
	stateFrame         int
	grid               [globalNumTilesX][globalNumTilesY]int
	p1Color            int
	p1ColorValidate    int
	p2Color            int
	p2CursorColor      int
	turn               int
	firstPlayer        int
	tokenPosition      int
	result             int
	serverAddress      string
	serverReady        bool
	connectionMessage  string
	errorConnection    string
	gameReady          bool
	conn               net.Conn
	errorMessage       string
	restartOk          bool
	nbJoueurConnecte   int
	messageWaitRematch string
	stateFrameIntro    int
	mouseReleased      bool
	debugMode          bool
	audioContext       *audio.Context
	audioPlayer        *audio.Player
	audioFile          *os.File
}

// Constantes pour représenter la séquence de jeu actuelle (écran titre,
// écran de sélection des couleurs, jeu, écran de résultats).
const (
	introStateLogo int = iota
	introStateTexte
	titleState
	colorSelectState
	playState
	resultState
	inputServerState
	waitingState
	waitingColorSelect
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
	// Réinitialiser la grille
	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			g.grid[x][y] = noToken
		}
	}

	// Réinitialiser les variables du jeu
	g.stateFrame = 0    // Réinitialiser le compteur d'états
	g.tokenPosition = 0 // Réinitialiser la position du jeton
	g.result = noToken  // Aucun gagnant pour la nouvelle partie

	log.Printf("Grille réinitialisée. Votre tour : %v\n", g.turn == p1Turn)
}
