package menu

import (
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/button"
	"github.com/ymohl-cl/game-builder/objects/image"
	"github.com/ymohl-cl/game-builder/objects/input"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/database"
)

func (m *Menu) createBlockToDefaultButton(x, y, w, h int32) (*block.Block, error) {
	var b *block.Block
	var err error

	if b, err = block.New(block.Filled); err != nil {
		return nil, err
	}
	// Set style fix and basic
	if err = b.SetVariantStyle(conf.ColorButtonRed, conf.ColorButtonGreen, conf.ColorButtonBlue, conf.ColorButtonOpacity, objects.SFix, objects.SBasic); err != nil {
		return nil, err
	}
	// Set style over
	if err = b.SetVariantStyle(conf.ColorOverButtonRed, conf.ColorOverButtonGreen, conf.ColorOverButtonBlue, conf.ColorOverButtonOpacity, objects.SOver); err != nil {
		return nil, err
	}
	// set style click
	if err = b.SetVariantStyle(conf.ColorClickButtonRed, conf.ColorClickButtonGreen, conf.ColorClickButtonBlue, conf.ColorClickButtonOpacity, objects.SClick); err != nil {
		return nil, err
	}
	// Set position
	b.UpdatePosition(x, y)
	// Set size
	b.UpdateSize(w, h)
	return b, nil
}

