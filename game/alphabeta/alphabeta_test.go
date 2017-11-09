package alphabeta

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	"github.com/ymohl-cl/gomoku/game/ruler"
)

func TestNewNode(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var n *Node

	b = boards.GetStartP1()
	state = New(b, ruler.Player2)

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
	if state.currentPlayer != n.rule.GetPlayer() && n.rule.GetPlayer() != ruler.Player2 {
		t.Error(t.Name() + " > test: 2")
	}

}
