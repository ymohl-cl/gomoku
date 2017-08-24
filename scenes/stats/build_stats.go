package stats

import (
	"github.com/ymohl-cl/game-builder/audio"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/button"
	"github.com/ymohl-cl/game-builder/objects/image"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/scenes/builder"
)

func (s *Stats) addMusic() error {
	var err error

	s.music, err = audio.New(conf.StatsMusic, 1)
	if err != nil {
		return err
	}

	if err = s.music.Init(s.renderer); err != nil {
		return err
	}
	return nil
}

func (s *Stats) addBackground() error {
	var i *image.Image
	var err error

	i = image.New(conf.StatsBackground, conf.OriginX, conf.OriginY, conf.WindowWidth, conf.WindowHeight)
	if err = i.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerBackground] = append(s.layers[layerBackground], i)
	return nil
}

func (s *Stats) addStructuresPage() error {
	var b *block.Block
	var err error
	var x, y int32

	// Create blockheader
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, conf.MarginTop)
	b.UpdateSize(conf.WindowWidth, conf.StatsHeaderHeight)
	if err = b.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerStructure] = append(s.layers[layerStructure], b)

	// Create blockLeft
	y = conf.MarginTop + conf.StatsHeaderHeight + conf.PaddingBlock
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.MarginLeft, y)
	b.UpdateSize(conf.StatsContentBlockWidth, conf.StatsContentLargeBlockHeight)
	if err = b.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerStructure] = append(s.layers[layerStructure], b)

	// Create blockTopRight
	x = conf.WindowWidth - conf.MarginRight - conf.StatsContentBlockWidth
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(x, y)
	b.UpdateSize(conf.StatsContentBlockWidth, conf.StatsContentMediumBlockHeight)
	if err = b.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerStructure] = append(s.layers[layerStructure], b)

	// Create blockBottomRight
	y = conf.MarginTop + conf.StatsHeaderHeight + conf.PaddingBlock + conf.StatsContentMediumBlockHeight + conf.PaddingBlock
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(x, y)
	b.UpdateSize(conf.StatsContentBlockWidth, conf.StatsContentMediumBlockHeight)
	if err = b.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerStructure] = append(s.layers[layerStructure], b)

	// Create blockFooter
	y = conf.WindowHeight - conf.MarginBot - conf.StatsFooterHeight
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, y)
	b.UpdateSize(conf.WindowWidth, conf.StatsHeaderHeight)
	if err = b.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerStructure] = append(s.layers[layerStructure], b)

	return nil
}

func (s *Stats) getButtonReturn() (*button.Button, error) {
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
	if t, err = builder.CreateTxtToButton(x+conf.ButtonWidth/2, y+conf.ButtonHeight/2, "Return !"); err != nil {
		return nil, err
	}

	b = button.New(bl, nil, t)
	b.SetAction(s.Return)

	return b, nil
}

func (s *Stats) addButtons() error {
	var err error
	var b *button.Button

	if b, err = s.getButtonReturn(); err != nil {
		return err
	}
	if err = b.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerButton] = append(s.layers[layerButton], b)

	return nil
}



func (s *Stats) addText() error {
	var t *text.Text
	var err error
	var x, y int32

	x = conf.WindowWidth / 2
	y = conf.MarginTop + (conf.StatsHeaderHeight / 2)
	// add title
	if t, err = text.New("GOMOKU - STATS", conf.TxtLarge, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	if err = t.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerText] = append(s.layers[layerText], t)

	// add signature
	y = conf.WindowHeight - (conf.MarginBot / 2)
	signature := "Gomuku is present to you by Anis (agadhgad) and MrPiou (ymohl-cl), Enjoy !"
	if t, err = text.New(signature, conf.TxtLittle, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	if err = t.Init(s.renderer); err != nil {
		return err
	}
	s.layers[layerText] = append(s.layers[layerText], t)

	return nil
}
