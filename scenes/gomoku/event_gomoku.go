package gomoku

import (
	"errors"
	"fmt"
	"time"

	"github.com/ymohl-cl/game-builder/objects"
	"github.com/ymohl-cl/gomoku/conf"
)

func (g *Gomoku) initMove() {
	var x, y uint8 = 9, 9
	var err error

	player := g.game.GetCurrentPlayer()
	g.game.AppliesMove(x, y)
	g.game.SwitchPlayer()
	durationToPlay := g.game.GetTimeToPlay()

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

	g.game.Playing()
}

func (g *Gomoku) selectToken(values ...interface{}) {
	var x, y uint8
	var err error
	var end bool

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

	if end, err = g.game.Move(x, y); err != nil {
		// setNotice
		g.setNotice("You can't make this move")
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

	if end == true {
		g.setNotice("WINNER YEAH BRAVO ! VOILA")
		time.Sleep(5 * time.Second)
		g.switcher(conf.SMenu, true)
		return
	}

	g.game.Playing()

	if g.game.GetCurrentPlayer().Name == "AI" {
		go g.DrawFilter()
		go func() {
			//g.game.Bot.PlayOpposing(y, x)
			c := make(chan uint8)
			go g.game.Bot.Play(g.game.GetBoard(), g.data.Current, c)
			yi, xi := <-c, <-c
			fmt.Println("AI play on x ", xi, " - y: ", yi)
			go g.selectToken(yi, xi)
		}()
	} else {
		go g.HideFilter()
	}

}

func (g *Gomoku) quit(values ...interface{}) {
	g.data.SaveSession()
	go func() {
		if err := g.switcher(conf.SMenu, true); err != nil {
			panic(err)
		}
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
	time.Sleep(3 * time.Second)
	if g.notice.GetIDSDL() == idSDL {
		g.notice.Close()
	}
}
