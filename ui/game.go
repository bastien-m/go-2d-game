package ui

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/bastien-m/mario/assets"
	"github.com/bastien-m/mario/engine"
	"github.com/bastien-m/mario/ui/components"
	"github.com/hajimehoshi/ebiten/v2"
)

type Screen int

type Game struct {
	screen Screen

	player     *Player
	playerIcon *ebiten.Image

	tileset   map[int]*ebiten.Image
	resources embed.FS

	newGameBtn *components.Button
	exitBtn    *components.Button
}

func buildGame(resources embed.FS) (*Game, error) {
	img, _, err := image.Decode(bytes.NewReader(assets.Sprites))
	if err != nil {
		log.Fatal(err)
	}
	spriteEbitenImg := ebiten.NewImageFromImage(img)
	if tilsetManager, err := engine.BuildTilesetManager(spriteEbitenImg, assetSize, 3, 3); err != nil {
		return nil, err
	} else {
		return &Game{
			screen:    SplashScreen,
			tileset:   tilsetManager.Tiles,
			resources: resources,
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

	mapFile, err := engine.BuildLevel(assets.Level01)
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
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{245, 131, 96, 255})
	switch g.screen {
	case SplashScreen:
		g.drawSplash(screen)
	case LevelScreen:
		g.drawLevelScreen(screen)
		g.drawPlayer(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	switch g.screen {
	case SplashScreen:
		return screenWidth, screenHeight
	case LevelScreen:
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
		log.Fatal(err)
	}
}
