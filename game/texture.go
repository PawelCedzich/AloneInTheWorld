package game

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
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
	ButtonStartT
	ButtonSaveT
	ButtonSettingsT
	ButtonContinueT
	ButtonExitT
	ButtonNoTextT
	ButtonSliderT
	ButtonOnT
	ButtonOffT
	CharacterT
	MortyVanillaT
	MortyJumpingT
	MortyMeditatingT
	MortyWalkingT
	MortyJoyFulT
	MortySadFulT
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
		panic("Texture id not defined ")
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

type TextureAnimation struct {
	frames []*ebiten.Image
	idx    float64
	fps    float64
	w, h   int
	pivot  Vec
}

func NewTextureAnimation(tex *ebiten.Image, atlas TextureAtlas, fps float64) *TextureAnimation {
	if len(atlas.SubTexture) == 0 {
		panic("invalid state subTexture cant be 0")
	}
	t := &TextureAnimation{
		w:   atlas.SubTexture[0].Width,
		h:   atlas.SubTexture[0].Height,
		fps: fps,
		pivot: Vec{
			x: atlas.SubTexture[0].PivotX,
			y: atlas.SubTexture[0].PivotY,
		},
	}

	for _, subTexture := range atlas.SubTexture {
		t.frames = append(t.frames, tex.SubImage(image.Rect(subTexture.X, subTexture.Y, subTexture.X+subTexture.Width, subTexture.Y+subTexture.Height)).(*ebiten.Image))
	}

	return t
}

func (t *TextureAnimation) Update() {
	delta := t.fps / 60
	t.idx += delta
	if t.idx > float64(len(t.frames)-1) {
		t.idx = 0
	}
}

func (t *TextureAnimation) Draw(dst *Canvas) {
	dst.DrawImage(t.frames[int(t.idx)], &ebiten.DrawImageOptions{})
}

func (t *TextureAnimation) Size() (x, y float64) {
	return float64(t.w), float64(t.h)
}

func (t *TextureAnimation) Pivot() Vec {
	return t.pivot
}

// =====================================================================================================================

func NewTextureAnimationFromFrames(fps float64, frames ...*ebiten.Image) *TextureAnimation {
	w, h := frames[0].Size()
	return &TextureAnimation{
		fps:    fps,
		frames: frames,
		idx:    0,
		w:      w,
		h:      h,
	}
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

// =====================================================================================================================

type AnimationGroup struct {
	animations []*TextureAnimation
	idx        int
}

func NewAnimationGroup(animations ...*TextureAnimation) *AnimationGroup {
	a := &AnimationGroup{
		animations: animations,
	}
	return a
}

func (a AnimationGroup) Update() {
	a.animations[a.idx].Update()
}

func (a AnimationGroup) Draw(dst *Canvas) {
	a.animations[a.idx].Draw(dst)
}

func (a AnimationGroup) Size() (x, y float64) {
	texW, texH := a.animations[a.idx].Size()
	return texW, texH
}

func (a AnimationGroup) Pivot() Vec {
	p := a.animations[a.idx].Pivot()
	return p
}

func CenterInside(w, h int, tex *ebiten.Image) *ebiten.Image {
	dst := ebiten.NewImage(w, h)
	op := &ebiten.DrawImageOptions{}
	texW, texH := tex.Size()
	op.GeoM.Translate(float64(w)/2-float64(texW)/2, float64(h)/2-float64(texH)/2)
	dst.DrawImage(tex, op)
	return dst
}

func mirrorY(tex *ebiten.Image) *ebiten.Image {
	dst := ebiten.NewImage(tex.Size())
	op := &ebiten.DrawImageOptions{}
	texW, _ := tex.Size()
	op.GeoM.Scale(-1, 1)
	op.GeoM.Translate(float64(texW), 0)
	dst.DrawImage(tex, op)
	return dst
}
