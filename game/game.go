package game

import (
	"fmt"

	"github.com/PawelCedzich/AloneInTheWorld/AloneInTheWorld/assets"
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

	g.LoadStartMenu()
	//g.LoadLevel1()

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

func (g *Game) LoadStartMenu() {
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(StartButtonT)), Rect{250, 250, 350, 300}, g.engine.Scale()))
}
