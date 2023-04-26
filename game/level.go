package game

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type LevelComponent rune

const (
	LevelSpace    LevelComponent = ' '
	LevelGround   LevelComponent = 'X'
	LevelPlatform LevelComponent = 'P'
	LevelPlayer   LevelComponent = '*'
	LevelGoal     LevelComponent = 'M'
	LevelRaccoon  LevelComponent = 'R'
	LevelMono     LevelComponent = 'B'
	LevelVampire  LevelComponent = 'V'
	LevelGold     LevelComponent = 'G'
	LevelOpossum  LevelComponent = 'O'
)

type Level struct {
	raster [][]LevelComponent
	res    []Renderable
	npcs   []*Npc
	player *MainCharacter
	goal   *RectObject
}

func ParseLevel(str string) (*Level, error) {
	lines := strings.Split(str, "\n")

	if len(lines) < 2 {
		return nil, fmt.Errorf("level must have at least 2 lines %d ", len(lines))
	}

	width := len(lines[0])
	if width < 1 {
		return nil, fmt.Errorf("level must have at least a single column to work %d", width)
	}

	raster := make([][]LevelComponent, len(lines), len(lines))

	hasPlayer := false
	for lineNo, line := range lines {
		if len(line) != width {
			return nil, fmt.Errorf("the width of the lvl must be same, the line %d has a wdith of %d expected %d", lineNo, len(line), width)
		}

		raster[lineNo] = make([]LevelComponent, width, width)
		for runeNo, char := range line {
			switch LevelComponent(char) {
			case LevelPlayer:
				if hasPlayer {
					return nil, fmt.Errorf("second player declaration in %d line and %d column", lineNo, runeNo)
				}
				hasPlayer = true
				fallthrough
			case LevelGround:
				fallthrough
			case LevelGold:
				fallthrough
			case LevelSpace:
				fallthrough
			case LevelRaccoon:
				fallthrough
			case LevelMono:
				fallthrough
			case LevelVampire:
				fallthrough
			case LevelGoal:
				fallthrough
			case LevelOpossum:
				fallthrough
			case LevelPlatform:
				raster[lineNo][runeNo] = LevelComponent(char)
			default:
				return nil, fmt.Errorf("invalid level component in %d line %d column %s", lineNo, runeNo, string(char))
			}
		}
	}

	if !hasPlayer {
		return nil, fmt.Errorf("level has no declared player need at least one")
	}

	return &Level{raster: raster}, nil
}

