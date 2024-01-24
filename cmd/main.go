package main

import (
	"flag"
	_ "image/png"

	"github.com/PawelCedzich/AloneInTheWorld/AloneInTheWorld/assets"
	"github.com/PawelCedzich/AloneInTheWorld/AloneInTheWorld/game"
)

func main() {
	if err := realMain(); err != nil {
		panic(err)
	}
}

func realMain() error {
	var cfg game.Config
	cfg.Reset()
	cfg.Configure(flag.CommandLine)
	flag.Parse()

	texture, err := game.NewTextureManager(map[game.Texture][]byte{
		game.GroundMidT:            assets.GroundMid,
		game.GroundSingleT:         assets.GroundSingle,
		game.GroundRightT:          assets.GroundRight,
		game.GroundLeftT:           assets.GroundLeft,
		game.GroundFillT:           assets.GroundFill,
		game.GroundFillSingleT:     assets.GroundFillSingle,
		game.GroundFillLeftT:       assets.GroundFillLeft,
		game.GroundFillRightT:      assets.GroundFillRight,
		game.GroundFillDeepT:       assets.GroundFillDeep,
		game.GroundFillDeepSingleT: assets.GroundFillDeepSingle,
		game.GroundFillDeepLeftT:   assets.GroundFillDeepLeft,
		game.GroundFillDeepRightT:  assets.GroundFillDeepRight,
		game.GroundFillDeeperT:     assets.GroundFillDeeper,
		game.GroundFillDeeper2T:    assets.GroundFillDeeper2,
		game.GroundFillDeeper3T:    assets.GroundFillDeeper3,
		game.BackgroundImageT:      assets.BackgroundImage,
		game.BackgroundTownT:       assets.BackgroundTown,
		game.BackgroundTownFrontT:  assets.BackgroundTownFront,
		game.ButtonStartT:          assets.ButtonStart,
		game.ButtonSaveT:           assets.ButtonSave,
		game.ButtonSettingsT:       assets.ButtonSettings,
		game.ButtonExitT:           assets.ButtonExit,
		game.ButtonContinueT:       assets.ButtonContinue,
		game.ButtonNoTextT:         assets.ButtonNoText,
		game.ButtonSliderT:         assets.ButtonSlider,
		game.ButtonOnT:             assets.ButtonOn,
		game.ButtonOffT:            assets.ButtonOff,
		game.CharacterT:            assets.Character,
		game.MortyVanillaT:         assets.MortyVanilla,
		game.MortyJumpingT:         assets.MortyJumping,
		game.MortyMeditatingT:      assets.MortyMeditating,
		game.MortyJoyFulT:          assets.MortyJoyFul,
		game.MortySadFulT:          assets.MortySadFul,
		game.MortyWalkingT:         assets.MortyWalking,
	})
	if err != nil {
		panic(err)
	}

	font, err := game.NewFontManager(map[game.Font][]byte{
		game.TusjF: assets.Tusj,
	})

	music, err := game.NewAudioManager(map[game.Audio][]byte{
		game.Music: assets.Music,
	})
	if err != nil {
		panic(err)
	}

	engine := game.NewEngine(cfg)
	game.NewGame(engine, texture, font, music)

	if err := engine.Start(); err != nil {
		panic(err)
	}
	return nil
}
