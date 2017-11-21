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

func TestUpdateScoreAlignment(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var scoreCurrent int8
	var scoreOpponent int8

	b = boards.GetEmpty()
	state = New(b, ruler.Player2)

	flag := false
	scoreCurrent = 12
	scoreOpponent = 10

	saveCurrent := scoreCurrent
	saveOpponent := scoreOpponent

	// score current
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 7)
	if saveCurrent != scoreCurrent {
		t.Error(t.Name() + " > test: 1")
	} else if saveOpponent != scoreOpponent {
		t.Error(t.Name() + " > test: 2")
	} else if flag == false {
		t.Error(t.Name() + " > test: 3")
	}

	// score opponent
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 15)
	if saveCurrent != scoreCurrent {
		t.Error(t.Name() + " > test: 4")
	} else if scoreOpponent != 15 {
		t.Error(t.Name()+" > test: 5 > ret: ", scoreOpponent)
	} else if flag == true {
		t.Error(t.Name() + " > test: 6")
	}

	// score current
	saveOpponent = scoreOpponent
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 0)
	if saveCurrent != scoreCurrent {
		t.Error(t.Name() + " > test: 7")
	} else if saveOpponent != scoreOpponent {
		t.Error(t.Name() + " > test: 8")
	} else if flag == false {
		t.Error(t.Name() + " > test: 9")
	}

	// score opponent
	flag = !flag
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 22)
	if scoreCurrent != 22 {
		t.Error(t.Name()+" > test: 10 > ret: ", scoreCurrent)
	} else if saveOpponent != scoreOpponent {
		t.Error(t.Name() + " > test: 11")
	} else if flag == false {
		t.Error(t.Name() + " > test: 12")
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

	ret := state.evalAlignment(n3, 0)
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

	ret = state.evalAlignment(n3, 0)
	// test: 3 > check value of eval
	if ret != 13 {
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

	ret = state.evalAlignment(n3, 0)
	// test: 5 > check value of eval
	if ret != 13 {
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

	ret = state.evalAlignment(n3, 0)
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

	ret = state.evalAlignment(n3, 0)
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

	ret = state.evalAlignment(n3, 0)
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

	// State board
	// - . . . . . . . . x o . . . . . . . . .

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

	// p1 capture: 1
	// p1 alignment: 1 | max: 3 free
	// p2 capture: 0
	// p2 alignment: 0 | max: 0

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval
	if ret != -100 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . o . . . . . . . .
	// - . . . . . . . . x . . x x . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
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

	// p1 capture: 2
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 1
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 3 > check value of eval to advantage P1 Caps and align
	if ret != -67 {
		t.Error(t.Name()+" > test: 3 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . o . . . . .
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

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

	// p1 capture: 3
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 1
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 5 > check value of eval to advantage P1 Caps and align
	if ret != -79 {
		t.Error(t.Name()+" > test: 5 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . . . . . . . x o . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

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

	// p1 capture: 3
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 2
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 7 > check value of eval
	if ret != -83 {
		t.Error(t.Name()+" > test: 7 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . . . o . . o o x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

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

	// p1 capture: 4
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 2
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 9 > check value of eval
	if ret != -75 {
		t.Error(t.Name()+" > test: 9 > resultat: ", ret)
	}

	// State board
	//                     |
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

	// p1 capture: 4
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 4
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 11 > check value of eval to score capture P2 == -25 and align == -24
	if ret != 120 {
		t.Error(t.Name()+" > test: 11 > resultat: ", ret)
	}

	// State board
	//                     |
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
	//   . . . . . o . . . . . . . . . . . . .
	//   . . . . . . o . . . . . . . . . . . .
	// - . . . . . . . o x x x x x o x . . . .
	//   . . . . . . . . . . o . o . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
}

func TestEval_situation1(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	b = boards.GetStartP2_1()

	// move ai
	state = New(b, ruler.Player1)
	n0 := state.newNode(9, 7)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(8, 9)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(11, 11)
	state.updateData(n3, n1)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval to align P1
	if ret != -44 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}
}
