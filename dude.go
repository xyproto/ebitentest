package main

import (
	_ "image/png"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Dude struct {
	X, Y, AX, AY     float64
	Image            *ebiten.Image
	jumping          bool
	jumpStarted      time.Time
	jumpBlocked      bool
	inAir            bool
	inAirJumpCounter int
	prevSpaceState   bool
	fire             bool
}

func (g *Game) UpdateDude() {
	g.dude.fire = ebiten.IsKeyPressed(ebiten.KeySpace)

	jumpKeyDown := ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW)

	// Jump logic
	const maxInAirJumps = 2
	const jumpTimeout = time.Millisecond * 100

	if !jumpKeyDown {
		g.dude.jumping = false
		g.dude.jumpBlocked = false
		if !g.dude.inAir {
			g.dude.inAirJumpCounter = 0
		}
	} else if !g.dude.jumpBlocked {
		// The jump key has either been pressed right now
		// or has been held down for a while, which is it?
		if !g.dude.jumping && (g.dude.inAirJumpCounter+1) < maxInAirJumps {
			g.dude.jumping = true
			g.dude.jumpStarted = time.Now()
		} else if time.Since(g.dude.jumpStarted) >= jumpTimeout { // jump timeout
			g.dude.jumping = false
			g.dude.jumpBlocked = true
		}
	}

	// Double jump counter
	spaceToggled := jumpKeyDown != g.dude.prevSpaceState
	if g.dude.inAir && spaceToggled && jumpKeyDown {
		g.dude.inAirJumpCounter++
	}
	g.dude.prevSpaceState = jumpKeyDown

	// Are we doing the jump acceleration right now?

	const jumpHigherWhenMovingFactor = 1.2
	const jumpSpeed = 8.0

	if g.dude.jumping {
		g.dude.AY = -jumpSpeed
		g.dude.AY -= math.Abs(g.dude.AX) * jumpHigherWhenMovingFactor
	}

	duckKey := ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown)
	if duckKey {
		g.dude.AY++
	}

	// "gravity"
	g.dude.AY++

	g.dude.Y += g.dude.AY

	if g.dude.AY < 0 {
		g.dude.inAir = true
	}

	// ceiling
	if g.dude.Y <= (0 + frameHeight/2) {
		g.dude.Y = (0 + frameHeight/2)
		g.dude.AY = 0
	}

	// floor
	if g.dude.Y >= (screenHeight - frameHeight/2) {
		g.dude.Y = screenHeight - frameHeight/2
		g.dude.AY = 0
		g.dude.inAir = false
	}

	const moveSpeed = 5.0
	const moveSpeedDampen = 0.7

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if !g.dude.inAir {
			g.dude.AX = -moveSpeed
		} else {
			g.dude.AX = (-moveSpeed)*0.7 + g.dude.AX*0.3
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if !g.dude.inAir {
			g.dude.AX = moveSpeed
		} else {
			g.dude.AX = moveSpeed*0.7 + g.dude.AX*0.3
		}
	}

	if duckKey {
		g.dude.AX *= 0.7
	}

	g.dude.X += g.dude.AX
	g.dude.AX *= moveSpeedDampen

	// left wall
	if g.dude.X <= (0 + frameWidth/2) {
		g.dude.X = (0 + frameWidth/2)
		g.dude.AX = 0
	}

	// right wall
	if g.dude.X >= (screenWidth - frameWidth/2) {
		g.dude.X = (screenWidth - frameWidth/2)
		g.dude.AX = 0
	}
}
