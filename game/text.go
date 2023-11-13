package game

import (
	"image/color"
	"strings"

	"golang.org/x/image/font/opentype"
)

type Text struct {
	fontFace *opentype.Font
	lines    []string
	fontSize float64
	w, h     float64
}

func NewText(fon *opentype.Font, lines string, fontSize, w, h float64) *Text {
	return &Text{
		fontFace: fon,
		lines:    strings.Split(lines, "\n"),
		fontSize: fontSize,
		w:        w,
		h:        h,
	}
}

func (t *Text) Draw(dst *Canvas) {
	sw, sh := dst.Size()
	offset := -t.fontSize
	linesScape := t.fontSize / 4

	dst.SetTextSize(t.fontSize)
	dst.SetFont(t.fontFace)

	for _, line := range t.lines {
		rect := dst.MeasureText(line)
		tw, th := rect.width(), rect.height()
		offset += th + linesScape

		dst.SetColor(color.Black)
		shadowWidth := t.fontSize / 12
		for x := -shadowWidth; x <= shadowWidth; x++ {
			for y := -shadowWidth; y <= shadowWidth; y++ {
				dst.DrawTexT(line, sw*t.w-float64(tw)/2+x, sh*t.h+offset+y)
			}
		}

		dst.SetColor(color.White)
		dst.DrawTexT(line, sw*t.w-float64(tw)/2, sh*t.h+offset)
	}
}

func (t *Text) Layout(sw, sh float64) {
}
