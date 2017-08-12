package builder

import (
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
)

func CreateBlockToDefaultButton(x, y, w, h int32) (*block.Block, error) {
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

func CreateTxtToButton(x, y int32, txt string) (*text.Text, error) {
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
