package game

type Npc struct {
	*RectObject
	grounds        []*Ground
	player         *Player
	scale          float64
	snapped        bool
	velocity       Vec
	speed          float64
	stunDuration   float64
	pushPower      float64
	min, max       float64
	turnOffDrawing bool
	walkLeft       bool
	target         bool
	canWalk        bool
}

func NewNpc(tex Drawable, area Rect, scale float64, speed float64, stunDuration float64, pushPower float64, min float64, max float64, target bool) *Npc {
	n := &Npc{
		RectObject: NewRectObject(tex, area),
		scale:      scale,
		speed:      speed * scale,
		velocity: Vec{
			x: speed * scale,
			y: 0,
		},
		stunDuration: stunDuration,
		pushPower:    pushPower,
		min:          min,
		max:          max,
		target:       target,
	}

	return n
}

func (n *Npc) Layout(sw, sh float64) {
	n.RectObject.Layout(sw, sh)

	npcBox := n.BoundingBox()
	for _, ground := range n.grounds {
		gBox := ground.BoundingBox()
		topRegion := gBox.TopRegion(n.scale * 20)
		if topRegion.Overlaps(npcBox) {
			delta := npcBox.Bottom - gBox.Top
			n.area = n.area.Offset(0, -delta)
		}
	}

	if n.walkLeft && n.CanLeft() {
		n.velocity.x = -n.speed
	} else if n.CanRight() {
		n.velocity.x = n.speed
		n.walkLeft = false
	} else {
		n.walkLeft = true
	}

	if n.target {
		playerIsLeft := (n.player.area.Left+n.player.area.width()/2)-(n.area.Left+n.area.width()/2) <= 1*n.scale
		n.walkLeft = playerIsLeft

		if playerIsLeft && !n.CanLeft() {
			n.velocity.x = 0
		}
		if !playerIsLeft && !n.CanRight() {
			n.velocity.x = 0
		}
	}

	n.mirrorYaxis = !n.walkLeft
	n.area = n.area.Offset(n.velocity.x, n.velocity.y)
}

func (n *Npc) Draw(dst *Canvas) {

	if !n.turnOffDrawing {
		n.RectObject.Draw(dst)
		//n.RectObject.BoundingBox().Draw(dst)
		//n.BoundingBox().Draw(dst)
	}
}

func (n *Npc) BoundingBox() Rect {
	box := n.RectObject.BoundingBox()
	density := n.scale
	x := 1.0

	box.Left += 7 * density * x
	box.Right -= 43 * density * x
	box.Top += 12 * density
	box.Bottom -= 10 * density

	return box
}

func (n *Npc) RegisterGround(ground *Ground) {
	n.grounds = append(n.grounds, ground)
}

func (n *Npc) PushPLayer(p *Player) {

	p.stun = n.stunDuration
	p.velocity = Vec{
		x: -p.velocity.x * 2,
		y: -n.pushPower * n.scale,
	}
}

func (n *Npc) register(player *Player) {
	n.player = player
}

func (n *Npc) TouchedGround(npcBox Rect) *Ground {
	delta := n.BoundingBox().width() / 2.5

	for _, ground := range n.grounds {
		gBox := ground.BoundingBox()

		if gBox.LeftRegion(delta).Overlaps(npcBox) || gBox.RightRegion(delta).Overlaps(npcBox) {
			return nil
		}

		topRegion := gBox.TopRegion(n.scale*20).withPadding(delta, 0)
		if topRegion.Overlaps(npcBox) {
			return ground
		}
	}
	return nil
}

func (n *Npc) CanLeft() bool {
	return n.TouchedGround(n.BoundingBox().Offset(-n.speed, 0)) != nil
}

func (n *Npc) CanRight() bool {
	return n.TouchedGround(n.BoundingBox().Offset(n.speed, 0)) != nil
}
