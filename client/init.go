package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"io"
	"log"
	"os"
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

	// Charger l'icône pour le mode fenêtré
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

	// Charger l'icône pour le mode fenêtré
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
}

func initBackground() {
	// Charger l'icône pour le mode plein écran
	background, err := os.Open("assets/images/background.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func(fullscreenFile *os.File) {
		err := fullscreenFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(background)
	globalBackgroundImage, _, err = ebitenutil.NewImageFromReader(background)
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
}

func (g *game) initAudio() {
	// Créer un contexte audio
	audioContext := audio.NewContext(sampleRate)
	g.audioContext = audioContext

	// Charger le fichier audio
	file, err := os.Open("assets/musics/bo.mp3") // Chemin vers votre fichier audio
	if err != nil {
		log.Fatalf("Impossible d'ouvrir le fichier audio: %v", err)
	}

	// Décoder le fichier audio
	stream, err := mp3.DecodeWithSampleRate(sampleRate, file)
	if err != nil {
		log.Fatalf("Impossible de décoder le fichier audio: %v", err)
	}

	// Créer un lecteur en boucle
	loop := audio.NewInfiniteLoop(stream, stream.Length())
	player, err := audioContext.NewPlayer(loop)
	if err != nil {
		log.Fatalf("Impossible de créer le lecteur audio: %v", err)
	}

	// Lancer la musique
	player.Play()

	// Stocker le lecteur dans la struct pour une utilisation future
	g.audioPlayer = player

	// Stocker le fichier pour le fermer proprement plus tard
	g.audioFile = file
}

func (g *game) closeAudio() {
	if g.audioFile != nil {
		g.audioFile.Close()
	}
	if g.audioPlayer != nil {
		g.audioPlayer.Close()
	}
}
