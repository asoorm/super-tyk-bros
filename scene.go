package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	time int64
	bg   *sdl.Texture
	tyks []*sdl.Texture
}

func NewScene(r *sdl.Renderer) (*Scene, error) {
	t, err := img.LoadTexture(r, "./assets/sprites/title-bg.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background: %v", err)
	}

	var tyks []*sdl.Texture
	for i := 1; i <= 3; i++ {
		path := fmt.Sprintf("./assets/sprites/tyk-lg-0%d.png", i)

		tyk, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load tyk-0%d: %v", i, err)
		}
		tyks = append(tyks, tyk)
	}

	return &Scene{bg: t, tyks: tyks}, nil
}

func (s *Scene) paint(r *sdl.Renderer) error {
	s.time++

	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	i := s.time / 10 % int64(len(s.tyks))

	tykRect := &sdl.Rect{X: 10, Y: 495, H: 80, W: 75}
	if err := r.Copy(s.tyks[i], nil, tykRect); err != nil {
		return fmt.Errorf("could not copy tyk: %v", err)
	}

	r.Present()

	return nil
}

func (s *Scene) Run(ctx context.Context, r *sdl.Renderer) <-chan error {

	errChan := make(chan error)

	go func() {
		defer close(errChan)

		for range time.Tick(time.Millisecond * 100) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := s.paint(r); err != nil {
					errChan <- err
				}
			}
		}
	}()

	return errChan
}

func (s *Scene) Destroy() {
	s.bg.Destroy()
}
