package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// Constantes définissant les paramètres généraux du programme.
const (
	globalWidth         = globalNumTilesX * globalTileSize
	globalHeight        = (globalNumTilesY + 1) * globalTileSize
	globalTileSize      = 100
	globalNumTilesX     = 7
	globalNumTilesY     = 6
	globalCircleMargin  = 5
	globalBlinkDuration = 60
	globalNumColorLine  = 3
	globalNumColorCol   = 3
	globalNumColor      = globalNumColorLine * globalNumColorCol
)

// Variables définissant les paramètres généraux du programme.
var (
	globalBackgroundColor    color.Color = color.NRGBA{R: 176, G: 196, B: 222, A: 255}
	globalGridColor          color.Color = color.NRGBA{R: 119, G: 136, B: 153, A: 255}
	globalTextColor          color.Color = color.NRGBA{R: 25, G: 25, B: 5, A: 255}
	globalSelectColor        color.Color = color.NRGBA{R: 25, G: 25, B: 5, A: 255}
	globalP2SelectColor      color.Color = color.RGBA{100, 100, 100, 255}
	globalValidatorColor     color.Color = color.NRGBA{R: 50, G: 205, B: 50, A: 255} // Vert clair
	globalValidatorDarkColor color.Color = color.NRGBA{R: 0, G: 100, B: 0, A: 255}   // Vert foncé
	smallFont                font.Face
	largeFont                font.Face
	globalTokenColors        [globalNumColor]color.Color = [globalNumColor]color.Color{
		color.NRGBA{R: 255, G: 239, B: 213, A: 255},
		color.NRGBA{R: 60, G: 179, B: 113, A: 255},
		color.NRGBA{R: 154, G: 205, B: 50, A: 255},
		color.NRGBA{R: 189, G: 183, B: 107, A: 255},
		color.NRGBA{R: 255, G: 127, B: 80, A: 255},
		color.NRGBA{R: 240, G: 128, B: 128, A: 255},
		color.NRGBA{R: 152, G: 251, B: 152, A: 255},
		color.NRGBA{R: 221, G: 160, B: 221, A: 255},
		color.NRGBA{R: 245, G: 255, B: 250, A: 255},
	}
	offScreenImage *ebiten.Image
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
