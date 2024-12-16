package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"image/color"
	"strings"
)

func getTextDimensions(text string, fontFace font.Face) (width, height int) {
	bounds := font.MeasureString(fontFace, text)
	width = bounds.Ceil()                     // Largeur en pixels
	height = fontFace.Metrics().Height.Ceil() // Hauteur en pixels
	return
}

func drawRoundedRectangle(screen *ebiten.Image, x, y, width, height, radius float32, fillColor color.Color) {
	// Dessiner le rectangle central (sans les coins)
	vector.DrawFilledRect(screen, x+radius, y, width-2*radius, height, fillColor, true)
	vector.DrawFilledRect(screen, x, y+radius, width, height-2*radius, fillColor, true)

	// Dessiner les coins arrondis
	vector.DrawFilledCircle(screen, x+radius, y+radius, radius, fillColor, true)              // Coin supérieur gauche
	vector.DrawFilledCircle(screen, x+width-radius, y+radius, radius, fillColor, true)        // Coin supérieur droit
	vector.DrawFilledCircle(screen, x+radius, y+height-radius, radius, fillColor, true)       // Coin inférieur gauche
	vector.DrawFilledCircle(screen, x+width-radius, y+height-radius, radius, fillColor, true) // Coin inférieur droit
}

// centerPosition calcule les coordonnées X et Y pour centrer un élément donné
// centerPosition calcule les coordonnées X et Y pour centrer un élément donné
func centerPosition(elementWidth, elementHeight, screenWidth, screenHeight int) (centerX, centerY int) {
	centerX = (screenWidth - elementWidth) / 2
	centerY = (screenHeight - elementHeight) / 2
	return
}

func (g *game) drawIntroLogo(screen *ebiten.Image) {
	// Préparer les options pour dessiner le logo
	op := &ebiten.DrawImageOptions{}

	// Redimensionner (par exemple, réduire à 20% de sa taille d'origine)
	scaleFactor := 0.15
	op.GeoM.Scale(scaleFactor, scaleFactor)

	// Obtenir les dimensions originales du logo
	logoWidth, logoHeight := iconIUT.Size()

	// Calculer la position centrée
	logoWidthScaled := int(float64(logoWidth) * scaleFactor)
	logoHeightScaled := int(float64(logoHeight) * scaleFactor)
	centerX, centerY := centerPosition(logoWidthScaled, logoHeightScaled, globalWidth, globalHeight)

	// Déplacer le logo au centre
	op.GeoM.Translate(float64(centerX), float64(centerY))

	// Définir les durées des différentes phases
	totalDuration := fadeInDuration + holdDuration + fadeOutDuration

	// Calculer l'opacité (alpha)
	alpha := 0.0
	if g.stateFrame < fadeInDuration { // Phase d'apparition
		alpha = float64(g.stateFrame) / float64(fadeInDuration)
	} else if g.stateFrame < fadeInDuration+holdDuration { // Phase de maintien
		alpha = 1.0
	} else if g.stateFrame < totalDuration { // Phase de disparition
		alpha = float64(totalDuration-g.stateFrame) / float64(fadeOutDuration)
	}

	// Afficher le logo uniquement si dans les phases fadeIn, hold, fadeOut ou invisibleHold
	if g.stateFrame < totalDuration+invisibleHoldDuration {
		// Appliquer l'opacité calculée
		op.ColorM.Scale(1, 1, 1, alpha)    // Opacité réelle
		if g.stateFrame >= totalDuration { // Phase d'invisibilité
			op.ColorM.Scale(1, 1, 1, 0) // Forcer l'opacité à 0
		}

		// Dessiner l'image
		screen.DrawImage(iconIUT, op)
	} else {
		// Une fois la phase invisibleHold terminée, passer à une étape suivante
		return // Exemple : changer l'état ou démarrer une autre animation
	}
}

