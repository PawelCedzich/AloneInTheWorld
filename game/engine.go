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
	Scale      float64
	Fullscreen bool
}

func (cfg *Config) Reset() {
	cfg.Width = 800
	cfg.Height = 600
	cfg.Scale = 0
	cfg.Fullscreen = false
}

func (cfg *Config) Configure(flags *flag.FlagSet) {
	flags.IntVar(&cfg.Width, "width", 800, "default width setting")
	flags.IntVar(&cfg.Height, "height", 600, "default height  settings")
	flags.Float64Var(&cfg.Scale, "scale", 0, "Default scale settings")
	flags.BoolVar(&cfg.Fullscreen, "Fullscreen", false, "default fullscreen settings")
}

// =====================================================================================================================

type Engine struct {
	renderables []Renderable
	camera      *Camera
	Cfg         Config
	windowSize  Vec
	stage       float64
	Close       float64
	textCounter float64
	stageBlock  bool
	mainScreen  func()
	level1      func()
}

func NewEngine(cfg Config) *Engine {
	e := &Engine{Cfg: cfg, camera: NewCamera(), stage: 0, Close: 0, textCounter: 0}
	return e
}

func (e *Engine) AddObject(obj Renderable) {
	e.renderables = append(e.renderables, obj)

}

func (e *Engine) Update() error {

	e.StageManager()

	for _, object := range e.renderables {
		object.Layout(e.windowSize.x, e.windowSize.y)
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {

	canvas := NewCanvas(screen)

	if e.renderables != nil {
		for _, object := range e.renderables {
			object.Draw(canvas)
		}

		if e.stage == 0 {
			if e.textCounter < 40 {
				e.renderables[len(e.renderables)-1].Draw(canvas)
			}
			if e.textCounter > 70 {
				e.textCounter = 0
			}
		}
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {

	scale := e.Scale()
	e.windowSize = Vec{x: float64(outsideWidth) * scale, y: float64(outsideHeight) * scale}
	return int(e.windowSize.x), int(e.windowSize.y)
}

func (e *Engine) ClearRenderables() {
	e.renderables = nil
}

func (e *Engine) ClearCamera() {
	e.camera = NewCamera()
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
	ebiten.SetWindowTitle("Morty game_v2")
	if err := ebiten.RunGame(e); err != nil {
		return fmt.Errorf("cant start game %w", err)
	}

	return nil
}

func (e *Engine) StageManager() {
	e.Close--
	e.textCounter++

	if e.stage == 0 {
		if !e.stageBlock {
			if e.mainScreen != nil {
				e.mainScreen()
				e.stageBlock = true
			}
		}
		if e.Close <= 0 && ebiten.IsKeyPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}
	}
	if e.stage == 0 && ebiten.IsKeyPressed(ebiten.KeySpace) {
		e.ClearRenderables()
		if e.level1 != nil {
			e.level1()
			e.stage = 1
		}
	}
	if e.stage == 1 {
		//e.camera.MainCharacter.onLose = func() {
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			e.ClearRenderables()
			e.ClearCamera()
			e.stageBlock = false
			e.stage = 0
			e.Close = 40
		}
		//}
	}

}
