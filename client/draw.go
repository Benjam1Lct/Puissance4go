package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g game) Draw(screen *ebiten.Image) {
	if globalBackgroundImage != nil && g.gameState != introStateLogo && g.gameState != introStateTexte {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(
			float64(globalWidth)/float64(globalBackgroundImage.Bounds().Dx()),
			float64(globalHeight)/float64(globalBackgroundImage.Bounds().Dy()),
		)
		screen.DrawImage(globalBackgroundImage, op)
	}

	switch g.gameState {
	case introStateLogo:
		g.drawIntroLogo(screen)
	case introStateTexte:
		g.drawIntroTexte(screen)
	case titleState:
		g.titleDraw(screen)
	case inputServerState:
		g.inputServerDraw(screen)
	case waitingState:
		g.waitingDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	}

	if g.gameState != introStateLogo && g.gameState != introStateTexte {
		g.topMenu(screen)
		g.drawFullscreenButton(screen)
	}

	if g.debugMode {
		// Créer le texte debug
		debug := fmt.Sprintf(
			"FPS: %0.2f\nTPS: %0.2f",
			ebiten.ActualFPS(),
			ebiten.ActualTPS())

		// Définir la position du texte en bas à gauche
		margin := 10     // Marge de 10 pixels par rapport au bord
		textHeight := 30 // Hauteur approximative pour 2 lignes (ajustez selon la taille réelle du texte)
		textX := margin
		textY := globalHeight - textHeight - margin

		// Dessiner le texte
		ebitenutil.DebugPrintAt(screen, debug, textX, textY)
	}

}

func (g game) titleDraw(screen *ebiten.Image) {
	// Définir la couleur du rectangle (#D9D9D9)
	rectColor := color.NRGBA{R: 217, G: 217, B: 217, A: 255}

	g.topMenuJoueur(screen)

	// Texte principal
	mainText := "Puissance 4"
	mainWidth, mainHeight := getTextDimensions(mainText, firstTitleFont)
	mainX := (globalWidth - mainWidth) / 2
	mainY := (globalHeight / 2) - 60 // Décalage vertical vers le haut
	text.Draw(screen, mainText, firstTitleFont, mainX, mainY, globalTextColorYellow)

	// Sous-titre 1
	subTitle1 := "PROJET DE PROGRAMMATION SYSTEM"
	subTitle1Width, subTitle1Height := getTextDimensions(subTitle1, smallFont)
	subTitle1X := (globalWidth - subTitle1Width) / 2
	subTitle1Y := mainY + mainHeight - 50 // Décalage sous le texte principal

	// Dessiner le rectangle pour le sous-titre 1 derrière le texte
	padding := 10 // Espace autour du texte
	rect1X := subTitle1X - padding
	rect1Y := subTitle1Y - subTitle1Height + 20
	rect1Width := subTitle1Width + 2*padding
	rect1Height := subTitle1Height - 5
	vector.DrawFilledRect(screen, float32(rect1X), float32(rect1Y), float32(rect1Width), float32(rect1Height), rectColor, true)

	// Dessiner le texte pour le sous-titre 1
	text.Draw(screen, subTitle1, smallFont, subTitle1X, subTitle1Y, globalTextColor)

	// Sous-titre 2
	subTitle2 := "ANNEE 2024-2025"
	subTitle2Width, subTitle2Height := getTextDimensions(subTitle2, smallFont)
	subTitle2X := (globalWidth - subTitle2Width) / 2
	subTitle2Y := subTitle1Y + subTitle1Height + 10 // Décalage sous le premier sous-titre

	// Dessiner le rectangle pour le sous-titre 2 derrière le texte
	rect2X := subTitle2X - padding
	rect2Y := subTitle2Y - subTitle2Height + 20
	rect2Width := subTitle2Width + 2*padding
	rect2Height := subTitle2Height - 5
	vector.DrawFilledRect(screen, float32(rect2X), float32(rect2Y), float32(rect2Width), float32(rect2Height), rectColor, true)

	// Dessiner le texte pour le sous-titre 2
	text.Draw(screen, subTitle2, smallFont, subTitle2X, subTitle2Y, globalTextColor)

	// Ajouter un message clignotant en bas de l'écran
	blinkMessage := "Appuyez sur Entrée pour commencer"
	blinkTextWidth, blinkTextHeight := getTextDimensions(blinkMessage, smallFont)
	blinkX := (globalWidth - blinkTextWidth) / 2
	blinkY := globalHeight - blinkTextHeight - 20 // 20 pixels de marge par rapport au bas

	// Dimensions du rectangle
	rectX := blinkX - padding - 20
	rectY := blinkY - padding - 25
	rectWidth := blinkTextWidth + padding*2 + 40
	rectHeight := blinkTextHeight

	// Clignotement basé sur l'état du frame
	if g.stateFrame%60 < 30 { // Clignote toutes les 30 frames (1/2 seconde à 60 FPS)
		vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(rectWidth), float32(rectHeight), globalTextColorGreen, true)
		text.Draw(screen, blinkMessage, smallFont, blinkX, blinkY, globalTextColor)

	}
}

