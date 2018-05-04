package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Tyk struct {
	time     int64
	textures []*sdl.Texture
	y        float64
	ySpeed   float64
	x        float64
	w        int32
	h        int32
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

	return &Tyk{textures: textures, y: floorY, x: startX, w: 75, h: 80}, nil
}

func (t *Tyk) paint(r *sdl.Renderer) error {
	t.time++

	i := t.time / 10 % int64(len(t.textures))

	//tykRect := &sdl.Rect{X: int32(t.x), Y: windowHeight - int32(t.y), H: t.h, W: t.w}
	tykRect := &sdl.Rect{X: int32(t.x), Y: int32(t.y) - t.h, H: t.h, W: t.w}
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

func (t *Tyk) jump() {
	t.ySpeed += 10
}
