package gomoku

import (
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/audio"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/gomoku/database"
)

const (
	// order layers of scene
	layerBackground = iota
	layerStructure
	layerText
	layerToken
	layerButton
	layerModal
)

// Gomoku is a scene which used when build other scene
type Gomoku struct {
	/* infos scene */
	initialized bool
	switcher    func(uint8, bool) error

	/* objects by layers */
	m      *sync.Mutex
	layers map[uint8][]objects.Object
	music  *audio.Audio

	/* specific objects */

	/* sdl ressources */
	renderer *sdl.Renderer
	/* data game */
	data *database.Data
}

/*
** constructor
 */

// New provide a new object
func New(d *database.Data, r *sdl.Renderer) (*Gomoku, error) {
	if r == nil {
		return nil, errors.New(objects.ErrorRenderer)
	}
	if d == nil {
		return nil, errors.New(objects.ErrorData)
	}

	g := Gomoku{renderer: r, data: d}
	g.layers = make(map[uint8][]objects.Object)
	g.m = new(sync.Mutex)
	return &g, nil
}
