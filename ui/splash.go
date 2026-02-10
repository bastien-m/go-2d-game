package ui

import (
	"os"

	"github.com/bastien-m/mario/ui/components"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	ebitenutil.DebugPrintAt(screen, "Mario", 155, 100)

	g.newGameBtn = &components.Button{X: 125, Y: 250, W: 150, H: 50, Label: "New Game"}
	g.exitBtn = &components.Button{X: 125, Y: 320, W: 150, H: 50, Label: "Exit"}

	g.newGameBtn.Draw(screen)
	g.exitBtn.Draw(screen)

}
