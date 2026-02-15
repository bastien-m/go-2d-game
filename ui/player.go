package ui

import (
	"fmt"
	"image/color"

	"github.com/bastien-m/mario/engine"
	"github.com/bastien-m/mario/engine/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	hitboxOffsetX = 5.0
	hitboxOffsetY = 6.0
)

type Player struct {
	x, y, vX, vY float64
}

func VMax() float64 {
	return constants.AssetSize * 1.0 / ebiten.ActualTPS() / 2
}

func (p *Player) moveForward(g *Game) {
	if !p.hit(g) {
		if p.vX <= VMax() {
			p.vX += VMax() * 0.1
		}
	}
}

func (p *Player) moveBackward(g *Game) {
	if !p.hit(g) {
		if p.vX >= -VMax() {
			p.vX -= VMax() * 0.1
		}
	}
}

func (p *Player) jump() {
	p.vY += 0.05
}

func (p *Player) hit(g *Game) bool {
	collision := false
	// handle horizontal collision
	if p.vX != 0 {
		if p.vX > 0 {
			fmt.Printf("[DEBUG] p.x: %f p.vX: %f\n", p.x, p.vX)
			tile := g.level.TileAt(p.x-(hitboxOffsetX/constants.AssetSize)+p.vX, p.y, engine.RIGHT)
			// we found a tile and it is not air and not our own player
			// collision
			fmt.Printf("[DEBUG] tile is %d\n", tile)
			if tile != -1 && tile != engine.Air && tile != engine.InitialPosition {
				p.x -= constants.AssetSize * 0.01
				p.vX = 0
				collision = true
			}
		} else if p.vX < 0 {
			tile := g.level.TileAt(p.x+(hitboxOffsetX/constants.AssetSize)-p.vX, p.y, engine.LEFT)
			// we found a tile and it is not air
			// collision
			if tile != -1 && tile != engine.Air && tile != engine.InitialPosition {
				p.vX = 0
				p.x += constants.AssetSize * 0.01
				collision = true
			}
		}
	}
	// handle vertical collision
	if p.vY != 0 {
	}
	return collision
}

func (g *Game) updatePlayer() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.jump()
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		fmt.Printf("hitbox: x: %f y: %f w: %f h: %f", g.player.x+hitboxOffsetX/constants.AssetSize, g.player.y+hitboxOffsetY, float64(constants.AssetSize)-(hitboxOffsetY/constants.AssetSize)*2, float64(constants.AssetSize)-(hitboxOffsetY/constants.AssetSize)*2)
		g.playerIcon = g.tileset[int(engine.GoingLeft)]
		g.player.moveBackward(g)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		fmt.Printf("hitbox: x: %f y: %f w: %f h: %f", g.player.x+hitboxOffsetX/constants.AssetSize, g.player.y+hitboxOffsetY, float64(constants.AssetSize)-(hitboxOffsetY/constants.AssetSize)*2, float64(constants.AssetSize)-(hitboxOffsetY/constants.AssetSize)*2)
		g.playerIcon = g.tileset[int(engine.GoingRight)]
		g.player.moveForward(g)
	}
}

func (g *Game) drawPlayer(screen *ebiten.Image) {
	tile := g.playerIcon

	imgOpts := &ebiten.DrawImageOptions{}

	if g.player.vX != 0 {
		if !g.player.hit(g) {
			g.player.x += g.player.vX
			if g.player.vX > 0 {
				if g.player.vX-VMax()*0.05 < 0 {
					g.player.vX = 0
				} else {
					g.player.vX -= VMax() * 0.05
				}
			} else {
				if g.player.vX+VMax()*0.05 > 0 {
					g.player.vX = 0
				} else {
					g.player.vX += VMax() * 0.05
				}
			}
		}
	}

	imgOpts.GeoM.Translate(float64(g.player.x)*constants.AssetSize+g.player.vX, float64(g.player.y)*constants.AssetSize)
	screen.DrawImage(tile, imgOpts)

	// draw hitbox
	drawRectOutline(screen,
		g.player.x*constants.AssetSize+hitboxOffsetX,
		g.player.y*constants.AssetSize+hitboxOffsetY,
		float64(tile.Bounds().Dx())-hitboxOffsetY*2,
		float64(tile.Bounds().Dy())-hitboxOffsetY*2,
		color.RGBA{255, 0, 0, 255})

	imgOpts.GeoM.Reset()
}

func drawRectOutline(screen *ebiten.Image, x, y, w, h float64, clr color.Color) {
	// Top
	drawLine(screen, x, y, w, 1, clr)
	// Bottom
	drawLine(screen, x, y+h-1, w, 1, clr)
	// Left
	drawLine(screen, x, y, 1, h, clr)
	// Right
	drawLine(screen, x+w-1, y, 1, h, clr)
}

func drawLine(screen *ebiten.Image, x, y, w, h float64, clr color.Color) {
	pixel := ebiten.NewImage(1, 1)
	pixel.Fill(color.White)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(w, h)
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	screen.DrawImage(pixel, op)
}
