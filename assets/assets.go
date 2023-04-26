package assets

import (
	_ "embed"
)

//go:embed level1.txt
var Level1 string

//go:embed TileLeft.png
var GroundLeft []byte

//go:embed TileFill.png
var GroundFill []byte

//go:embed TileMiddle.png
var GroundMid []byte

//go:embed TileRight.png
var GroundRight []byte

//go:embed BackgroundImage.png
var BackgroundImage []byte

//go:embed BackgroundTown.png
var BackgroundTown []byte

//go:embed BackgroundTownFront.png
var BackgroundTownFrotn []byte
