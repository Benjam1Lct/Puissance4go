package main

import (
	"io"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func loadFont(path string, size float64, scaleFactor float64) font.Face {
	// Ouvrir le fichier de police
	fontFile, err := os.Open(path)
	if err != nil {
		log.Fatalf("Erreur lors du chargement de la police %s: %v", path, err)
	}
	defer fontFile.Close()

	// Lire le contenu du fichier
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier %s: %v", path, err)
	}

	// Parser les données de la police
	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Erreur lors du parsing de la police %s: %v", path, err)
	}

	dpi := 72.0

	// Créer une Face avec la taille ajustée
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size: size * scaleFactor, // Taille ajustée
		DPI:  dpi,
	})
	if err != nil {
		log.Fatalf("Erreur lors de la création de la Face pour %s: %v", path, err)
	}

	return face
}

func initFonts() {

	// Calculez un facteur d'échelle basé sur la largeur de la fenêtre
	scaleFactor := 1.0

	baseFontSizeError = 20.0
	baseFontSizeSmall = 30.0
	baseFontSizeLarge = 50.0 // Taille de base pour les grandes polices
	baseFontSizeBig = 75.0
	baseFontTitle = 1.0
	if globalWidth > 1920 {
		baseFontTitle = 140.0
	} else {
		baseFontTitle = 100
	}

	// Initialiser chaque style avec deux tailles
	lightFontSmall = loadFont("assets/fonts/Hind/HindMysuru-Light.ttf", baseFontSizeSmall, scaleFactor)
	lightFontLarge = loadFont("assets/fonts/Hind/HindMysuru-Light.ttf", baseFontSizeLarge, scaleFactor)

	regularFontSmall = loadFont("assets/fonts/Hind/HindMysuru-Regular.ttf", baseFontSizeSmall, scaleFactor)
	regularFontLarge = loadFont("assets/fonts/Hind/HindMysuru-Regular.ttf", baseFontSizeLarge, scaleFactor)

	mediumFontSmall = loadFont("assets/fonts/Hind/HindMysuru-Medium.ttf", baseFontSizeSmall, scaleFactor)
	mediumFontLarge = loadFont("assets/fonts/Hind/HindMysuru-Medium.ttf", baseFontSizeLarge, scaleFactor)
	mediumFontError = loadFont("assets/fonts/Hind/HindMysuru-Medium.ttf", baseFontSizeError, scaleFactor)

	semiboldFontSmall = loadFont("assets/fonts/Hind/HindMysuru-SemiBold.ttf", baseFontSizeSmall, scaleFactor)
	semiboldFontLarge = loadFont("assets/fonts/Hind/HindMysuru-SemiBold.ttf", baseFontSizeLarge, scaleFactor)

	boldFontSmall = loadFont("assets/fonts/Hind/HindMysuru-Bold.ttf", baseFontSizeSmall, scaleFactor)
	boldFontLarge = loadFont("assets/fonts/Hind/HindMysuru-Bold.ttf", baseFontSizeLarge, scaleFactor)
	bigFontResult = loadFont("assets/fonts/Hind/HindMysuru-Bold.ttf", baseFontSizeBig, scaleFactor)

	firstTitleMinusFont = loadFont("assets/fonts/Valorant.ttf", baseFontSizeSmall, scaleFactor)
	firstTitleFont = loadFont("assets/fonts/Valorant.ttf", baseFontTitle, scaleFactor)
	firstTitleSmallFont = loadFont("assets/fonts/Valorant.ttf", baseFontSizeBig, scaleFactor)
	firstTitleSmallerFont = loadFont("assets/fonts/Valorant.ttf", baseFontSizeLarge, scaleFactor)

	smallFont = mediumFontSmall
	largeFont = boldFontLarge
}

