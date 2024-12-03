package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g game) Draw(screen *ebiten.Image) {
	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	case titleState:
		g.titleDraw(screen)
	case inputServerState:
		g.inputServerDraw(screen)
	case waitingState:
		g.waitingDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case waitingColorSelect:
		g.waitingColorSelectDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}
}

// Affichage des graphismes de l'écran titre.
func (g game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)

	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Appuyez sur entrée", smallFont, 210, 500, globalTextColor)
	}
}
func (g *game) colorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "Quelle couleur pour vos pions ?", smallFont, 110, 80, globalTextColor)

	line := 0
	col := 0
	for numColor := 0; numColor < globalNumColor; numColor++ {
		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line

		centerX := float32(globalTileSize/2 + xPos*globalTileSize)
		centerY := float32(globalTileSize/2 + yPos*globalTileSize)

		// Dessiner un cercle vert foncé si la couleur est validée pour p1
		if g.p1ColorValidate != -1 && g.p1ColorValidate == numColor {
			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				(globalTileSize / 2), // Cercle légèrement plus grand
				globalValidatorColor, // Couleur vert foncé
				true,
			)
		}

		// Dessiner un cercle vert foncé si la couleur est validée pour p1
		if g.p2Color != -1 && g.p2Color == numColor {
			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				(globalTileSize / 2),     // Cercle légèrement plus grand
				globalValidatorDarkColor, // Couleur vert foncé
				true,
			)
		}

		// Chevauchement : même position pour p1 et p2
		if numColor == g.p1Color && numColor == g.p2CursorColor {
			// Cercle extérieur pour le joueur 2
			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				(globalTileSize/2)+5, // Cercle légèrement plus grand
				globalP2SelectColor,
				true,
			)
			// Cercle intérieur pour le joueur 1
			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				(globalTileSize/2)-1, // Cercle légèrement plus petit
				globalSelectColor,
				true,
			)
		} else {
			// Dessiner le cercle pour le joueur 1
			if numColor == g.p1Color {
				vector.DrawFilledCircle(
					screen,
					centerX,
					centerY,
					(globalTileSize / 2),
					globalSelectColor,
					true,
				)
			}

			// Dessiner le cercle pour le joueur 2
			if numColor == g.p2CursorColor {
				vector.DrawFilledCircle(
					screen,
					centerX,
					centerY,
					(globalTileSize / 2),
					globalP2SelectColor,
					true,
				)
			}
		}

		// Dessiner le cercle pour représenter la couleur de fond
		vector.DrawFilledCircle(
			screen,
			centerX,
			centerY,
			(globalTileSize/2)-globalCircleMargin-2,
			globalTokenColors[numColor],
			true,
		)

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}

	// Afficher les messages d'erreur, le cas échéant
	if g.errorMessage != "" {
		text.Draw(screen, g.errorMessage, smallFont, 20, 40, color.RGBA{255, 0, 0, 255}) // En rouge
	}
}

// Affichage des graphismes durant le jeu.
func (g game) playDraw(screen *ebiten.Image) {
	g.drawGrid(screen)

	// Définir la couleur du contour en fonction de l'état du joueur
	var borderColor color.Color
	if g.turn == p1Turn { // Ajoutez un champ `isPlayerTurn` dans la structure `game` pour suivre l'état
		borderColor = color.RGBA{0, 255, 0, 255} // Vert pour le tour du joueur
	} else {
		borderColor = color.RGBA{255, 0, 0, 255} // Rouge sinon
	}

	// Dessiner le contour autour de l'écran
	vector.DrawFilledRect(screen, 0, 0, float32(screen.Bounds().Dx()), 10, borderColor, true)                                // Haut
	vector.DrawFilledRect(screen, 0, float32(screen.Bounds().Dy()-10), float32(screen.Bounds().Dx()), 10, borderColor, true) // Bas
	vector.DrawFilledRect(screen, 0, 0, 10, float32(screen.Bounds().Dy()), borderColor, true)                                // Gauche
	vector.DrawFilledRect(screen, float32(screen.Bounds().Dx()-10), 0, 10, float32(screen.Bounds().Dy()), borderColor, true) // Droite

	if !g.restartOk {
		// Afficher le message d'attente par-dessus
		message := "En attente de l'autre joueur pour commencer..."
		text.Draw(screen, message, smallFont, 10, 50, globalTextColor)
		return // Ne pas dessiner les autres éléments tant que restartOk n'est pas validé
	}

	// Si restartOk est vrai, continuer à dessiner les éléments du jeu
	vector.DrawFilledCircle(
		screen,
		float32(globalTileSize/2+g.tokenPosition*globalTileSize),
		float32(globalTileSize/2),
		globalTileSize/2-globalCircleMargin,
		globalTokenColors[g.p1Color],
		true,
	)
}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Égalité"
	if g.result == p1wins {
		message = "Gagné !"
	} else if g.result == p2wins {
		message = "Perdu…"
	}
	text.Draw(screen, message, smallFont, 300, 350, globalTextColor)

	if g.messageWaitRematch != "" {
		text.Draw(screen, "L'autre joueur est en attente de rematch...", smallFont, 50, 400, redColor)
	}

}

// Affichage de la grille de puissance 4, incluant les pions déjà joués.
func (g game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {
			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.p2Color]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
