package assets

import (
	_ "embed"
)

var (
	//go:embed sprites-v2.png
	Sprites []byte

	//go:embed level-01.json
	Level01 []byte
)
