package game

type Ground struct {
	*RectObject
	solid bool
}

func NewGround(tex Drawable, area Rect, solid bool, scale float64) *Ground {

	return &Ground{
		RectObject: NewRectObject(tex, area.Scale(scale)),
		solid:      solid,
	}
}

func (g *Ground) Layout(sw, sh float64) {

}

func (g *Ground) Draw(dst *Canvas) {

	g.RectObject.Draw(dst)
	//g.area.Draw(dst)
}

func (g *Ground) Solid() bool {
	return g.solid
}
