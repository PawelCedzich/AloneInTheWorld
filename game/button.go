package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	*RectObject
	scale  float64
	action func()
}

func NewButton(tex Drawable, area Rect, scale float64) *Button {
	b := &Button{
		RectObject: NewRectObject(tex, area),
		scale:      scale,
	}

	return b
}

func (b *Button) Layout(sw, sh float64) {
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if x > int(b.area.Left) && x < int(b.area.Right) && y > int(b.area.Top) && y < int(b.area.Bottom) {
			log.Println("wykonano akcje")
			//b.action()
		}
	}
}
