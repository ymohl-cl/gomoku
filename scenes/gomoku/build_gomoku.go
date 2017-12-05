package gomoku

import (
	"errors"

	"github.com/ymohl-cl/game-builder/audio"
	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/button"
	"github.com/ymohl-cl/game-builder/objects/image"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/scenes/builder"
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

	// create block left
	y = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(conf.MarginLeft, y)
	b.UpdateSize(conf.GameContentBlockWidth, conf.MenuContentLargeBlockHeight)
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerStructure] = append(g.layers[layerStructure], b)

	// create block right
	x = conf.WindowWidth - conf.MarginRight - conf.GameContentBlockWidth
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SFix)
	b.UpdatePosition(x, y)
	b.UpdateSize(conf.GameContentBlockWidth, conf.MenuContentLargeBlockHeight)
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

	nameP2 := database.Bot
	nameP1 := g.data.Current.P1.Name
	if g.data.Current.P2 != nil {
		nameP2 = g.data.Current.P2.Name
	}

	x = conf.WindowWidth / 2
	y = conf.MarginTop + (conf.MenuHeaderHeight / 2)
	if t, err = text.New(nameP1+" VS "+nameP2, conf.TxtLarge, conf.Font, x, y); err != nil {
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
	return nil
}

// RestoreToken restore the initial token
func (g *Gomoku) RestoreToken(x, y uint8, p *database.Player) error {
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

	img = image.New(conf.GameMarkTokenBlack, posX, posY, 22, 22)
	img.SetVariantStyle(conf.GameMarkTokenBlack, conf.GameTokenBlack, conf.GameMarkTokenBlack)
	img.SetAction(g.selectToken, y, x)
	if err = img.Init(g.renderer); err != nil {
		panic(err)
	}

	tmp := g.layers[layerToken][iY*19+iX]
	g.m.Lock()
	g.layers[layerToken][iY*19+iX] = img
	g.m.Unlock()
	if err = tmp.Close(); err != nil {
		panic(err)
	}

	return nil
}

// ChangeToken update token status and color selected by player
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

func (g *Gomoku) addButtons() error {
	var err error
	var b *button.Button

	if b, err = g.getButtonQuit(); err != nil {
		return err
	}
	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerButton] = append(g.layers[layerButton], b)
	return nil
}

func (g *Gomoku) getButtonQuit() (*button.Button, error) {
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
	b.SetAction(g.quit)

	return b, nil
}

func (g *Gomoku) addBlockTime() error {
	var b *block.Block
	var t *text.Text
	var x, y int32
	var err error

	y = conf.WindowHeight - conf.MarginBot - (conf.MenuFooterHeight / 2) - (conf.ButtonHeight / 2)
	x = conf.MarginLeft

	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	// Set style fix and basic
	if err = b.SetVariantStyle(conf.ColorOverButtonRed, conf.ColorOverButtonGreen, conf.ColorOverButtonBlue, conf.ColorOverButtonOpacity, objects.SFix); err != nil {
		return err
	}
	// Set position
	b.UpdatePosition(x, y)
	// Set size
	b.UpdateSize(conf.ButtonWidth, conf.ButtonHeight)

	if err = b.Init(g.renderer); err != nil {
		return err
	}
	g.layers[layerBlockTime] = append(g.layers[layerBlockTime], b)

	y += conf.ButtonHeight / 2
	x += conf.ButtonWidth / 2
	if t, err = text.New("00:00:00", conf.TxtMedium, conf.Font, x, y); err != nil {
		return err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	if err = t.Init(g.renderer); err != nil {
		return err
	}

	g.layers[layerText] = append(g.layers[layerText], t)
	g.timer = t
	return nil
}

func (g *Gomoku) addNotice() error {
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
	g.notice = t
	g.layers[layerNotice] = append(g.layers[layerNotice], g.notice)
	return nil
}

func (g *Gomoku) addModal() error {
	var b *block.Block
	var err error

	// create blockheader
	if b, err = block.New(block.Filled); err != nil {
		return err
	}
	b.SetVariantStyle(conf.ColorBlockRed, conf.ColorBlockGreen, conf.ColorBlockBlue, conf.ColorBlockOpacity, objects.SBasic, objects.SClick, objects.SOver)
	b.UpdatePosition(conf.OriginX, conf.OriginY)
	b.UpdateSize(conf.WindowWidth, conf.WindowHeight-conf.MenuHeaderHeight)
	b.SetAction(g.ModalFunction)
	g.modal = b
	g.layers[layerModal] = append(g.layers[layerModal], b)
	return nil
}
