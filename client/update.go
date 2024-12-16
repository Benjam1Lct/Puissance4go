package main

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"log"
	"strings"
	"time"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleMouseClick(x, y)
	}

	g.UpdateFullscreen()
	g.UpdateDebug()
	g.updateChat()
	g.updateCloseButton()
	g.updateChatButton()

	g.stateFrame++

	switch g.gameState {
	case introStateLogo:
		if g.introStateLogo() {
			g.gameState = introStateTexte
			g.stateFrame = 0
		}
	case introStateTexte:
		if g.introStateTexte() {
			g.gameState = titleState
			g.stateFrame = 0
		}
	case titleState:
		if g.titleUpdate() {
			g.gameState = inputServerState
		}
	case themeState:
		if g.UpdateThemesPage() {
			g.gameState = titleState
		}
	case inputServerState:
		if g.inputServerUpdate() {
			g.gameState = waitingState
			go connectToServer(g) // Lancer la connexion au serveur
		}
	case waitingState:

		if g.serverReady {
			// La transition vers colorSelectState est déjà gérée par connectToServer
		}
	case colorSelectState:
		if g.colorSelectUpdate() {
			g.gameState = waitingColorSelect
			if g.gameState == waitingColorSelect {
			}
		}
	case waitingColorSelect:
		if g.gameReady {
			// La transition vers colorSelectState est déjà gérée par connectToServer
		}
	case shifumiState:
		g.UpdateShifumi()
	case playState:
		if !g.restartOk {
			return nil // Bloquer les mises à jour
		}
		g.p1Color = g.p1ColorValidate
		g.tokenPosUpdate()
		var lastXPositionPlayed, lastYPositionPlayed int
		if g.turn == p1Turn {
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update()
		}
		if lastXPositionPlayed >= 0 {
			finished, result, posWinnerCheck := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				if posWinnerCheck != nil {
					g.posWinner = posWinnerCheck
					log.Println("posWinnerCheck", g.posWinner)
				}
				if g.result == p1wins {
					g.nbPartieWin++
					g.turn = p2Turn
				} else if g.result == p2wins {
					g.nbPartieAdversaireWin++
					g.turn = p1Turn
				} else {
					g.nbPartieWin++
					g.nbPartieAdversaireWin++
					g.turn = g.firstPlayer
				}
				g.gameState = resultState
				g.restartOk = false
			}
		}
	case resultState:
		if g.resultUpdate() {
			g.reset()
			g.gameState = playState
		}
	case replayState:
		if g.replayDrawUpdate() {
			g.gameState = resultState
			g.blinking = false
		}
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.mouseReleased = true // Souris relâchée
	}

	return nil
}

func (g *game) UpdateFullscreen() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		var icon *ebiten.Image
		if ebiten.IsFullscreen() {
			icon = iconWindowed
		} else {
			icon = iconFullscreen
		}
		// Facteur d'échelle pour redimensionner le logo
		iconScale := 0.06 // Réduire à 50% de la taille d'origine

		// Calculer la taille de l'icône redimensionnée
		originalWidth, originalHeight := icon.Size()
		scaledWidth := float64(originalWidth) * iconScale
		scaledHeight := float64(originalHeight) * iconScale

		// Définir la position du bouton (centré horizontalement en haut à droite)
		buttonX := globalWidth - int(scaledWidth) - 35
		buttonY := 11

		// Vérifier si le clic est sur le bouton
		if x >= buttonX && x <= buttonX+int(scaledWidth) && y >= buttonY && y <= buttonY+int(scaledHeight) {
			// Basculer entre plein écran et fenêtré
			ebiten.SetFullscreen(!ebiten.IsFullscreen())
			// Si en plein écran, ajustez la résolution cible
			// Ajuster les dimensions globales
			adjustGlobalDimensions()
			initFonts() // Recalculer les polices après le changement

			// Redimensionner la fenêtre en conséquence
			ebiten.SetWindowSize(globalWidth, globalHeight)

		}
	}
	return nil
}

