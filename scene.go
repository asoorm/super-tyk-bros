package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	time int64
	bg   *sdl.Texture
	tyk  *Tyk
}

func NewScene(r *sdl.Renderer) (*Scene, error) {
	t, err := img.LoadTexture(r, "./assets/sprites/title-bg.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background: %v", err)
	}

	tyk, err := NewTyk(r)
	if err != nil {
		return nil, fmt.Errorf("could not load tyk: %v", err)
	}

	return &Scene{bg: t, tyk: tyk}, nil
}

func (s *Scene) paint(r *sdl.Renderer) error {
	s.time++

	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.tyk.paint(r); err != nil {
		return fmt.Errorf("could not paint tyk: %v", err)
	}

	if flagDebug {
		// draw floor
		gfx.LineColor(r, 0, floorY, windowWidth, floorY, sdl.Color{255, 0, 0, 255})
	}

	r.Present()

	return nil
}

func (s *Scene) Run(ctx context.Context, r *sdl.Renderer) <-chan error {

	errChan := make(chan error)

	go func() {
		defer close(errChan)

		for range time.Tick(time.Millisecond * 10) {
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
	s.tyk.destroy()
}
