package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Region int

const (
	None Region = iota
	Left
	Top
	Right
	Bottom
)

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

func (r Rect) Overlaps(box Rect) bool {

	return r.Left < box.Right && r.Right > box.Left && r.Top < box.Bottom && r.Bottom >= box.Top
}

func (r Rect) Scale(v float64) Rect {
	return Rect{
		Left:   r.Left * v,
		Top:    r.Top * v,
		Right:  r.Right * v,
		Bottom: r.Bottom * v,
	}
}

func (r Rect) Apply(m ebiten.GeoM) Rect {
	var t Rect
	t.Left, t.Top = m.Apply(r.Left, r.Top)
	t.Right, t.Bottom = m.Apply(r.Right, r.Bottom)
	return t
}

func (r Rect) Offset(x, y float64) Rect {
	return Rect{
		Left:   r.Left + x,
		Top:    r.Top + y,
		Right:  r.Right + x,
		Bottom: r.Bottom + y,
	}
}

func (r Rect) Padding(x, y float64) Rect {

	return Rect{
		Left:   r.Left + x,
		Top:    r.Top + y,
		Right:  r.Right - x,
		Bottom: r.Bottom - y,
	}
}

func (r Rect) ChangeX(width float64) Rect {
	n := Rect{
		Left:   r.Right - width/2,
		Top:    r.Top,
		Right:  r.Left - width/2,
		Bottom: r.Bottom,
	}
	return n
}

func (r Rect) OverlapsRegion(other Rect, delta float64) Region {
	if r.LeftRegion(delta).Overlaps(other) {
		return Left
	}

	if r.TopRegion(delta).Overlaps(other) {
		return Top
	}

	if r.RightRegion(delta).Overlaps(other) {
		return Right
	}

	if r.BottomRegion(delta).Overlaps(other) {
		return Bottom
	}

	return None
}

func (r Rect) LeftRegion(delta float64) Rect {
	return Rect{
		Left:   r.Left - delta,
		Top:    r.Top + delta,
		Right:  r.Left + delta,
		Bottom: r.Bottom - delta,
	}
}

func (r Rect) TopRegion(delta float64) Rect {
	return Rect{
		Left:   r.Left,
		Top:    r.Top - delta,
		Right:  r.Right,
		Bottom: r.Top + delta,
	}
}

func (r Rect) RightRegion(delta float64) Rect {
	return Rect{
		Left:   r.Right - delta,
		Top:    r.Top + delta,
		Right:  r.Right + delta,
		Bottom: r.Bottom - delta,
	}
}

func (r Rect) BottomRegion(delta float64) Rect {
	return Rect{
		Left:   r.Left,
		Top:    r.Bottom - delta,
		Right:  r.Right,
		Bottom: r.Bottom + delta,
	}
}

func (r Rect) withPadding(x float64, y float64) Rect {
	t := Rect{
		Left:   r.Left + x,
		Top:    r.Top - y,
		Right:  r.Right - x,
		Bottom: r.Bottom + y,
	}
	return t
}
