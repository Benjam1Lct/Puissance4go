package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g game) Draw(screen *ebiten.Image) {
	if g.gameState != introStateLogo && g.gameState != introStateTexte {
		var image *ebiten.Image
		if g.nbBackground == 0 {
			image = background1
		} else if g.nbBackground == 1 {
			image = background3
		} else if g.nbBackground == 2 {
			image = background2
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(
			float64(globalWidth)/float64(image.Bounds().Dx()),
			float64(globalHeight)/float64(image.Bounds().Dy()),
		)
		screen.DrawImage(image, op)
		if g.chatIsFocus {
			g.errorMessageDisplay(screen, "Fermez le chat pour reprendre la partie")
		}
	}

	if g.isReset {
		g.errorMessageDisplay(screen, "L'autre joueur s'est déconnecté")
	}

	if g.debugMode {
		// Exemple de changement d'état via une entrée utilisateur
		if ebiten.IsKeyPressed(ebiten.KeyD) {
			g.gameState = shifumiState // Changer vers l'état titre
		} else if ebiten.IsKeyPressed(ebiten.Key2) {
			g.gameState = playState // Changer vers l'état de jeu
		} else if ebiten.IsKeyPressed(ebiten.Key3) {
			g.gameState = resultState // Changer vers l'état de résultat
		}
	}

	switch g.gameState {
	case introStateLogo:
		g.drawIntroLogo(screen)
	case introStateTexte:
		g.drawIntroTexte(screen)
	case titleState:
		g.titleDraw(screen)
	case themeState:
		g.themeDraw(screen)
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
	case replayState:
		g.drawReplay(screen)
	case shifumiState:
		g.DrawShifumi(screen)
	}

	if g.gameState != introStateLogo && g.gameState != introStateTexte {
		g.topMenu(screen)
		g.drawFullscreenButton(screen)
		if g.gameState != themeState {
			g.drawNbPlayer(screen)
		}
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
	g.drawThemeButton(screen)
	g.topMenuButton(screen, "CONNEXION")

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
	vector.DrawFilledRect(screen, float32(rect1X), float32(rect1Y), float32(rect1Width), float32(rect1Height), globalTextColorBright, true)

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
	vector.DrawFilledRect(screen, float32(rect2X), float32(rect2Y), float32(rect2Width), float32(rect2Height), globalTextColorBright, true)

	// Dessiner le texte pour le sous-titre 2
	text.Draw(screen, subTitle2, smallFont, subTitle2X, subTitle2Y, globalTextColor)

	// Ajouter un message clignotant en bas de l'écran
	blinkMessage := "Appuyez sur Entrée pour vous connectez"
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
			(float32(globalTileSize/2-globalCircleMargin) - 2),
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

	if g.chatIsFocus {
		g.drawChat(screen)
		g.drawCloseButton(screen)
	} else {
		g.drawChatButton(screen)
	}
}

// Affichage des graphismes durant le jeu.
func (g game) playDraw(screen *ebiten.Image) {
	if g.chatIsFocus {
		g.drawChat(screen)
		g.drawCloseButton(screen)
	} else {
		g.drawChatButton(screen)
	}

	g.drawScore(screen)

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

	// Afficher le pion du joueur actif (au-dessus de la grille)
	pionAdversaireX := float32(startX + globalTileSize/2 + g.adversaryTokenPosition*globalTileSize)
	pionAdversaireY := float32(startY-globalTileSize/2) - 20 // Juste au-dessus de la grille

	if pionX == pionAdversaireX {
		sizeP1 = -1.0
		sizeP2 = 5.0
	} else {
		sizeP1 = 0.0
		sizeP2 = 0.0
	}

	vector.DrawFilledCircle(
		screen,
		pionAdversaireX,
		pionAdversaireY,
		float32(globalTileSize/2-globalCircleMargin)+float32(sizeP2),
		globalTokenColors[g.p2Color],
		true,
	)

	vector.DrawFilledCircle(
		screen,
		pionX,
		pionY,
		float32(globalTileSize/2-globalCircleMargin)+float32(sizeP1),
		globalTextColor,
		true,
	)

	vector.DrawFilledCircle(
		screen,
		pionX,
		pionY,
		float32(globalTileSize/2-globalCircleMargin)+float32(sizeP1)-5,
		globalTokenColors[g.p1Color],
		true,
	)

	// Afficher les messages d'erreur, le cas échéant
	if !g.restartOk {
		g.errorMessageDisplay(screen, "En attente de l'autre joueur")
	}

	if g.chatIsFocus {
		g.drawChat(screen)
		g.drawCloseButton(screen)
	} else {
		g.drawChatButton(screen)
	}
}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {

	g.drawScore(screen)
	g.drawReplayButton(screen)
	// Centrer la grille
	g.drawGrid(offScreenImage)
	g.topMenuButton(screen, "REJOUER")

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

	if g.chatIsFocus {
		g.drawChat(screen)
		g.drawCloseButton(screen)
	} else {
		g.drawChatButton(screen)
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

func (g *game) DrawShifumi(screen *ebiten.Image) {
	if pierreImg == nil || papierImg == nil || ciseauxImg == nil {
		log.Println("Une ou plusieurs images ne sont pas chargées.")
		return
	}

	// Calculer la position de départ pour centrer les images
	scale := 0.5                                           // Facteur d'échelle pour les images
	imageWidth := float64(pierreImg.Bounds().Dx()) * scale // Largeur d'une image mise à l'échelle
	spacing := 50.0                                        // Espace entre les images
	totalWidth := imageWidth*3 + spacing*2                 // Largeur totale de toutes les images
	startX := (float64(globalWidth) - totalWidth) / 2      // Position X de départ
	startY := float64(globalHeight - 450)                  // Position Y de départ pour les images

	// Dessiner Pierre
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(startX, startY)
	screen.DrawImage(pierreImg, op)

	// Dessiner Papier
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(startX+imageWidth+spacing, startY)
	screen.DrawImage(papierImg, op)

	// Dessiner Ciseaux
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(startX+2*(imageWidth+spacing), startY)
	screen.DrawImage(ciseauxImg, op)

	message := "Choix du coup"
	textWidth, _ := getTextDimensions(message, firstTitleSmallFont)
	textX := (globalWidth - textWidth) / 2
	textY := globalHeight/2 - 200
	text.Draw(screen, message, firstTitleSmallFont, textX, textY, globalTextColorYellow)

	// Afficher le losange au-dessus de l'image sélectionnée
	if g.selected != "" {
		var selectedX float64

		// Déterminer la position X du centre de l'image sélectionnée
		switch g.selected {
		case "Pierre":
			selectedX = startX + (imageWidth / 2)
		case "Papier":
			selectedX = startX + imageWidth + spacing + (imageWidth / 2)
		case "Ciseaux":
			selectedX = startX + 2*(imageWidth+spacing) + (imageWidth / 2)
		}

		// Calculer la position du losange
		losangeScale := scale
		losangeWidth := float64(losangeImg.Bounds().Dx()) * losangeScale
		losangeX := selectedX - (losangeWidth / 2)
		losangeY := float64(globalHeight - 500) // 10 pixels d'espace entre le losange et l'image

		// Dessiner le losange
		opLosange := &ebiten.DrawImageOptions{}
		opLosange.GeoM.Scale(losangeScale, losangeScale)
		opLosange.GeoM.Translate(losangeX, losangeY)
		screen.DrawImage(losangeImg, opLosange)
	}

	// Afficher le résultat du shifumi si disponible
	if g.showShifumiResult {
		var textColor color.Color
		switch g.shifumiResult {
		case "Gagné":
			textColor = globalTextColorGreen
		case "Perdu":
			textColor = globalTextRed
		default:
			textColor = globalTextColorYellow
		}

		// Afficher le choix de l'adversaire
		adversaryText := "Choix de l'adversaire : " + g.adversaryChoice
		bounds := text.BoundString(smallFont, adversaryText)
		x := (globalWidth - bounds.Dx()) / 2
		text.Draw(screen, adversaryText, smallFont, x, int(startY)-50, globalTextColorBright)

		// Afficher le résultat
		resultText := g.shifumiResult
		bounds = text.BoundString(largeFont, resultText)
		x = (globalWidth - bounds.Dx()) / 2
		text.Draw(screen, resultText, largeFont, x, int(startY)-20, textColor)
	}
}

func (g *game) handleMouseClick(x, y int) {
	if g.gameState == shifumiState {
		scale := 0.5                                           // Facteur d'échelle pour les images
		imageWidth := float64(pierreImg.Bounds().Dx()) * scale // Largeur d'une image mise à l'échelle
		spacing := 50.0
		startX := (float64(globalWidth) - (imageWidth*3 + spacing*2)) / 2
		startY := globalHeight - 450 // Position Y de départ pour les images

		// Vérifier si l'utilisateur a cliqué sur l'image Pierre
		if x >= int(startX) && x <= int(startX+imageWidth) && y >= startY && y <= startY+int(float64(pierreImg.Bounds().Dy())) {
			g.selected = "Pierre"
		}

		// Vérifier si l'utilisateur a cliqué sur l'image Papier
		if x >= int(startX+imageWidth+spacing) && x <= int(startX+imageWidth+spacing+imageWidth) && y >= startY && y <= startY+int(float64(papierImg.Bounds().Dy())) {
			g.selected = "Papier"
		}

		// Vérifier si l'utilisateur a cliqué sur l'image Ciseaux
		if x >= int(startX+2*(imageWidth+spacing)) && x <= int(startX+2*(imageWidth+spacing)+imageWidth) && y >= startY && y <= startY+int(float64(ciseauxImg.Bounds().Dy())) {
			g.selected = "Ciseaux"
		}
		log.Print(g.selected)
	}
}