func (g *game) UpdateThemesPage() bool {
	mouseX, mouseY := ebiten.CursorPosition()

	// Vérifier si le clic est dans les limites du rectangle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseReleased && !g.chatIsFocus {
		// Dimensions du rectangle
		textWidth, _ := getTextDimensions("VALIDER", firstTitleMinusFont)

		// Ajouter du padding à gauche et à droite
		horizontalPadding := 60
		smallRectWidth := textWidth + horizontalPadding*2
		smallRectHeight := largeRectHeight

		// Facteur d'échelle pour redimensionner le logo
		iconScaleArrow := 0.1 // Réduire à 50% de la taille d'origine

		// Calculer la taille de l'icône redimensionnée
		originalWidthArrow, originalHeightArrow := background2.Size()
		scaledWidthArrow := float64(originalWidthArrow) * iconScaleArrow
		scaledHeightArrow := float64(originalHeightArrow) * iconScaleArrow

		// Définir la position du bouton (centré horizontalement en haut à droite)
		buttonXArrow := globalWidth/2 - (int(scaledWidthArrow) / 2) - 45
		buttonYArrow := (globalHeight / 2) + 220

		buttonXArrow2 := globalWidth/2 + (int(scaledWidthArrow) / 2) - 100

		// Calculer la position centrée du rectangle
		smallRectX, smallRectY := centerPosition(smallRectWidth, smallRectHeight, globalWidth, largeRectHeight)
		// Vérification des coordonnées
		if mouseX >= smallRectX && mouseX <= smallRectX+smallRectWidth &&
			mouseY >= smallRectY && mouseY <= smallRectY+smallRectHeight {
			// Passer à l'état suivant
			g.mouseReleased = false
			g.gameState = titleState
			g.nbBackground = g.nbBackgroundTheme // Remplacez "nextState" par l'état que vous voulez
		} else if mouseX >= buttonXArrow && mouseX <= buttonXArrow+int(scaledWidthArrow)/2-50 &&
			mouseY >= buttonYArrow && mouseY <= buttonYArrow+int(scaledHeightArrow)/2-15 {

			g.nbBackgroundTheme--
			if g.nbBackgroundTheme < 0 {
				g.nbBackgroundTheme = 2
			}
			g.mouseReleased = false

		} else if mouseX >= buttonXArrow2 && mouseX <= buttonXArrow2+int(scaledWidthArrow)/2-50 &&
			mouseY >= buttonYArrow && mouseY <= buttonYArrow+int(scaledHeightArrow)/2-15 {
			g.nbBackgroundTheme++
			if g.nbBackgroundTheme > 2 {
				g.nbBackgroundTheme = 0
			}
			g.mouseReleased = false

		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.chatIsFocus {
		g.gameState = titleState
		g.nbBackground = g.nbBackgroundTheme
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && !g.chatIsFocus {
		g.nbBackgroundTheme++
		if g.nbBackgroundTheme > 2 {
			g.nbBackgroundTheme = 0
		}

	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && !g.chatIsFocus {
		g.nbBackgroundTheme--
		if g.nbBackgroundTheme < 0 {
			g.nbBackgroundTheme = 2
		}

	}
	return false
}

func (g *game) UpdateThemes() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.chatIsFocus {
		x, y := ebiten.CursorPosition()
		var icon *ebiten.Image = Themes
		// Facteur d'échelle pour redimensionner le logo
		iconScale := 0.055 // Réduire à 50% de la taille d'origine

		// Calculer la taille de l'icône redimensionnée
		originalWidth, originalHeight := icon.Size()
		scaledWidth := float64(originalWidth) * iconScale
		scaledHeight := float64(originalHeight) * iconScale

		// Définir la position du bouton (centré horizontalement en haut à droite)
		buttonX := globalWidth - int(scaledWidth) - 155
		buttonY := 13

		// Vérifier si le clic est sur le bouton
		if x >= buttonX && x <= buttonX+int(scaledWidth) && y >= buttonY && y <= buttonY+int(scaledHeight) {
			g.gameState = themeState
		}
	}
	return nil
}

func (g *game) UpdateReplayButton() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.chatIsFocus {
		x, y := ebiten.CursorPosition()
		var icon *ebiten.Image = replay
		// Facteur d'échelle pour redimensionner le logo
		iconScale := 0.055 // Réduire à 50% de la taille d'origine

		// Calculer la taille de l'icône redimensionnée
		originalWidth, originalHeight := icon.Size()
		scaledWidth := float64(originalWidth) * iconScale
		scaledHeight := float64(originalHeight) * iconScale

		// Définir la position du bouton (centré horizontalement en haut à droite)
		buttonX := globalWidth - int(scaledWidth) - 155
		buttonY := 13

		// Vérifier si le clic est sur le bouton
		if x >= buttonX && x <= buttonX+int(scaledWidth) && y >= buttonY && y <= buttonY+int(scaledHeight) {
			g.gameState = replayState
			g.resetGrid()
			err := requestHistory(g.conn)
			if err != nil {
				return err
			}
			g.blinking = false
		}
	}
	return nil
}

func (g *game) UpdateDebug() error {
	// ... votre code update existant ...

	// Activer/désactiver le mode debug avec F3 par exemple
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) && !g.chatIsFocus {
		g.debugMode = !g.debugMode
	}

	return nil
}

