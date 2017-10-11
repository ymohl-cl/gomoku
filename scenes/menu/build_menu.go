package menu

import (
	"errors"

	"github.com/ymohl-cl/game-builder/audio"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/button"
	"github.com/ymohl-cl/game-builder/objects/image"
	"github.com/ymohl-cl/game-builder/objects/input"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/scenes/builder"
)

func (m *Menu) addMusic() error {
	var err error

	m.music, err = audio.New(conf.MenuMusic, 1)
	if err != nil {
		return err
	}

	if err = m.music.Init(m.renderer); err != nil {
		return err
	}
	return nil
}

func (m *Menu) addBackground() error {
	var i *image.Image
	var err error

	i = image.New(conf.MenuBackground, conf.OriginX, conf.OriginY, conf.WindowWidth, conf.WindowHeight)
	if err = i.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerBackground] = append(m.layers[layerBackground], i)
	return nil
}

func (m *Menu) addStructuresPage() error {
	var b *block.Block
	var err error
	var x, y int32

	// Create blockheader
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, conf.MarginTop)
	b.UpdateSize(conf.WindowWidth, conf.MenuHeaderHeight)
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerStructure] = append(m.layers[layerStructure], b)

	// Create blockLeft
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.MarginLeft, y)
	b.UpdateSize(conf.MenuContentBlockWidth, conf.MenuContentLargeBlockHeight)
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerStructure] = append(m.layers[layerStructure], b)

	// Create blockTopRight
	x = conf.WindowWidth - conf.MarginRight - conf.MenuContentBlockWidth
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(x, y)
	b.UpdateSize(conf.MenuContentBlockWidth, conf.MenuContentMediumBlockHeight)
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerStructure] = append(m.layers[layerStructure], b)

	// Create blockBottomRight
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight + conf.PaddingBlock
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(x, y)
	b.UpdateSize(conf.MenuContentBlockWidth, conf.MenuContentMediumBlockHeight)
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerStructure] = append(m.layers[layerStructure], b)

	// Create blockFooter
	y = conf.WindowHeight - conf.MarginBot - conf.MenuFooterHeight
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, y)
	b.UpdateSize(conf.WindowWidth, conf.MenuHeaderHeight)
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerStructure] = append(m.layers[layerStructure], b)

	return nil
}

func (m *Menu) addButtons() error {
	var err error
	var b *button.Button

	if b, err = m.getButtonNewPlayer(); err != nil {
		return err
	}
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerButton] = append(m.layers[layerButton], b)

	if b, err = m.getButtonResetName(); err != nil {
		return err
	}
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerButton] = append(m.layers[layerButton], b)

	if b, err = m.getButtonPlay(); err != nil {
		return err
	}
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerButton] = append(m.layers[layerButton], b)

	if b, err = m.getButtonDefaultPlayers(); err != nil {
		return err
	}
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerButton] = append(m.layers[layerButton], b)

	if b, err = m.getButtonQuit(); err != nil {
		return err
	}
	if err = b.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerButton] = append(m.layers[layerButton], b)

	return nil
}

func (m *Menu) getButtonQuit() (*button.Button, error) {
	var x, y int32
	var t *text.Text
	var bl *block.Block
	var b *button.Button
	var err error

	y = conf.WindowHeight - conf.MarginBot - (conf.MenuFooterHeight / 2) - (conf.ButtonHeight / 2)
	x = conf.WindowWidth - conf.MarginRight - conf.ButtonWidth

	// create block
	if bl, err = builder.CreateBlockToDefaultButton(x, y, conf.ButtonWidth, conf.ButtonHeight); err != nil {
		return nil, err
	}

	// create txt
	if t, err = builder.CreateTxtToButton(x+conf.ButtonWidth/2, y+conf.ButtonHeight/2, "QUIT !"); err != nil {
		return nil, err
	}

	b = button.New(bl, nil, t)
	b.SetAction(m.QuitGame)

	return b, nil
}

