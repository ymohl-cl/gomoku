package menu

import (
	"errors"
	"fmt"
	"time"

	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/database"
)

/*
** Endpoint action from objects click
 */

// LoadGame : start a new party
func (m *Menu) LoadGame(values ...interface{}) {
	fmt.Println("Load Game")
}

// DeletePlayer : _
func (m *Menu) DeletePlayer(values ...interface{}) {
	var p *database.Player
	var err error
	var ok bool
	var id int

	if len(values) == 1 {
		p, ok = values[0].(*database.Player)
		if !ok {
			panic(errorInterface)
		}
	} else {
		panic(errorValuesEmpty)
	}

	if id, err = m.data.DeletePlayer(p); err != nil {
		go m.setNotice(err.Error())
		return
	}
	go func() {
		if err = m.removeUIPlayer(id); err != nil {
			panic(err)
		}
		m.updateVS()
	}()
}

// DrawStat : on user selected
func (m *Menu) DrawStat(values ...interface{}) {
	fmt.Println("Draw stat")
}

// SelectPlayer : to the futur game
func (m *Menu) SelectPlayer(values ...interface{}) {
	var p *database.Player
	var err error
	var ok bool

	if len(values) == 1 {
		p, ok = values[0].(*database.Player)
		if !ok {
			panic(errorInterface)
		}
	} else {
		panic(errorValuesEmpty)
	}
	if err = m.data.UpdateCurrent(p); err != nil {
		go m.setNotice(err.Error())
		return
	}
	m.updateVS()
}

// NewPlayer : on the database
func (m *Menu) NewPlayer(values ...interface{}) {
	var name string
	var nbPlayer int
	var err error

	// check number players
	nbPlayer = len(m.data.Players)
	if nbPlayer >= playerMax {
		go m.setNotice(noticeMaxPlayer)
		return
	}
	// get name player
	name = m.input.GetTxt()
	if len(name) == 0 {
		go m.setNotice(noticeNameEmpty)
		return
	}
	// check if player already exist
	for _, p := range m.data.Players {
		if p.Name == name {
			go m.setNotice(noticeNameExist)
			return
		}
	}
	// create player on protocol data
	p := database.CreatePlayer(name)
	// add player on protocol data
	m.data.AddPlayer(p)

	// Reset input
	m.input.Reset(m.renderer)

	// add player on layers players
	go func() {
		if err = m.addUIPlayer(nbPlayer, p); err != nil {
			panic(err)
		}
	}()
}

// Play start the game
func (m *Menu) Play(values ...interface{}) {
	var err error

	go func() {
		if err = m.switcher(conf.SGame, true); err != nil {
			panic(err)
		}
	}()
}

// ResetName : reset the input value
func (m *Menu) ResetName(values ...interface{}) {
	m.input.Reset(m.renderer)
}

// DefaultPlayer : init the defaults player to the game
func (m *Menu) DefaultPlayer(values ...interface{}) {
	var err error

	if err = m.data.DefaultPlayers(); err != nil {
		panic(err.Error)
	}
	m.updateVS()
}

// setNotice allow draw information to the player
func (m *Menu) setNotice(str string) {
	idSDL := m.notice.NewIDSDL()
	if m.notice.IsInit() == true {
		m.notice.Close()
	}
	m.notice.UpdateText(str, m.renderer)
	if err := m.notice.Init(m.renderer); err != nil {
		panic(errors.New(objects.ErrorRenderer))
	}
	time.Sleep(3 * time.Second)
	if m.notice.GetIDSDL() == idSDL {
		m.notice.Close()
	}
}