func (g *game) introStateLogo() bool {
	// Incrémenter le compteur de frames
	g.stateFrame++
	totalDuration := fadeInDuration + holdDuration + fadeOutDuration + invisibleHoldDuration

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.chatIsFocus {
		return true
	}

	if g.stateFrame >= totalDuration {
		return true
	}

	return false
}

func (g *game) introStateTexte() bool {
	// Incrémenter le compteur de frames
	g.stateFrame++
	totalDuration := fadeInDuration + holdDuration + fadeOutDuration + invisibleHoldDuration

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.chatIsFocus {
		return true
	}

	if g.stateFrame >= totalDuration {
		return true
	}

	return false
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.UpdateThemes()

	g.stateFrame = g.stateFrame % globalBlinkDuration
	// Obtenir les coordonnées de la souris
	mouseX, mouseY := ebiten.CursorPosition()

	// Vérifier si le clic est dans les limites du rectangle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseReleased && !g.chatIsFocus {
		// Dimensions du rectangle "JOUER"
		smallRectWidth := 200
		smallRectHeight := 80
		smallRectX := (globalWidth - smallRectWidth) / 2
		smallRectY := (80 - smallRectHeight) / 2

		// Vérification des coordonnées
		if mouseX >= smallRectX && mouseX <= smallRectX+smallRectWidth &&
			mouseY >= smallRectY && mouseY <= smallRectY+smallRectHeight {
			// Passer à l'état suivant
			g.mouseReleased = false
			g.gameState = inputServerState // Remplacez "nextState" par l'état que vous voulez
		}
	}
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.chatIsFocus

}

