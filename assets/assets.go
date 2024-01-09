package assets

import (
	_ "embed"
)

//go:embed level1.txt
var Level1 string

//go:embed TileMiddle.png
var GroundMid []byte

//go:embed TileSingle.png
var GroundSingle []byte

//go:embed TileRight.png
var GroundRight []byte

//go:embed TileLeft.png
var GroundLeft []byte

//go:embed TileFill.png
var GroundFill []byte

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

//go:embed BackgroundImage.png
var BackgroundImage []byte

//go:embed BackgroundTown.png
var BackgroundTown []byte

//go:embed BackgroundTownFront.png
var BackgroundTownFront []byte

//go:embed ButtonStart.png
var ButtonStart []byte

//go:embed ButtonSettings.png
var ButtonSettings []byte

//go:embed ButtonContinue.png
var ButtonContinue []byte

//go:embed ButtonExit.png
var ButtonExit []byte

//go:embed ButtonNoText.png
var ButtonNoText []byte

//go:embed ButtonSlider.png
var ButtonSlider []byte

//go:embed ButtonOn.png
var ButtonOn []byte

//go:embed ButtonOff.png
var ButtonOff []byte

//go:embed charac.png
var Character []byte

//go:embed Tusj.ttf
var Tusj []byte

//go:embed mainscreen_bgm.mp3
var Music []byte
