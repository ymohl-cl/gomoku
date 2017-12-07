package gomoku

import (
	"errors"
	"time"

	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/gomoku/conf"
)

func (g *Gomoku) initMove() {
	var y, x uint8 = 9, 9

	g.game.Start = true
	go g.selectToken(y, x)
}

func (g *Gomoku) selectToken(values ...interface{}) {
	var x, y uint8
	var err error

	if g.game.End {
		return
	}

	if len(values) != 2 {
		panic(errors.New("More two values specified on select token"))
	}
	for id, v := range values {
		switch res := v.(type) {
		case uint8:
			if id == 0 {
				y = res
			} else if id == 1 {
				x = res
			}
		default:
			panic(errors.New("type to indexes cells of board is wrong. must be uint8"))
		}
	}

	PrevWin := true
	if g.game.RulerWin != nil {
		for _, s := range g.game.RulerWin.CapturableWin {
			if y == uint8(s.Y) && x == uint8(s.X) {
				PrevWin = false
				g.game.RulerWin = nil
				break
			}
		}
	} else {
		PrevWin = false
	}

	player := g.game.GetCurrentPlayer()

	if ok, message := g.game.Move(x, y); !ok {
		// setNotice
		g.setNotice("You can't make this move: " + message)
		return
	}

	durationToPlay := g.game.GetTimeToPlay()

	// nb is saved to GetWinCapturabled
	nb := 0
	for _, cap := range g.game.GetCaptures() {
		nb++
		xt := cap.X
		yt := cap.Y
		go func() {
			if err = g.RestoreToken(uint8(xt), uint8(yt), player); err != nil {
				panic(err)
			}
		}()
	}

	go func() {
		if err = g.ChangeToken(x, y, player); err != nil {
			panic(err)
		}
	}()

	go func() {
		if err = g.addHistory(x, y, durationToPlay, player); err != nil {
			panic(err)
		}
	}()

	if r := g.game.GetRules(); r != nil && len(r.CapturableWin) != 0 {
		if g.game.RulerWin == nil {
			g.game.RulerWin = r
		}
	}

	if PrevWin == true {
		g.game.End = true
		g.setNotice("PREV WINNER " + g.game.GetCurrentPlayer().Name + " " + g.game.RulerWin.Info)
		g.game.RulerWin = nil
		return
	} else if ok, message := g.game.IsWin(); ok {
		g.game.End = true
		g.setNotice("CURRENT WINNER " + player.Name + " " + message)
		return
	}

	g.game.Playing()

	if g.game.IsBot(g.game.GetCurrentPlayer()) {
		go g.DrawFilter()
		go func() {
			g.game.Bot.PlayOpponnent(int8(y), int8(x))
			c := make(chan uint8)
			if g.game.RulerWin != nil {
				go g.game.Bot.Play(g.game.GetBoard(), g.data.Current, c, g.game.RulerWin.CapturableWin)
			} else {
				go g.game.Bot.Play(g.game.GetBoard(), g.data.Current, c, nil)
			}
			yi, xi := <-c, <-c
			go g.selectToken(yi, xi)
		}()
	} else if g.game.Bot != nil {
		go g.HideFilter()
	}

}

func (g *Gomoku) quit(values ...interface{}) {
	go func() {
		if err := g.switcher(conf.SMenu, true); err != nil {
			panic(err)
		}
		g.data.SaveSession()
	}()
}

// setNotice allow draw informations to the player
func (g *Gomoku) setNotice(str string) {
	idSDL := g.notice.NewIDSDL()
	if g.notice.IsInit() == true {
		g.notice.Close()
	}
	g.notice.UpdateText(str, g.renderer)
	if err := g.notice.Init(g.renderer); err != nil {
		panic(errors.New(objects.ErrorRenderer))
	}
	time.Sleep(8 * time.Second)
	if g.notice.GetIDSDL() == idSDL {
		g.notice.Close()
	}
}

// ModalFunction mock action
func (g *Gomoku) ModalFunction(values ...interface{}) {
	return
}
