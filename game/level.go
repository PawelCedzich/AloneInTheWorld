package game

import (
	"fmt"
	"strings"
)

type LevelComponent rune

const (
	LevelSpace  LevelComponent = ' '
	LevelGround LevelComponent = 'X'
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

	for lineNo, line := range lines {
		if len(line) != width {
			return nil, fmt.Errorf("the width of the lvl must be same, the line %d has a wdith of %d expected %d", lineNo, len(line), width)
		}

		raster[lineNo] = make([]LevelComponent, width, width)
		for runeNo, char := range line {
			switch LevelComponent(char) {
			case LevelGround:
				raster[lineNo][runeNo] = LevelComponent(char)
			case LevelSpace:
				raster[lineNo][runeNo] = LevelComponent(char)
			default:
				//return nil, fmt.Errorf("invalid level component in %d line %d column %s", lineNo, runeNo, string(char))
			}
		}
	}

	return &Level{raster: raster}, nil
}

func (l *Level) Build(g *Game) []Renderable {

	l.res = append(l.res, NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundImageT))))
	l.res = append(l.res, NewBackground(NewDrawableTexture(g.texture.LoadTexture(BackgroundTownT))))

	for y, line := range l.raster {
		for x, component := range line {
			pos := Vec{
				x: float64(x) * 100,
				y: float64(y) * 100,
			}
			cell := Rect{pos.x, pos.y, pos.x + 100, pos.y + 100}
			switch component {
			case LevelGround:
				tex := GroundFillT
				if l.Get(x, y-1) != LevelGround {
					tex = GroundMidT
					if l.Get(x, y+1) != LevelGround {
						if l.Get(x-1, y) == LevelSpace && l.Get(x+1, y) == LevelGround {
							tex = GroundLeftT
						} else if l.Get(x+1, y) == LevelSpace && l.Get(x-1, y) == LevelGround {
							tex = GroundRightT
						}
					}
				}

				ground := NewGround(NewDrawableTexture(g.texture.LoadTexture(tex)), cell, true, g.engine.Scale())
				l.res = append(l.res, ground)
			default:
				//panic(fmt.Errorf("illegal state"))
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
