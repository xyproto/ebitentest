package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"strings"
	"time"

	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/keyboard/keyboard"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/images"
	rkeyboard "github.com/hajimehoshi/ebiten/v2/examples/resources/images/keyboard"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
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

var (
	dudeX  = float64(screenWidth) / 2.0
	dudeY  = float64(screenHeight) / 2.0
	dudeAX = 0.0
	dudeAY = 0.0

	runnerImage   *ebiten.Image
	keyboardImage *ebiten.Image
)

type Game struct {
	count            int
	keys             []ebiten.Key
	jumping          bool
	jumpStarted      time.Time
	jumpBlocked      bool
	inAir            bool
	inAirJumpCounter int
	prevSpaceState   bool
	fire             bool
}

func (g *Game) UpdateDude() {
	g.fire = ebiten.IsKeyPressed(ebiten.KeySpace)

	jumpKeyDown := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)

	// Jump logic
	const maxInAirJumps = 2
	const jumpTimeout = time.Millisecond * 100

	if !jumpKeyDown {
		g.jumping = false
		g.jumpBlocked = false
		if !g.inAir {
			g.inAirJumpCounter = 0
		}
	} else if !g.jumpBlocked {
		// The jump key has either been pressed right now
		// or has been held down for a while, which is it?
		if !g.jumping && (g.inAirJumpCounter+1) < maxInAirJumps {
			g.jumping = true
			g.jumpStarted = time.Now()
		} else if time.Since(g.jumpStarted) >= jumpTimeout { // jump timeout
			g.jumping = false
			g.jumpBlocked = true
		}
	}

	// Double jump counter
	spaceToggled := jumpKeyDown != g.prevSpaceState
	if g.inAir && spaceToggled && jumpKeyDown {
		g.inAirJumpCounter++
	}
	g.prevSpaceState = jumpKeyDown

	// Are we doing the jump acceleration right now?

	const jumpHigherWhenMovingFactor = 1.2
	const jumpSpeed = 8.0

	if g.jumping {
		dudeAY = -jumpSpeed
		dudeAY -= math.Abs(dudeAX) * jumpHigherWhenMovingFactor
	}

	duckKey := ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)
	if duckKey {
		dudeAY++
	}

	// "gravity"
	dudeAY++

	dudeY += dudeAY

	if dudeAY < 0 {
		g.inAir = true
	}

	// ceiling
	if dudeY <= (0 + frameHeight/2) {
		dudeY = (0 + frameHeight/2)
		dudeAY = 0
	}

	// floor
	if dudeY >= (screenHeight - frameHeight/2) {
		dudeY = screenHeight - frameHeight/2
		dudeAY = 0
		g.inAir = false
	}

	const moveSpeed = 5.0
	const moveSpeedDampen = 0.7

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if !g.inAir {
			dudeAX = -moveSpeed
		} else {
			dudeAX = (-moveSpeed)*0.7 + dudeAX*0.3
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if !g.inAir {
			dudeAX = moveSpeed
		} else {
			dudeAX = moveSpeed*0.7 + dudeAX*0.3
		}
	}

	if duckKey {
		dudeAX *= 0.7
	}

	dudeX += dudeAX
	dudeAX *= moveSpeedDampen

	// left wall
	if dudeX <= (0 + frameWidth/2) {
		dudeX = (0 + frameWidth/2)
		dudeAX = 0
	}

	// right wall
	if dudeX >= (screenWidth - frameWidth/2) {
		dudeX = (screenWidth - frameWidth/2)
		dudeAX = 0
	}
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
	op.GeoM.Translate(dudeX, dudeY)
	i := (g.count / 5) % frameCount
	sx, sy := frameOffsetX+i*frameWidth, frameOffsetY
	screen.DrawImage(runnerImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)

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

	if g.jumping {
		text.Draw(screen, "JUMP", bitmapfont.Face, 8, 24, color.RGBA{0xff, 0, 0, 0xff})
	}

	if g.fire {
		text.Draw(screen, "FIRE", bitmapfont.Face, 8, 36, color.RGBA{0xff, 0, 0, 0xff})
	}

	if g.inAirJumpCounter > 0 {
		text.Draw(screen, fmt.Sprintf("DOUBLE JUMP %d", g.inAirJumpCounter), bitmapfont.Face, 8, 48, color.RGBA{0xff, 0xff, 0, 0xff})
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// The animated dude
	img, _, err := image.Decode(bytes.NewReader(images.Runner_png))
	if err != nil {
		log.Fatal(err)
	}
	runnerImage = ebiten.NewImageFromImage(img)

	// The keyboard image
	img, _, err = image.Decode(bytes.NewReader(rkeyboard.Keyboard_png))
	if err != nil {
		log.Fatal(err)
	}
	keyboardImage = ebiten.NewImageFromImage(img)

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation + Keyboard")

	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
