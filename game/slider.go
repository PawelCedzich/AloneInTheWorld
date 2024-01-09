package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Slider struct {
	*RectObject
	scale         float64
	value         float64
	sliderElement *SliderElement // Dodajemy element suwaka
	action        func(float64)
}

type SliderElement struct {
	*RectObject
	dragging bool
}

func NewSliderElement(tex Drawable, area Rect) *SliderElement {
	s := &SliderElement{
		RectObject: NewRectObject(tex, area),
		dragging:   false,
	}

	return s
}

func NewSlider(tex Drawable, sliderTex Drawable, area Rect, scale float64, action func(float64)) *Slider {

	left := area.Left + (area.width() / 2) - 10
	right := area.Left + (area.width() / 2) + 10
	top := area.Top
	bottom := area.Bottom

	cell := Rect{left, top, right, bottom}

	s := &Slider{
		RectObject:    NewRectObject(tex, area),
		scale:         scale,
		sliderElement: NewSliderElement(sliderTex, cell),
		action:        action,
	}

	return s
}

func (s *Slider) Layout(sw, sh float64) {
	x, y := ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		if x > int(s.sliderElement.area.Left) && x < int(s.sliderElement.area.Right) && y > int(s.sliderElement.area.Top) && y < int(s.sliderElement.area.Bottom) {
			s.sliderElement.dragging = true
		}
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.sliderElement.dragging = false
	}

	if s.sliderElement.dragging {
		mouseX, _ := ebiten.CursorPosition()
		if mouseX >= int(s.area.Left) && mouseX <= int(s.area.Right) {
			delta := float64(mouseX) - s.sliderElement.area.Left - 10
			s.sliderElement.area = s.sliderElement.area.Offset(delta, 0)
		}

		lenght := s.area.Right - s.area.Left
		currentSliderPosition := s.sliderElement.area.Left + 10 - s.area.Left
		s.value = currentSliderPosition / lenght * 100

		s.action(s.value)
	}
}

func (s *Slider) Draw(dst *Canvas) {

	s.RectObject.Draw(dst)
	s.sliderElement.RectObject.Draw(dst)

}
