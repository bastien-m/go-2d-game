package ui

import (
	"embed"
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/bastien-m/mario/engine"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Screen int

const (
	SplashScreen Screen = iota
	LevelScreen
	GameOverScreen
)

const (
	screenWidth  = 600
	screenHeight = 500
	assetSize    = 16
	gravity      = 0.5
)

type Button struct {
	X, Y, W, H int
	Label      string
}

func (b *Button) Contains(x, y int) bool {
	return x >= b.X && x <= b.X+b.W && y >= b.Y && y <= b.Y+b.H
}

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

type Game struct {
	screen Screen

	newGameBtn *Button
	exitBtn    *Button
	backBtn    *Button

	player     *Player
	playerIcon *ebiten.Image

	tileset   map[int]*ebiten.Image
	resources embed.FS
}

func buildGame(resources embed.FS) (*Game, error) {
	if tilsetManager, err := engine.BuildTilesetManager(resources, "assets/sprites-v2.png", assetSize, 3, 3); err != nil {
		return nil, err
	} else {
		return &Game{
			screen:     SplashScreen,
			newGameBtn: &Button{X: 125, Y: 250, W: 150, H: 50, Label: "New Game"},
			exitBtn:    &Button{X: 125, Y: 320, W: 150, H: 50, Label: "Exit"},
			backBtn:    &Button{X: 125, Y: 400, W: 150, H: 50, Label: "Back to Menu"},
			tileset:    tilsetManager.Tiles,
			resources:  resources,
			player: &Player{
				x:  0,
				y:  0,
				vX: 0,
				vY: 0,
			},
			playerIcon: tilsetManager.Tiles[int(engine.Steady)],
		}, nil
	}
}

func (g *Game) startNewGame() {
	g.screen = LevelScreen

	mapFile, err := engine.BuildLevel("assets/level-01.json")
	if err != nil {
		fmt.Printf("Error occured while fetching assets %v", err)
	}

	for _, layer := range mapFile.Layers {
		for _, chunk := range layer.Chunks {
			for i, data := range chunk.Data {
				if data == engine.InitialPosition {
					g.player.x = float64(i % 16)
					g.player.y = float64(i / 16)
				}
			}
		}
	}
}

func (g *Game) Update() error {
	switch g.screen {
	case SplashScreen:
		g.updateSplashScreen()
	case LevelScreen:
		g.updateLevelScreen()
	case GameOverScreen:
		g.updateGameOverScreen()
	}
	return nil
}

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

func (g *Game) updateLevelScreen() {
	ebiten.KeyName(ebiten.KeyQ)
	// here is the game interface handling code
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		// jump
		g.player.jump()
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		// left
		g.playerIcon = g.tileset[int(engine.GoingLeft)]
		g.player.moveBackward()
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		// right
		g.playerIcon = g.tileset[int(engine.GoingRight)]
		g.player.moveForward()
	}
}

func (g *Game) updateGameOverScreen() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if g.backBtn.Contains(x, y) {
			g.screen = SplashScreen
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{245, 131, 96, 255})
	switch g.screen {
	case SplashScreen:
		g.drawSplash(screen)
	case LevelScreen:
		g.drawLevelScreen(screen)
		g.drawPlayer(screen)
	case GameOverScreen:
		g.drawLevelScreen(screen)
		g.drawGameOverScreen(screen)
	}
}

func (g *Game) drawSplash(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, "Mario", 155, 100)

	g.drawButton(screen, g.newGameBtn)
	g.drawButton(screen, g.exitBtn)

}

func (g *Game) drawButton(screen *ebiten.Image, btn *Button) {
	// Background
	vector.FillRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), color.RGBA{100, 100, 200, 255}, true)
	// Border
	vector.StrokeRect(screen, float32(btn.X), float32(btn.Y), float32(btn.W), float32(btn.H), 2, color.RGBA{50, 50, 150, 255}, true)

	labelX := btn.X + btn.W/2 - len(btn.Label)*3
	labelY := btn.Y + btn.H/2 - 8
	ebitenutil.DebugPrintAt(screen, btn.Label, labelX, labelY)

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

func (g *Game) drawLevelScreen(screen *ebiten.Image) {
	mapFile, err := engine.BuildLevel("assets/level-01.json")
	if err != nil {
		fmt.Printf("Error occured while fetching assets %v", err)
	}

	imgOpts := &ebiten.DrawImageOptions{}

	for _, layer := range mapFile.Layers {
		for _, chunk := range layer.Chunks {
			startXPx := chunk.X * assetSize
			for i, tileType := range chunk.Data {
				if tileType == engine.Air || tileType == engine.InitialPosition {
					continue
				}
				imgOpts.GeoM.Reset()

				xPx := startXPx + ((i % chunk.Width) * assetSize)
				yPx := (i / chunk.Width) * assetSize

				tile := g.tileset[int(tileType)]

				imgOpts.GeoM.Translate(float64(xPx), float64(yPx))

				// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("type %d, %d at x:%d y:%d", tileType, tile.Bounds().Dx(), xPx, yPx), xPx, yPx)
				screen.DrawImage(tile, imgOpts)
			}
		}
	}

}

func (g *Game) drawGameOverScreen(screen *ebiten.Image) {
	overlay := ebiten.NewImage(screenWidth, screenHeight)
	overlay.Fill(color.RGBA{0, 0, 0, 150})
	screen.DrawImage(overlay, nil)

	ebitenutil.DebugPrintAt(screen, "Game Over", 170, 300)

	g.drawButton(screen, g.backBtn)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	switch g.screen {
	case SplashScreen:
		return screenWidth, screenHeight
	case LevelScreen:
	case GameOverScreen:
		return int(math.Abs((screenWidth + 100) / 3)), int(math.Abs(screenHeight / 3))
	}
	return int(math.Abs(screenWidth / 1.7)), int(math.Abs(screenHeight / 1.7))
}

func Run(resources embed.FS) {
	// ebiten.SetWindowIcon()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mario")

	game, err := buildGame(resources)
	if err != nil {
		fmt.Printf("Error while initializing game %v\n", err)
		os.Exit(1)
	}
	if err = ebiten.RunGame(game); err != nil {
		fmt.Printf("Error while initializing game %v\n", err)
		os.Exit(1)
	}
}
