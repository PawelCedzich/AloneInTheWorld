package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Region int

type Rect struct {
	Left, Top, Right, Bottom float64
}

func (r Rect) width() float64 {

	return r.Right - r.Left
}

func (r Rect) height() float64 {
	return r.Bottom - r.Top
}

func (r Rect) Draw(dst *Canvas) {
	dst.DrawRect(r, color.RGBA{R: 255, G: 0, B: 0, A: 100})
}

func (r Rect) Scale(v float64) Rect {
	return Rect{
		Left:   r.Left * v,
		Top:    r.Top * v,
		Right:  r.Right * v,
		Bottom: r.Bottom * v,
	}
}

func (r Rect) Offset(x, y float64) Rect {
	return Rect{
		Left:   r.Left + x,
		Top:    r.Top + y,
		Right:  r.Right + x,
		Bottom: r.Bottom + y,
	}
}

func (r Rect) Apply(m ebiten.GeoM) Rect {
	var t Rect
	t.Left, t.Top = m.Apply(r.Left, r.Top)
	t.Right, t.Bottom = m.Apply(r.Right, r.Bottom)
	return t
}

func (r Rect) Overlaps(box Rect) bool {

	return r.Left < box.Right && r.Right > box.Left && r.Top < box.Bottom && r.Bottom >= box.Top
}
