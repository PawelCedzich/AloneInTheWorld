package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/PawelCedzich/AloneInTheWorld/AloneInTheWorld/assets"
)

type Game struct {
	engine     *Engine
	texture    *TextureManager
	font       *FontManager
	music      *AudioManager
	playerArea Rect
	player     *Player
	dataLoaded bool
}

func NewGame(e *Engine, tex *TextureManager, font *FontManager, music *AudioManager) *Game {
	g := &Game{
		engine:     e,
		texture:    tex,
		font:       font,
		music:      music,
		dataLoaded: false,
	}

	e.mainScreen = func() {
		g.LoadStartMenu()
	}

	e.level = func(lvl int, newGame bool) {
		g.LoadLevel(lvl, newGame)
	}

	e.levelSettings = func() {
		g.LoadSettings()
	}

	e.levelGameSettings = func() {
		g.LoadGameSettings()
	}

	e.levelGameOver = func() {
		g.LoadGameOver()
	}

	e.updateCam = func() {
		g.UpdateCamera()
	}

	//g.music.LoadAudio(Music).Play()

	return g
}

func (g *Game) LoadLevel(lvl int, newGame bool) {
	var level *Level
	var err error
	switch lvl {
	case 1:
		level, err = ParseLevel(assets.Level1)
	case 2:
		level, err = ParseLevel(assets.Level2)
	}

	if err != nil {
		panic(fmt.Errorf("illegall state, embeded level 1 is invalid %w", err))
	}
	for _, renderable := range level.Build(g) {
		g.engine.AddObject(renderable)
	}

	g.player = level.player
	g.UpdateCamera()

	if !newGame {
		deltaX := g.playerArea.Left - level.player.area.Left
		deltaY := g.playerArea.Top - level.player.area.Top
		level.player.area = level.player.area.Offset(deltaX, deltaY)
	}
}

func (g *Game) LoadStartMenu() {

	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundImageT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownFrontT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	//new game
	cell := Rect{250, 250, 350, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonStartT)), cell, g.engine.Scale(), func() { g.engine.NewGameBool() }))

	//continue
	cell = Rect{400, 250, 500, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonContinueT)), cell, g.engine.Scale(), func() { g.load() }))

	//settings
	cell = Rect{550, 250, 650, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonSettingsT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(10) }))

	//exit
	cell = Rect{700, 250, 800, 300}
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

	//back to start
	cell := Rect{250, 250, 350, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonStartT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(0) }))

	//to do
	cell = Rect{250, 350, 350, 400}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonSettingsT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(0) }))

	//exit
	cell = Rect{250, 450, 350, 500}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonExitT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(-1) }))

	//fullscreen
	cell = Rect{250, 550, 350, 600}
	g.engine.AddObject(NewButtonOnOff(NewDrawableTexture(g.texture.LoadTexture(ButtonOnT)), NewDrawableTexture(g.texture.LoadTexture(ButtonOffT)), cell, g.engine.Scale(), func() { g.engine.ChangeFullscreen() }))

	//volume
	cell = Rect{200, 650, 400, 675}
	g.engine.AddObject(NewSlider(NewDrawableTexture(g.texture.LoadTexture(ButtonNoTextT)), NewDrawableTexture(g.texture.LoadTexture(ButtonSliderT)), cell, g.engine.Scale(), g.music.ChangeVolume))
}

func (g *Game) LoadGameSettings() {
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundImageT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownFrontT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	//back to start
	cell := Rect{250, 250, 350, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonStartT)), cell, g.engine.Scale(), func() { g.engine.moveBackToStart() }))

	//Save
	cell = Rect{250, 350, 350, 400}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonSaveT)), cell, g.engine.Scale(), func() { g.save() }))

	//fullscreen
	cell = Rect{250, 550, 350, 600}
	g.engine.AddObject(NewButtonOnOff(NewDrawableTexture(g.texture.LoadTexture(ButtonOnT)), NewDrawableTexture(g.texture.LoadTexture(ButtonOffT)), cell, g.engine.Scale(), func() {}))

	//volume
	cell = Rect{200, 650, 400, 675}
	g.engine.AddObject(NewSlider(NewDrawableTexture(g.texture.LoadTexture(ButtonNoTextT)), NewDrawableTexture(g.texture.LoadTexture(ButtonSliderT)), cell, g.engine.Scale(), g.music.ChangeVolume))

	//exit
	cell = Rect{250, 450, 350, 500}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonExitT)), cell, g.engine.Scale(), func() { g.engine.ChangeStage(-1) }))
}

func (g *Game) LoadGameOver() {
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundImageT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownFrontT))))
	g.engine.AddObject(NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	g.engine.AddObject(
		NewText(
			g.font.LoadFont(TusjF),
			"Game over!",
			24*g.engine.Scale(),
			0.2,
			0.2,
		),
	)

	//back to start
	cell := Rect{250, 250, 350, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonStartT)), cell, g.engine.Scale(), func() { g.engine.moveBackToStart() }))

	//continue
	cell = Rect{400, 250, 500, 300}
	g.engine.AddObject(NewButton(NewDrawableTexture(g.texture.LoadTexture(ButtonContinueT)), cell, g.engine.Scale(), func() { g.load() }))
}

func (g *Game) save() {
	saveToJSON("playerData.json", g.player.area)
	saveToJSON("playerLevel.json", g.engine.playerLevel)
}

func (g *Game) load() {
	var levelData int
	loadFromJSON("playerData.json", &g.playerArea)
	loadFromJSON("playerLevel.json", &levelData)

	g.engine.playerLevel = levelData

	g.engine.SavedGameBool()

	g.dataLoaded = true
}

func saveToJSON(filename string, data interface{}) {
	file, _ := json.MarshalIndent(data, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
}

func loadFromJSON(filename string, result interface{}) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Nie można odczytać pliku: %v", err)
	}
	json.Unmarshal([]byte(file), result)
}

func (g *Game) UpdateCamera() {
	g.engine.camera.UpdateMainCharacter(g.player)
}
