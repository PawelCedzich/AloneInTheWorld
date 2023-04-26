package game

import (
	"fmt"

	"./assets"
)

type Game struct {
	engine  *Engine
	texture *TextureManager
}

func NewGame(e *Engine, tex *TextureManager) *Game {
	g := &Game{
		e,
		tex,
	}

	g.LoadLevel1()

	return g
}

func (g *Game) LoadLevel1() {
	level, err := ParseLevel(assets.Level1)
	if err != nil {
		panic(fmt.Errorf("illegall state, embeded level 1 is invalid %w", err))
	}

	for _, renderable := range level.Build(g) {
		g.engine.AddObject(renderable)
	}

}
