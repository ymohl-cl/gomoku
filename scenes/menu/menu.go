package menu

import (
	"errors"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/ymohl-cl/game-builder/audio"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/input"
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
	layerVS
	layerInput
	layerPlayers

	nbrLayers = 8

	// Configuration menu
	playerMax      = 10
	buttonByPlayer = 6

	// notice message
	noticeMaxPlayer = "You can't save more players"
	noticeNameEmpty = "New player's name is empty"
	noticeNameExist = "New player's name already exist"

	// erreur message
	errorValuesEmpty = "Function polymorphic call without valid argument"
	errorInterface   = "Error interface type insertion"
	errorPlayer      = "Player not found"
	errorData        = "Data not corrupted"
)

// Menu is a scene
type Menu struct {
	/* infos scene */
	initialized bool
	closer      chan (uint8)
	switcher    func(uint8, bool) error
	quit        func()

	/* objects by layers */
	m      *sync.Mutex
	layers map[uint8][]objects.Object

	/* specific objects */
	input   *input.Input
	notice  *text.Text
	music   *audio.Audio
	player1 *text.Text
	player2 *text.Text

	/* sdl ressources */
	renderer *sdl.Renderer
	/* data game */
	data *database.Data
}

/*
** constructor
 */

// New provide a new object
func New(d *database.Data, r *sdl.Renderer) (*Menu, error) {
	if r == nil {
		return nil, errors.New(objects.ErrorRenderer)
	}
	if d == nil {
		return nil, errors.New(objects.ErrorData)
	}

	m := Menu{renderer: r, data: d}
	m.m = new(sync.Mutex)
	return &m, nil
}
