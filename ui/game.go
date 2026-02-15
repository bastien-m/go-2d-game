package ui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bastien-m/mario/assets"
	"github.com/bastien-m/mario/engine"
	"github.com/bastien-m/mario/engine/constants"
	"github.com/bastien-m/mario/ui/components"
	"github.com/hajimehoshi/ebiten/v2"
)

type Screen int

const (
	SplashScreen Screen = iota
	LevelScreen
)

const (
	screenWidth  = 600
	screenHeight = 500
)

type Game struct {
	screen Screen

	player     *Player
	playerIcon *ebiten.Image

	tileset map[int]*ebiten.Image

	newGameBtn *components.Button
	exitBtn    *components.Button
	tileSize   float64

	level *engine.MapLevel
}

func buildGame() (*Game, error) {
	img, _, err := image.Decode(bytes.NewReader(assets.Sprites))
	if err != nil {
		log.Fatal(err)
	}
	spriteEbitenImg := ebiten.NewImageFromImage(img)
	if tilsetManager, err := engine.BuildTilesetManager(spriteEbitenImg, constants.AssetSize, 3, 3); err != nil {
		return nil, err
	} else {
		return &Game{
			screen:  SplashScreen,
			tileset: tilsetManager.Tiles,
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

	level, err := engine.GetLevel(assets.Level01)
	g.level = level
	if err != nil {
		fmt.Printf("Error occured while fetching assets %v", err)
	}

	for _, layer := range level.Layers {
		for _, chunk := range layer.Chunks {
			for i, data := range chunk.Data {
				if data == engine.InitialPosition {
					g.player.x = float64(i % chunk.Width)
					g.player.y = float64(i / chunk.Height)
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
	screen.Fill(HexToRGBA("#68D0E3"))
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
		return screenWidth / 2, screenHeight / 2
	}
	return screenWidth, screenHeight
}

func Run() {
	// ebiten.SetWindowIcon()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mario")

	game, err := buildGame()
	if err != nil {
		fmt.Printf("Error while initializing game %v\n", err)
		os.Exit(1)
	}
	if err = ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func HexToRGBA(hex string) color.RGBA {
	// 1. Nettoyage : on enlève le # et les éventuels espaces
	hex = strings.TrimPrefix(hex, "#")

	// 2. Conversion de la string (base 16) en un nombre entier
	// On utilise 64 bits pour être tranquille
	value, err := strconv.ParseUint(hex, 16, 64)
	if err != nil {
		// En cas d'erreur (string invalide), on retourne du noir ou du blanc par défaut
		return color.RGBA{0, 0, 0, 255}
	}

	// 3. Application du Bitwise
	// Si la string fait 6 caractères (RRGGBB)
	r := uint8(value >> 16)         // On décale de 2 octets (16 bits)
	g := uint8((value >> 8) & 0xFF) // On décale de 1 octet et on masque
	b := uint8(value & 0xFF)        // On masque les 8 derniers bits

	return color.RGBA{r, g, b, 255}
}
