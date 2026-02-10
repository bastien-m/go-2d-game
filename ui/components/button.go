package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Button struct {
	X, Y, W, H int
	Label      string
}

func (b *Button) Contains(x, y int) bool {
	return x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H
}

func (b *Button) Draw(screen *ebiten.Image) {
	// Background
	vector.FillRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), color.RGBA{100, 100, 200, 255}, true)
	// Border
	vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 2, color.RGBA{50, 50, 150, 255}, true)

	labelX := b.X + b.W/2 - len(b.Label)*3
	labelY := b.Y + b.H/2 - 8
	ebitenutil.DebugPrintAt(screen, b.Label, labelX, labelY)

}
