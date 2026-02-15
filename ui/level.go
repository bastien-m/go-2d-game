package ui

import (
	"fmt"

	"github.com/bastien-m/mario/assets"
	"github.com/bastien-m/mario/engine"
	"github.com/bastien-m/mario/engine/constants"
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) drawLevelScreen(screen *ebiten.Image) {
	mapFile, err := engine.GetLevel(assets.Level01)
	if err != nil {
		fmt.Printf("Error occured while fetching assets %v", err)
	}

	imgOpts := &ebiten.DrawImageOptions{}

	for _, layer := range mapFile.Layers {
		for _, chunk := range layer.Chunks {
			startXPx := chunk.X * constants.AssetSize
			for i, tileType := range chunk.Data {
				if tileType == engine.Air || tileType == engine.InitialPosition {
					continue
				}
				imgOpts.GeoM.Reset()

				xPx := startXPx + ((i % chunk.Width) * constants.AssetSize)
				yPx := (i / chunk.Width) * constants.AssetSize

				tile := g.tileset[int(tileType)]

				imgOpts.GeoM.Translate(float64(xPx), float64(yPx))

				// ebitenutil.DebugPrintAt(screen, fmt.Sprintf("type %d, %d at x:%d y:%d", tileType, tile.Bounds().Dx(), xPx, yPx), xPx, yPx)
				screen.DrawImage(tile, imgOpts)
			}
		}
	}

}

func (g *Game) updateLevelScreen() {
	g.updatePlayer()
}
