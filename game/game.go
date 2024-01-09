package game

import (
	"fmt"

	"github.com/PawelCedzich/AloneInTheWorld/AloneInTheWorld/assets"
)

type Game struct {
	engine  *Engine
	texture *TextureManager
	font    *FontManager
	music   *AudioManager
}

func NewGame(e *Engine, tex *TextureManager, font *FontManager, music *AudioManager) *Game {
	g := &Game{
		e,
		tex,
		font,
		music,
	}

	e.mainScreen = func() {
		g.LoadStartMenu()
	}

	e.level1 = func() {
		g.LoadLevel1()
	}

	e.levelSettings = func() {
		g.LoadSettings()
	}

	g.music.LoadAudio(Music).Play()

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

	g.engine.camera.UpdateMainCharacter(level.player)

}

func (g *Game) LoadStartMenu() {

	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundImageT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownFrontT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	cell := Rect{250, 250, 350, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonStartT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(1) }))

	cell = Rect{400, 250, 500, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonContinueT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(3) }))

	cell = Rect{550, 250, 650, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonSettingsT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(10) }))

	g.engine.AddObject(
		NewButton(
			NewDrawableTexture(g.texture.LoadTexture(ButtonExitT)),
			cell,
			g.engine.Scale(),
			func() { g.engine.ChangeStage(-1) },
		),
	)

	g.engine.AddObject(
		NewText(
			g.font.LoadFont(TusjF),
			"Hello",
			24*g.engine.Scale(),
			0.3,
			0.3,
		),
	)
}

func (g *Game) LoadSettings() {
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundImageT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownFrontT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	cell := Rect{250, 250, 350, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonStartT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(1) }))

	cell = Rect{250, 350, 350, 400}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonSettingsT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(10) }))

	cell = Rect{250, 450, 350, 500}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonExitT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(-1) }))

	cell = Rect{250, 550, 350, 600}
	g.engine.AddObject(NewButtonOnOff(NewDrawableTexture(g.texture.LoadTexture(ButtonOnT)), NewDrawableTexture(g.texture.LoadTexture(ButtonOffT)), cell, g.engine.Scale(), func() {}))

	cell = Rect{200, 650, 400, 675}
	g.engine.AddObject(NewSlider(NewDrawableTexture(g.texture.LoadTexture(ButtonNoTextT)), NewDrawableTexture(g.texture.LoadTexture(ButtonSliderT)), cell, g.engine.Scale(), g.music.ChangeVolume))
}
