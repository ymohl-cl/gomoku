package stats

import (
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/gomoku/scenes"
)

/*
** interface functions
 */

// Build describe the scene with objects needest
func (s *Stats) Build() error {
	var err error

	if err = s.addMusic(); err != nil {
		return err
	}
	if err = s.addBackground(); err != nil {
		return err
	}
	if err = s.addStructuresPage(); err != nil {
		return err
	}
	if err = s.addButtons(); err != nil {
		return err
	}
	if err = s.addText(); err != nil {
		return err
	}
	return nil
}

// Init the scene
func (s *Stats) Init() error {
	if s.renderer == nil {
		return errors.New(objects.ErrorRenderer)
	}
	if s.data == nil {
		return errors.New(objects.ErrorData)
	}
	if s.music == nil {
		return errors.New(scenes.ErrorMissing)
	}

	s.initialized = true
	return nil
}

// IsInit return status initialize
func (s Stats) IsInit() bool {
	return s.initialized
}

// Run the scene
func (s Stats) Run() error {
	var wg sync.WaitGroup

	if ok := s.music.IsInit(); !ok {
		if err := s.music.Init(s.renderer); err != nil {
			return err
		}
	}
	wg.Add(1)
	go s.music.Play(&wg, s.renderer)
	wg.Wait()
	return nil
}

// Stop the scene
func (s *Stats) Stop() {
	if ok := s.music.IsInit(); ok {
		if err := s.music.Close(); err != nil {
			panic(err)
		}
	}
}

// Close the scene
func (s *Stats) Close() error {
	var err error

	s.initialized = false

	s.m.Lock()
	defer s.m.Unlock()
	if err = objects.Closer(s.layers); err != nil {
		return err
	}
	return nil
}

// GetLayers provide all scene's objects to draw them
func (s Stats) GetLayers() (map[uint8][]objects.Object, *sync.Mutex) {
	return s.layers, s.m
}

// KeyDownEvent provide key down to the scene
func (s *Stats) KeyDownEvent(keyDown *sdl.KeyDownEvent) {
	return
}

// SetSwitcher can be call to change scene with index scene and flag closer
func (s *Stats) SetSwitcher(f func(uint8, bool) error) {
	s.switcher = f
}

// Update : called on each frame
func (s *Stats) Update() {
	return
}
