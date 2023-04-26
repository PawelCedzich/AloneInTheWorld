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
		game.GroundLeftT:      assets.GroundLeft,
		game.GroundFillT:      assets.GroundFill,
		game.GroundMidT:       assets.GroundMid,
		game.GroundRightT:     assets.GroundRight,
		game.BackgroundImageT: assets.BackgroundImage,
		game.BackgroundTownT:  assets.BackgroundTown,
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
