package alphabeta

import (
	"fmt"
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

func TestGetTotalCapture(t *testing.T) {
	b := boards.GetEmpty()
	s := New(b, rdef.Player1)

	// test: 0 > just init, so capture == 0 to p1 and p2
	if s.getTotalCapture(rdef.Player1) != 0 {
		t.Error(t.Name() + " > test: 0-0")
	}
	if s.getTotalCapture(rdef.Player2) != 0 {
		t.Error(t.Name() + " > test: 0-1")
	}

	s.addTotalCapture(rdef.Player1, 1)
	s.addTotalCapture(rdef.Player2, 3)

	// test: 1 > p1 == 1 and p2 == 3
	if s.getTotalCapture(rdef.Player1) != 1 {
		t.Error(t.Name() + " > test: 1-0")
	}
	if s.getTotalCapture(rdef.Player2) != 3 {
		t.Error(t.Name() + " > test: 1-1")
	}
}

func TestAddTotalCapture(t *testing.T) {
	b := boards.GetEmpty()
	s := New(b, rdef.Player1)

	s.addTotalCapture(rdef.Player1, 1)
	s.addTotalCapture(rdef.Player2, 2)

	// test: 0 > p1 == 1 and p2 == 2
	if s.getTotalCapture(rdef.Player1) != 1 {
		t.Error(t.Name() + " > test: 0-0")
	}
	if s.getTotalCapture(rdef.Player2) != 2 {
		t.Error(t.Name() + " > test: 0-1")
	}

	s.addTotalCapture(rdef.Player2, 2)

	// test: 1 > p1 == 1 and p2 == 4
	if s.getTotalCapture(rdef.Player1) != 1 {
		t.Error(t.Name() + " > test: 1-0")
	}
	if s.getTotalCapture(rdef.Player2) != 4 {
		t.Error(t.Name() + " > test: 1-1")
	}

	s.addTotalCapture(rdef.Player1, 1)

	// test: 2 > p1 == 2 and p2 == 4
	if s.getTotalCapture(rdef.Player1) != 2 {
		t.Error(t.Name() + " > test: 2-0")
	}
	if s.getTotalCapture(rdef.Player2) != 4 {
		t.Error(t.Name() + " > test: 2-1")
	}
}

func TestSubTotalCapture(t *testing.T) {
	b := boards.GetEmpty()
	s := New(b, rdef.Player1)

	s.addTotalCapture(rdef.Player1, 1)
	s.addTotalCapture(rdef.Player2, 2)

	// test: 0 > p1 == 1 and p2 == 2
	if s.getTotalCapture(rdef.Player1) != 1 {
		t.Error(t.Name() + " > test: 0-0")
	}
	if s.getTotalCapture(rdef.Player2) != 2 {
		t.Error(t.Name() + " > test: 0-1")
	}

	s.subTotalCapture(rdef.Player2, 1)

	// test: 1 > p1 == 1 and p2 == 1
	if s.getTotalCapture(rdef.Player1) != 1 {
		t.Error(t.Name() + " > test: 1-0")
	}
	if s.getTotalCapture(rdef.Player2) != 1 {
		t.Error(t.Name() + " > test: 1-1")
	}

	s.subTotalCapture(rdef.Player1, 1)

	// test: 2 > p1 == 0 and p2 == 1
	if s.getTotalCapture(rdef.Player1) != 0 {
		t.Error(t.Name() + " > test: 2-0")
	}
	if s.getTotalCapture(rdef.Player2) != 1 {
		t.Error(t.Name() + " > test: 2-1")
	}
}

func TestNewNode(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var n *Node

	b = boards.GetStartP1_1()
	state = New(b, rdef.Player2)

	// test: 0 > simple invalid move
	n = state.newNode(9, 0)
	if n != nil {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > simple allowed move
	n = state.newNode(9, 7)
	if n == nil {
		t.Error(t.Name() + " > test: 1")
	}
	// test: 2 > check the id player
	if state.currentPlayer != n.rule.GetPlayer() && n.rule.GetPlayer() != rdef.Player2 {
		t.Error(t.Name() + " > test: 2")
	}
}

func ExampleUpdateData() {
	var b *[19][19]uint8
	var state *State
	var n *Node

	b = boards.GetStartP2_1()
	state = New(b, rdef.Player1)
	n = state.newNode(9, 7)
	if n == nil {
		fmt.Println("nil node")
	}

	state.Print()
	state.updateData(n, n)
	state.Print()

	// Output:
	//maxDepth:  4  - player:  1
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . o o x . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//p1 nb capture:  0
	//p2 nb capture:  0
	//maxDepth:  4  - player:  2
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . x . . x . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//p1 nb capture:  1
	//p2 nb capture:  0
}

func ExampleRestoreData() {
	var b *[19][19]uint8
	var state *State
	var n *Node

	b = boards.GetStartP2_1()
	state = New(b, rdef.Player1)
	n = state.newNode(9, 7)
	if n == nil {
		fmt.Println("nil node")
	}

	state.Print()
	state.updateData(n, n)
	state.restoreData(n, n)
	state.Print()

	// Output:
	//maxDepth:  4  - player:  1
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . o o x . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//p1 nb capture:  0
	//p2 nb capture:  0
	//maxDepth:  4  - player:  1
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . o o x . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//  . . . . . . . . . . . . . . . . . . .
	//p1 nb capture:  0
	//p2 nb capture:  0
}

/*func ExampleAlphabetaNegaScout_1() {
	var b *[19][19]uint8
	var state *State

	b = boards.GetStartP2_1()
	state = New(b, ruler.Player1)
	state.maxDepth = 1

	state.addTotalCapture(ruler.Player1, uint8(0))
	state.addTotalCapture(ruler.Player2, uint8(0))

	ret := state.alphabetaNegaScout(math.MinInt8+1, math.MaxInt8, state.maxDepth, nil)
	fmt.Println("ret: ", ret)

	tmp := int8(math.MinInt8)
	for n := state.lst; n != nil; n = n.next {
		if tmp < n.weight {
			tmp = n.weight
			state.save = n
		}
		y, x := n.rule.GetPosition()
		fmt.Println("Node weight: ", n.weight, " y et x: ", y, " - ", x)
	}
	y, x := state.save.rule.GetPosition()
	fmt.Println("select: weight: ", state.save.weight, " y et x: ", y, " - ", x)

	// Output:
	// ?
}*/
