package game

import (
	"fmt"
	"strings"
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

	l.res = append(l.res, NewBackground(NewDrawableTexture(g.texture.LoadTexture(Blueshroom))))

	for y, line := range l.raster {
		for x, component := range line {
			pos := Vec{
				x: float64(x) * 100,
				y: float64(y) * 100,
			}
			cell := Rect{pos.x, pos.y, pos.x + 100, pos.y + 100}
			switch component {
			case LevelGround:
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
				//no-op
			case LevelPlatform:
				platformCell := cell
				platformCell.Bottom = platformCell.Top + platformCell.height()/2
				platform := NewGround(NewDrawableTexture(g.texture.LoadTexture(Platform)), platformCell, false, g.engine.Scale())
				l.res = append(l.res, platform)
			default:
				panic(fmt.Errorf("illegal state"))
			}
		}
	}

	return l.res
}

func (l *Level) Get(x, y int) LevelComponent {
	if x < 0 || y < 0 || y >= len(l.raster) || x >= len(l.raster[0]) {
		return LevelSpace
	}
	return l.raster[y][x]
}
