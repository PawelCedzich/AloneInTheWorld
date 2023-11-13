package game

import (
	"fmt"

	"golang.org/x/image/font/opentype"
)

type Font int

const (
	TusjF Font = iota + 1
)

type FontManager struct {
	fonts map[Font]*opentype.Font
}

func NewFontManager(resources map[Font][]byte) (*FontManager, error) {
	parsedFont := map[Font]*opentype.Font{}
	for id, buf := range resources {
		tt, err := opentype.Parse(buf)
		if err != nil {
			return nil, fmt.Errorf("cant initialiye fontmanager %w", err)
		}
		parsedFont[id] = tt
	}
	return &FontManager{fonts: parsedFont}, nil
}

func (f *FontManager) LoadFont(id Font) *opentype.Font {
	font, ok := f.fonts[id]
	if !ok {
		panic("Font id not defined ")
	}
	return font
}
