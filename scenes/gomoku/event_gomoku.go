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

	player := g.game.GetCurrentPlayer()

	if ok, message := g.game.Move(x, y); !ok {
		// setNotice
		g.setNotice("You can't make this move: " + message)
		return
	}

	durationToPlay := g.game.GetTimeToPlay()

	for _, cap := range g.game.GetCaptures() {
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

	if ok, message := g.game.IsWin(); ok {
		g.game.End = true
		g.setNotice("WINNER YEAH BRAVO ! " + message)
		// save the game
		//time.Sleep(15 * time.Second)

		go func() {
			if err := g.switcher(conf.SMenu, true); err != nil {
				panic(err)
			}
			g.data.SaveSession()
		}()
		return
	}

	g.game.Playing()

	if g.game.IsBot(g.game.GetCurrentPlayer()) {
		go g.DrawFilter()
		go func() {
			g.game.Bot.PlayOpponnent(int8(y), int8(x))
			c := make(chan uint8)
			go g.game.Bot.Play(g.game.GetBoard(), g.data.Current, c)
			yi, xi := <-c, <-c
			go g.selectToken(yi, xi)
		}()
	} else {
		go g.HideFilter()
	}

}

func (g *Gomoku) quit(values ...interface{}) {
	if !g.game.End {
		go func() {
			if err := g.switcher(conf.SMenu, true); err != nil {
				panic(err)
			}
			g.data.SaveSession()
		}()
	}
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
	time.Sleep(5 * time.Second)
	if g.notice.GetIDSDL() == idSDL {
		g.notice.Close()
	}
}
