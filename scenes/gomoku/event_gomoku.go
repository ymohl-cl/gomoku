package gomoku

import (
	"errors"
	"fmt"

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

	if err, end = g.game.Move(x, y); err != nil {
		// setNotice
		println("crash", err.Error())
		println("x: ", x, "y: ", y)
		return
	}

	if end == true {
		fmt.Println("WINNER YEAH BRAVO ! VOILA")
		g.switcher(conf.SMenu, true)
		return
	}

	durationToPlay := g.game.GetTimeToPlay()

	caps := g.game.GetCaptures()
	for i, _ := range caps {
		xt := caps[i].X
		yt := caps[i].Y
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

	g.game.Playing()
	/*
		 if p must play and p == ia {
		 	setNotice(it's not you which must play)
			return
		}
		get y, x
		if checkMove == false {
			setNotive(you can't do this move)
			return
		}
		applie move
		applie draw
		----
		if p1 play then p2 must play
		else p1 must play
		if futur player == ia {
			timeSet
			ia.Play
			timeStop
			if futur player == p1 {
				futur player = p2
			} else {
				futur player = p1
			}
			applie draw
		}
	*/
}

func (g *Gomoku) quit(values ...interface{}) {
	g.switcher(conf.SMenu, true)
}
