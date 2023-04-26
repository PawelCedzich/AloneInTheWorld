package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font/opentype"
)

type canvasState struct {
	transformation ebiten.GeoM
	font           *opentype.Font
	color          color.Color
	textSize       float64
}

type Canvas struct {
	dst   *ebiten.Image
	stack []canvasState
}

func NewCanvas(dst *ebiten.Image) *Canvas {
	return &Canvas{dst: dst, stack: []canvasState{{}}}
}

func (c *Canvas) SetTransformation(m ebiten.GeoM) {
	c.stack[len(c.stack)-1].transformation = m
}

func (c *Canvas) Transformation() ebiten.GeoM {
	return c.stack[len(c.stack)-1].transformation
}

func (c *Canvas) Translate(x, y float64) {
	m := c.Transformation()
	m.Translate(x, y)
	c.SetTransformation(m)
}

func (c *Canvas) Size() (float64, float64) {
	x, y := c.dst.Size()
	return float64(x), float64(y)

}

func (c *Canvas) Restore() {
	c.stack = c.stack[:len(c.stack)-1]
}

func (c *Canvas) Concat(other ebiten.GeoM) {
	m := c.Transformation()
	other.Concat(m)
	c.SetTransformation(other)
}

func (c *Canvas) Save() {
	c.stack = append(c.stack, c.lastState())
}

func (c *Canvas) lastState() canvasState {
	return c.stack[len(c.stack)-1]
}

func (c *Canvas) DrawImage(img *ebiten.Image, op *ebiten.DrawImageOptions) {
	copyOp := &ebiten.DrawImageOptions{
		GeoM:          op.GeoM,
		ColorM:        op.ColorM,
		CompositeMode: op.CompositeMode,
		Filter:        op.Filter,
	}
	copyOp.GeoM.Concat(c.Transformation())
	c.dst.DrawImage(img, copyOp)
}

func (c *Canvas) DrawRect(r Rect, color color.RGBA) {
	r = r.Apply(c.Transformation())
	ebitenutil.DrawRect(c.dst, r.Left, r.Top, r.width(), r.height(), color)
}
