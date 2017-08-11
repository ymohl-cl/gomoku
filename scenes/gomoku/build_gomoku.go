package gomoku

import (
	"errors"
	"fmt"

	"github.com/ymohl-cl/game-builder/audio"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/image"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/database"
)

const (
	errorNotInit = "you can't do this action because gomoku scene is not initialized"
)

func (g *Gomoku) addMusic() error {
	var err error

	g.music, err = audio.New(conf.GameMusic, 3)
	if err != nil {
		return err
	}

	if err = g.music.Init(g.renderer); err != nil {
		return err
	}
	return nil
}

func (g *Gomoku) addBackground() error {
	var i *image.Image
	var err error

	i = image.New(conf.MenuBackground, conf.OriginX, conf.OriginY, conf.WindowWidth, conf.WindowHeight)
	if err = i.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerBackground] = append(g.layers[layerBackground], i)
	return nil
}

func (g *Gomoku) addStructures() error {
	var b *block.Block
	var err error
	var x, y int32

	// create blockheader
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, conf.MarginTop)
	b.UpdateSize(conf.WindowWidth, conf.MenuHeaderHeight)
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerStructure] = append(g.layers[layerStructure], b)

	// create block top left
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.MarginLeft, y)
	b.UpdateSize(conf.GameContentBlockWidth, conf.MenuContentMediumBlockHeight)
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerStructure] = append(g.layers[layerStructure], b)

	// create block top right
	x = conf.WindowWidth - conf.MarginRight - conf.GameContentBlockWidth
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(x, y)
	b.UpdateSize(conf.GameContentBlockWidth, conf.MenuContentMediumBlockHeight)
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerStructure] = append(g.layers[layerStructure], b)

	// create block botton right
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock + conf.MenuContentMediumBlockHeight + conf.PaddingBlock
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(x, y)
	b.UpdateSize(conf.GameContentBlockWidth, conf.MenuContentMediumBlockHeight)
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerStructure] = append(g.layers[layerStructure], b)

	// create block bottom left
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.MarginLeft, y)
	b.UpdateSize(conf.GameContentBlockWidth, conf.MenuContentMediumBlockHeight)
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerStructure] = append(g.layers[layerStructure], b)

	// create block footer
	y = conf.WindowHeight - conf.MarginBot - conf.MenuFooterHeight
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.OriginX, y)
	b.UpdateSize(conf.WindowWidth, conf.MenuHeaderHeight)
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerStructure] = append(g.layers[layerStructure], b)

	return nil
}

func (g *Gomoku) addText() error {
	var t *text.Text
	var err error
	var y, x int32

	p1 := g.data.Current.P1
	p2 := g.data.Current.P2
	x = conf.WindowWidth / 2
	y = conf.MarginTop + (conf.MenuHeaderHeight / 2)
	if t, err = text.New(p1.Name+" VS "+p2.Name, conf.TxtLarge, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	t.SetVariantUnderStyle(conf.ColorUnderTxtRed, conf.ColorUnderTxtGreen, conf.ColorUnderTxtBlue, conf.ColorUnderTxtOpacity, objects.SFix)
	t.SetUnderPosition(x-conf.TxtUnderPadding, y-conf.TxtUnderPadding)
	if err = t.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerText] = append(g.layers[layerText], t)

	return nil
}

func (g *Gomoku) addTokens() error {
	var img *image.Image
	var err error
	var x, y int32

	x = (conf.WindowWidth / 2) - (250 + 5)
	y = (conf.WindowHeight / 2) - (250 + 5)

	for yi := 0; yi < 19; yi++ {
		for xi := 0; xi < 19; xi++ {
			img = image.New(conf.GameMarkTokenBlack, x, y, 22, 22)
			img.SetVariantStyle(conf.GameMarkTokenBlack, conf.GameTokenBlack, conf.GameMarkTokenBlack)
			img.SetAction(g.selectToken, uint8(yi), uint8(xi))
			if err = img.Init(g.renderer); err != nil {
				panic(err)
			}
			g.layers[layerToken] = append(g.layers[layerToken], img)
			x += 22 + 5
		}
		x = (conf.WindowWidth / 2) - (250 + 5)
		y += 22 + 5
	}
	fmt.Println("Len graphic board: ", len(g.layers[layerToken]))
	return nil
}

func (g *Gomoku) ChangeToken(x, y uint8, p *database.Player) error {
	var posX, posY int32
	var img *image.Image
	var err error
	var iX, iY int32

	iX = int32(x)
	iY = int32(y)
	if !g.IsInit() {
		return errors.New(errorNotInit)
	}
	posX, posY = g.layers[layerToken][iY*19+iX].GetPosition()
	if p == g.data.Current.P1 {
		img = image.New(conf.GameTokenWhite, posX, posY, 22, 22)
	} else {
		img = image.New(conf.GameTokenGold, posX, posY, 22, 22)
	}

	tmp := g.layers[layerToken][iY*19+iX]
	if err = img.Init(g.renderer); err != nil {
		return err
	}
	g.m.Lock()
	g.layers[layerToken][iY*19+iX] = img
	g.m.Unlock()
	if err = tmp.Close(); err != nil {
		return err
	}

	return nil
}
