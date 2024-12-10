package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"image/color"
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

func (g game) topMenu(screen *ebiten.Image) {
	// Dimensions du grand rectangle en haut
	vector.DrawFilledRect(screen, 0, 0, float32(globalWidth), float32(largeRectHeight), color.NRGBA{R: 0, G: 0, B: 0, A: 0}, true) // Rectangle transparent

	// Dessiner le bord inférieur du grand rectangle
	borderHeight := 2
	vector.DrawFilledRect(screen, 0, float32(largeRectHeight-borderHeight), float32(globalWidth), float32(borderHeight), globalTextColorBright, true)

	// Dessiner un petit cercle blanc
	circleRadius := float32(5)                          // Rayon du cercle
	circleX := float32(globalWidth) - circleRadius - 45 // marge depuis la droite
	circleY := float32(80)                              // 80 pixels depuis le haut

	// Dessiner le cercle
	vector.DrawFilledCircle(screen, circleX, circleY, circleRadius, globalTextColorBright, true)
}

func (g game) topMenuJoueur(screen *ebiten.Image) {
	// Dimensions du petit rectangle centré dans le grand rectangle
	smallRectWidth := 200
	smallRectHeight := 80
	smallRectX, smallRectY := centerPosition(smallRectWidth, smallRectHeight, globalWidth, largeRectHeight)

	// Dessiner le petit rectangle
	vector.DrawFilledRect(screen, float32(smallRectX), float32(smallRectY), float32(smallRectWidth), float32(smallRectHeight), globalTextColorBright, true)

	// Texte "JOUER" dans le petit rectangle
	buttonText := "JOUER"
	textWidth, textHeight := getTextDimensions(buttonText, smallFont)
	textX, textY := centerPosition(textWidth, textHeight, smallRectWidth, smallRectHeight)
	textX += smallRectX
	textY += smallRectY + 35
	text.Draw(screen, buttonText, smallFont, textX, textY, globalTextColor)
}

func (g game) topMenuRejouer(screen *ebiten.Image) {
	// Couleur pour le bord et le petit rectangle
	backgroundColor := color.NRGBA{R: 217, G: 217, B: 217, A: 255} // #D9D9D9 pour le petit rectangle

	// Dimensions du grand rectangle en haut
	largeRectHeight := 80

	// Dimensions du petit rectangle centré dans le grand rectangle
	smallRectWidth := 200
	smallRectHeight := 80
	smallRectX, smallRectY := centerPosition(smallRectWidth, smallRectHeight, globalWidth, largeRectHeight)

	// Dessiner le petit rectangle
	vector.DrawFilledRect(screen, float32(smallRectX), float32(smallRectY), float32(smallRectWidth), float32(smallRectHeight), backgroundColor, true)

	// Texte "JOUER" dans le petit rectangle
	buttonText := "REJOUER"
	textWidth, textHeight := getTextDimensions(buttonText, smallFont)
	textX, textY := centerPosition(textWidth, textHeight, smallRectWidth, smallRectHeight)
	textX += smallRectX
	textY += smallRectY + 35
	text.Draw(screen, buttonText, smallFont, textX, textY, globalTextColor)
}

func (g game) inputServerDraw(screen *ebiten.Image) {

	g.topMenuJoueur(screen)

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
	iconScale := 0.07 // Réduire à 50% de la taille d'origine

	// Calculer la taille de l'icône redimensionnée
	originalWidth, _ := icon.Size()
	scaledWidth := float64(originalWidth) * iconScale

	// Définir la position du bouton (centré horizontalement en haut à droite)
	buttonX := globalWidth - int(scaledWidth) - 19
	buttonY := 8

	// Dessiner l'icône redimensionnée
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(iconScale, iconScale)                   // Appliquer le facteur d'échelle
	op.GeoM.Translate(float64(buttonX), float64(buttonY)) // Positionner l'icône
	screen.DrawImage(icon, op)

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
