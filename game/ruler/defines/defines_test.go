package defines

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
)

func TestGetOtherPlayer(t *testing.T) {
	// test: 0
	if GetOtherPlayer(Player1) != Player2 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	if GetOtherPlayer(Player2) != Player1 {
		t.Error(t.Name() + " > test: 1")
	}
}

func TestIsOnTheBoard(t *testing.T) {
	var y, x int8

	// test: 0 > top left of the board
	y = 0
	x = 0
	if !IsOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > top right of the board
	y = 0
	x = 18
	if !IsOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > bot right of the board
	y = 18
	x = 18
	if !IsOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3 > bot left of the board
	y = 18
	x = 0
	if !IsOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 3")
	}

	// test: 4 > out of board
	y = 19
	x = 0
	if IsOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 4")
	}
}

func TestIsEmpty(t *testing.T) {
	var b *[19][19]uint8

	// test: 0
	b = boards.GetEmpty()
	if !IsEmpty(b, 9, 9) {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	b = boards.GetStartP1_1()
	if IsEmpty(b, 9, 9) {
		t.Error(t.Name() + " > test: 1")
	}
}
