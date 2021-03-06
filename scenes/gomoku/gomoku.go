package gomoku

import (
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/go-ui/audio"
	"github.com/ymohl-cl/go-ui/objects"
	"github.com/ymohl-cl/go-ui/objects/block"
	"github.com/ymohl-cl/go-ui/objects/text"
	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game"
)

const (
	// order layers of scene
	layerBackground = iota
	layerStructure
	layerBlockTime
	layerToken
	layerButton
	layerText
	layerHistoryBlock
	layerHistoryText
	layerNotice
	layerModal
)

// Gomoku is a scene which used when build other scene
type Gomoku struct {
	/* infos scene */
	initialized bool
	switcher    func(indexScene string) error
	closeScene  func(indexScene string) error
	game        *game.Game

	/* objects by layers */
	m      *sync.Mutex
	layers map[uint8][]objects.Object
	music  *audio.Audio
	notice *text.Text
	timer  *text.Text
	modal  *block.Block
	// input

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
	g.m = new(sync.Mutex)
	return &g, nil
}
