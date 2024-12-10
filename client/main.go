package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

// Création, paramétrage et lancement du jeu.
func main() {

	initResolution(true)

	// Configurer la fenêtre pour utiliser cette résolution
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowSize(globalWidth, globalHeight)
	ebiten.SetWindowTitle("Programmation système : projet puissance 4")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	// Passer en mode plein écran

	g := game{}
	g.gameState = introStateLogo
	g.stateFrame = 0
	g.restartOk = true
	g.p2Color = -1
	g.p1ColorValidate = -1
	g.playerID = -1

	// Initialiser les ressources nécessaires
	g.initAudio()
	defer g.closeAudio()
	initFonts()
	initImage()
	initOffScreen()
	initBackground()

	// Lancer le jeu
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

}
