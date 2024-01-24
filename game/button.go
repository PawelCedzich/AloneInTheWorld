package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Relative float64

type RelativeObject struct {
	*RenderObject
	x, y          Relative
	offset, scale float64
}

type Button struct {
	*RenderObject
	x, y   Relative
	scale  float64
	action func()
}

func NewButton(tex *ebiten.Image, x, y Relative, scale float64, action func()) *Button {
	b := &Button{
		NewObject(tex),
		x,
		y,
		scale,
		action,
	}

	return b
}

func (b *Button) Layout(sw, sh float64) {
	imgW, imgH := b.Texture.Size()
	scale := sh / float64(imgH) * b.scale
	b.Scale = Vec{
		x: scale,
		y: scale,
	}
	b.Position.x = sw * float64(b.x)
	b.Position.y = sh * float64(b.y)

	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if x > int(b.Position.x) && x < int(b.Position.x+float64(imgW)) && y > int(b.Position.y) && y < int(b.Position.y+float64(imgH)) {
			b.action()
		}
	}
}

func (b *Button) Draw(dst *Canvas) {

}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type ButtonOnOff struct {
	*Button
	x, y       Relative
	onTex      *ebiten.Image
	offTex     *ebiten.Image
	state      bool
	clickDelay int
}

func NewButtonOnOff(onTex, offTex *ebiten.Image, x, y Relative, scale float64, action func()) *ButtonOnOff {
	b := &ButtonOnOff{
		Button: &Button{
			NewObject(onTex),
			x,
			y,
			scale,
			action,
		},
		onTex:  onTex,
		offTex: offTex,
		state:  true,
	}

	return b
}

func (b *ButtonOnOff) Layout(sw, sh float64) {
	imgW, imgH := b.Texture.Size()
	b.clickDelay++
	if b.clickDelay > 60 {
		b.clickDelay = 30
	}

	x, y := ebiten.CursorPosition()

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
			if x > int(b.Position.x) && x < int(b.Position.x+float64(imgW)) && y > int(b.Position.y) && y < int(b.Position.y+float64(imgH)) {
				if b.state == true && b.clickDelay >= 30 {
					b.Button.RenderObject.Texture = b.offTex
					b.state = false
					b.clickDelay = 0
				}
				if b.state == false && b.clickDelay >= 30 {
					//do action
					b.Button.RenderObject.Texture = b.onTex
					b.state = true
					b.clickDelay = 0
				}
				b.action()
			}
		}
	}
}

func (b *ButtonOnOff) Draw(dst *Canvas) {

}
