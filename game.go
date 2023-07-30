package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"strings"

	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/keyboard/keyboard"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	count int
	keys  []ebiten.Key
	dude  Dude
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("quit")
	}

	g.count++

	g.keys = inpututil.AppendPressedKeys(g.keys[:0])

	g.UpdateDude()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the little dude
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Translate(g.dude.X, g.dude.Y)
	i := (g.count / 5) % frameCount
	sx, sy := frameOffsetX+i*frameWidth, frameOffsetY
	screen.DrawImage(g.dude.Image.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

	const (
		offsetX = 24
		offsetY = 40
	)

	// Draw the base (grayed) keyboard image.
	op = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(offsetX, offsetY)
	op.ColorScale.Scale(0.5, 0.5, 0.5, 1)
	screen.DrawImage(keyboardImage, op)

	// Draw the highlighted keys.
	op = &ebiten.DrawImageOptions{}
	for _, p := range g.keys {
		op.GeoM.Reset()
		r, ok := keyboard.KeyRect(p)
		if !ok {
			continue
		}
		op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(keyboardImage.SubImage(r).(*ebiten.Image), op)
	}

	var keyStrs []string
	var keyNames []string
	for _, k := range g.keys {
		keyStrs = append(keyStrs, k.String())
		if name := ebiten.KeyName(k); name != "" {
			keyNames = append(keyNames, name)
		}
	}

	// Use bitmapfont.Face instead of ebitenutil.DebugPrint, since some key names might not be printed with DebugPrint.
	text.Draw(screen, strings.Join(keyStrs, ", ")+"\n"+strings.Join(keyNames, ", "), bitmapfont.Face, 8, 12, color.White)

	if g.dude.jumping {
		text.Draw(screen, "JUMP", bitmapfont.Face, 8, 24, color.RGBA{0xff, 0, 0, 0xff})
	}

	if g.dude.fire {
		text.Draw(screen, "FIRE", bitmapfont.Face, 8, 36, color.RGBA{0xff, 0, 0, 0xff})
	}

	if g.dude.inAirJumpCounter > 0 {
		text.Draw(screen, fmt.Sprintf("DOUBLE JUMP %d", g.dude.inAirJumpCounter), bitmapfont.Face, 8, 48, color.RGBA{0xff, 0xff, 0, 0xff})
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
