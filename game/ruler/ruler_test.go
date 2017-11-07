package ruler

import (
	"fmt"
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
)

func TestIsOnTheBoard(t *testing.T) {
	var y, x int8

	// test top left of the board
	y = 0
	x = 0
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " failed with y and x: " + string(y) + " - " + string(x))
	}

	// test top right of the board
	y = 0
	x = 18
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " failed with y and x: " + string(y) + " - " + string(x))
	}

	// test bot right of the board
	y = 18
	x = 18
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " failed with y and x: " + string(y) + " - " + string(x))
	}

	// test bot left of the board
	y = 18
	x = 0
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " failed with y and x: " + string(y) + " - " + string(x))
	}

	// test out of board
	y = 19
	x = 0
	if isOnTheBoard(y, x) {
		t.Error(t.Name() + " failed with y and x: " + string(y) + " - " + string(x))
	}
}

func TestIsAvailablePosition(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetEmptyBoard()
	r = New(TokenP2, 9, 8)
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " failed with y and x: " + string(r.y) + " - " + string(r.x))
	}

	b = boards.GetSimpleBoardP2()
	r = New(TokenP1, 9, 6)
	if !r.isAvailablePosition(b) {
		t.Error(t.Name() + " failed with y and x: " + string(r.y) + " - " + string(r.x))
	}

	r = New(TokenP1, 9, 5)
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " failed with y and x: " + string(r.y) + " - " + string(r.x))
	}

	r = New(TokenP1, 9, 0)
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " failed with y and x: " + string(r.y) + " - " + string(r.x))
	}

	r = New(TokenP1, 9, 9)
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " failed with y and x: " + string(r.y) + " - " + string(r.x))
	}

	r = New(TokenP1, 9, 9)
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " failed with y and x: " + string(r.y) + " - " + string(r.x))
	}

	r = New(TokenP1, -1, 6)
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " failed with y and x: " + string(r.y) + " - " + string(r.x))
	}
}

func ExampleRules_getMaskFromBoard() {
	var b *[19][19]uint8
	var r *Rules
	var mask [11]uint8

	b = boards.GetSimpleBoardP2()

	// simple test_part1
	r = New(TokenP1, 9, 6)
	r.getMaskFromBoard(b, -1, -1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 0, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, 0, -1, &mask)
	fmt.Println(mask)

	// simple test_part2
	r = New(TokenP1, 8, 9)
	r.getMaskFromBoard(b, -1, -1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 0, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, 0, -1, &mask)
	fmt.Println(mask)

	// test out of board
	b = boards.GetBoardFilledRightP2()
	r = New(TokenP1, 10, 18)
	r.getMaskFromBoard(b, -1, -1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 0, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, 0, -1, &mask)
	fmt.Println(mask)

	// Output:
	// [0 0 0 0 0 1 0 0 0 0 0]
	// [0 0 0 0 0 1 0 0 0 0 0]
	// [0 0 0 0 0 1 0 0 0 0 0]
	// [0 0 0 0 0 1 1 1 2 0 0]
	// [0 0 0 0 0 1 0 0 0 0 0]
	// [0 0 0 0 0 1 2 0 0 0 0]
	// [0 0 0 0 0 1 1 0 0 0 0]
	// [0 0 0 0 0 1 2 0 0 0 0]
	// [0 0 0 0 1 1 3 3 3 3 3]
	// [0 0 0 1 2 1 0 0 0 0 0]
	// [3 3 3 3 3 1 0 0 0 0 0]
	// [0 0 0 2 1 1 3 3 3 3 3]
}

func prepareAnalyzeCapture(b *[19][19]uint8, r *Rules) *Rules {
	var mask [11]uint8

	r.getMaskFromBoard(b, -1, -1, &mask)
	r.analyzeCapture(&mask, -1, -1)

	r.getMaskFromBoard(b, -1, 0, &mask)
	r.analyzeCapture(&mask, -1, 0)

	r.getMaskFromBoard(b, -1, 1, &mask)
	r.analyzeCapture(&mask, -1, 1)

	r.getMaskFromBoard(b, 0, -1, &mask)
	r.analyzeCapture(&mask, 0, -1)

	return r
}

