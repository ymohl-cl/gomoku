package gomoku

import (
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/gomoku/controller"
	"github.com/ymohl-cl/gomoku/game"
	"github.com/ymohl-cl/gomoku/scenes"
)

/*
** interface functions
 */

// Build describe the scene with objects needest
func (g *Gomoku) Build() error {
	var err error

	g.layers = make(map[uint8][]objects.Object)

	if err = g.addMusic(); err != nil {
		return err
	}
	if err = g.addBackground(); err != nil {
		return err
	}
	if err = g.addStructures(); err != nil {
		return err
	}
	if err = g.addBlockTime(); err != nil {
		return err
	}
	if err = g.addText(); err != nil {
		return err
	}
	if err = g.addTokens(); err != nil {
		return err
	}
	if err = g.addButtons(); err != nil {
		return err
	}
	if err = g.addNotice(); err != nil {
		return err
	}
	/*
		if err = g.addModal(); err != nil {
			return err
		}
	*/
	g.layers[layerHistoryBlock] = nil
	g.layers[layerHistoryText] = nil
	return nil
}

// Init the scene
func (g *Gomoku) Init() error {
	if g.renderer == nil {
		return errors.New(objects.ErrorRenderer)
	}
	if g.layers == nil {
		return errors.New(scenes.ErrorLayers)
	}
	if g.music == nil {
		return errors.New(scenes.ErrorMissing)
	}
	if g.notice == nil {
		return errors.New(scenes.ErrorMissing)
	}

	g.initialized = true
	return nil
}

// IsInit return status initialize
func (g *Gomoku) IsInit() bool {
	return g.initialized
}

// Run the scene
func (g *Gomoku) Run() error {
	var err error
	var wg sync.WaitGroup

	if g.game, err = game.New(g.data); err != nil {
		return err
	}

	if ok := g.music.IsInit(); ok {
		wg.Add(1)
		go g.music.Play(&wg, g.renderer)
		wg.Wait()
	}
	g.game.Run()
	g.initMove()
	return nil
}

// Stop the scene
func (g *Gomoku) Stop() {
	if g.music.IsInit() {
		if err := g.music.Close(); err != nil {
			panic(err)
		}
	}
}

// Close the scene
func (g *Gomoku) Close() error {
	var err error

	g.initialized = false
	g.m.Lock()
	defer g.m.Unlock()
	if err = objects.Closer(g.layers); err != nil {
		return err
	}

	g.layers = nil
	g.game = nil
	return nil
}

// GetLayers provide all scene's objects to draw them
func (g *Gomoku) GetLayers() (map[uint8][]objects.Object, *sync.Mutex) {
	return g.layers, g.m
}

// KeyDownEvent provide key down to the scene
func (g *Gomoku) KeyDownEvent(keyDown *sdl.KeyDownEvent) {
	return
}

// SetSwitcher can be call to change scene with index scene and flag closer
func (g *Gomoku) SetSwitcher(f func(uint8, bool) error) {
	g.switcher = f
}

// Update : called on each frame
func (g *Gomoku) Update() {
	var err error

	duration := g.game.GetTimeGame()
	str := controller.TimeToString(duration)
	// check if need to change
	if g.timer.GetTxt() == str {
		return
	}

	if err = g.timer.UpdateText(str, g.renderer); err != nil {
		panic(err)
	}
}

// SetCloser : allow quit the application
func (g *Gomoku) SetCloser(f func()) {
	// nothing to do
}
