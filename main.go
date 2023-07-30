package main

import (
	"bytes"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	rkeyboard "github.com/hajimehoshi/ebiten/v2/examples/resources/images/keyboard"
)

const (
	screenWidth  = 320
	screenHeight = 240

	frameOffsetX = 0
	frameOffsetY = 32

	frameWidth  = 32
	frameHeight = 32
	frameCount  = 8
)

var keyboardImage *ebiten.Image

func main() {
	var dude Dude
	dude.X = float64(screenWidth) / 2.0
	dude.Y = float64(screenHeight) / 2.0
	dude.AX = 0.0
	dude.AY = 0.0

	// The animated dude
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	dude.Image = ebiten.NewImageFromImage(img)

	// The keyboard image
	img, _, err = image.Decode(bytes.NewReader(rkeyboard.Keyboard_png))
	if err != nil {
		log.Fatal(err)
	}
	keyboardImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation + Keyboard")

	var game Game
	game.dude = dude

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
