package ui

import (
	"image/color"
	"os"

	"github.com/bastien-m/mario/ui/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (g *Game) updateSplashScreen() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.newGameBtn.Contains(x, y) {
			g.startNewGame()
		}
		if g.exitBtn.Contains(x, y) {
			os.Exit(0)
		}
	}
}

func (g *Game) drawSplash(screen *ebiten.Image) {
	mX := screenWidth / 2
	mY := screenHeight / 2

	face := &text.GoTextFace{
		Source: components.FontSource,
		Size:   36,
	}

	title := "Mario"
	w, h := text.Measure(title, face, 0)

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(mX)-w/2, float64(mY)-100-h/2)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, title, face, op)

	btnW := 150
	btnH := 50

	g.newGameBtn = &components.Button{X: mX - btnW/2, Y: mY - btnH/2, W: btnW, H: btnH, Label: "New Game"}
	g.exitBtn = &components.Button{X: mX - btnW/2, Y: mY + btnH, W: btnW, H: btnH, Label: "Exit"}

	g.newGameBtn.Draw(screen)
	g.exitBtn.Draw(screen)
}
