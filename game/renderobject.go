package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type degree float64

func (d degree) Rad() float64 {
	return math.Pi / 180 * float64(d)
}

// =====================================================================================================================

type Vec struct {
	x, y float64
}

// =====================================================================================================================

type RenderObject struct {
	Position Vec
	pivot    Vec
	Scale    Vec
	Texture  Drawable
	rotation degree
}

func NewObject(texture Drawable) *RenderObject {
	x, y := texture.Size()
	r := &RenderObject{
		pivot: Vec{
			x: -float64(x) / 2,
			y: -float64(y) / 2,
		},
		Texture: texture,
		Scale: Vec{
			x: 1,
			y: 1,
		},
		rotation: 0,
	}
	return r
}

func (r *RenderObject) Draw(dst *Canvas) {

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(r.pivot.x, r.pivot.y)
	op.GeoM.Scale(r.Scale.x, r.Scale.y)
	op.GeoM.Rotate(r.rotation.Rad())
	op.GeoM.Translate(r.Position.x, r.Position.y)
	r.Draw(dst)
}

func (r *RenderObject) Layout(w, h float64) {
}

// =====================================================================================================================

type RectObject struct {
	texture     Drawable
	area        Rect
	mirrorYaxis bool
}

func NewRectObject(tex Drawable, area Rect) *RectObject {
	return &RectObject{
		texture: tex,
		area:    area,
	}
}

func (r *RectObject) Draw(dst *Canvas) {
	texW, texH := r.texture.Size()
	position := Vec{
		x: r.area.Left,
		y: r.area.Top,
	}
	scale := Vec{
		x: r.area.width() / texW,
		y: r.area.height() / texH,
	}

	op := &ebiten.DrawImageOptions{}

	if r.mirrorYaxis {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(texW*2/3, 0)
	}
	op.GeoM.Scale(scale.x, scale.y)
	op.GeoM.Translate(position.x, position.y)

	dst.Save()
	dst.Concat(op.GeoM)
	r.texture.Draw(dst)

	dst.Restore()
}

func (r *RectObject) BoundingBox() Rect {
	return r.area
}

func (r *RectObject) Layout(sw, sh float64) {
	r.texture.Update()
}
