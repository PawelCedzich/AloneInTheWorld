package game

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type canvasState struct {
	transformation ebiten.GeoM
	font           *opentype.Font
	color          color.Color
	textSize       float64
}

type fontCacheKey struct {
	face *opentype.Font
	size float64
}

type Canvas struct {
	dst           *ebiten.Image
	stack         []canvasState
	fontFaceCache map[fontCacheKey]font.Face
}

func NewCanvas(dst *ebiten.Image) *Canvas {
	return &Canvas{dst: dst, stack: []canvasState{{}}, fontFaceCache: map[fontCacheKey]font.Face{}}
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

func (c *Canvas) Fill(color color.RGBA) {
	c.dst.Fill(color)
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

func (c *Canvas) SetColor(color color.Color) {
	c.stack[len(c.stack)-1].color = color
}

func (c *Canvas) SetTextSize(size float64) {
	c.stack[len(c.stack)-1].textSize = size
}

func (c *Canvas) SetFont(fnt *opentype.Font) {
	c.stack[len(c.stack)-1].font = fnt
}

func (c *Canvas) MeasureText(str string) Rect {
	r := text.BoundString(c.fontFace(), str)
	return Rect{
		Left:   float64(r.Min.X),
		Top:    float64(r.Min.Y),
		Right:  float64(r.Max.X),
		Bottom: float64(r.Max.Y),
	}
}

func (c *Canvas) fontFace() font.Face {
	state := c.lastState()
	key := fontCacheKey{
		face: state.font,
		size: state.textSize,
	}

	face, ok := c.fontFaceCache[key]
	if ok {
		return face
	}

	const dpi = 76
	ff, err := opentype.NewFace(key.face, &opentype.FaceOptions{
		Size:    key.size,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	if err != nil {
		panic(fmt.Errorf("invalid state %w", err))
	}

	c.fontFaceCache[key] = ff

	return ff
}

func (c *Canvas) DrawTexT(str string, x float64, y float64) {
	state := c.lastState()
	m := c.Transformation()
	x, y = m.Apply(x, y)
	text.Draw(c.dst, str, c.fontFace(), int(x), int(y), state.color)
}
