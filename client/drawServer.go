package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
)

var (
	redColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Rouge opaque
)

func (g game) inputServerDraw(screen *ebiten.Image) {
	text.Draw(screen, "Entrez l'adresse du serveur :", smallFont, 100, 100, globalTextColor)

	// Afficher l'adresse actuelle ou une suggestion par défaut
	addressToShow := g.serverAddress
	if addressToShow == "" {
		addressToShow = "localhost:8080 (par défaut)"
	}
	text.Draw(screen, addressToShow, smallFont, 100, 140, globalTextColor)

	// Afficher le message de connexion ou d'erreur
	if g.errorConnection != "" {
		text.Draw(screen, g.errorConnection, smallFont, 100, 200, redColor)
	}
}

func (g game) waitingDraw(screen *ebiten.Image) {
	text.Draw(screen, "En attente de l'autre joueur...", smallFont, 100, 100, globalTextColor)
	text.Draw(screen, fmt.Sprintf("Joueurs connectés : %d / 2", g.nbJoueurConnecte), smallFont, 100, 200, globalTextColor)
}
func (g game) waitingColorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "En attente de selection de la couleur", smallFont, 100, 100, globalTextColor)
	text.Draw(screen, fmt.Sprintf("Joueurs connectés : %d / 2", g.nbJoueurConnecte), smallFont, 100, 200, globalTextColor)

}