func (g *game) inputServerUpdate() bool {

	// Réinitialiser l'adresse si une erreur est survenue
	if g.connectionMessage == "Erreur : Adresse incorrecte. Entrez une nouvelle adresse." {
		g.serverAddress = ""
		g.connectionMessage = ""
	}

	// Ajouter les caractères saisis
	if !g.chatIsFocus {
		for _, char := range ebiten.AppendInputChars(nil) {
			g.errorConnection = ""
			g.serverAddress += string(char)
		}
	}

	// Gérer le maintien de Backspace
	backspacePressed := inpututil.KeyPressDuration(ebiten.KeyBackspace) > 0
	if backspacePressed && !g.chatIsFocus {
		if (inpututil.KeyPressDuration(ebiten.KeyBackspace) == 1 || (inpututil.KeyPressDuration(ebiten.KeyBackspace) > 30 && inpututil.KeyPressDuration(ebiten.KeyBackspace)%3 == 0)) && !g.chatIsFocus {
			// Supprimer le dernier caractère si l'adresse n'est pas vide
			if len(g.serverAddress) > 0 {
				g.serverAddress = g.serverAddress[:len(g.serverAddress)-1]
			}
		}
	}

	// Obtenir les coordonnées de la souris
	mouseX, mouseY := ebiten.CursorPosition()

	// Vérifier si le clic est dans les limites du rectangle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseReleased && !g.chatIsFocus {
		// Dimensions du rectangle "JOUER"
		smallRectWidth := 200
		smallRectHeight := 80
		smallRectX := (globalWidth - smallRectWidth) / 2
		smallRectY := (80 - smallRectHeight) / 2

		// Vérification des coordonnées
		if mouseX >= smallRectX && mouseX <= smallRectX+smallRectWidth &&
			mouseY >= smallRectY && mouseY <= smallRectY+smallRectHeight {
			// Passer à l'état suivant
			if g.serverAddress == "" {
				g.serverAddress = "localhost:8080"
			}
			g.mouseReleased = false
			g.isReset = false
			return true
		}
	}

	// Valider l'adresse (Entrée)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.chatIsFocus {
		if g.serverAddress == "" {
			g.serverAddress = "localhost:8080"
		}
		g.isReset = false
		return true
	}
	return false
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {

	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && !g.chatIsFocus {
		g.errorMessage = ""
		col = (col + 1) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && !g.chatIsFocus {
		g.errorMessage = ""
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) && !g.chatIsFocus {
		g.errorMessage = ""
		line = (line + 1) % globalNumColorLine
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) && !g.chatIsFocus {
		g.errorMessage = ""
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
	}

	g.p1Color = line*globalNumColorLine + col
	sendCursorUpdateToServer(g.conn, g.p1Color)

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.chatIsFocus {
		if g.p1Color == g.p2Color {
			log.Println("Erreur : Couleur déjà sélectionnée par l'autre joueur. Choisissez une autre couleur.")
			g.errorMessage = "Couleur deja choisi par l'autre joueur"
			return false // Ne pas permettre de continuer
		}
		// Envoyer la couleur au serveur
		if g.conn != nil {
			sendColorToServer(g.conn, g.p1Color)
			g.p1ColorValidate = g.p1Color
		}
		if g.p2Color != -1 {
			return true
		} else {
			return false
		}
	}
	return false

}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate() {
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) && !g.chatIsFocus {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
		sendTokenUpdateToServer(g.conn, g.tokenPosition)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) && !g.chatIsFocus {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
		sendTokenUpdateToServer(g.conn, g.tokenPosition)
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if (inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter)) && !g.chatIsFocus {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			g.turn = p2Turn
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos

			// Envoyer la position au serveur
			if g.conn != nil {
				sendMoveToServer(g.conn, lastXPositionPlayed, lastYPositionPlayed)
			}
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion de la position du prochain pion joué par le joueur 2 et
// du moment où ce pion est joué.
func (g *game) p2Update() (int, int) {
	// Ne fait rien, attend que le serveur mette à jour la grille
	return -1, -1
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g *game) resultUpdate() bool {
	g.UpdateReplayButton()
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && !g.chatIsFocus {
		return true
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Vérifier si le clic est dans les limites du rectangle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseReleased && !g.chatIsFocus {
		// Dimensions du rectangle "JOUER"
		smallRectWidth := 200
		smallRectHeight := 80
		smallRectX := (globalWidth - smallRectWidth) / 2
		smallRectY := (80 - smallRectHeight) / 2

		// Vérification des coordonnées
		if mouseX >= smallRectX && mouseX <= smallRectX+smallRectWidth &&
			mouseY >= smallRectY && mouseY <= smallRectY+smallRectHeight {
			// Passer à l'état suivant

			g.mouseReleased = false
			return true // Remplacez "nextState" par l'état que vous voulez
		}
	}
	return false
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

func (g *game) updateGridExtend(id, positionX int, positionY int) (updated bool) {
	token := noToken
	if id == g.playerID {
		token = p1Token
	} else {
		token = p2Token
	}

	if g.grid[positionX][positionY] == noToken {
		updated = true
		g.grid[positionX][positionY] = token
		return
	}
	return
}

func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int, winningPositions [][2]int) {

	tokenType := g.grid[xPos][yPos]

	// Horizontal
	count := 0
	var tempPositions [][2]int

	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
		tempPositions = append(tempPositions, [2]int{x, yPos})
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
		tempPositions = append(tempPositions, [2]int{x, yPos})
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins, tempPositions
		}
		return true, p2wins, tempPositions
	}

	// Vertical
	count = 0
	tempPositions = nil
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
		tempPositions = append(tempPositions, [2]int{xPos, y})
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins, tempPositions
		}
		return true, p2wins, tempPositions
	}

	// Diagonal haut gauche / bas droit
	count = 0
	tempPositions = nil
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
		tempPositions = append(tempPositions, [2]int{x, y})
	}
	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
		tempPositions = append(tempPositions, [2]int{x, y})
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins, tempPositions
		}
		return true, p2wins, tempPositions
	}

	// Diagonal haut droit / bas gauche
	count = 0
	tempPositions = nil
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
		tempPositions = append(tempPositions, [2]int{x, y})
	}
	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
		tempPositions = append(tempPositions, [2]int{x, y})
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins, tempPositions
		}
		return true, p2wins, tempPositions
	}

	// Égalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality, nil
	}

	return
}

