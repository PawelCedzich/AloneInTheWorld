package game

import (
	"fmt"
	"math/rand"
	"strings"
)

type LevelComponent rune

const (
	LevelSpace  LevelComponent = ' '
	LevelGround LevelComponent = 'X'
	LevelPlayer LevelComponent = 'P'
	LevelGoal   LevelComponent = 'G'
)

type Level struct {
	raster [][]LevelComponent
	res    []Renderable
	player *Player
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
			case LevelGoal:
				fallthrough
			case LevelSpace:
				raster[lineNo][runeNo] = LevelComponent(char)
			}
		}
	}

	return &Level{raster: raster}, nil
}

func (l *Level) Build(g *Game) []Renderable {

	l.res = append(l.res, NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundImageT))))
	l.res = append(l.res, NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownFrontT))))
	l.res = append(l.res, NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	for y, line := range l.raster {
		for x, component := range line {
			pos := Vec{
				x: float64(x) * 50,
				y: float64(y) * 50,
			}
			cell := Rect{pos.x, pos.y, pos.x + 50, pos.y + 50}
			switch component {
			case LevelPlayer:
				l.player = NewPlayer(NewDrawableTexture(g.texture.LoadTexture(CharacterT)), cell, g.engine.Scale())
			case LevelGoal:
				l.goal = NewRectObject(NewDrawableTexture(g.texture.LoadTexture(ButtonNoTextT)), cell)
			case LevelGround:
				var tex Texture
				var texList []Texture
				aboveGroundCount := 0

				for i := 1; i <= 3; i++ {
					if l.Get(x, y-i) == LevelGround {
						aboveGroundCount++
					} else if l.Get(x, y-i) == LevelSpace {
						break
					}
				}

				if aboveGroundCount == 3 {
					randomNumber := rand.Intn(30) + 1
					if randomNumber >= 1 && randomNumber <= 19 {
						tex = GroundFillDeeperT
					} else if randomNumber >= 20 && randomNumber <= 27 {
						tex = GroundFillDeeper2T
					} else if randomNumber >= 28 {
						tex = GroundFillDeeper3T
					}
				} else if aboveGroundCount == 2 {
					texList = append(texList, GroundFillDeepLeftT, GroundFillDeepRightT, GroundFillDeepSingleT, GroundFillDeepT)
					tex = l.CheckTilePlacement(x, y, texList)
				} else if aboveGroundCount == 1 {
					texList = append(texList, GroundFillLeftT, GroundFillRightT, GroundFillSingleT, GroundFillT)
					tex = l.CheckTilePlacement(x, y, texList)

				} else {
					texList = append(texList, GroundLeftT, GroundRightT, GroundSingleT, GroundMidT)
					tex = l.CheckTilePlacement(x, y, texList)
				}

				ground := NewGround(NewDrawableTexture(g.texture.LoadTexture(tex)), cell, true, g.engine.Scale())
				l.res = append(l.res, ground)
			}
		}
	}

	if l.goal == nil {
		panic(fmt.Errorf("invalid state, no goal"))
	}

	for _, r := range l.res {
		if ground, ok := r.(*Ground); ok {
			l.player.AppendGround(ground)
		}
	}
	loseText := NewText(g.font.LoadFont(TusjF), "You Lost\nTHE END", 24*g.engine.Scale(), 0.5, 0.4)
	winText := NewText(g.font.LoadFont(TusjF), "You Won\nTHE END", 24*g.engine.Scale(), 0.5, 0.4)

	l.res = append(l.res, l.goal)

	l.res = append(l.res, l.player)
	l.res = append(l.res, NewGoal(l.player, loseText, winText, l.goal, g.engine.Scale()))

	return l.res
}

func (l *Level) Get(x, y int) LevelComponent {
	if x < 0 || y < 0 || y >= len(l.raster) || x >= len(l.raster[0]) {
		return LevelSpace
	}
	return l.raster[y][x]
}

func (l *Level) CheckLeftSpace(x, y int) bool {
	if l.Get(x-1, y) == LevelSpace && l.Get(x+1, y) == LevelGround {
		return true
	} else {
		return false
	}
}

func (l *Level) CheckRightSpace(x, y int) bool {
	if l.Get(x-1, y) == LevelGround && l.Get(x+1, y) == LevelSpace {
		return true
	} else {
		return false
	}
}

func (l *Level) CheckTilePlacement(x, y int, texList []Texture) Texture {
	if l.CheckLeftSpace(x, y) {
		return texList[0]
	} else if l.CheckRightSpace(x, y) {
		return texList[1]
	} else if l.Get(x+1, y) == LevelSpace && l.Get(x-1, y) == LevelSpace {
		return texList[2]
	}
	return texList[3]
}
