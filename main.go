package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	windowWidth  = 800
	windowHeight = 600
	floorY       = 575
	startX       = 10
	gravity      = 0.3
)

var (
	gameTitle  = "Super Tyk Bros."
	colorTitle = sdl.Color{
		R: 0, G: 0, B: 255, A: 0,
	}

	flagDebug = false

	log = logrus.WithField("prefix", "main")
)

func init() {
	flag.BoolVar(&flagDebug, "debug", false, "set debug mode")

	flag.Parse()

	if flagDebug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {

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

	time.Sleep(time.Second * 3)

	s, err := NewScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.Destroy()

	events := make(chan sdl.Event)

	errChan := s.Run(events, r)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errChan:
			return err
		}
	}
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	f, err := ttf.OpenFont("assets/fonts/Roboto-Regular.ttf", 20)
	if err != nil {
		return err
	}
	defer f.Close()

	s, err := f.RenderUTF8Solid(gameTitle, colorTitle)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	var width, height = int32(400), int32(300)
	rect := &sdl.Rect{X: windowWidth/2 - width/2, Y: windowHeight/2 - height/2, W: width, H: height}
	if err := r.Copy(t, nil, rect); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()

	return nil
}