func (g *game) replayDrawUpdate() bool {
	mouseX, mouseY := ebiten.CursorPosition()

	// Vérifier si le clic est dans les limites du rectangle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseReleased && !g.chatIsFocus {
		// Dimensions du rectangle "JOUER"
		smallRectWidth := 200
		smallRectHeight := 80
		smallRectX := (globalWidth - smallRectWidth) / 2
		smallRectY := (80 - smallRectHeight) / 2

		// Vérification des coordonnées
		if mouseX >= smallRectX && mouseX <= smallRectX+smallRectWidth &&
			mouseY >= smallRectY && mouseY <= smallRectY+smallRectHeight {
			// Passer à l'état suivant
			g.mouseReleased = false
			return true
		}
	}

	// Gestion du délai entre les mises à jour
	const delay = 350 // Délai en millisecondes (500 ms = 0,5 s)
	if g.lastUpdateTime == 0 {
		g.lastUpdateTime = int(time.Now().UnixMilli())
	}

	currentTime := int(time.Now().UnixMilli())
	if currentTime-g.lastUpdateTime >= delay {
		// Identifier la clé la plus basse
		var minKey int
		found := false
		for key := range history {
			if !found || key < minKey {
				minKey = key
				found = true
			}
		}

		// Effectuer la mise à jour si une clé est trouvée
		if found {
			value := history[minKey]
			g.updateGridExtend(value.ID, value.X, value.Y)
			delete(history, minKey) // Supprimer après traitement
		} else if len(history) == 0 && g.posWinner != nil { // Si tout est affiché et pas encore clignotant
			g.blinking = true // Lancer le clignotement des gagnants
		}

		g.lastUpdateTime = currentTime
	}

	return false

}

func (g *game) updateChat() {

	if !g.chatIsFocus {
		return // Ignorer les événements clavier si le chat n'est pas en focus
	}

	for _, key := range ebiten.InputChars() {
		// Prévisualiser la saisie avec le caractère ajouté
		testInput := g.chatInput + string(key)
		if len(strings.Split(testInput, "\n")) <= 2 && text.BoundString(smallFont, testInput).Dx() <= inputMaxWidth {
			g.chatInput = testInput
		}
	}

	// Gérer le maintien de Backspace
	backspaceDuration := inpututil.KeyPressDuration(ebiten.KeyBackspace)
	if backspaceDuration > 0 {
		// Supprimer immédiatement un caractère si la touche vient d'être pressée
		if backspaceDuration == 1 || (backspaceDuration > 30 && backspaceDuration%3 == 0) {
			// Supprimer un caractère si l'entrée n'est pas vide
			if len(g.chatInput) > 0 {
				g.chatInput = g.chatInput[:len(g.chatInput)-1]
			}
		}
	}

	// Envoyer le message si Enter est pressé
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(g.chatInput) > 0 && g.chatIsFocus {
		sendChatMessage(g.conn, g.chatInput)
		g.chatInput = "" // Réinitialiser l'entrée
	}

}

