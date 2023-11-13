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
	flags.BoolVar(&cfg.Fullscreen, "fullscreen", true, "default fullscreen settings")
}

// =====================================================================================================================

type Engine struct {
	renderables   []Renderable
	Cfg           Config
	camera        *Camera
	windowSize    Vec
	stage         int
	blockStage    bool
	mainScreen    func()
	levelSettings func()
	level1        func()
}

func NewEngine(cfg Config) *Engine {
	e := &Engine{Cfg: cfg, camera: NewCamera(), stage: 0, blockStage: false}
	return e
}

func (e *Engine) AddObject(obj Renderable) {
	e.renderables = append(e.renderables, obj)

}

func (e *Engine) Update() error {

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

	if e.blockStage == false {
		e.ClearCamera()

		switch e.stage {
		case -1:
			os.Exit(0)
		case 0:
			if e.mainScreen != nil {
				e.ClearRenderables()
				e.mainScreen()
			}
		case 1:
			if e.level1 != nil {
				e.ClearRenderables()
				e.level1()
			}
		case 10:
			if e.levelSettings != nil {
				e.ClearRenderables()
				e.levelSettings()
			}
		}
		e.blockStage = true
	}

	//dubug tools
	if ebiten.IsKeyPressed(ebiten.Key1) {
		e.ChangeStage(1)
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		e.ChangeStage(0)
	}

	if ebiten.IsKeyPressed(ebiten.Key0) {
		e.ChangeStage(-1)
	}
	return nil
}

func (e *Engine) ClearRenderables() {
	e.renderables = nil
}

func (e *Engine) ChangeStage(value int) {
	e.stage = value
	e.blockStage = false
}

func (e *Engine) ClearCamera() {
	e.camera = NewCamera()
}
