package game

type BoundingBoxer interface {
	BoundingBox() Rect
}

type Goal struct {
	player       *Player
	lostText     *Text
	winText      *Text
	winZone      BoundingBoxer
	scale        float64
	onWinAction  func()
	onLoseAction func()
}

func NewGoal(player *Player, winZone BoundingBoxer, scale float64, onwin func(), onLose func()) *Goal {
	return &Goal{
		player:       player,
		winZone:      winZone,
		scale:        scale,
		onWinAction:  onwin,
		onLoseAction: onLose,
	}
}

func (g *Goal) Draw(canvas *Canvas) {
}

func (g *Goal) Layout(sw, sh float64) {
	if g.player.area.Bottom > 1000*g.scale {
		g.onLoseAction()
	}

	playerBox := g.player.BoundingBox()
	if g.player.lookingDirR {
		playerBox = playerBox.ChangeX(0)
	}

	if playerBox.Overlaps(g.winZone.BoundingBox()) {
		g.onWinAction()
	}
}