func (g *game) drawIntroTexte(screen *ebiten.Image) {
	// Préparer les options pour dessiner le logo
	op := &ebiten.DrawImageOptions{}

	// Redimensionner le logo
	scaleFactor := 0.35
	op.GeoM.Scale(scaleFactor, scaleFactor)

	// Obtenir les dimensions originales du logo
	logoWidth, logoHeight := logoStudio.Size()

	// Calculer la position centrée pour le logo
	logoWidthScaled := int(float64(logoWidth) * scaleFactor)
	logoHeightScaled := int(float64(logoHeight) * scaleFactor)
	centerX, centerY := centerPosition(logoWidthScaled, logoHeightScaled, globalWidth, globalHeight)
	centerY -= 50

	// Déplacer le logo au centre
	op.GeoM.Translate(float64(centerX), float64(centerY))

	// Définir les durées des différentes phases
	totalDuration := fadeInDuration + holdDuration + fadeOutDuration

	// Calculer l'opacité (alpha)
	alpha := 0.0
	if g.stateFrame < fadeInDuration { // Phase d'apparition
		alpha = float64(g.stateFrame) / float64(fadeInDuration)
	} else if g.stateFrame < fadeInDuration+holdDuration { // Phase de maintien
		alpha = 1.0
	} else if g.stateFrame < totalDuration { // Phase de disparition
		alpha = float64(totalDuration-g.stateFrame) / float64(fadeOutDuration)
	}

	// Afficher le logo uniquement si dans les phases fadeIn, hold, fadeOut ou invisibleHold
	if g.stateFrame < totalDuration+invisibleHoldDuration {
		// Appliquer l'opacité calculée
		op.ColorM.Scale(1, 1, 1, alpha)    // Opacité réelle
		if g.stateFrame >= totalDuration { // Phase d'invisibilité
			op.ColorM.Scale(1, 1, 1, 0) // Forcer l'opacité à 0
		}

		// Dessiner le logo
		screen.DrawImage(logoStudio, op)

		// Dessiner le texte en dessous du logo
		textStudio := "a BenCheikh Production"

		// Obtenir les dimensions du texte pour centrer
		textWidth, _ := getTextDimensions(textStudio, largeFont)
		textX, textY := centerPosition(textWidth, logoHeightScaled, globalWidth, globalHeight)
		textY += 210

		// Configurer la couleur du texte avec opacité
		textColor := color.NRGBA{255, 255, 255, uint8(alpha * 255)}

		// Dessiner le texte
		text.Draw(screen, textStudio, largeFont, textX, textY, textColor)
	} else {
		// Une fois la phase invisibleHold terminée, passer à une étape suivante
		return
	}
}

