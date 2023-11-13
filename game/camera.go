package game

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	player *Player
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) UpdateMainCharacter(player *Player) {
	c.player = player
}

func (c *Camera) Transformation(sw, sh float64) ebiten.GeoM {

	var m ebiten.GeoM
	if c.player != nil {
		m.Translate(-c.player.BoundingBox().Left+sw/2-c.player.BoundingBox().width()/2, -c.player.BoundingBox().Top+sh/2-c.player.BoundingBox().height()/2)
	}
	return m
}
