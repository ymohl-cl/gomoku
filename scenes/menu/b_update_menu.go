package menu

import (
	"errors"

	"github.com/ymohl-cl/go-ui/objects/button"
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
	m.m.Lock()
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b1)
	m.m.Unlock()

	if err = m.addButtonChoicePlayer(x, y, p); err != nil {
		return err
	}

	if b2, err = m.addButtonDeletePlayer(x, y, p); err != nil {
		return err
	}
	if err = b2.Init(m.renderer); err != nil {
		return err
	}
	m.m.Lock()
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b2)
	m.m.Unlock()

	if b3, err = m.addButtonStat(x, y, p); err != nil {
		return err
	}
	if err = b3.Init(m.renderer); err != nil {
		return err
	}
	m.m.Lock()
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b3)
	m.m.Unlock()

	if b4, err = m.addLoadGame(x, y, p); err != nil {
		return err
	}
	if err = b4.Init(m.renderer); err != nil {
		return err
	}
	m.m.Lock()
	m.layers[layerPlayers] = append(m.layers[layerPlayers], b4)
	m.m.Unlock()

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
	if err = m.closeUIPlayer(id + 4); err != nil {
		return err
	}
	if err = m.closeUIPlayer(id + 5); err != nil {
		return err
	}
	// Update position next elements
	for _, b := range m.layers[layerPlayers][id+buttonByPlayer:] {
		x, y := b.GetPosition()
		y -= (conf.MenuElementPlayerHeight + conf.MenuElementPadding)
		b.UpdatePosition(x, y)
	}

	m.m.Lock()
	m.layers[layerPlayers] = append(m.layers[layerPlayers][:id], m.layers[layerPlayers][id+buttonByPlayer:]...)
	m.m.Unlock()
	return nil
}

func (m *Menu) updateVS() {
	var err error
	nameP1 := "Unknow"
	nameP2 := "Unknow"

	p1 := m.data.Current.P1
	p2 := m.data.Current.P2
	if p1 != nil {
		nameP1 = p1.Name
	}
	if p2 != nil {
		nameP2 = p2.Name
	}

	if m.player1.IsInit() {
		if err = m.player1.Close(); err != nil {
			panic(err)
		}
	}
	if err = m.player1.UpdateText(nameP1, m.renderer); err != nil {
		panic(err)
	}

	if m.player2.IsInit() {
		if err = m.player2.Close(); err != nil {
			panic(err)
		}
	}
	if err = m.player2.UpdateText(nameP2, m.renderer); err != nil {
		panic(err)
	}
}
