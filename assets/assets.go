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

//go:embed TileSingle.png
var GroundSingle []byte

//go:embed TileFillSingle.png
var GroundFillSingle []byte

//go:embed TileFillLeft.png
var GroundFillLeft []byte

//go:embed TileFillRight.png
var GroundFillRight []byte

//go:embed TileFillDeep.png
var GroundFillDeep []byte

//go:embed TileFillDeepSingle.png
var GroundFillDeepSingle []byte

//go:embed TileFillDeepLeft.png
var GroundFillDeepLeft []byte

//go:embed TileFillDeepRight.png
var GroundFillDeepRight []byte

//go:embed TileFillDeeper.png
var GroundFillDeeper []byte

//go:embed TileFillDeeper2.png
var GroundFillDeeper2 []byte

//go:embed TileFillDeeper3.png
var GroundFillDeeper3 []byte

//go:embed TileMiddle.png
var GroundMid []byte

//go:embed TileRight.png
var GroundRight []byte

//go:embed BackgroundImage.png
var BackgroundImage []byte

//go:embed BackgroundTown.png
var BackgroundTown []byte

//go:embed BackgroundTownFront.png
var BackgroundTownFront []byte

//go:embed charac.png
var Character []byte
