package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	*RectObject
	scale        float64
	grounds      []*Ground
	velocity     Vec
	jumpDuration int
	jumping      bool
	snapped      bool
}

func NewPlayer(tex Drawable, area Rect, scale float64) *Player {
	area.Top = area.Top - 50
	area.Left = area.Left - 50
	p := &Player{
		RectObject: NewRectObject(tex, area),
		scale:      scale,
	}

	return p
}

func (p *Player) Layout(sw, sh float64) {

	scale := p.scale

	//inputs
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.velocity.x = 8 * scale
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.velocity.x = -8 * scale
	} else {
		p.velocity.x = 0
	}

	if p.jumping {
		p.jumpDuration++
		if ebiten.IsKeyPressed(ebiten.KeyUp) {
			p.velocity.y -= 2 * scale
		}
		if p.jumpDuration > 5 {
			p.jumping = false
		}
	} else {
		if p.Grounded() {
			p.velocity.y = 0
			if ebiten.IsKeyPressed(ebiten.KeyUp) {
				p.jumping = true
				p.jumpDuration = 0
			}
		} else {
			if p.velocity.y < 20*scale {
				p.velocity.y += 1 * scale
			}
		}
	}

	p.area = p.area.Offset(p.velocity.x, p.velocity.y)

	p.snapped = false
	p.Snap()

}

func (p *Player) Snap() {

	scale := p.scale

	playerBox := p.BoundingBox()

	for _, ground := range p.grounds {
		if !p.snapped {
			groundBox := ground.BoundingBox()
			if playerBox.Overlaps(groundBox) {
				if math.Abs(groundBox.Top-playerBox.Bottom) < 20*scale {
					delta := playerBox.Bottom - groundBox.Top
					p.area = p.area.Offset(0, -delta)
					p.snapped = true
				} else if ground.Solid() {
					bottomHalf := groundBox
					bottomHalf.Left = groundBox.Left + 10*scale
					bottomHalf.Right = groundBox.Right - 10*scale
					bottomHalf.Top = groundBox.Bottom + groundBox.height()/2
					if playerBox.Overlaps(bottomHalf) {
						delta := playerBox.Top - groundBox.Bottom
						p.area = p.area.Offset(0, -delta)
						p.velocity.y += 3 * scale
					} else {
						leftHalf := groundBox
						leftHalf.Right = groundBox.Left + groundBox.width()/2
						if playerBox.Overlaps(leftHalf) {
							delta := groundBox.Left - playerBox.Right
							p.area = p.area.Offset(delta, 0)
							p.snapped = true
						} else {
							delta := groundBox.Right - playerBox.Left
							p.area = p.area.Offset(delta, 0)
							p.snapped = true
						}
					}
				}
			}
		}
	}
}

func (p *Player) Draw(dst *Canvas) {

	p.RectObject.Draw(dst)

	//hitbox debug
	//p.BoundingBox().Draw(dst)

	//debug sprite
	//p.RectObject.BoundingBox().Draw(dst)

}

func (p *Player) BoundingBox() Rect {
	playerBox := p.RectObject.BoundingBox()
	density := p.scale
	x := 1.0
	playerBox.Left += 25 * density * x
	playerBox.Right -= 35 * density * x
	playerBox.Top += 5 * density
	playerBox.Bottom -= 2 * density

	return playerBox
}

func (p *Player) AppendGround(ground *Ground) {
	p.grounds = append(p.grounds, ground)
}

func (p *Player) Grounded() bool {
	if p.velocity.y < 0 {
		return false
	}

	scale := p.scale

	playerBox := p.BoundingBox()

	for _, ground := range p.grounds {
		groundBox := ground.BoundingBox()
		if playerBox.Overlaps(groundBox) {
			if math.Abs(groundBox.Top-playerBox.Bottom) < 20*scale {
				return true
			}
		}
	}

	return false
}
