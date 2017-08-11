package gomoku

import (
	"errors"
	"fmt"
)

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

	fmt.Println("y: ", y)
	fmt.Println("x: ", x)
	if err = g.game.Move(x, y); err != nil {
		// setNotice
		panic(err)
	}

	go func() {
		if err = g.ChangeToken(x, y, g.game.GetCurrentPlayer()); err != nil {
			panic(err)
		}
	}()
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
