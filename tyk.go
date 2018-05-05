package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Tyk struct {
	time      int64
	textures  []*sdl.Texture
	yPos      float64
	xPos      float64
	xVelocity float64
	yVelocity float64
	width     int32
	height    int32
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

	return &Tyk{textures: textures, yPos: floorY, xPos: startX, width: 75, height: 80}, nil
}

func (t *Tyk) paint(r *sdl.Renderer) error {
	t.time++

	// which tyk to render
	i := t.time / 10 % int64(len(t.textures))

	t.yPos += t.yVelocity

	// dont let tyk go above the ceiling
	if t.yPos < float64(t.height) {
		t.yPos = float64(t.height)
		t.yVelocity = 0
	}

	// apply gravity
	t.yVelocity += gravity

	// dont let tyk go below the floor.
	if t.yPos > float64(int32(floorY)) {
		t.yPos = float64(int32(floorY))
		t.yVelocity = 0
	}

	// restrict horizontal velocity
	if t.xVelocity > 40 {
		t.xVelocity = 40
	}
	if t.xVelocity < -40 {
		t.xVelocity = -40
	}

	// move x dependent on velocity
	t.xPos = t.xPos + t.xVelocity

	// dont let tyk go to off left of screen
	if t.xPos < 0 {
		t.xPos = 0
		t.xVelocity = 0
	}

	// dont let tyk go to off right of screen
	if t.xPos > float64(windowWidth-t.width) {
		t.xPos = float64(windowWidth - t.width)
		t.xVelocity = 0
	}

	tykRect := &sdl.Rect{X: int32(t.xPos), Y: int32(t.yPos) - t.height, H: t.height, W: t.width}

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