func (g game) themeDraw(screen *ebiten.Image) {
	g.topMenuButton(screen, "VALIDER")

	image1 := background2
	image2 := background3
	image3 := background1

	if g.nbBackgroundTheme == 0 {
		image1 = background2
		image2 = background3
		image3 = background1
	} else if g.nbBackgroundTheme == 1 {
		image1 = background1
		image2 = background2
		image3 = background3
	} else if g.nbBackgroundTheme == 2 {
		image1 = background3
		image2 = background1
		image3 = background2
	}

	// Facteur d'échelle pour redimensionner le logo
	iconScale := 0.15 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée
	originalWidth, originalHeight := background2.Size()
	scaledWidth := float64(originalWidth) * iconScale
	scaledHeight := float64(originalHeight) * iconScale

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX := globalWidth/2 - (int(scaledWidth) / 2)
	buttonY := globalHeight/2 - (int(scaledHeight) / 2)

	// Facteur d'échelle pour redimensionner le logo
	iconScale2 := 0.1 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée
	originalWidth2, originalHeight2 := background1.Size()
	scaledWidth2 := float64(originalWidth2) * iconScale2
	scaledHeight2 := float64(originalHeight2) * iconScale2

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX2 := globalWidth/2 - (int(scaledWidth2) / 2) - (int(scaledWidth) / 2)
	buttonX3 := globalWidth/2 - (int(scaledWidth2) / 2) + (int(scaledWidth) / 2)
	buttonY2 := globalHeight/2 - (int(scaledHeight2) / 2)

	vector.DrawFilledRect(screen, float32(buttonX2-5), float32(buttonY2-5), float32(scaledWidth2+10), float32(scaledHeight2+10), globalTextColorBright, true) // Rectangle transparent
	vector.DrawFilledRect(screen, float32(buttonX3-5), float32(buttonY2-5), float32(scaledWidth2+10), float32(scaledHeight2+10), globalTextColorBright, true) // Rectangle transparent

	// Dessiner l'icône redimensionnée
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(iconScale2, iconScale2)                   // Appliquer le facteur d'échelle
	op2.GeoM.Translate(float64(buttonX2), float64(buttonY2)) // Positionner l'icône
	screen.DrawImage(image1, op2)

	// Dessiner l'icône redimensionnée
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Scale(iconScale2, iconScale2)                   // Appliquer le facteur d'échelle
	op3.GeoM.Translate(float64(buttonX3), float64(buttonY2)) // Positionner l'icône
	screen.DrawImage(image2, op3)

	vector.DrawFilledRect(screen, float32(buttonX-5), float32(buttonY-5), float32(scaledWidth+10), float32(scaledHeight+10), globalTextColorGreen, true) // Rectangle transparent

	// Dessiner l'icône redimensionnée
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
	op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
	screen.DrawImage(image3, op)

	// Facteur d'échelle pour redimensionner le logo
	iconScaleArrow := 0.1 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée
	originalWidthArrow, originalHeightArrow := background2.Size()
	scaledWidthArrow := float64(originalWidthArrow) * iconScaleArrow
	scaledHeightArrow := float64(originalHeightArrow) * iconScaleArrow

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonXArrow := globalWidth/2 - (int(scaledWidthArrow) / 2) - 45
	buttonYArrow := (globalHeight / 2) + 220

	vector.DrawFilledRect(screen, float32(buttonXArrow), float32(buttonYArrow), float32(scaledWidthArrow/2-50), float32(scaledHeightArrow/2-15), globalTextColorBright, true) // Rectangle transparent

	opArrow := &ebiten.DrawImageOptions{}
	opArrow.GeoM.Scale(iconScaleArrow, iconScaleArrow)                      // Appliquer le facteur d'échelle
	opArrow.GeoM.Translate(float64(buttonXArrow+25), float64(buttonYArrow)) // Positionner l'icône
	screen.DrawImage(leftArrowImage, opArrow)

	buttonXArrow2 := globalWidth/2 + (int(scaledWidthArrow) / 2) - 100

	vector.DrawFilledRect(screen, float32(buttonXArrow2), float32(buttonYArrow), float32(scaledWidthArrow/2-50), float32(scaledHeightArrow/2-15), globalTextColorBright, true) // Rectangle transparent

	opArrow2 := &ebiten.DrawImageOptions{}
	opArrow2.GeoM.Scale(iconScaleArrow, iconScaleArrow)                       // Appliquer le facteur d'échelle
	opArrow2.GeoM.Translate(float64(buttonXArrow2+25), float64(buttonYArrow)) // Positionner l'icône
	screen.DrawImage(rightArrowImage, opArrow2)

}

func (g game) topMenu(screen *ebiten.Image) {
	// Dimensions du grand rectangle en haut
	vector.DrawFilledRect(screen, 0, 0, float32(globalWidth), float32(largeRectHeight), color.NRGBA{R: 0, G: 0, B: 0, A: 0}, true) // Rectangle transparent

	// Dessiner le bord inférieur du grand rectangle
	borderHeight := 2
	vector.DrawFilledRect(screen, 0, float32(largeRectHeight-borderHeight), float32(globalWidth), float32(borderHeight), globalTextColorBright, true)
}

