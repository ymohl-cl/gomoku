package stats

import (
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/audio"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/database"
)

const (
	// order layers of scene
	layerBackground = iota
	layerStructure
	layerButton
	layerNotice
	layerText
)

// Stats is a scene
type Stats struct {
	/* infos scene */
	initialized bool
	closer      chan (uint8)
	switcher    func(uint8, bool) error

	/* objects by layers */
	m      *sync.Mutex
	layers map[uint8][]objects.Object

	/* specific objects */
	notice *text.Text
	music  *audio.Audio

	/* sdl ressources */
	renderer *sdl.Renderer
	/* data game */
	data *database.Data
}

/*
** constructor
 */

// New provide a new object
func New(d *database.Data, r *sdl.Renderer) (*Stats, error) {
	if r == nil {
		return nil, errors.New(objects.ErrorRenderer)
	}
	if d == nil {
		return nil, errors.New(objects.ErrorData)
	}

	s := Stats{renderer: r, data: d}
	s.layers = make(map[uint8][]objects.Object)
	s.m = new(sync.Mutex)
	return &s, nil
}
