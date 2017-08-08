package loader

import (
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/image"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
)

func (l *Load) addBackground() error {
	var i *image.Image
	var err error

	i = image.New(conf.MenuBackground, conf.OriginX, conf.OriginY, conf.WindowWidth, conf.WindowHeight)
	if err = i.Init(l.renderer); err != nil {
		return err
	}
	l.layers[layerBackground] = append(l.layers[layerBackground], i)
	return nil
}

func (l *Load) addBlockLoading() error {
	var b *block.Block
	var err error

	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorButtonRed, conf.ColorButtonGreen, conf.ColorButtonBlue, conf.LoadBlockOpacity, objects.SFix)
	b.UpdateSize(conf.LoadBlockWidth, conf.LoadBlockHeight)
	b.UpdatePosition(conf.OriginX+conf.LoadBlockWidth, conf.WindowHeight-(conf.MarginBot*2)-conf.LoadFooterHeight)

	if err = b.Init(l.renderer); err != nil {
		return err
	}
	l.layers[layerLoadingBar] = append(l.layers[layerLoadingBar], b)
	l.lastLoadBlock = b
	return nil
}

func (l *Load) addStructure() error {
	var b *block.Block
	var err error
	var y int32

	// Create blockheader
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, conf.MarginTop)
	b.UpdateSize(conf.WindowWidth, conf.MenuHeaderHeight)
	if err = b.Init(l.renderer); err != nil {
		return nil
	}
	l.layers[layerStructure] = append(l.layers[layerStructure], b)

	// Create blockFooter
	y = conf.WindowHeight - conf.MarginBot - conf.MenuFooterHeight
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, y)
	b.UpdateSize(conf.WindowWidth, conf.MenuHeaderHeight)
	if err = b.Init(l.renderer); err != nil {
		return nil
	}
	l.layers[layerStructure] = append(l.layers[layerStructure], b)

	return nil
}

func (l *Load) addTxt() error {
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
	if err = t.Init(l.renderer); err != nil {
		return err
	}
	l.layers[layerText] = append(l.layers[layerText], t)

	// add signature
	y = conf.WindowHeight / 2
	if t, err = text.New("LOADING ...", conf.TxtLarge, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	if err = t.Init(l.renderer); err != nil {
		return err
	}
	l.layers[layerText] = append(l.layers[layerText], t)

	return nil
}