func (g game) topMenuButton(screen *ebiten.Image, buttonText string) {
	// Calculer les dimensions du texte
	textWidth, _ := getTextDimensions(buttonText, firstTitleMinusFont)

	// Ajouter du padding à gauche et à droite
	horizontalPadding := 60
	smallRectWidth := textWidth + horizontalPadding*2
	smallRectHeight := largeRectHeight

	// Calculer la position centrée du rectangle
	smallRectX, smallRectY := centerPosition(smallRectWidth, smallRectHeight, globalWidth, largeRectHeight)

	// Dessiner le rectangle avec un padding horizontal et une hauteur fixe
	vector.DrawFilledRect(screen, float32(smallRectX), float32(smallRectY), float32(smallRectWidth), float32(smallRectHeight), globalTextColorBright, true)

	// Calculer la position centrée du texte à l'intérieur du rectangle
	textX := smallRectX + horizontalPadding
	textY := (largeRectHeight / 2) + 6

	// Dessiner le texte
	text.Draw(screen, buttonText, firstTitleMinusFont, textX, textY, globalTextColor)
}

func (g game) inputServerDraw(screen *ebiten.Image) {

	g.topMenuButton(screen, "JOUER")

	// Texte principal
	mainText := "Entrez l'adresse du serveur"
	mainWidth, mainHeight := getTextDimensions(mainText, firstTitleSmallerFont)
	mainX, mainY := centerPosition(mainWidth, 0, globalWidth, globalHeight)
	mainY -= 50 // Décalage vertical vers le haut
	text.Draw(screen, mainText, firstTitleSmallerFont, mainX, mainY, globalTextColorYellow)

	// Afficher l'adresse actuelle ou une suggestion par défaut
	addressToShow := g.serverAddress
	if addressToShow == "" {
		addressToShow = "localhost:8080 (par défaut)"
	}
	// Calcul des dimensions et position
	subTitle1Width, subTitle1Height := getTextDimensions(addressToShow, smallFont)
	subTitle1X := (globalWidth - subTitle1Width) / 2
	subTitle1Y := mainY + mainHeight + 30

	// Rectangle classique
	rectX := subTitle1X - 40
	rectY := subTitle1Y - 45
	rectWidth := subTitle1Width + 80
	rectHeight := subTitle1Height + 20

	// Dessiner le rectangle classique
	vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(rectWidth), float32(rectHeight), globalTextColor, true)

	// Dessiner le texte de l'adresse
	text.Draw(screen, addressToShow, smallFont, subTitle1X, subTitle1Y, globalTextColorBright)

	// Afficher le message de connexion ou d'erreur
	if g.errorConnection != "" {
		g.errorMessageDisplay(screen, "Adresse introuvable")
	}

}

func (g game) errorMessageDisplay(screen *ebiten.Image, msg string) {
	subTitle2 := msg

	// Dimensions du texte
	textWidth, _ := getTextDimensions(subTitle2, mediumFontError)
	padding := 15 // Espace entre le texte et le bord du rectangle
	rectWidth := textWidth + 2*padding
	rectHeight := 80

	rectX := 0.0                   // Position au bord de l'ecran
	rectY := (80 - rectHeight) / 2 // Centré verticalement dans le grand rectangle (80px de hauteur)
	textX := int(rectX) + padding
	textY := rectY + padding + 30

	const visibleFrames = 70                              // Nombre de frames où le message est visible
	const invisibleFrames = 30                            // Nombre de frames où le message est invisible
	const cycleDuration = visibleFrames + invisibleFrames // Durée totale du cycle

	if g.stateFrame%cycleDuration < visibleFrames {
		// Le message est visible
		// Dessiner le rectangle rouge
		vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(rectWidth), float32(rectHeight), globalTextRed, true)
		// Dessiner le texte centré dans le rectangle
		text.Draw(screen, subTitle2, mediumFontError, textX, textY, globalTextColorBright)
	}
}