func (g *game) colorSelectDraw(screen *ebiten.Image) {
	// Calculer la position de départ pour centrer la grille
	gridWidth := globalNumColorCol * globalTileSize
	gridHeight := globalNumColorLine * globalTileSize
	startX := (globalWidth - gridWidth) / 2
	startY := (globalHeight-gridHeight)/2 + 40

	// Dessiner le texte titre au-dessus de la grille
	title := "Choisissez la couleur de vos pions"
	titleWidth, titleHeight := getTextDimensions(title, firstTitleSmallerFont)
	titleX := (globalWidth - titleWidth) / 2
	titleY := startY - titleHeight + 20 // 20 pixels au-dessus de la grille
	text.Draw(screen, title, firstTitleSmallerFont, titleX, titleY, globalTextColorYellow)

	line := 0
	col := 0
	for numColor := 0; numColor < globalNumColor; numColor++ {
		// Calculer la position des cercles
		xPos := startX + col*globalTileSize
		yPos := startY + line*globalTileSize

		centerX := float32(xPos + globalTileSize/2)
		centerY := float32(yPos + globalTileSize/2)

		// Dessiner un cercle vert foncé si la couleur est validée pour p1
		if g.p1ColorValidate != -1 && g.p1ColorValidate == numColor {
			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				(float32(globalTileSize) / 2),
				globalValidatorColor, // Couleur vert foncé
				true,
			)
		}

		// Dessiner un cercle vert foncé si la couleur est validée pour p2
		if g.p2Color != -1 && g.p2Color == numColor {
			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				(float32(globalTileSize) / 2),
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
				(float32(globalTileSize)/2)+5,
				globalP2SelectColor,
				true,
			)
			// Cercle intérieur pour le joueur 1
			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				(float32(globalTileSize)/2)-1,
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
					(float32(globalTileSize) / 2),
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
					(float32(globalTileSize) / 2),
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
			(float32(globalTileSize)/2)-float32(globalCircleMargin)-2,
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
		g.errorMessageDisplay(screen, g.errorMessage)
	}
}