func (l *Level) Build(g *Game) []Renderable {
	var platformGroup *PlatformGroup
	platformMoveLeft := false

	l.res = append(l.res, NewBackground(NewDrawableTexture(g.texture.LoadTexture(Blueshroom))))

	for y, line := range l.raster {
		for x, component := range line {
			pos := Vec{
				x: float64(x) * 100,
				y: float64(y) * 100,
			}
			cell := Rect{pos.x, pos.y, pos.x + 100, pos.y + 100}
			switch component {
			case LevelPlayer:
				platformGroup = nil
				l.player = l.NewPlayer(g, cell)
			case LevelGround:
				platformGroup = nil
				tex := GroundFill
				if l.Get(x, y-1) != LevelGround {
					tex = GroundMid
					if l.Get(x, y+1) != LevelGround {
						if l.Get(x-1, y) == LevelSpace && l.Get(x+1, y) == LevelGround {
							tex = GroundLeft
						} else if l.Get(x+1, y) == LevelSpace && l.Get(x-1, y) == LevelGround {
							tex = GroundRight
						}
					}
				}

				ground := NewGround(NewDrawableTexture(g.texture.LoadTexture(tex)), cell, true, g.engine.Scale())
				l.res = append(l.res, ground)
			case LevelSpace:
				platformGroup = nil
				//no-op
			case LevelPlatform:
				platformCell := cell
				platformCell.Bottom = platformCell.Top + platformCell.height()/2
				platform := NewGround(NewDrawableTexture(g.texture.LoadTexture(Platform)), platformCell, false, g.engine.Scale())
				l.res = append(l.res, platform)

				if platformGroup == nil {
					platformGroup = NewPlatformGroup(g.engine.Scale(), platformMoveLeft)
					platformMoveLeft = !platformMoveLeft
					l.res = append(l.res, platformGroup)
				}
				platformGroup.Add(platform)
			case LevelGoal:
				l.goal = NewRectObject(NewMortimer(g.texture.LoadTexture(ProfessorHead), g.texture.LoadTexture(ProfessorBody)), cell.Scale(g.engine.Scale()))
			case LevelRaccoon:
				min, max := l.GroundGroup(x, y, cell.Scale(g.engine.Scale()))
				anim := NewAnimationGroup(l.MustParseAnimation(g, MortyRaccoonWalking, assets.MortyRaccoonWalkingAtlas))
				r := NewNpc(anim, cell.Scale(g.engine.Scale()), g.engine.Scale(), 2, 15, 10, min, max, false)
				l.npcs = append(l.npcs, r)
			case LevelMono:
				min, max := l.GroundGroup(x, y, cell.Scale(g.engine.Scale()))
				anim := NewAnimationGroup(l.MustParseAnimation(g, MortySadMono, assets.MortySadMonoAtlas))
				r := NewNpc(anim, cell.Scale(g.engine.Scale()), g.engine.Scale(), 1, 60, 20, min, max, false)
				l.npcs = append(l.npcs, r)
			case LevelVampire:
				min, max := l.GroundGroup(x, y, cell.Scale(g.engine.Scale()))
				anim := NewAnimationGroup(l.MustParseAnimation(g, MortyAlbino, assets.MortyAlbinoAtlas))
				r := NewNpc(anim, cell.Scale(g.engine.Scale()), g.engine.Scale(), 6, 30, 15, min, max, false)
				l.npcs = append(l.npcs, r)
			case LevelOpossum:
				min, max := l.GroundGroup(x, y, cell.Scale(g.engine.Scale()))
				anim := NewAnimationGroup(l.MustParseAnimation(g, MortyOpossum, assets.MortyOpossumAtlas))
				r := NewNpc(anim, cell.Scale(g.engine.Scale()), g.engine.Scale(), 2, 15, 10, min, max, true)
				l.npcs = append(l.npcs, r)
			case LevelGold:
				tex := NewTextureAnimationFromFrames(10,
					CenterInside(84, 84, g.texture.LoadTexture(Gold1)),
					CenterInside(84, 84, g.texture.LoadTexture(Gold2)),
					CenterInside(84, 84, g.texture.LoadTexture(Gold3)),
					CenterInside(84, 84, g.texture.LoadTexture(Gold4)),
					mirrorY(CenterInside(84, 84, g.texture.LoadTexture(Gold3))),
					mirrorY(CenterInside(84, 84, g.texture.LoadTexture(Gold2))),
				)
				gold := NewCollectable(tex, cell.Scale(g.engine.Scale()))
				l.res = append(l.res, gold)
			default:
				platformGroup = nil
				panic(fmt.Errorf("illegal state"))
			}
		}
	}

	if l.player == nil {
		panic(fmt.Errorf("invalid state,  no player"))
	}

	if l.goal == nil {
		panic(fmt.Errorf("invalid state, no goal"))
	}

	for _, npc := range l.npcs {
		l.res = append(l.res, npc)
		npc.register(l.player)
		l.player.npcs = append(l.player.npcs, npc)
		l.player.projectiles.RegisterNpc(npc)
	}

	for _, r := range l.res {
		if ground, ok := r.(*Ground); ok {
			for _, npc := range l.npcs {
				npc.RegisterGround(ground)
			}
			l.player.RegisterGround(ground)
			l.player.projectiles.RegisterGround(ground)
		}
	}
	for _, r := range l.res {
		if c, ok := r.(*Collectable); ok {
			l.player.RegisterCollectable(c)
		}
	}

	l.res = append(l.res, l.player)

	loseText := NewRelativeText(g.font.LoadFont(OpenSans), "You Lost\nTHE END", 24*g.engine.Scale(), 0.5, 0.4)
	winText := NewRelativeText(g.font.LoadFont(OpenSans), "You Won\nTHE END", 24*g.engine.Scale(), 0.5, 0.4)

	h := NewPlayerHud(g.engine.Scale(), g.font.LoadFont(OpenSans), g.engine.Scale(),
		NewIconTextView(NewDrawableTexture(g.texture.LoadTexture(Nut)), " x ", "3", g.engine.Scale()),
		NewIconTextView(NewDrawableTexture(g.texture.LoadTexture(Gold1)), " x ", "0", g.engine.Scale()),
		NewIconTextView(NewDrawableTexture(g.texture.LoadTexture(ProfessorHead)), " x ", "0", g.engine.Scale()))
	l.res = append(l.res, h)
	l.res = append(l.res, l.goal)

	l.player.updateStats = func() {
		h.HudElements[0].value = strconv.Itoa(l.player.projectiles.available)
		h.HudElements[1].value = strconv.Itoa(l.player.goldCounter)
		h.HudElements[2].value = strconv.Itoa(l.player.score)
	}

	l.res = append(l.res, NewGoal(l.player, loseText, winText, l.goal, g.engine.Scale()))
	return l.res
}

