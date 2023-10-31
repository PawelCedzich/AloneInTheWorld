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
		game.CharacterT:            assets.Character,
	})
	if err != nil {
		panic(err)
	}

	engine := game.NewEngine(cfg)
	game.NewGame(engine, texture)

	if err := engine.Start(); err != nil {
		panic(err)
	}
	return nil
}
