package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Tyk struct {
	time     int64
	textures []*sdl.Texture
}

func NewTyk(r *sdl.Renderer) (*Tyk, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 3; i++ {
		path := fmt.Sprintf("./assets/sprites/tyk-lg-0%d.png", i)

		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load tyk-0%d: %v", i, err)
		}
		textures = append(textures, texture)
	}

	return &Tyk{textures: textures}, nil
}

func (t *Tyk) paint(r *sdl.Renderer) error {
	t.time++

	i := t.time / 10 % int64(len(t.textures))

	tykRect := &sdl.Rect{X: 10, Y: 495, H: 80, W: 75}
	if err := r.Copy(t.textures[i], nil, tykRect); err != nil {
		return fmt.Errorf("could not copy tyk: %v", err)
	}
	return nil
}

func (t *Tyk) destroy() {
	for _, t := range t.textures {
		t.Destroy()
	}
}
