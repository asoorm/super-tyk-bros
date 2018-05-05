package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	time   int64
	bg     *sdl.Texture
	tyk    *Tyk
	baddie *Baddie
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

	baddie, err := NewBaddie(r)
	if err != nil {
		return nil, fmt.Errorf("could not load baddie: %v", err)
	}

	return &Scene{bg: t, tyk: tyk, baddie: baddie}, nil
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

	if err := s.baddie.paint(r); err != nil {
		return fmt.Errorf("could not paint baddie: %v", err)
	}

	//if flagDebug {
	//	// draw floor
	//	gfx.LineColor(r, 0, floorY, windowWidth, floorY, sdl.Color{255, 0, 0, 255})
	//}

	r.Present()

	return nil
}

func (s *Scene) Run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {

	errChan := make(chan error)

	go func() {
		defer close(errChan)

		tick := time.Tick(time.Millisecond * 10)

		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				if err := s.paint(r); err != nil {
					errChan <- err
				}
			}
		}
	}()

	return errChan
}

func (s *Scene) handleEvent(event sdl.Event) bool {

	log = log.WithField("prefix", "event")

	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseMotionEvent, *sdl.TouchFingerEvent, *sdl.MouseButtonEvent, *sdl.WindowEvent,
		*sdl.AudioDeviceEvent:
		return false
	case *sdl.KeyboardEvent:
		switch e.Type {
		case sdl.KEYDOWN:
			switch e.Keysym.Sym {
			case sdl.K_LEFT:
				log.Debug("pressed: left")
				restriction := 1.0
				// if in air, restrict movement
				if s.tyk.yPos != floorY {
					restriction = 0.5
				}
				s.tyk.xVelocity = -5 * restriction
			case sdl.K_RIGHT:
				log.Debug("pressed: right")
				restriction := 1.0
				// if in air, restrict movement
				if s.tyk.yPos != floorY {
					restriction = 0.5
				}
				s.tyk.xVelocity = 5 * restriction
			case sdl.K_UP:
				log.Debug("pressed: up")
				if s.tyk.yPos == floorY {
					log.Debug("onFloor")
					s.tyk.yVelocity -= 10
				}
				//case sdl.K_DOWN:
				//	log.Info("pressed: down")
				//	s.tyk.yVelocity += 5
			}
		case sdl.KEYUP:
			switch e.Keysym.Sym {
			case sdl.K_LEFT:
				log.Debug("released: left")
				// if tyk is moving to the left, we zero the velocity
				if s.tyk.xVelocity < 0 {
					s.tyk.xVelocity = 0
				}
			case sdl.K_RIGHT:
				log.Debug("released: right")
				// if tyk is moving to the right, we zero the velocity
				if s.tyk.xVelocity > 0 {
					s.tyk.xVelocity = 0
				}
				//case sdl.K_UP:
				//	log.Info("released: up")
				//	if s.tyk.yVelocity > 0 || s.tyk.yPos < float64(floorY-s.tyk.height*4) {
				//		s.tyk.yVelocity = 0
				//	}
				//case sdl.K_DOWN:
				//	log.Info("released: down")
			}
		}
	default:
		log.Debug("other event: %T", event)
	}

	return false
}

func (s *Scene) Destroy() {
	s.bg.Destroy()
	s.tyk.destroy()
}
