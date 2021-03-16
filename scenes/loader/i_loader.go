package loader

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
func (l *Load) Build() error {
	var err error

	l.layers = make(map[uint8][]objects.Object)

	if err = l.addBackground(); err != nil {
		return err
	}
	if err = l.addStructure(); err != nil {
		return err
	}
	if err = l.addTxt(); err != nil {
		return err
	}
	if err = l.addBlockLoading(); err != nil {
		return err
	}
	return nil
}

// Init the scene
func (l *Load) Init() error {
	if l.renderer == nil {
		return errors.New(scenes.ErrorRenderer)
	}

	if len(l.layers) != nbrLayers {
		return errors.New(scenes.ErrorLayers)
	}
	if l.loadBlock == nil {
		return errors.New(scenes.ErrorLayers)
	}

	l.initialized = true
	return nil
}

// IsInit return status initialize
func (l Load) IsInit() bool {
	return l.initialized
}

// Run the scene
func (l *Load) Run() error {
	//	go l.addLoadingBar()
	return nil
}

// Stop scene
func (l *Load) Stop() {
	l.resetLoadingBlock()
	//	l.closer <- true
}

// Close the scene
func (l *Load) Close() error {
	var err error

	l.initialized = false

	l.m.Lock()
	defer l.m.Unlock()
	if err = objects.Closer(l.layers); err != nil {
		return err
	}

	l.layers = nil
	return nil
}

// GetLayers provide all scene's objects to draw them
func (l Load) GetLayers() (map[uint8][]objects.Object, *sync.Mutex) {
	return l.layers, l.m
}

// KeyboardEvent provide key down to the scene
func (l *Load) KeyboardEvent(keyboard *sdl.KeyboardEvent) {
	return
}

// SetSwitcher can be call to change scene with index scene and flag closer
func (l *Load) SetSwitcher(f func(uint8, bool) error) {
	return
}

// Update : called on each frame
func (l *Load) Update() {
	l.addLoadingBar()
	return
}

// SetCloser : allow quit the application
func (l *Load) SetCloser(f func()) {
	// nothing to do
}