func (g *game) updateChatButton() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.chatIsFocus {
		x, y := ebiten.CursorPosition()
		var icon *ebiten.Image = chat
		// Facteur d'échelle pour redimensionner le logo
		iconScale := 0.06 // Réduire à 50% de la taille d'origine

		// Calculer la taille de l'icône redimensionnée
		originalWidth, originalHeight := icon.Size()
		scaledWidth := float64(originalWidth) * iconScale
		scaledHeight := float64(originalHeight) * iconScale

		// Définir la position du bouton (centré horizontalement en haut à droite)
		buttonX := 40
		buttonY := globalHeight - 85

		// Vérifier si le clic est sur le bouton
		if x >= buttonX && x <= buttonX+int(scaledWidth) && y >= buttonY && y <= buttonY+int(scaledHeight) {
			g.chatIsFocus = true
			g.chatNewMessage = false
		}
	}
	return nil
}

func (g *game) updateCloseButton() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.chatIsFocus {
		x, y := ebiten.CursorPosition()
		var icon *ebiten.Image = close
		// Facteur d'échelle pour redimensionner le logo
		iconScale := 0.06 // Réduire à 50% de la taille d'origine

		// Calculer la taille de l'icône redimensionnée
		originalWidth, originalHeight := icon.Size()
		scaledWidth := float64(originalWidth) * iconScale
		scaledHeight := float64(originalHeight) * iconScale

		// Définir la position du bouton (centré horizontalement en haut à droite)
		buttonX := inputMaxWidth - 120
		buttonY := globalHeight - 390

		// Vérifier si le clic est sur le bouton
		if x >= buttonX && x <= buttonX+int(scaledWidth) && y >= buttonY && y <= buttonY+int(scaledHeight) {
			g.chatIsFocus = false
		}
	}
	return nil
}

func (g *game) UpdateShifumi() {
	// Gestion du clic de souris pour la sélection
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.handleMouseClick(x, y)

		// Si un choix a été fait, envoyer au serveur
		if g.selected != "" {
			message := Message{
				Type: "selected",
				Payload: SelectedPayload{
					Selected: strings.ToLower(g.selected),
				},
			}
			jsonData, err := json.Marshal(message)
			if err == nil {
				fmt.Fprintf(g.conn, string(jsonData)+"\n")
			}
		}
	}

	// Gestion du timer pour l'affichage du résultat
	if g.showShifumiResult {
		g.shifumiResultTimer--
		if g.shifumiResultTimer <= 0 {
			g.gameState = playState
			g.showShifumiResult = false
			g.selected = ""
			g.adversaryChoice = ""
			g.shifumiResult = ""
		}
	}
}

func (g *game) handleShifumiResult(message Message) {
	payload, ok := message.Payload.(map[string]interface{})
	if !ok {
		return
	}

	// Récupérer les informations des joueurs
	player1Data := payload["player1"].(map[string]interface{})
	player2Data := payload["player2"].(map[string]interface{})

	// Récupérer l'ID et la sélection de l'adversaire
	var adversarySelection string
	if int(player1Data["id"].(float64)) == g.playerID {
		adversarySelection = player2Data["selection"].(string)
	} else {
		adversarySelection = player1Data["selection"].(string)
	}
	g.adversaryChoice = adversarySelection

	// Déterminer le résultat
	if message.Type == "shifumi_result" {
		g.shifumiResult = "Égalité"
	} else if message.Type == "shifumi_complete" {
		winnerID := int(payload["winner"].(float64))
		if winnerID == g.playerID {
			g.shifumiResult = "Gagné"
		} else {
			g.shifumiResult = "Perdu"
		}
	}

	g.showShifumiResult = true
	g.shifumiResultTimer = 120 // 2 secondes à 60 FPS
}

const globalStatePlay = playState