func (g *game) connectingDraw(screen *ebiten.Image) {
	// Message fixe
	message := "Connexion au serveur"
	textWidth, textHeight := getTextDimensions(message, firstTitleSmallerFont)

	// Calculer les positions pour centrer le texte
	textX, textY := centerPosition(textWidth, textHeight, globalWidth, globalHeight)
	textY = textY - 30

	// Dessiner le texte
	text.Draw(screen, message, firstTitleSmallerFont, textX, textY, globalTextColorYellow)

	// Dimensions de la barre de chargement
	barWidth := 300
	barHeight := 10
	//position de la barre de chargement
	barX := (globalWidth - barWidth) / 2
	barY := textY + textHeight

	// Calcul de la progression (cycle toutes les 120 frames)
	progress := (g.stateFrame % 120) * barWidth / 120

	// Fond de la barre
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), float64(barWidth), float64(barHeight), globalTextColor)

	// Barre de progression en mouvement
	ebitenutil.DrawRect(screen, float64(barX), float64(barY), float64(progress), float64(barHeight), globalTextColorYellow)
}

func (g game) waitingDraw(screen *ebiten.Image) {
	if g.playerID == -1 {
		g.connectingDraw(screen)
		return
	}
	// Liste des messages avec différents points de suspension
	baseMessage := "En attente de l'autre joueur"
	dots := []string{"", ".", "..", "..."}       // Cycle des points
	dotsIndex := (g.stateFrame / 30) % len(dots) // Change toutes les 30 frames (1/2 seconde à 60 FPS)

	// Construire le message actuel
	message := baseMessage + dots[dotsIndex]

	// Calculer les dimensions du texte
	textWidth, textHeight := getTextDimensions(message, firstTitleSmallerFont)

	// Calculer les positions pour centrer le texte
	textX, textY := centerPosition(textWidth, textHeight, globalWidth, globalHeight)

	// Dessiner le texte centré
	text.Draw(screen, message, firstTitleSmallerFont, textX, textY, globalTextColorYellow)
}

func (g game) drawFullscreenButton(screen *ebiten.Image) {
	// Déterminer l'icône à afficher
	var icon *ebiten.Image
	if ebiten.IsFullscreen() {
		icon = iconWindowed
	} else {
		icon = iconFullscreen
	}

	// Facteur d'échelle pour redimensionner le logo
	iconScale := 0.06 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée
	originalWidth, _ := icon.Size()
	scaledWidth := float64(originalWidth) * iconScale

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX := globalWidth - int(scaledWidth) - 35
	buttonY := 11

	// Dessiner l'icône redimensionnée
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
	op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
	screen.DrawImage(icon, op)

	// Dessiner un petit cercle blanc
	circleRadius := float32(5)                          // Rayon du cercle
	circleX := float32(globalWidth) - circleRadius - 57 // marge depuis la droite
	circleY := float32(80)                              // 80 pixels depuis le haut

	// Dessiner le cercle
	vector.DrawFilledCircle(screen, circleX, circleY, circleRadius, globalTextColorBright, true)

}

func (g game) drawNbPlayer(screen *ebiten.Image) {

	// Ajouter le texte pour afficher le nombre de joueurs connectés
	playerText := fmt.Sprintf("JOUEURS CONNECTES : %d", g.nbJoueurConnecte)
	textWidth, textHeight := getTextDimensions(playerText, mediumFontError)

	// Calculer les coordonnées pour placer le texte en bas à droite
	margin := 20                                   // Marge par rapport aux bords
	textX := globalWidth - textWidth - margin - 20 // Aligné au bord droit
	textY := globalHeight - 42                     // Aligné au bord bas

	// Dimensions et position du fond
	padding := 10 // Espace autour du texte
	rectX := textX - padding
	rectY := textY - 25
	rectWidth := textWidth + 20
	rectHeight := textHeight

	// Dessiner le rectangle derrière le texte

	// Dessiner le texte par-dessus le fond
	if g.gameState != titleState {
		vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(rectWidth), float32(rectHeight), globalTextColorBright, true)

		text.Draw(screen, playerText, mediumFontError, textX, textY, globalTextColor)
	}
}

