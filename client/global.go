package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Constantes définissant les paramètres généraux du programme.
const (
	globalTileSize        = 70
	globalNumTilesX       = 7
	globalNumTilesY       = 6
	globalCircleMargin    = 5
	globalBlinkDuration   = 60
	globalNumColorLine    = 3
	globalNumColorCol     = 3
	globalNumColor        = globalNumColorLine * globalNumColorCol
	iconScale             = 0.20
	fadeInDuration        = 200
	holdDuration          = 120
	fadeOutDuration       = 250
	invisibleHoldDuration = 100
	sampleRate            = 44100
	largeRectHeight       = 80
)

// Variables définissant les paramètres généraux du programme.
var (
	globalBackgroundColor                color.Color = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
	globalGridColor                      color.Color = color.NRGBA{R: 119, G: 136, B: 153, A: 255}
	globalTextColor                      color.Color = color.NRGBA{R: 25, G: 25, B: 25, A: 255}
	globalTextColorBright                color.Color = color.NRGBA{R: 217, G: 217, B: 217, A: 255} // Gris clair
	globalTextRed                        color.Color = color.NRGBA{R: 231, G: 65, B: 80, A: 255}   // Rouge avec transparence pleine
	globalTextDarkRed                    color.Color = color.NRGBA{R: 161, G: 45, B: 56, A: 255}   // Rouge sombre
	globalTextColorYellow                color.Color = color.NRGBA{R: 235, G: 234, B: 169, A: 255}
	globalTextColorGreen                 color.Color = color.NRGBA{R: 193, G: 255, B: 58, A: 255}
	globalSelectColor                    color.Color = color.NRGBA{R: 25, G: 25, B: 5, A: 255}
	globalP2SelectColor                  color.Color = color.NRGBA{50, 50, 50, 255}
	globalValidatorColor                 color.Color = color.NRGBA{R: 50, G: 205, B: 50, A: 255} // Vert clair
	globalValidatorDarkColor             color.Color = color.NRGBA{R: 0, G: 150, B: 0, A: 255}
	lightFontSmall, lightFontLarge       font.Face
	regularFontSmall, regularFontLarge   font.Face
	mediumFontSmall, mediumFontLarge     font.Face
	semiboldFontSmall, semiboldFontLarge font.Face
	boldFontSmall, boldFontLarge         font.Face
	largeFont                            font.Face
	smallFont                            font.Face
	mediumFontError                      font.Face
	bigFontResult                        font.Face
	firstTitleFont                       font.Face
	firstTitleSmallFont                  font.Face
	firstTitleSmallerFont                font.Face
	globalTokenColors                    [globalNumColor]color.Color = [globalNumColor]color.Color{
		color.NRGBA{R: 255, G: 239, B: 213, A: 255}, // Pêche pastel
		color.NRGBA{R: 119, G: 221, B: 153, A: 255}, // Vert menthe pastel
		color.NRGBA{R: 174, G: 238, B: 152, A: 255}, // Vert clair pastel
		color.NRGBA{R: 230, G: 220, B: 170, A: 255}, // Jaune sable pastel
		color.NRGBA{R: 255, G: 178, B: 156, A: 255}, // Orange saumon pastel
		color.NRGBA{R: 255, G: 182, B: 193, A: 255}, // Rose clair pastel
		color.NRGBA{R: 202, G: 255, B: 219, A: 255}, // Turquoise clair pastel
		color.NRGBA{R: 219, G: 178, B: 255, A: 255}, // Violet lavande pastel
		color.NRGBA{R: 245, G: 255, B: 250, A: 255}, // Blanc cassé pastel
	}

	globalBackgroundImage *ebiten.Image
	offScreenImage        *ebiten.Image
	iconFullscreen        *ebiten.Image
	iconWindowed          *ebiten.Image
	iconIUT               *ebiten.Image
	logoStudio            *ebiten.Image
	globalWidth           = 1920
	globalHeight          = 1080
	baseFontSizeError     float64
	baseFontSizeSmall     float64
	baseFontSizeLarge     float64
	baseFontSizeBig       float64
	baseFontTitle         float64
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type ColorPayload struct {
	Color int `json:"color"`
}

type MovePayload struct {
	X int `json:"x"`
	Y int `json:"y"`
}
