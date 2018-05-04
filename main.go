package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"flag"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 800
	windowHeight = 600
	floorY       = 575
	startX       = 10
	gravity      = 9.8
)

var (
	colorTitle = sdl.Color{
		R: 0, G: 0, B: 255, A: 0,
	}

	flagDebug = false
)

func main() {

	flag.BoolVar(&flagDebug, "debug", false, "set debug mode")
	flag.Parse()

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not intitialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not intitialize ttf: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window and renderer: %v", err)
	}
	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}

	time.Sleep(time.Second * 1)

	s, err := NewScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.Destroy()

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(5*time.Second, cancel)

	return <-s.Run(ctx, r)
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	f, err := ttf.OpenFont("assets/fonts/Roboto-Regular.ttf", 20)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := f.RenderUTF8Solid("Super Tyk Bros.", colorTitle)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()

	return nil
}
