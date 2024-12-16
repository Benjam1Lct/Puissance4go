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

	// Initialiser le jeu
	g := game{}
	g.initGame()

	// Initialiser les ressources nécessaires
	initFonts()
	initImage()
	initOffScreen()
	initBackground()

	// Lancer le jeu
	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}

}
