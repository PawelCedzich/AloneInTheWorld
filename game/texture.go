package game

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Texture int

const (
	GroundMidT Texture = iota + 1
	GroundLeftT
	GroundRightT
	GroundSingleT
	GroundFillT
	GroundFillSingleT
	GroundFillLeftT
	GroundFillRightT
	GroundFillDeepT
	GroundFillDeepSingleT
	GroundFillDeepLeftT
	GroundFillDeepRightT
	GroundFillDeeperT
	GroundFillDeeper2T
	GroundFillDeeper3T
	BackgroundImageT
	BackgroundTownT
	BackgroundTownFrontT
	StartButtonT
	CharacterT
)

type TextureManager struct {
	textures map[Texture]*ebiten.Image
}

func NewTextureManager(resources map[Texture][]byte) (*TextureManager, error) {
	t := &TextureManager{textures: map[Texture]*ebiten.Image{}}
	for id, buf := range resources {
		img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(buf))
		if err != nil {
			return nil, fmt.Errorf("cant decode Texture %w", err)
		}
		t.textures[id] = img
	}
	return t, nil
}

func (t *TextureManager) LoadTexture(id Texture) *ebiten.Image {
	tex, ok := t.textures[id]
	if !ok {
		panic(fmt.Sprintf("Texture id not defined %d ", id))
	}
	return tex
}

// =====================================================================================================================

type TextureAtlas struct {
	ImagePath  string       `xml:"imagePath,attr"`
	SubTexture []SubTexture `xml:"SubTexture"`
}

type SubTexture struct {
	Name   string  `xml:"name,attr"`
	X      int     `xml:"x,attr"`
	Y      int     `xml:"y,attr"`
	Width  int     `xml:"w,attr"`
	Height int     `xml:"h,attr"`
	PivotX float64 `xml:"pivotX,attr"`
	PivotY float64 `xml:"pivotY,attr"`
}

func ParseTextureAtlas(in io.Reader) (TextureAtlas, error) {
	var atlas TextureAtlas

	dec := xml.NewDecoder(in)
	if err := dec.Decode(&atlas); err != nil {
		return atlas, fmt.Errorf("cant parse texture atlas: %w", err)
	}

	return atlas, nil
}

// =====================================================================================================================

type DrawableTexture struct {
	tex *ebiten.Image
}

func NewDrawableTexture(tex *ebiten.Image) DrawableTexture {
	d := DrawableTexture{tex: tex}
	return d
}

func (d DrawableTexture) Update() {
}

func (d DrawableTexture) Draw(dst *Canvas) {
	dst.DrawImage(d.tex, &ebiten.DrawImageOptions{})
}

func (d DrawableTexture) Size() (x, y float64) {
	texW, texH := d.tex.Size()
	return float64(texW), float64(texH)
}

func (d DrawableTexture) Pivot() Vec {
	v := Vec{
		x: 0,
		y: 0,
	}
	return v
}