// Mise en place des polices d'écritures utilisées pour l'affichage.
func initImage() {

	// Charger l'icône pour le mode plein écran
	fullscreenFile, err := os.Open("assets/images/fullscreen.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(fullscreenFile *os.File) {
		err := fullscreenFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(fullscreenFile)
	iconFullscreen, _, err = ebitenutil.NewImageFromReader(fullscreenFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône pour le mode fenêtré
	windowedFile, err := os.Open("assets/images/windowed.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(windowedFile *os.File) {
		err := windowedFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(windowedFile)
	iconWindowed, _, err = ebitenutil.NewImageFromReader(windowedFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône de l'iut
	iutFile, err := os.Open("assets/images/logoIUT-Blanc.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(iutFile *os.File) {
		err := iutFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(iutFile)
	iconIUT, _, err = ebitenutil.NewImageFromReader(iutFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône du studio
	logoStudioFile, err := os.Open("assets/images/logoStudio.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(logoStudioFile *os.File) {
		err := logoStudioFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(logoStudioFile)
	logoStudio, _, err = ebitenutil.NewImageFromReader(logoStudioFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône de son on
	replayFile, err := os.Open("assets/images/replay.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(replayFile *os.File) {
		err := replayFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(replayFile)
	replay, _, err = ebitenutil.NewImageFromReader(replayFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône des themes
	ThemesFile, err := os.Open("assets/images/theme.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(ThemesFile *os.File) {
		err := ThemesFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(ThemesFile)
	Themes, _, err = ebitenutil.NewImageFromReader(ThemesFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône des themes
	leftArrowImageFile, err := os.Open("assets/images/left.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(leftArrowImageFile *os.File) {
		err := leftArrowImageFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(leftArrowImageFile)
	leftArrowImage, _, err = ebitenutil.NewImageFromReader(leftArrowImageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône des themes
	rightArrowImageFile, err := os.Open("assets/images/right.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(rightArrowImageFile *os.File) {
		err := rightArrowImageFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rightArrowImageFile)
	rightArrowImage, _, err = ebitenutil.NewImageFromReader(rightArrowImageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône du chat
	chatFile, err := os.Open("assets/images/chat.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(chatFile *os.File) {
		err := chatFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(chatFile)
	chat, _, err = ebitenutil.NewImageFromReader(chatFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône du chat
	chatWarningFile, err := os.Open("assets/images/chatWarning.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(chatWarningFile *os.File) {
		err := chatWarningFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(chatWarningFile)
	chatWarning, _, err = ebitenutil.NewImageFromReader(chatWarningFile)
	if err != nil {
		log.Fatal(err)
	}

	// Charger l'icône pour close
	closeFile, err := os.Open("assets/images/close.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(closeFile *os.File) {
		err := closeFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(closeFile)
	close, _, err = ebitenutil.NewImageFromReader(closeFile)
	if err != nil {
		log.Fatal(err)
	}
	pierreImage, err := os.Open("assets/images/pierre.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(pierreImage *os.File) {
		err := pierreImage.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(pierreImage)
	pierreImg, _, err = ebitenutil.NewImageFromReader(pierreImage)
	if err != nil {
		log.Fatal(err)
	}
	papierImage, err := os.Open("assets/images/papier.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(papierImage *os.File) {
		err := papierImage.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(papierImage)
	papierImg, _, err = ebitenutil.NewImageFromReader(papierImage)
	if err != nil {
		log.Fatal(err)
	}
	// Charger ciseaux image
	ciseauxImage, err := os.Open("assets/images/ciseaux.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(ciseauxImage *os.File) {
		err := ciseauxImage.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(ciseauxImage)
	ciseauxImg, _, err = ebitenutil.NewImageFromReader(ciseauxImage)
	if err != nil {
		log.Fatal(err)
	}
}

func initBackground() {
	// Charger un fond d'ecran
	background1File, err := os.Open("assets/images/background/background1.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(background1File *os.File) {
		err := background1File.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(background1File)
	background1, _, err = ebitenutil.NewImageFromReader(background1File)
	if err != nil {
		log.Fatal(err)
	}

	// Charger un fond d'ecran
	background2File, err := os.Open("assets/images/background/background2.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(background2File *os.File) {
		err := background2File.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(background2File)
	background2, _, err = ebitenutil.NewImageFromReader(background2File)
	if err != nil {
		log.Fatal(err)
	}

	// Charger un fond d'ecran
	background3File, err := os.Open("assets/images/background/background3.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(background3File *os.File) {
		err := background3File.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(background3File)
	background3, _, err = ebitenutil.NewImageFromReader(background3File)
	if err != nil {
		log.Fatal(err)
	}
	//
	losangeImage, err := os.Open("assets/images/losange.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(losangeImage *os.File) {
		err := losangeImage.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(losangeImage)
	losangeImg, _, err = ebitenutil.NewImageFromReader(losangeImage)
	if err != nil {
		log.Fatal(err)
	}
}

// Création d'une image annexe pour l'affichage des résultats.
func initOffScreen() {
	offScreenImage = ebiten.NewImage(globalWidth, globalHeight)
}

func initResolution(isFullscreen bool) {
	// Récupérer la résolution de l'écran
	if isFullscreen {
		screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()
		globalHeight = screenHeight
		globalWidth = screenWidth
		if globalWidth >= 1920 && globalHeight >= 1080 {
			globalTileSize = 100
		} else {
			globalTileSize = 70
		}
		ebiten.SetFullscreen(true)
	} else {
		adjustGlobalDimensions()
	}

}

func adjustGlobalDimensions() {
	screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()

	if ebiten.IsFullscreen() {
		globalHeight = screenHeight
		globalWidth = screenWidth
	} else {
		// En mode fenêtré, ajustez en fonction de la résolution de l'écran
		if screenWidth >= 2560 && screenHeight >= 1440 {
			// Si l'écran est 2K ou plus, utilisez 1920x1080
			globalWidth = 1920
			globalHeight = 1080
		} else {
			// Si l'écran est Full HD, utilisez 1280x720
			globalWidth = 1280
			globalHeight = 720
		}

		screenWidth, screenHeight := ebiten.ScreenSizeInFullscreen()
		windowWidth, windowHeight := globalWidth, globalHeight

		// Calculer les coordonnées pour centrer la fenêtre
		centerX := (screenWidth - windowWidth) / 2
		centerY := (screenHeight - windowHeight) / 2

		// Positionner la fenêtre
		ebiten.SetWindowPosition(centerX, centerY)
	}

	if globalWidth >= 1920 && globalHeight >= 1080 {
		globalTileSize = 100
	} else {
		globalTileSize = 70
	}
}

func (g *game) initGame() {
	g.gameState = introStateLogo
	g.chatNewMessage = false
	g.stateFrame = 0
	g.restartOk = true
	g.p2Color = -1
	g.p1ColorValidate = -1
	g.playerID = -1
	g.adversaryTokenPosition = g.tokenPosition
	g.isReset = false
	g.mouseReleased = true
}
