package menu

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
func (m *Menu) Build() error {
	var err error

	m.layers = make(map[uint8][]objects.Object)

	if err = m.addMusic(); err != nil {
		return err
	}
	if err = m.addBackground(); err != nil {
		return err
	}
	if err = m.addStructuresPage(); err != nil {
		return err
	}
	if err = m.addButtons(); err != nil {
		return err
	}
	if err = m.addNotice(); err != nil {
		return err
	}
	if err = m.addText(); err != nil {
		return err
	}

	if err = m.addVS(); err != nil {
		return err
	}

	if err = m.addJS(); err != nil {
		return err
	}

	if err = m.addInput(); err != nil {
		return err
	}
	if err = m.addPlayers(); err != nil {
		return err
	}
	return nil
}

// Init the scene
func (m *Menu) Init() error {
	if m.renderer == nil {
		return errors.New(objects.ErrorRenderer)
	}
	if m.data == nil {
		return errors.New(objects.ErrorData)
	}
	if len(m.layers) == 0 {
		return errors.New(scenes.ErrorLayers)
	}
	if m.input == nil {
		return errors.New(scenes.ErrorMissing)
	}
	if m.notice == nil {
		return errors.New(scenes.ErrorMissing)
	}
	if m.music == nil {
		return errors.New(scenes.ErrorMissing)
	}
	if m.player1 == nil || m.player2 == nil {
		return errors.New(scenes.ErrorMissing)
	}
	/*	if m.vs == nil {
		return errors.New(scenes.ErrorMissing)
	}*/

	m.initialized = true
	return nil
}

// IsInit return status initialize
func (m Menu) IsInit() bool {
	return m.initialized
}

// Run the scene
func (m Menu) Run() error {
	var wg sync.WaitGroup

	if ok := m.music.IsInit(); !ok {
		if err := m.music.Init(m.renderer); err != nil {
			return err
		}
	}
	wg.Add(1)
	go m.music.Play(&wg, m.renderer)
	wg.Wait()
	return nil
}

// Stop the scene
func (m *Menu) Stop() {
	if ok := m.music.IsInit(); ok {
		if err := m.music.Close(); err != nil {
			panic(err)
		}
	}
}

// Close the scene
func (m *Menu) Close() error {
	var err error

	m.initialized = false

	m.m.Lock()
	defer m.m.Unlock()
	if err = objects.Closer(m.layers); err != nil {
		return err
	}

	m.layers = nil
	return nil
}

// GetLayers provide all scene's objects to draw them
func (m Menu) GetLayers() (map[uint8][]objects.Object, *sync.Mutex) {
	return m.layers, m.m
}

// KeyboardEvent provide key down to the scene
func (m *Menu) KeyboardEvent(keyboard *sdl.KeyboardEvent) {
	var err error

	if keyboard.Type == sdl.KEYDOWN {
		if err = m.input.SetNewRune(keyboard.Keysym, m.renderer); err != nil {
			go m.setNotice(err.Error())
		}
	}
}

// SetSwitcher can be call to change scene with index scene and flag closer
func (m *Menu) SetSwitcher(f func(uint8, bool) error) {
	m.switcher = f
}

// Update : called on each frame
func (m *Menu) Update() {
	return
}

// SetCloser : allow quit the application
func (m *Menu) SetCloser(f func()) {
	m.quit = f
}
