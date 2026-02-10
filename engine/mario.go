package engine

import (
	"embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"

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

func loadEmbeddedImage(resource embed.FS, name string) (*ebiten.Image, error) {
	// 1. Lire le fichier depuis l'embed (chemin racine)
	data, err := resource.Open(name)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	// 2. Décoder l'image (PNG/JPG)
	// C'est ici que l'import _ "image/png" est crucial !
	img, _, err := image.Decode(data)
	if err != nil {
		return nil, err
	}

	// 3. Convertir l'image standard Go en image Ebitengine
	return ebiten.NewImageFromImage(img), nil
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

type JsonMapFile struct {
	Layers []struct {
		Chunks []ChunkJson `json:"chunks"`
	} `json:"layers"`
}

func BuildLevel(level []byte) (*JsonMapFile, error) {
	var mapFile JsonMapFile

	err := json.Unmarshal(level, &mapFile)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing JSON file : %w", err)
	}

	return &mapFile, nil
}