func ExampleRules_analyzeCapture() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetCaptureBoardP1_1()
	r = New(TokenP2, 0, 18)
	r = prepareAnalyzeCapture(b, r)
	r.printCaptures()

	b = boards.GetCaptureBoardP2_1()
	r = New(TokenP1, 7, 9)
	r = prepareAnalyzeCapture(b, r)
	r.printCaptures()

	// Unordered output:
	// (y - x):  9  -  9
	// (y - x):  8  -  9
	// (y - x):  9  -  11
	// (y - x):  8  -  10
	// (y - x):  7  -  11
	// (y - x):  7  -  10
	// (y - x):  5  -  11
	// (y - x):  6  -  10
	// (y - x):  5  -  9
	// (y - x):  6  -  9
	// (y - x):  5  -  7
	// (y - x):  6  -  8
	// (y - x):  7  -  7
	// (y - x):  7  -  8
	// (y - x):  8  -  8
	// (y - x):  9  -  7
	// (y - x):  0  -  17
	// (y - x):  0  -  16
}

func prepareAnalyzeAlign(b *[19][19]uint8, r *Rules) *Rules {
	var mask [11]uint8

	r.getMaskFromBoard(b, -1, -1, &mask)
	r.analyzeAlign(&mask, b, -1, -1)

	r.getMaskFromBoard(b, -1, 0, &mask)
	r.analyzeAlign(&mask, b, -1, 0)

	r.getMaskFromBoard(b, -1, 1, &mask)
	r.analyzeAlign(&mask, b, -1, 1)

	r.getMaskFromBoard(b, 0, -1, &mask)
	r.analyzeAlign(&mask, b, 0, -1)

	return r
}

func ExampleRules_analyzeAlign_noAlign() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetBoardNoAlignP2_1()
	r = New(TokenP1, 2, 0)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()

	b = boards.GetBoardNoAlignP2_2()
	r = New(TokenP1, 9, 10)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()

	// Output:
	// nb alignment  0
	// nb alignment  0
}

func ExampleRules_analyzeAlign_moreAlign() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetBoardAlignFlankedP2()
	r = New(TokenP1, 9, 11)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()

	b = boards.GetBoardAlignFreeP2()
	r = New(TokenP1, 10, 10)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()

	b = boards.GetBoardAlignHalfP2()
	r = New(TokenP1, 7, 7)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()

	// Output:
	// nb alignment  1
	// [size: 4, style: flanked]
	// nb alignment  1
	// [size: 3, style: free]
	// nb alignment  1
	// [size: 4, style: half]
}

func ExampleRules_analyzeAlign_tree() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetTreeBoardP1()
	r = New(TokenP2, 7, 7)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()
	r.printThrees()

	b = boards.GetTreeBoardP1()
	r = New(TokenP2, 7, 13)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()
	r.printThrees()

	b = boards.GetTreeBoardP1()
	r = New(TokenP2, 8, 9)
	r = prepareAnalyzeAlign(b, r)
	r.printAlignments()
	r.printThrees()

	// Output:
	// nb alignment  2
	// [size: 3, style: free]
	// [size: 2, style: free]
	// nb threes  1
	// nb alignment  1
	// [size: 3, style: free]
	// nb threes  1
	// nb alignment  2
	// [size: 2, style: free]
	// [size: 3, style: free]
	// nb threes  1
}

func ExampleRules_analyzeAlign_doubleTree() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetTreeBoardP1()
	r = New(TokenP2, 10, 10)
	r = prepareAnalyzeAlign(b, r)
	r.printThrees()
	r.analyzeMoveCondition()
	fmt.Println(r.MovedStr)

	// Output:
	// nb threes  2
	// Double three
}

func ExampleRules_analyzeAlign_winNoCapturable() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetWinNoCapturableBoardP2()
	r = New(TokenP1, 9, 11)
	r = prepareAnalyzeAlign(b, r)
	r.analyzeWinCondition(0)
	fmt.Println(r.MessageWin)

	// Output:
	// Great, win by align: 5
}

func ExampleRules_analyzeAlign_winCapturable() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetWinCapturableBoardP2()
	r = New(TokenP1, 9, 11)
	r = prepareAnalyzeAlign(b, r)
	r.analyzeWinCondition(0)
	fmt.Println(r.MessageWin)

	// Output:
	//
}
