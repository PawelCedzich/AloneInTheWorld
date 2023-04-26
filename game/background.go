package game

import "github.com/hajimehoshi/ebiten/v2"

type Background struct {
	*RectObject
}

func NewBackground(tex Drawable) *Background {
	return &Background{
		RectObject: NewRectObject(tex, Rect{}),
	}
}

func (b *Background) Draw(canvas *Canvas) {
	sw, sh := canvas.Size()
	b.area = Rect{Right: sw, Bottom: sh}
	savedMatrix := canvas.Transformation()
	canvas.SetTransformation(ebiten.GeoM{})
	b.RectObject.Draw(canvas)
	canvas.SetTransformation(savedMatrix)
}

func (s *Background) Layout(_, _ float64) {}
