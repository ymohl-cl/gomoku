package menu

import (
	"errors"

	"github.com/ymohl-cl/game-builder/objects/button"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/database"
)

func (m *Menu) addUIPlayer(nb int, p *database.Player) error {
	var err error
	var x, y int32
	var b1, b2, b3, b4 *button.Button

	x = conf.MarginLeft
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock
	y += (conf.MenuElementPlayerHeight + conf.MenuElementPadding) * int32(nb)

	if b1, err = m.addButtonPlayer(x, y, p); err != nil {
		return err
	}
	if err = b1.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b1)

	if b2, err = m.addButtonDeletePlayer(x, y, p); err != nil {
		return err
	}
	if err = b2.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b2)

	if b3, err = m.addButtonStat(x, y, p); err != nil {
		return err
	}
	if err = b3.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b3)

	if b4, err = m.addLoadGame(x, y, p); err != nil {
		return err
	}
	if err = b4.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b4)

	return nil
}

func (m Menu) closeUIPlayer(idx int) error {
	var err error
	var size int

	size = len(m.layers[layerPlayers])
	if size <= idx {
		return errors.New("id object not found")
	}
	button := m.layers[layerPlayers][idx]
	if err = button.Close(); err != nil {
		return err
	}
	return nil
}

func (m *Menu) removeUIPlayer(idData int) error {
	var err error
	var id int

	id = idData * buttonByPlayer
	if err = m.closeUIPlayer(id); err != nil {
		return err
	}
	if err = m.closeUIPlayer(id + 1); err != nil {
		return err
	}
	if err = m.closeUIPlayer(id + 2); err != nil {
		return err
	}
	if err = m.closeUIPlayer(id + 3); err != nil {
		return err
	}

	// Update position next elements
	for _, b := range m.layers[layerPlayers][id+4:] {
		x, y := b.GetPosition()
		y -= (conf.MenuElementPlayerHeight + conf.MenuElementPadding)
		b.UpdatePosition(x, y)
	}

	m.layers[layerPlayers] = append(m.layers[layerPlayers][:id], m.layers[layerPlayers][id+4:]...)
	return nil
}

func (m *Menu) updateVS() {
	//	var err error
	//	p1 := m.data.Current.P1
	//	p2 := m.data.Current.P2

	/*	if p1 == nil || p2 == nil {
			panic(errors.New("Players is nil"))
		}
	*/
	/*	if m.vs.IsInit() {
		if err = m.vs.Close(); err != nil {
			panic(err)
		}
	}*/
	/*
		if err = m.vs.UpdateText(p1.Name+" VS "+p2.Name, m.renderer); err != nil {
			panic(err)
		}
	*/
}
