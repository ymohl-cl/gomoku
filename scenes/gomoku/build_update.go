package gomoku

import (
	"strconv"
	"time"

	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/game-builder/objects/block"
	"github.com/ymohl-cl/game-builder/objects/text"
	"github.com/ymohl-cl/gomoku/conf"
	"github.com/ymohl-cl/gomoku/controller"
	"github.com/ymohl-cl/gomoku/database"
)

func (g *Gomoku) createBlockHistory(x, y int32) (*block.Block, error) {
	var b *block.Block
	var err error

	if b, err = block.New(block.Filled); err != nil {
		return nil, err
	}
	// Set style fix and basic
	if err = b.SetVariantStyle(conf.ColorButtonRed, conf.ColorButtonGreen, conf.ColorButtonBlue, conf.ColorButtonOpacity, objects.SFix); err != nil {
		return nil, err
	}
	// Set position
	b.UpdatePosition(x, y)
	// Set size
	b.UpdateSize(conf.GameHistoryWidth, conf.GameHistoryHeight)
	return b, nil
}

func (g *Gomoku) createTextHistory(x, y int32, txt string) (*text.Text, error) {
	var t *text.Text
	var err error

	x += conf.GameHistoryWidth / 2
	y += conf.GameHistoryHeight / 2
	if t, err = text.New(txt, conf.TxtLittle, conf.Font, x, y); err != nil {
		return nil, err
	}
	t.SetVariantStyle(conf.ColorTxtRed, conf.ColorTxtGreen, conf.ColorTxtBlue, conf.ColorTxtOpacity, objects.SFix)
	return t, nil
}

func (g *Gomoku) addHistory(x, y uint8, d time.Duration, p *database.Player) error {
	var posX, posY int32
	var block *block.Block
	var text *text.Text
	var err error

	// get side player to draw history
	if g.data.Current.P1 == p {
		posX = conf.MarginLeft
	} else {
		posX = conf.WindowWidth - conf.MarginRight - conf.GameContentBlockWidth
	}
	posY = conf.MarginTop + conf.MenuHeaderHeight + conf.PaddingBlock

	// get txt to draw it on new history element
	str := "x: " + strconv.FormatUint(uint64(x), 10)
	str += " - y: " + strconv.FormatUint(uint64(y), 10)
	str += " - time: " + controller.TimeToString(d)

	// create block element
	if block, err = g.createBlockHistory(posX, posY); err != nil {
		return err
	}
	if err = block.Init(g.renderer); err != nil {
		return err
	}

	// create txt element
	if text, err = g.createTextHistory(posX, posY, str); err != nil {
		return err
	}
	if err = text.Init(g.renderer); err != nil {
		return err
	}

	g.m.Lock()
	defer g.m.Unlock()
	// move next history
	g.moveHistory()

	// add new elements on layer history
	g.layers[layerHistoryBlock] = append(g.layers[layerHistoryBlock], block)
	g.layers[layerHistoryText] = append(g.layers[layerHistoryText], text)
	return nil
}

func (g *Gomoku) moveHistory() {
	var mod int
	idDelete := -1

	if len(g.layers[layerHistoryBlock])%2 == 0 {
		mod = 0
	} else {
		mod = 1
	}

	for id, v := range g.layers[layerHistoryBlock] {
		if (id)%2 == mod {
			x, y := v.GetPosition()
			// calcul new position
			y += conf.GameHistoryHeight + conf.GameHistoryPadding

			// check if new position stay on the space history
			if y+conf.GameHistoryHeight >= conf.MarginTop+conf.MenuHeaderHeight+conf.PaddingBlock+conf.MenuContentLargeBlockHeight {
				// delete object outside space history
				idDelete = id
			} else {
				// update block position
				v.UpdatePosition(x, y)
				// update text position
				x, y = g.layers[layerHistoryText][id].GetPosition()
				y += conf.GameHistoryHeight + conf.GameHistoryPadding
				g.layers[layerHistoryText][id].UpdatePosition(x, y)
			}
		}
	}

	if idDelete >= 0 {
		b := g.layers[layerHistoryBlock][idDelete]
		t := g.layers[layerHistoryText][idDelete]
		go func() {
			if err := b.Close(); err != nil {
				panic(err)
			}
			if err := t.Close(); err != nil {
				panic(err)
			}
		}()
		g.layers[layerHistoryBlock] = append(g.layers[layerHistoryBlock][:idDelete], g.layers[layerHistoryBlock][idDelete+1:]...)
		g.layers[layerHistoryText] = append(g.layers[layerHistoryText][:idDelete], g.layers[layerHistoryText][idDelete+1:]...)
	}
}
