package ui

import (
	"github.com/bastien-m/mario/engine"
	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	x, y, vX, vY float64
}

func (p *Player) moveForward() {
	if p.vX <= 0.2 {
		p.vX += 0.01
	}
}

func (p *Player) moveBackward() {
	if p.vX >= -0.2 {
		p.vX -= 0.01
	}
}

func (p *Player) jump() {
	p.vY += 0.05
}

func (g *Game) updatePlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.jump()
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.playerIcon = g.tileset[int(engine.GoingLeft)]
		g.player.moveBackward()
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.playerIcon = g.tileset[int(engine.GoingRight)]
		g.player.moveForward()
	}
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	tile := g.playerIcon

	imgOpts := &ebiten.DrawImageOptions{}

	if g.player.vX != 0 {
		if g.player.vX > 0 {
			g.player.x += g.player.vX
			if g.player.vX-0.001 < 0 {
				g.player.vX = 0
			} else {
				g.player.vX -= 0.001
			}
		} else {
			g.player.x += g.player.vX
			if g.player.vX+0.001 > 0 {
				g.player.vX = 0
			} else {
				g.player.vX += 0.001
			}
		}
	}

	imgOpts.GeoM.Translate(float64(g.player.x)*assetSize, float64(g.player.y)*assetSize)
	screen.DrawImage(tile, imgOpts)

	imgOpts.GeoM.Reset()
}
