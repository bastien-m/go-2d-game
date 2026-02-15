package engine

import (
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Tile int

const (
	Air Tile = iota
	Ground
	Crate
	Tree
	Cactus
	GoingRight
	GoingLeft
	Steady
	InitialPosition = 99
)

type TilesetManager struct {
	Tiles map[int]*ebiten.Image
}

func BuildTilesetManager(img *ebiten.Image, tileSize, rows, columns int) (*TilesetManager, error) {

	manager := &TilesetManager{
		Tiles: make(map[int]*ebiten.Image),
	}

	for i := range rows {
		for j := range columns {
			rect := image.Rect(tileSize*j, tileSize*i, tileSize*(j+1), tileSize*(i+1))

			index := (i * rows) + j + 1
			manager.Tiles[index] = img.SubImage(rect).(*ebiten.Image)
		}
	}

	return manager, nil
}

type ChunkJson struct {
	Data   []Tile `json:"data"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
}

type MapLevel struct {
	Layers []struct {
		Chunks []ChunkJson `json:"chunks"`
	} `json:"layers"`
}

func GetLevel(level []byte) (*MapLevel, error) {
	var mapFile MapLevel

	err := json.Unmarshal(level, &mapFile)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing JSON file : %w", err)
	}

	return &mapFile, nil
}

type Direction int

const (
	UP Direction = iota
	RIGHT
	LEFT
	DOWN
)

func (m *MapLevel) TileAt(x, y float64, direction Direction) Tile {
	fmt.Printf("[DEBUG] tileX: %f tileY: %f\n", x, y)

	for i := range m.Layers {
		for _, chunk := range m.Layers[i].Chunks {
			fcx := float64(chunk.X)
			fcy := float64(chunk.Y)
			fwidth := float64(chunk.Width)
			fheight := float64(chunk.Height)
			if x >= fcx && x < fcx+fwidth && y >= fcy && y < fcy+fheight {
				return chunk.Data[int(math.Floor(x))+int(math.Floor(y))*int(math.Floor(fwidth))]
			}
		}
	}

	return -1
}
