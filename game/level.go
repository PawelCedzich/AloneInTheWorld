package game

import (
	"fmt"
	"strings"
)

type LevelComponent rune

const (
	LevelSpace  LevelComponent = ' '
	LevelGround LevelComponent = 'X'
	LevelPlayer LevelComponent = 'P'
)

type Level struct {
	raster [][]LevelComponent
	res    []Renderable
	player *Player
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
			}
		}
	}

	for _, r := range l.res {
		if ground, ok := r.(*Ground); ok {
			l.player.AppendGround(ground)
		}
	}

	l.res = append(l.res, l.player)

	return l.res
}

func (l *Level) Get(x, y int) LevelComponent {
	if x < 0 || y < 0 || y >= len(l.raster) || x >= len(l.raster[0]) {
		return LevelSpace
	}
	return l.raster[y][x]
}