// Affichage des graphismes durant le jeu.
func (g game) playDraw(screen *ebiten.Image) {
	// Calculer la position et les dimensions de la grille
	gridWidth := globalTileSize * globalNumTilesX
	gridHeight := globalTileSize * globalNumTilesY
	startX := (globalWidth - gridWidth) / 2
	startY := (globalHeight-gridHeight)/2 + 50
	cornerRadius := float32(25) // Rayon des coins arrondis pour le contour (un peu plus grand)

	// Définir la couleur du contour en fonction de l'état du joueur
	var borderColor color.Color
	if g.turn == p1Turn {
		borderColor = globalTextColorGreen // Vert pour le joueur 1
	} else {
		borderColor = globalTextRed // Rouge pour le joueur 2
	}

	// Dessiner le contour autour de la grille
	drawRoundedRectangle(
		screen,
		float32(startX)-7.5, // Décaler pour entourer la grille
		float32(startY)-7.5,
		float32(gridWidth)+15, // Ajouter de l'espace autour de la grille
		float32(gridHeight)+15,
		cornerRadius,
		borderColor,
	)

	// Dessiner la grille
	g.drawGrid(screen)

	// Afficher le pion du joueur actif (au-dessus de la grille)
	pionX := float32(startX + globalTileSize/2 + g.tokenPosition*globalTileSize)
	pionY := float32(startY-globalTileSize/2) - 20 // Juste au-dessus de la grille
	vector.DrawFilledCircle(
		screen,
		pionX,
		pionY,
		float32(globalTileSize/2-globalCircleMargin),
		globalTokenColors[g.p1Color],
		true,
	)

	// Afficher les messages d'erreur, le cas échéant
	if !g.restartOk {
		g.errorMessageDisplay(screen, "En attente de l'autre joueur")
	}
}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {
	// Centrer la grille
	g.drawGrid(offScreenImage)
	g.topMenuRejouer(screen)

	// Appliquer un effet d'atténuation
	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.1)
	screen.DrawImage(offScreenImage, options)

	// Afficher le message de résultat
	message := "Il y a Egalite"
	if g.result == p1wins {
		message = "Vous avez Gagne !"
	} else if g.result == p2wins {
		message = "Vous avez Perdu"
	}
	textWidth, _ := getTextDimensions(message, firstTitleSmallFont)
	textX := (globalWidth - textWidth) / 2
	textY := globalHeight/2 - 50
	text.Draw(screen, message, firstTitleSmallFont, textX, textY, globalTextColorYellow)

	if (g.stateFrame/30)%2 == 0 { // Clignotement toutes les 30 frames (0.5s à 60 FPS)
		blinkMessage := "Appuyez sur Entrée pour rejouer"
		blinkTextWidth, blinkTextHeight := getTextDimensions(blinkMessage, smallFont)
		blinkTextX := (globalWidth - blinkTextWidth) / 2
		blinkTextY := globalHeight - 100 // Position en bas de l'écran

		padding := 10
		rectX := blinkTextX - padding - 20
		rectY := blinkTextY - padding - 25
		rectWidth := blinkTextWidth + padding*2 + 40
		rectHeight := blinkTextHeight

		vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(rectWidth), float32(rectHeight), globalTextColorGreen, true)
		text.Draw(screen, blinkMessage, smallFont, blinkTextX, blinkTextY, globalTextColor)
	}

	// Afficher le message d'attente pour rematch au-dessus
	if g.messageWaitRematch != "" {
		msg := "L'adversaire est en attente de rematch !"
		textWidth, blinkTextHeight := getTextDimensions(msg, smallFont)
		textX := (globalWidth - textWidth) / 2
		textY := globalHeight/2 + 20

		padding := 10
		rectX := textX - padding - 20
		rectY := textY - padding - 25
		rectWidth := textWidth + padding*2 + 40
		rectHeight := blinkTextHeight

		vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(rectWidth), float32(rectHeight), globalTextRed, true)
		text.Draw(screen, msg, smallFont, textX, textY, globalTextColorBright)
	}

}

func (g game) drawGrid(screen *ebiten.Image) {
	// Calculer la position et les dimensions de la grille
	gridWidth := globalTileSize * globalNumTilesX
	gridHeight := globalTileSize * globalNumTilesY
	startX := (globalWidth - gridWidth) / 2
	startY := (globalHeight-gridHeight)/2 + 50
	cornerRadius := float32(20) // Rayon des coins arrondis

	// Dessiner la grille noire avec des coins arrondis
	drawRoundedRectangle(
		screen,
		float32(startX),
		float32(startY),
		float32(gridWidth),
		float32(gridHeight),
		cornerRadius,
		color.Black, // Couleur noire pour la grille
	)

	// Dessiner les cercles pour chaque cellule
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

			centerX := float32(startX + globalTileSize/2 + x*globalTileSize)
			centerY := float32(startY + globalTileSize/2 + y*globalTileSize)

			vector.DrawFilledCircle(
				screen,
				centerX,
				centerY,
				float32(globalTileSize/2-globalCircleMargin),
				tileColor,
				true,
			)
		}
	}
}