func (g game) drawScore(screen *ebiten.Image) {

	// Ajouter le texte pour afficher le nombre de joueurs connectés
	playerText := fmt.Sprintf("PERSONAL WIN : %d", g.nbPartieWin)
	textWidth, textHeight := getTextDimensions(playerText, mediumFontError)

	// Calculer les coordonnées pour placer le texte en bas à droite
	margin := 20                                   // Marge par rapport aux bords
	textX := globalWidth - textWidth - margin - 20 // Aligné au bord droit
	textY := globalHeight / 2                      // Aligné au bord bas

	// Dimensions et position du fond
	padding := 10 // Espace autour du texte
	rectX := textX - padding
	rectY := textY - 25
	rectWidth := textWidth + 20
	rectHeight := textHeight

	// Dessiner le rectangle derrière le texte
	// Dessiner le texte par-dessus le fond
	vector.DrawFilledRect(screen, float32(rectX), float32(rectY), float32(rectWidth), float32(rectHeight), globalTextColorBright, true)
	text.Draw(screen, playerText, mediumFontError, textX, textY, globalTextColor)

	// Ajouter le texte pour afficher le nombre de joueurs connectés
	playerText2 := fmt.Sprintf("OPPONENT WIN : %d", g.nbPartieAdversaireWin)
	textWidth2, textHeight2 := getTextDimensions(playerText2, mediumFontError)

	// Calculer les coordonnées pour placer le texte en bas à droite
	textX2 := textX      // Aligné au bord droit
	textY2 := textY + 50 // Aligné au bord bas

	// Dimensions et position du fond
	rectX2 := textX2 - padding
	rectY2 := textY2 - 25
	rectWidth2 := textWidth2 + 20
	rectHeight2 := textHeight2

	// Dessiner le rectangle derrière le texte
	// Dessiner le texte par-dessus le fond
	vector.DrawFilledRect(screen, float32(rectX2), float32(rectY2), float32(rectWidth2), float32(rectHeight2), globalTextColorBright, true)
	text.Draw(screen, playerText2, mediumFontError, textX2, textY2, globalTextColor)
}

func (g game) drawThemeButton(screen *ebiten.Image) {
	// Déterminer l'icône à afficher
	var icon *ebiten.Image = Themes

	// Facteur d'échelle pour redimensionner le logo
	iconScale := 0.055 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée
	originalWidth, _ := icon.Size()
	scaledWidth := float64(originalWidth) * iconScale

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX := globalWidth - int(scaledWidth) - 155
	buttonY := 13

	// Dessiner l'icône redimensionnée
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
	op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
	screen.DrawImage(icon, op)

	// Dessiner un petit cercle blanc
	circleRadius := float32(5)                           // Rayon du cercle
	circleX := float32(globalWidth) - circleRadius - 174 // marge depuis la droite
	circleY := float32(80)

	vector.DrawFilledCircle(screen, circleX, circleY, circleRadius, globalTextColorBright, true)

}

func (g game) drawReplayButton(screen *ebiten.Image) {
	// Déterminer l'icône à afficher
	var icon *ebiten.Image = replay

	// Facteur d'échelle pour redimensionner le logo
	iconScale := 0.055 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée
	originalWidth, _ := icon.Size()
	scaledWidth := float64(originalWidth) * iconScale

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX := globalWidth - int(scaledWidth) - 155
	buttonY := 13

	// Dessiner l'icône redimensionnée
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
	op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
	screen.DrawImage(icon, op)

	// Dessiner un petit cercle blanc
	circleRadius := float32(5)                           // Rayon du cercle
	circleX := float32(globalWidth) - circleRadius - 178 // marge depuis la droite
	circleY := float32(80)

	vector.DrawFilledCircle(screen, circleX, circleY, circleRadius, globalTextColorBright, true)

}

