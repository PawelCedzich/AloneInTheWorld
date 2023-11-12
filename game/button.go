package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Button struct {
	*RectObject
	scale  float64
	action func()
}

func NewButton(tex Drawable, area Rect, scale float64, action func()) *Button {
	b := &Button{
		RectObject: NewRectObject(tex, area),
		scale:      scale,
		action:     action,
	}

	return b
}

func (b *Button) Layout(sw, sh float64) {
	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if x > int(b.area.Left) && x < int(b.area.Right) && y > int(b.area.Top) && y < int(b.area.Bottom) {
			b.action()
		}
	}
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type ButtonOnOff struct {
	*Button
	onTex      Drawable
	offTex     Drawable
	state      bool
	clickDelay int
}

func NewButtonOnOff(onTex, offTex Drawable, area Rect, scale float64, action func()) *ButtonOnOff {
	b := &ButtonOnOff{
		Button: &Button{
			RectObject: &RectObject{
				texture: onTex,
				area:    area,
			},
			scale:  scale,
			action: action,
		},
		onTex:  onTex,
		offTex: offTex,
		state:  true,
	}

	return b
}

func (b *ButtonOnOff) Layout(sw, sh float64) {
	b.clickDelay++
	if b.clickDelay > 60 {
		b.clickDelay = 30
	}

	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if x > int(b.area.Left) && x < int(b.area.Right) && y > int(b.area.Top) && y < int(b.area.Bottom) {
			if b.state == true && b.clickDelay >= 30 {
				//do action
				b.Button.texture = b.offTex
				b.state = false
				b.clickDelay = 0
			}
			if b.state == false && b.clickDelay >= 30 {
				//do action
				b.Button.texture = b.onTex
				b.state = true
				b.clickDelay = 0
			}
		}
	}
}

func (b *ButtonOnOff) Draw(dst *Canvas) {
	b.RectObject.Draw(dst)
}
