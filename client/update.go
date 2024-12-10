package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		var icon *ebiten.Image
		if ebiten.IsFullscreen() {
			icon = iconWindowed
		} else {
			icon = iconFullscreen
		}
		// Facteur d'échelle pour redimensionner le logo
		iconScale := 0.07 // Réduire à 50% de la taille d'origine

		// Calculer la taille de l'icône redimensionnée
		originalWidth, originalHeight := icon.Size()
		scaledWidth := float64(originalWidth) * iconScale
		scaledHeight := float64(originalHeight) * iconScale

		// Définir la position du bouton (centré horizontalement en haut à droite)
		buttonX := globalWidth - int(scaledWidth) - 19
		buttonY := 8

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

	g.UpdateDebug()
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
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result

				if g.result == p1wins {
					g.turn = p2Turn
				} else if g.result == p2wins {
					g.turn = p1Turn
				} else {
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
			endGame(g)
		}
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.mouseReleased = true // Souris relâchée
	}

	return nil
}

func (g *game) UpdateDebug() error {
	// ... votre code update existant ...

	// Activer/désactiver le mode debug avec F3 par exemple
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		println("c'est fait")

		g.debugMode = !g.debugMode
		print(g.debugMode)
	}

	return nil
}

func (g *game) introStateLogo() bool {
	// Incrémenter le compteur de frames
	g.stateFrame++
	totalDuration := fadeInDuration + holdDuration + fadeOutDuration + invisibleHoldDuration

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}

	if g.stateFrame >= totalDuration {
		return true
	}

	return false
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	g.stateFrame = g.stateFrame % globalBlinkDuration
	// Obtenir les coordonnées de la souris
	mouseX, mouseY := ebiten.CursorPosition()

	// Vérifier si le clic est dans les limites du rectangle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseReleased {
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
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate() bool {

	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.errorMessage = ""
		col = (col + 1) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.errorMessage = ""
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		g.errorMessage = ""
		line = (line + 1) % globalNumColorLine
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.errorMessage = ""
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
	}

	g.p1Color = line*globalNumColorLine + col
	sendCursorUpdateToServer(g.conn, g.p1Color)

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		return true
	}

	mouseX, mouseY := ebiten.CursorPosition()

	// Vérifier si le clic est dans les limites du rectangle
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.mouseReleased {
		// Dimensions du rectangle "JOUER"
		smallRectWidth := 200
		smallRectHeight := 80
		smallRectX := (globalWidth - smallRectWidth) / 2
		smallRectY := (80 - smallRectHeight) / 2

		// Vérification des coordonnées
		if mouseX >= smallRectX && mouseX <= smallRectX+smallRectWidth &&
			mouseY >= smallRectY && mouseY <= smallRectY+smallRectHeight {
			// Passer à l'état suivant
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

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
