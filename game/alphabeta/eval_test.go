package alphabeta

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	"github.com/ymohl-cl/gomoku/game/ruler"
)

func TestEvalCapture(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	b = boards.GetStartP1()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 11)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(8, 10)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 12)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	ret := state.evalCapture(n3, n3.rule.GetPlayer(), ruler.GetOtherPlayer(n3.rule.GetPlayer()))
	// test: 1 > check value of eval
	if ret != 6 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . o . . . . . . . .
	// - . . . . . . . . x . . x x . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(9, 13)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(7, 10)
	state.updateData(n3, n2)

	// test: 2 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 2")
	}

	ret = state.evalCapture(n3, n3.rule.GetPlayer(), ruler.GetOtherPlayer(n3.rule.GetPlayer()))
	// test: 3 > check value of eval
	if ret != 4 {
		t.Error(t.Name()+" > test: 3 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . o . . . . .
	//   . . . . . . . . . . x . . . . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(7, 13)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 13)
	state.updateData(n3, n2)

	// test: 4 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 4")
	}

	ret = state.evalCapture(n3, n3.rule.GetPlayer(), ruler.GetOtherPlayer(n3.rule.GetPlayer()))
	// test: 5 > check value of eval
	if ret != 16 {
		t.Error(t.Name()+" > test: 5 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . . . . . . . x o . x . . . . .

	// move ai
	n0 = state.newNode(10, 12)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 9)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 8)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(9, 13)
	state.updateData(n3, n2)

	// test: 6 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 6")
	}

	ret = state.evalCapture(n3, n3.rule.GetPlayer(), ruler.GetOtherPlayer(n3.rule.GetPlayer()))
	// test: 7 > check value of eval
	if ret != -11 {
		t.Error(t.Name()+" > test: 7 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . . . o . . o o x . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 7)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 6)
	state.updateData(n3, n2)

	// test: 8 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 8")
	}

	ret = state.evalCapture(n3, n3.rule.GetPlayer(), ruler.GetOtherPlayer(n3.rule.GetPlayer()))
	// test: 9 > check value of eval
	if ret != 25 {
		t.Error(t.Name()+" > test: 9 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . x o o . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(11, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 5)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 4)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 11)
	state.updateData(n3, n2)

	// test: 10 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 10")
	}

	ret = state.evalCapture(n3, n3.rule.GetPlayer(), ruler.GetOtherPlayer(n3.rule.GetPlayer()))
	// test: 11 > check value of eval
	if ret != -25 {
		t.Error(t.Name()+" > test: 11 > resultat: ", ret)
	}

}

func TestEvalAlignment(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	b = boards.GetStartP1()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 11)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(8, 10)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 12)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	ret := state.evalAlignment(n3)
	// test: 1 > check value of eval
	if ret != 24 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . o . . . . . . . .
	// - . . . . . . . . x . . x x . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(9, 13)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(7, 10)
	state.updateData(n3, n2)

	// test: 2 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 2")
	}

	ret = state.evalAlignment(n3)
	// test: 3 > check value of eval
	if ret != 11 {
		t.Error(t.Name()+" > test: 3 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . o . . . . .
	//   . . . . . . . . . . x . . . . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(7, 13)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 13)
	state.updateData(n3, n2)

	// test: 4 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 4")
	}

	ret = state.evalAlignment(n3)
	// test: 5 > check value of eval
	if ret != 11 {
		t.Error(t.Name()+" > test: 5 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . . . . . . . x o . x . . . . .

	// move ai
	n0 = state.newNode(10, 12)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 9)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 8)
	state.updateData(n2, n1)
	// move player
	//	boards.Print(state.board)
	n3 = state.newNode(9, 13)
	state.updateData(n3, n2)

	// test: 6 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 6")
	}

	ret = state.evalAlignment(n3)
	// test: 7 > check value of eval
	if ret != 24 {
		t.Error(t.Name()+" > test: 7 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . . . o . . o o x . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 7)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 6)
	state.updateData(n3, n2)

	// test: 8 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 8")
	}

	ret = state.evalAlignment(n3)
	// test: 9 > check value of eval
	if ret != 0 {
		t.Error(t.Name()+" > test: 9 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . x o o . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(11, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 5)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 4)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 11)
	state.updateData(n3, n2)

	// test: 10 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 10")
	}

	ret = state.evalAlignment(n3)
	// test: 11 > check value of eval
	if ret != -24 {
		t.Error(t.Name()+" > test: 11 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . o . . o o . x x . . . . . . .
	//   . . . . . . . . . . . . . o . . . . .
}

func TestEval_noWin(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	b = boards.GetStartP1()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 11)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(8, 10)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 12)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval to align P1
	if ret != -119 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . o . . . . . . . .
	// - . . . . . . . . x . . x x . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(9, 13)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(7, 10)
	state.updateData(n3, n2)

	// test: 2 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 2")
	}

	ret = state.eval(n3, 0)
	// test: 3 > check value of eval to advantage P1 Caps and align
	if ret != -65 {
		t.Error(t.Name()+" > test: 3 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . o . . . . .
	//   . . . . . . . . . . x . . . . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(7, 13)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 13)
	state.updateData(n3, n2)

	// test: 4 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 4")
	}

	ret = state.eval(n3, 0)
	// test: 5 > check value of eval to advantage P1 Caps and align
	if ret != -77 {
		t.Error(t.Name()+" > test: 5 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . . . . . . . x o . x . . . . .

	// move ai
	n0 = state.newNode(10, 12)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 9)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 8)
	state.updateData(n2, n1)
	// move player
	//	boards.Print(state.board)
	n3 = state.newNode(9, 13)
	state.updateData(n3, n2)

	// test: 6 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 6")
	}

	ret = state.eval(n3, 0)
	// test: 7 > check value of eval to score align P1 == 24
	if ret != -119 {
		t.Error(t.Name()+" > test: 7 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . . . o . . o o x . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 7)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 6)
	state.updateData(n3, n2)

	// test: 8 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 8")
	}

	ret = state.eval(n3, 0)
	// test: 9 > check value of eval to score capture P1 == 24
	if ret != -121 {
		t.Error(t.Name()+" > test: 9 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . x o o . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(11, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 5)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 4)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 11)
	state.updateData(n3, n2)

	// test: 10 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 10")
	}

	ret = state.eval(n3, 0)
	// test: 11 > check value of eval to score capture P2 == -25 and align == -24
	if ret != 121 {
		t.Error(t.Name()+" > test: 11 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . o . . o o . x x . . . . . . .
	//   . . . . . . . . . . . . . o . . . . .
}

func TestEval_win(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	b = boards.GetWinNoCapturableP2()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(8, 7)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 15)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(7, 6)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 11)
	state.updateData(n3, n1)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > check rule
	if !n3.rule.Win {
		t.Error(t.Name() + " > test: 1")
	}

	ret := state.eval(n3, 0)
	// test: 2 > check value of eval to align P1
	if ret != -123 {
		t.Error(t.Name()+" > test: 2 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . o . . . . . . . . . . . .
	// - . . . . . . . o x x x x x o x . . . .
	//   . . . . . . . . . . o . o . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
}
