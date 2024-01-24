package game

import (
	"flag"
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type Renderable interface {
	Layout(w, h float64)
	Draw(dst *Canvas)
}

// =====================================================================================================================

type Drawable interface {
	Update()
	Draw(dst *Canvas)
	Size() (x, y float64)
	Pivot() Vec
}

// =====================================================================================================================

type Config struct {
	Width      int
	Height     int
	Stage      int
	Scale      float64
	Fullscreen bool
}

func (cfg *Config) Reset() {
	cfg.Width = 800
	cfg.Height = 600
	cfg.Scale = 0
	cfg.Stage = 1
	cfg.Fullscreen = true
}

func (cfg *Config) Configure(flags *flag.FlagSet) {
	flags.IntVar(&cfg.Width, "width", 800, "default width setting")
	flags.IntVar(&cfg.Height, "height", 600, "default height settings")
	flags.IntVar(&cfg.Stage, "stage", 0, "default stage settings")
	flags.Float64Var(&cfg.Scale, "scale", 0, "Default scale settings")
	flags.BoolVar(&cfg.Fullscreen, "fullscreen", false, "default fullscreen settings")
}

// =====================================================================================================================

type Engine struct {
	renderables       []Renderable
	renderablesStack  [][]Renderable
	stageStack        []int
	Cfg               Config
	camera            *Camera
	windowSize        Vec
	changingStage     bool
	newGame           bool
	blockCountDown    int
	playerLevel       int
	stage             int
	mainScreen        func()
	levelSettings     func()
	levelGameSettings func()
	levelGameOver     func()
	level             func(int, bool)
	updateCam         func()
}

func NewEngine(cfg Config) *Engine {
	e := &Engine{Cfg: cfg, camera: NewCamera(), stage: 0, changingStage: true, blockCountDown: 10, playerLevel: 1, newGame: true}
	return e
}

func (e *Engine) Update() error {
	e.blockCountDown++
	if err := e.StageManager(); err != nil {
		return fmt.Errorf("cant load Stage manager %w", err)
	}

	for _, object := range e.renderables {
		object.Layout(e.windowSize.x, e.windowSize.y)
	}

	return nil

}

func (e *Engine) Draw(screen *ebiten.Image) {

	canvas := NewCanvas(screen)
	sw, sh := canvas.Size()
	canvas.SetTransformation(e.camera.Transformation(sw, sh))

	if e.renderables != nil {
		for _, object := range e.renderables {
			object.Draw(canvas)
		}
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	scale := e.Scale()
	e.windowSize = Vec{x: float64(outsideWidth) * scale, y: float64(outsideHeight) * scale}
	return int(e.windowSize.x), int(e.windowSize.y)
}

func (e *Engine) AddObject(obj Renderable) {
	e.renderables = append(e.renderables, obj)

}

func (e *Engine) Scale() float64 {
	scale := ebiten.DeviceScaleFactor()
	if e.Cfg.Scale != 0 {
		scale = e.Cfg.Scale
	}

	return scale
}

func (e *Engine) Start() error {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(e.Cfg.Width, e.Cfg.Height)
	ebiten.SetFullscreen(e.Cfg.Fullscreen)
	ebiten.SetWindowTitle("Alone In The World")
	if err := ebiten.RunGame(e); err != nil {
		return fmt.Errorf("cant start game %w", err)
	}

	return nil
}

func (e *Engine) StageManager() error {

	//dubugging tools
	if ebiten.IsKeyPressed(ebiten.Key1) {
		e.ChangeStage(4)
	}

	if e.changingStage == true {
		e.ClearCamera()
		switch e.stage {
		case 0:
			if e.mainScreen != nil {
				if len(e.renderablesStack) == 0 {
					e.mainScreen()
					e.PushStage(0, e.renderables)
				}
			}
		case 1:
			if e.level != nil {
				if len(e.stageStack) <= 1 {
					e.level(e.playerLevel, e.newGame)
					e.PushStage(e.stage, e.renderables)
				}
			}
			if e.updateCam != nil {
				e.updateCam()
			}
		case 2:
			if e.levelGameSettings != nil {
				e.levelGameSettings()
				e.PushStage(e.stage, e.renderables)
			}
		case 4:
			if e.levelGameOver != nil {
				e.levelGameOver()
				e.PushStage(e.stage, e.renderables)
			}
		case 10:
			if e.levelSettings != nil {
				e.levelSettings()
				e.PushStage(10, e.renderables)
			}
		}
		e.changingStage = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if e.stage == 10 {
			e.ChangeStage(0)
		} else if e.stage == 1 {
			e.ChangeStage(2)
		} else if e.stage == 2 {
			e.ChangeStage(1)
		} else {
			e.ChangeStage(10)
		}
	}

	if ebiten.IsKeyPressed(ebiten.Key0) {
		e.ChangeStage(-1)
	}
	return nil
}

func (e *Engine) PushStage(stage int, renderables []Renderable) {
	e.stageStack = append(e.stageStack, stage)
	e.renderablesStack = append(e.renderablesStack, renderables)
}

func (e *Engine) PopStage() {
	length := len(e.stageStack)
	if length == 0 {
		return
	}

	e.stageStack = e.stageStack[:length-1]
	if length == 2 {
		e.stage = e.stageStack[0]
	} else {
		e.stage = e.stageStack[length-2]
	}

	e.ClearRenderables()
	e.renderables = e.renderablesStack[length-1]
	e.renderablesStack = e.renderablesStack[:length-1]
	if length == 2 {
		e.renderables = e.renderablesStack[0]
	} else {
		e.renderables = e.renderablesStack[length-2]
	}
}

func (e *Engine) ClearRenderables() {
	e.renderables = nil
}

func (e *Engine) ChangeStage(value int) {
	popping := false
	if e.blockCountDown > 15 {
		if value == -1 {
			os.Exit(0)
		}
		if value < e.stage {
			e.PopStage()
			popping = true
		} else if value == e.stage {
			e.PopStage()
		}
		if popping != true {
			e.stage = value
		}
		e.changingStage = true
		e.blockCountDown = 0

	} else {
		return
	}
}

func (e *Engine) moveBackToStart() {
	e.PopStage()
	e.ChangeStage(0)
	e.playerLevel = 1
}

func (e *Engine) ClearCamera() {
	e.camera = NewCamera()
}

func (e *Engine) ChangePlayerLvL(val int) {

	e.ChangeStage(1)
	e.playerLevel = val
}

func (e *Engine) NewGameBool() {
	e.newGame = true
	e.ChangeStage(1)
}

func (e *Engine) SavedGameBool() {
	e.newGame = false
	e.ChangeStage(1)
}

func (e *Engine) ChangeFullscreen() {
	if ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
	} else {
		ebiten.SetFullscreen(true)
	}
}
