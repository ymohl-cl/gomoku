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
	var mess string

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

	if end, mess, err = g.game.Move(x, y); err != nil {
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
		g.setNotice("WINNER YEAH BRAVO ! VOILA | " + mess)

		time.Sleep(5 * time.Second)
		g.switcher(conf.SMenu, true)
		return
	}

	g.game.Playing()

	if g.game.GetCurrentPlayer().Name == "AI" {
		go g.DrawFilter()
		go func() {
			g.game.Bot.PlayOpposing(y, x)
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

/*function AlphaBetaWithMemory(n : node_type; alpha , beta , d : integer) : integer;

  if retrieve(n) == OK then // Transposition table lookup
      if n.lowerbound >= beta then return n.lowerbound;
      if n.upperbound <= alpha then return n.upperbound;
      alpha := max(alpha, n.lowerbound);
      beta := min(beta, n.upperbound);
  if d == 0 then g := evaluate(n); // leaf node
  else if n == MAXNODE then
      g := -INFINITY; a := alpha; // save original alpha value
      c := firstchild(n);
      while (g < beta) and (c != NOCHILD) do
          g := max(g, AlphaBetaWithMemory(c, a, beta, d - 1));
          a := max(a, g);
          c := nextbrother(c);
  else // n is a MINNODE
      g := +INFINITY; b := beta; // save original beta value
      c := firstchild(n);
      while (g > alpha) and (c != NOCHILD) do
          g := min(g, AlphaBetaWithMemory(c, alpha, b, d - 1));
          b := min(b, g);
          c := nextbrother(c);
  // Traditional transposition table storing of bounds
  // Fail low result implies an upper bound
  if g <= alpha then n.upperbound := g; store n.upperbound;
  // Found an accurate minimax value - will not occur if called with zero window
  if g >  alpha and g < beta then
      n.lowerbound := g; n.upperbound := g; store n.lowerbound, n.upperbound;
  // Fail high result implies a lower bound
  if g >= beta then n.lowerbound := g; store n.lowerbound;
  return g;
*/

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
