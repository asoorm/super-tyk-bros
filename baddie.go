package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Baddie struct {
	time      int64
	texture   *sdl.Texture
	yPos      float64
	xPos      float64
	xVelocity float64
	yVelocity float64
	width     int32
	height    int32
}

func NewBaddie(r *sdl.Renderer) (*Baddie, error) {

	path := "./assets/sprites/banana-01.png"

	texture, err := img.LoadTexture(r, path)
	if err != nil {
		return nil, fmt.Errorf("could not load baddie image: %v", err)
	}

	return &Baddie{texture: texture, yPos: floorY, xPos: windowWidth, width: 75, height: 80, xVelocity: -0.01}, nil
}

func (k *Baddie) paint(r *sdl.Renderer) error {
	k.time++

	k.xPos -= k.xVelocity

	k.xVelocity += 0.001

	if k.xPos < float64(-k.width-100) {
		k.xPos = windowWidth
	}

	rect := &sdl.Rect{X: int32(k.xPos), Y: int32(k.yPos) - k.height, H: k.height, W: k.width}

	if err := r.Copy(k.texture, nil, rect); err != nil {
		return fmt.Errorf("could not copy tyk: %v", err)
	}
	return nil
}

func (k *Baddie) destroy() {
	k.texture.Destroy()
}