func (g game) drawReplay(screen *ebiten.Image) {

	g.topMenuButton(screen, "RETOUR")
	g.drawGrid(screen)

	message := "Il y a Egalite"
	if g.result == p1wins {
		message = "Vous avez Gagne !"
	} else if g.result == p2wins {
		message = "Vous avez Perdu"
	}
	textWidth, _ := getTextDimensions(message, firstTitleSmallerFont)
	textX := (globalWidth - textWidth) / 2
	textY := globalHeight/2 - ((globalNumTilesY * globalTileSize) / 2) + 20
	text.Draw(screen, message, firstTitleSmallerFont, textX, textY, globalTextColorYellow)

	// Calculer la position et les dimensions de la grille
	gridWidth := globalTileSize * globalNumTilesX
	gridHeight := globalTileSize * globalNumTilesY
	startX := (globalWidth - gridWidth) / 2
	startY := (globalHeight-gridHeight)/2 + 50

	if g.blinking {
		// Dessiner les cercles pour chaque cellule
		for _, pos := range g.posWinner {
			x, y := pos[0], pos[1]

			centerX := float32(startX + globalTileSize/2 + x*globalTileSize)
			centerY := float32(startY + globalTileSize/2 + y*globalTileSize)

			if g.stateFrame%60 < 10 { // Clignote toutes les 30 frames (1/2 seconde à 60 FPS)
				vector.DrawFilledCircle(
					screen,
					centerX,
					centerY,
					float32(globalTileSize/2-globalCircleMargin),
					globalTextRed,
					true,
				)
			}

		}
	}

	if g.chatIsFocus {
		g.drawChat(screen)
		g.drawCloseButton(screen)
	} else {
		g.drawChatButton(screen)
	}

}

func (g game) drawChat(screen *ebiten.Image) {

	// Dessiner une bordure autour de la zone de saisie
	ebitenutil.DrawRect(screen, 5, float64(globalHeight)-5, inputMaxWidth-80, -400, globalTextColor)

	// Dessiner la zone des messages
	yOffset := globalHeight - 50 // Point de départ en bas de l'écran
	for i := len(g.chatMessages) - 1; i >= 0; i-- {
		message := g.chatMessages[i]

		// Découper le message en lignes
		lines := strings.Split(message, "\n")
		totalLinesHeight := len(lines)*(20+10) - 10 // Hauteur totale pour ce message (sans double espacement)

		// Ajuster yOffset pour inclure l'espacement entre messages
		yOffset -= totalLinesHeight + 10
		if yOffset < 50 { // Ne pas dépasser la zone visible
			break
		}

		// Dessiner chaque ligne du message
		for j, line := range lines {
			text.Draw(screen, line, mediumFontError, 10, yOffset+(j*(20+5)), globalTextColorBright)
		}
	}

	// Afficher le champ de saisie
	text.Draw(screen, "Message: "+g.chatInput, mediumFontError, 10, globalHeight-20, globalTextColorGreen)

}

func (g game) drawChatButton(screen *ebiten.Image) {
	// Déterminer l'icône à afficher
	var icon *ebiten.Image = chat
	var iconWarn = chatWarning

	// Facteur d'échelle pour redimensionner le logo
	iconScale := 0.06 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX := 40
	buttonY := globalHeight - 85

	if !g.chatNewMessage {
		// Dessiner l'icône redimensionnée
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
		op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
		screen.DrawImage(icon, op)
	} else {
		// Dessiner l'icône redimensionnée
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
		op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
		screen.DrawImage(icon, op)

		if g.stateFrame%60 < 10 { // Clignote toutes les 30 frames (1/2 seconde à 60 FPS)
			// Dessiner l'icône redimensionnée
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
			op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
			screen.DrawImage(iconWarn, op)
		}
	}

}

func (g game) drawCloseButton(screen *ebiten.Image) {
	// Déterminer l'icône à afficher
	var icon *ebiten.Image = close
	// Facteur d'échelle pour redimensionner le logo
	iconScale := 0.035 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX := inputMaxWidth - 120
	buttonY := globalHeight - 390

	// Dessiner l'icône redimensionnée
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
	op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
	screen.DrawImage(icon, op)

}
