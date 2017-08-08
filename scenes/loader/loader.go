package loader

import (
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/gomoku/database"
)

const (
	// order layers of scene
	layerBackground = iota
	layerStructure
	layerText
	layerLoadingBar
	nbrLayers = 4
)

// Load is a scene which used when build other scene
type Load struct {
	/* infos scene */
	initialized bool
	closer      chan (bool)
	refresh     bool

	/* objects by layers */
	m             *sync.Mutex
	layers        map[uint8][]objects.Object
	lastLoadBlock *block.Block

	/* specific objects */

	/* sdl ressources */
	renderer *sdl.Renderer
}

/*
** constructor
 */

// New provide a new object
func New(d *database.Data, r *sdl.Renderer) (*Load, error) {
	if r == nil {
		return nil, errors.New(objects.ErrorRenderer)
	}

	l := Load{renderer: r}
	l.layers = make(map[uint8][]objects.Object)
	l.closer = make(chan (bool))
	l.m = new(sync.Mutex)
	return &l, nil
}