func (l *Level) Get(x, y int) LevelComponent {
	if x < 0 || y < 0 || y >= len(l.raster) || x >= len(l.raster[0]) {
		return LevelSpace
	}
	return l.raster[y][x]
}

func (l *Level) NewPlayer(g *Game, cell Rect) *MainCharacter {
	mortyCell := cell
	mortyCell.Right *= 1.05

	anim := NewAnimationGroup(
		l.MustParseAnimation(g, MortyVanilla, assets.MortyVanillaAtlas),
		l.MustParseAnimation(g, MortyJumping, assets.MortyJumpingAtlas),
		l.MustParseAnimation(g, MortyWalking, assets.MortyWalkingAtlas),
		l.MustParseAnimation(g, MortyMeditating, assets.MortyMeditatingAtlas),
		l.MustParseAnimation(g, MortySadFul, assets.MortySadFulAtlas),
		l.MustParseAnimation(g, MortyJoyFul, assets.MortyJoyFulAtlas),
	)

	const (
		Idle = iota
		Jumping
		Walking
		Meditating
		SadFul
		JoyFul
	)

	p := NewMainCharacter(anim, mortyCell.Scale(g.engine.Scale()), g.engine.Scale(), NewDrawableTexture(g.texture.LoadTexture(Nut)), 3, 3)

	p.onMovement = func() {
		anim.idx = Walking
	}

	p.onJump = func() {
		anim.idx = Jumping
	}

	p.onIdle = func() {
		anim.idx = Idle
	}

	p.onBoring = func() {
		anim.idx = Meditating
	}

	start := time.Now()
	onShotWin := false
	p.onWin = func() {
		anim.idx = JoyFul
		if !onShotWin {
			g.music.Play(PetAha)
			onShotWin = true
			completionTime := time.Now().Sub(start)
			max := time.Minute
			var points int
			if completionTime < max {
				pt := completionTime.Seconds() / max.Seconds()
				points = int(1000 * (1 - pt))
			}
			p.score += points
		}

	}

	p.onCollect = func() {
		g.music.Play(OpeningChest)
	}

	p.onHit = func() {
		g.music.Play(PetDisappointed)
	}

	onShotLose := false
	p.onLose = func() {
		anim.idx = SadFul
		if !onShotLose {
			g.music.Play(PetUpset)
			onShotLose = true
		}
	}

	return p
}

func (l *Level) MustParseAnimation(g *Game, id Texture, xml []byte) *TextureAnimation {
	atl, err := ParseTextureAtlas(bytes.NewReader(xml))
	if err != nil {
		panic(fmt.Errorf("invalid state by parsing texture atlas %w", err))
	}

	return NewTextureAnimation(g.texture.LoadTexture(id), atl, 30)
}

func (l *Level) GroundGroup(x, y int, cell Rect) (float64, float64) {
	min := float64(x) * 100
	max := float64(x)*100 + cell.width()

	for i := 1; i < 10; i++ {
		if l.Get(x-i, y) != LevelGround {
			if l.Get(x-i, y+1) == LevelGround {
				min -= cell.width()
			} else if l.Get(x-i, y+1) == LevelSpace {
				break
			}
		}
	}

	for i := 1; i < 10; i++ {
		if l.Get(x+i, y) != LevelGround {
			if l.Get(x+i, y+1) == LevelGround {
				max += cell.width()
			} else if l.Get(x+i, y+1) == LevelSpace {
				break
			}
		}
	}

	return min, max
}
