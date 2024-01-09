package game

import (
	"bytes"
	"fmt"

	audio2 "github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

type Audio int

const (
	Music Audio = iota + 1
)

type AudioManager struct {
	audios       map[Audio]*audio2.Player
	audioContext *audio2.Context
}

func NewAudioManager(resources map[Audio][]byte) (*AudioManager, error) {
	p := map[Audio]*audio2.Player{}
	c := audio2.NewContext(32000)
	for id, buf := range resources {
		s, err := mp3.DecodeWithSampleRate(32000, bytes.NewReader(buf))
		if err != nil {
			return nil, fmt.Errorf("cant decode music %w", err)
		}
		p[id], err = c.NewPlayer(s)
		if err != nil {
			return nil, fmt.Errorf("cant create new player %w", err)
		}
	}
	a := &AudioManager{
		audios:       p,
		audioContext: c,
	}
	return a, nil
}

func (a *AudioManager) LoadAudio(id Audio) *audio2.Player {
	audio, ok := a.audios[id]
	if !ok {
		panic("Audio id not defined ")
	}
	return audio
}

func (a *AudioManager) Play(audio Audio) {
	player, ok := a.audios[audio]

	if !ok {
		panic(fmt.Errorf("cant find audio %v", audio))
	}
	if err := player.Rewind(); err != nil {
		panic(fmt.Errorf("cant rewind player invalid state %v", player))
	}

	player.Play()
}

func (a *AudioManager) ChangeVolume(value float64) {

	for _, player := range a.audios {
		player.SetVolume(value)
	}
}
