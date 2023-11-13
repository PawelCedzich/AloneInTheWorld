package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type BoundingBoxer interface {
	BoundingBox() Rect
}

type Goal struct {
	player    *Player
	lost, win bool
	lostText  *Text
	winText   *Text
	winZone   BoundingBoxer
	scale     float64
	onWin     func()
	onLose    func()
}

func NewGoal(player *Player, lostText *Text, winText *Text, winZone BoundingBoxer, scale float64) *Goal {
	return &Goal{
		player:   player,
		lost:     false,
		win:      false,
		lostText: lostText,
		winText:  winText,
		winZone:  winZone,
		scale:    scale,
	}
}

func (g *Goal) Draw(canvas *Canvas) {
	if g.lost {
		savedMatrix := canvas.Transformation()
		var identity ebiten.GeoM
		canvas.SetTransformation(identity)
		g.lostText.Draw(canvas)
		canvas.SetTransformation(savedMatrix)

	} else if g.win {
		savedMatrix := canvas.Transformation()
		var identity ebiten.GeoM
		canvas.SetTransformation(identity)
		g.winText.Draw(canvas)
		canvas.SetTransformation(savedMatrix)

	}
}

func (g *Goal) Layout(sw, sh float64) {
	if g.player.area.Bottom > 1000*g.scale {
		g.lost = true
		//g.player.Lose()
	}

	playerBox := g.player.BoundingBox()
	// if g.player.lookingDirR {
	// 	playerBox = playerBox.ChangeX(0)
	// }

	if playerBox.Overlaps(g.winZone.BoundingBox()) {
		g.win = true
		//g.player.Win()
	}
}