func (m *Menu) createTxtToButton(x, y int32, txt string) (*text.Text, error) {
	var t *text.Text
	var err error

	if t, err = text.New(txt, conf.TxtMedium, conf.Font, x, y); err != nil {
		return nil, err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	return t, nil
}

func (m *Menu) addButtonPlayer(x, y int32, p *database.Player) (*button.Button, error) {
	var t *text.Text
	var bl *block.Block
	var b *button.Button
	var err error

	if bl, err = m.createBlockToDefaultButton(x, y, conf.MenuElementPlayerWidth, conf.MenuElementPlayerHeight); err != nil {
		return nil, err
	}
	w, h := bl.GetSize()
	if t, err = m.createTxtToButton(x+(w/2), y+(h/2), p.Name); err != nil {
		return nil, err
	}
	b = button.New(bl, nil, t)

	return b, nil
}

func (m *Menu) addButtonChoicePlayer(x, y int32, p *database.Player) error {
	var t *text.Text
	var b *button.Button
	var bl *block.Block
	var err error
	var str string

	str = "J1"
	x += conf.MenuElementPlayerWidth + conf.MenuElementPadding

	for i := 0; i < 2; i++ {
		if i == 1 {
			str = "J2"
			x += conf.MenuElementPadding + conf.MenuIconWidth
		}
		if bl, err = m.createBlockToDefaultButton(x, y, conf.MenuIconWidth, conf.MenuIconWidth); err != nil {
			return err
		}

		w, h := bl.GetSize()
		if t, err = m.createTxtToButton(x+(w/2), y+(h/2), str); err != nil {
			return err
		}
		b = button.New(bl, nil, t)
		b.SetAction(m.SelectPlayer, p, i)

		if err = b.Init(m.renderer); err != nil {
			return err
		}
		m.layers[layerPlayers] = append(m.layers[layerPlayers], b)
	}

	return nil
}

func (m *Menu) addButtonDeletePlayer(x, y int32, p *database.Player) (*button.Button, error) {
	var i *image.Image
	var b *button.Button

	x += conf.MenuElementPlayerWidth + (conf.MenuElementPadding * 3) + (conf.MenuIconWidth * 2)
	i = image.New(conf.MenuIconDelete, x, y, conf.MenuIconWidth, conf.MenuElementPlayerHeight)
	i.SetVariantStyle(conf.MenuIconDelete, conf.MenuIconOverDelete, conf.MenuIconOverDelete)
	b = button.New(nil, i, nil)
	b.SetAction(m.DeletePlayer, p)

	return b, nil
}

func (m *Menu) addLoadGame(x, y int32, p *database.Player) (*button.Button, error) {
	var i *image.Image
	var b *button.Button

	x += conf.MenuElementPlayerWidth + (conf.MenuElementPadding * 4) + (conf.MenuIconWidth * 3)
	i = image.New(conf.MenuIconLoad, x, y, conf.MenuIconWidth, conf.MenuElementPlayerHeight)
	i.SetVariantStyle(conf.MenuIconLoad, conf.MenuIconOverLoad, conf.MenuIconOverLoad)
	b = button.New(nil, i, nil)
	b.SetAction(m.LoadGame, p)

	return b, nil
}

func (m *Menu) addButtonStat(x, y int32, p *database.Player) (*button.Button, error) {
	var i *image.Image
	var b *button.Button

	x += conf.MenuElementPlayerWidth + (conf.MenuElementPadding * 5) + (conf.MenuIconWidth * 4)

	i = image.New(conf.MenuIconTrophy, x, y, conf.MenuIconWidth, conf.MenuElementPlayerHeight)
	i.SetVariantStyle(conf.MenuIconTrophy, conf.MenuIconOverTrophy, conf.MenuIconOverTrophy)
	b = button.New(nil, i, nil)
	b.SetAction(m.DrawStat, p)

	return b, nil
}

func (m *Menu) createBlockToInput(x, y int32) (*block.Block, error) {
	var b *block.Block
	var err error

	if b, err = block.New(block.Filled); err != nil {
		return nil, err
	}
	// Set style fix and basic
	if err = b.SetVariantStyle(conf.ColorInputRed, conf.ColorInputGreen, conf.ColorInputBlue, conf.ColorInputOpacity, objects.SFix, objects.SBasic); err != nil {
		return nil, err
	}
	// Set style over
	if err = b.SetVariantStyle(conf.ColorOverInputRed, conf.ColorOverInputGreen, conf.ColorOverInputBlue, conf.ColorOverInputOpacity, objects.SOver); err != nil {
		return nil, err
	}
	// set style click
	if err = b.SetVariantStyle(conf.ColorClickInputRed, conf.ColorClickInputGreen, conf.ColorClickInputBlue, conf.ColorClickInputOpacity, objects.SClick); err != nil {
		return nil, err
	}
	// Set position
	b.UpdatePosition(x, y)
	// Set size
	b.UpdateSize((conf.ButtonWidth*2)+conf.PaddingBlock, conf.ButtonHeight)
	return b, nil
}

func (m *Menu) createInput() (*input.Input, error) {
	var x, y int32
	var w, h int32
	var b *block.Block
	var i *input.Input
	var err error

	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight/2 - (conf.PaddingBlock/2 + conf.ButtonHeight)
	interval := int32((conf.MenuContentBlockWidth - (conf.ButtonWidth*2 + conf.PaddingBlock)) / 2)
	x = conf.WindowWidth - conf.MarginRight - conf.MenuContentBlockWidth + interval

	if b, err = m.createBlockToInput(x, y); err != nil {
		return nil, err
	}

	if i, err = input.New(conf.TxtMedium, conf.Font, b); err != nil {
		return nil, err
	}
	i.Txt.SetVariantStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	w = (conf.ButtonWidth * 2) + conf.PaddingBlock
	h = (conf.ButtonHeight)
	i.Txt.UpdatePosition(x+(w/2), y+(h/2))

	return i, nil
}

func (m *Menu) getButtonDefaultPlayers() (*button.Button, error) {
	var x, y int32
	var t *text.Text
	var bl *block.Block
	var b *button.Button
	var err error

	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight + conf.MenuContentMediumBlockHeight/2 + conf.PaddingBlock + conf.PaddingBlock/2
	interval := int32((conf.MenuContentBlockWidth - (conf.ButtonWidth*2 + conf.PaddingBlock)) / 2)
	x = conf.WindowWidth - conf.MarginRight - interval - conf.ButtonWidth

	// create block
	if bl, err = m.createBlockToDefaultButton(x, y, conf.ButtonWidth, conf.ButtonHeight); err != nil {
		return nil, err
	}

	// create txt
	x += conf.ButtonWidth / 2
	y += conf.ButtonHeight / 2
	if t, err = m.createTxtToButton(x, y, "J1 VS IA"); err != nil {
		return nil, err
	}

	// create button
	b = button.New(bl, nil, t)
	b.SetAction(m.Play)

	return b, nil
}

func (m *Menu) getButtonPlay() (*button.Button, error) {
	var x, y int32
	var t *text.Text
	var bl *block.Block
	var b *button.Button
	var err error

	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight + conf.MenuContentMediumBlockHeight/2 + conf.PaddingBlock + conf.PaddingBlock/2
	interval := int32((conf.MenuContentBlockWidth - (conf.ButtonWidth*2 + conf.PaddingBlock)) / 2)
	x = conf.WindowWidth - conf.MarginRight - conf.MenuContentBlockWidth + interval

	// create block
	if bl, err = m.createBlockToDefaultButton(x, y, conf.ButtonWidth, conf.ButtonHeight); err != nil {
		return nil, err
	}

	// create txt
	if t, err = m.createTxtToButton(x+conf.ButtonWidth/2, y+conf.ButtonHeight/2, "J1 VS J2"); err != nil {
		return nil, err
	}

	b = button.New(bl, nil, t)
	b.SetAction(m.Play)

	return b, nil
}

func (m *Menu) getButtonResetName() (*button.Button, error) {
	var x, y int32
	var t *text.Text
	var bl *block.Block
	var b *button.Button
	var err error

	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight/2 + conf.PaddingBlock/2
	interval := int32((conf.MenuContentBlockWidth - (conf.ButtonWidth*2 + conf.PaddingBlock)) / 2)
	x = conf.WindowWidth - conf.MarginRight - interval - conf.ButtonWidth

	// create block
	if bl, err = m.createBlockToDefaultButton(x, y, conf.ButtonWidth, conf.ButtonHeight); err != nil {
		return nil, err
	}

	// create txt
	if t, err = m.createTxtToButton(x+conf.ButtonWidth/2, y+conf.ButtonHeight/2, "RESET NAME"); err != nil {
		return nil, err
	}

	b = button.New(bl, nil, t)
	b.SetAction(m.ResetName)
	return b, nil
}

func (m *Menu) getButtonNewPlayer() (*button.Button, error) {
	var x, y int32
	var t *text.Text
	var bl *block.Block
	var b *button.Button
	var err error

	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight/2 + conf.PaddingBlock/2
	interval := int32((conf.MenuContentBlockWidth - (conf.ButtonWidth*2 + conf.PaddingBlock)) / 2)
	x = conf.WindowWidth - conf.MarginRight - conf.MenuContentBlockWidth + interval

	// create block
	if bl, err = m.createBlockToDefaultButton(x, y, conf.ButtonWidth, conf.ButtonHeight); err != nil {
		return nil, err
	}

	// create txt
	if t, err = m.createTxtToButton(x+conf.ButtonWidth/2, y+conf.ButtonHeight/2, "NEW PLAYER"); err != nil {
		return nil, err
	}

	b = button.New(bl, nil, t)
	b.SetAction(m.NewPlayer)

	return b, nil
}