func (m *Menu) addNotice() error {
	var t *text.Text
	var err error
	var x, y int32

	x = conf.WindowWidth / 2
	y = conf.WindowHeight - conf.MarginBot - (conf.MenuFooterHeight / 2)
	if t, err = text.New("", conf.TxtLittle, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	m.notice = t
	m.layers[layerNotice] = append(m.layers[layerNotice], m.notice)
	return nil
}

func (m *Menu) addText() error {
	var t *text.Text
	var err error
	var x, y int32

	x = conf.WindowWidth / 2
	y = conf.MarginTop + (conf.MenuHeaderHeight / 2)
	// add title
	if t, err = text.New("GOMOKU", conf.TxtLarge, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	if err = t.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerText] = append(m.layers[layerText], t)

	// add signature
	y = conf.WindowHeight - (conf.MarginBot / 2)
	signature := "Gomuku is present to you by Val (vroche) and MrPiou (ymohl-cl), Enjoy !"
	if t, err = text.New(signature, conf.TxtLittle, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	if err = t.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerText] = append(m.layers[layerText], t)

	return nil
}

func (m *Menu) addJS() error {
	var bl *block.Block
	var t *text.Text
	var b *button.Button
	var err error
	var str string
	var x, y int32

	str = "J1"
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight + conf.MenuContentMediumBlockHeight/4 - (conf.PaddingBlock / 2)
	interval := int32((conf.MenuContentBlockWidth - (conf.ButtonWidth*2 + conf.PaddingBlock)) / 2)
	x = conf.WindowWidth - conf.MarginRight - conf.MenuContentBlockWidth + interval

	for i := 0; i < 2; i++ {
		if i == 1 {
			str = "J2"
			y += conf.MenuElementPadding + conf.MenuIconWidth
		}

		if bl, err = block.New(block.Filled); err != nil {
			return err
		}
		// Set style fix and basic
		if err = bl.SetVariantStyle(conf.ColorButtonRed, conf.ColorButtonGreen, conf.ColorButtonBlue, conf.ColorButtonOpacity, objects.SFix); err != nil {
			return err
		}
		// Set position
		bl.UpdatePosition(x, y)
		// Set size
		bl.UpdateSize(conf.MenuIconWidth, conf.MenuIconWidth)

		w, h := bl.GetSize()
		if t, err = m.createTxtToButton(x+(w/2), y+(h/2), str); err != nil {
			return err
		}
		b = button.New(bl, nil, t)

		if err = b.Init(m.renderer); err != nil {
			return err
		}
		m.layers[layerVS] = append(m.layers[layerVS], b)
	}

	return nil
}

func (m *Menu) addVS() error {
	var t *text.Text
	var err error
	var y, x int32
	var str string

	// add name player 1
	//str = m.data.Current.P1.Name
	str = "Unknow"

	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight + conf.MenuContentMediumBlockHeight/4 - (conf.PaddingBlock / 2)
	x = conf.WindowWidth - conf.MarginRight - conf.MenuContentBlockWidth/2

	for i := 0; i < 2; i++ {
		if i == 1 {
			//str = m.data.Current.P2.Name
			y += conf.MenuElementPadding + conf.MenuIconWidth
		}
		if t, err = text.New(str, conf.TxtMedium, conf.Font, x, y+(conf.MenuIconWidth/2)); err != nil {
			return err
		}
		t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
		t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
		t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
		if err = t.Init(m.renderer); err != nil {
			return err
		}

		m.layers[layerVS] = append(m.layers[layerVS], t)
		if i == 1 {
			m.player2 = t
		} else {
			m.player1 = t
		}
	}

	return nil
}

func (m *Menu) addInput() error {
	var err error
	var i *input.Input

	if i, err = m.createInput(); err != nil {
		return err
	}
	if err = i.Init(m.renderer); err != nil {
		return err
	}
	m.layers[layerInput] = append(m.layers[layerInput], i)

	m.input = i
	return nil
}

func (m *Menu) addPlayers() error {
	var err error
	var x, y int32
	var b1, b2, b3, b4 *button.Button

	if m.data == nil {
		return errors.New(errorData)
	}

	x = conf.MarginLeft
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock
	for _, p := range m.data.Players {
		if b1, err = m.addButtonPlayer(x, y, p); err != nil {
			return err
		}
		if err = b1.Init(m.renderer); err != nil {
			return err
		}
		m.layers[layerPlayers] = append(m.layers[layerPlayers], b1)

		if err = m.addButtonChoicePlayer(x, y, p); err != nil {
			return err
		}

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

		y += conf.MenuElementPlayerHeight + conf.MenuElementPadding
	}

	return nil
}
