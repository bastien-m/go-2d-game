package engine

import (
	"embed"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"os"

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

func BuildTilesetManager(resources embed.FS, filename string, tileSize, rows, columns int) (*TilesetManager, error) {

	fullImage, err := loadEmbeddedImage(resources, filename)
	if err != nil {
		return nil, fmt.Errorf("Cant open file %s %w", filename, err)
	}
	manager := &TilesetManager{
		Tiles: make(map[int]*ebiten.Image),
	}

	for i := range rows {
		for j := range columns {
			rect := image.Rect(tileSize*j, tileSize*i, tileSize*(j+1), tileSize*(i+1))

			index := (i * rows) + j + 1
			manager.Tiles[index] = fullImage.SubImage(rect).(*ebiten.Image)
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

func BuildLevel(levelDescriptorPath string) (*JsonMapFile, error) {
	data, err := os.ReadFile(levelDescriptorPath)

	if err != nil {
		return nil, fmt.Errorf("Unable to read file %s %w", levelDescriptorPath, err)
	}

	var mapFile JsonMapFile

	err = json.Unmarshal(data, &mapFile)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing JSON file : %w", err)
	}

	return &mapFile, nil
}
