package ruler

import (
	"fmt"
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
)

// prepareAnalyzeCapture mock the check rules to operate on captures
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

// prepareAnalyzeCapture mock the check rules to operate on alignment
func prepareAnalyzeAlign(b *[19][19]uint8, r *Rules) *Rules {
	var mask [11]uint8

	r.getMaskFromBoard(b, -1, -1, &mask)
	r.analyzeAlign(&mask, -1, -1)

	r.getMaskFromBoard(b, -1, 0, &mask)
	r.analyzeAlign(&mask, -1, 0)

	r.getMaskFromBoard(b, -1, 1, &mask)
	r.analyzeAlign(&mask, -1, 1)

	r.getMaskFromBoard(b, 0, -1, &mask)
	r.analyzeAlign(&mask, 0, -1)

	return r
}

func TestIsOnTheBoard(t *testing.T) {
	var y, x int8

	// test: 0 > top left of the board
	y = 0
	x = 0
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > top right of the board
	y = 0
	x = 18
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > bot right of the board
	y = 18
	x = 18
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3 > bot left of the board
	y = 18
	x = 0
	if !isOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 3")
	}

	// test: 4 > out of board
	y = 19
	x = 0
	if isOnTheBoard(y, x) {
		t.Error(t.Name() + " > test: 4")
	}
}

func TestIsAvailablePosition(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetEmpty()
	r = New(Player2, 9, 8)
	// test: 0 > already used 1/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 0")
	}

	b = boards.GetSimpleP2()
	r = New(Player1, 9, 6)
	// test: 1 > allowed position
	if !r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 1")
	}

	r = New(Player1, 9, 5)
	// test: 2 > no neighbour 1/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 2")
	}

	r = New(Player1, 9, 0)
	// test: 3 > neighbour 2/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 3")
	}

	r = New(Player1, 9, 9)
	// test: 4 > already used 2/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 4")
	}

	r = New(Player1, -1, 6)
	// test: 5 > out of board
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 5")
	}
}

func ExampleRules_ApplyMove() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetCaptureP1_1()
	r = New(Player2, 0, 18)
	r = prepareAnalyzeCapture(b, r)

	fmt.Println(b[0])
	r.ApplyMove(b)
	fmt.Println(b[0])
	r.RestoreMove(b)
	fmt.Println(b[0])

	// Output:
	// [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 2 1 1 0]
	// [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 2 0 0 2]
	// [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 2 1 1 0]
}

func ExampleRules_getMaskFrom() {
	var b *[19][19]uint8
	var r *Rules
	var mask [11]uint8

	b = boards.GetSimpleP2()

	// simple test_part1
	r = New(Player1, 9, 6)
	r.getMaskFromBoard(b, -1, -1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 0, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, 0, -1, &mask)
	fmt.Println(mask)

	// simple test_part2
	r = New(Player1, 8, 9)
	r.getMaskFromBoard(b, -1, -1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 0, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, 0, -1, &mask)
	fmt.Println(mask)

	// test out of board
	b = boards.GetFilledRightP2()
	r = New(Player1, 10, 18)
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

func ExampleRules_analyzeCapture() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetCaptureP1_1()
	r = New(Player2, 0, 18)
	r = prepareAnalyzeCapture(b, r)
	r.PrintCaptures()

	b = boards.GetCaptureP2_1()
	r = New(Player1, 7, 9)
	r = prepareAnalyzeCapture(b, r)
	r.PrintCaptures()

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

func ExampleRules_analyzeAlign_noAlign() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetNoAlignP2_1()
	r = New(Player1, 2, 0)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()

	b = boards.GetNoAlignP2_2()
	r = New(Player1, 9, 10)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()

	// Output:
	// nb alignment  0
	// nb alignment  0
}

func ExampleRules_analyzeAlign_moreAlign() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetAlignFlankedP2()
	r = New(Player1, 9, 11)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()

	b = boards.GetAlignFreeP2()
	r = New(Player1, 10, 10)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()

	b = boards.GetAlignHalfP2()
	r = New(Player1, 7, 7)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()

	b = boards.GetDoubleAlignmentP2()
	r = New(Player1, 7, 10)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()

	// Output:
	// nb alignment  1
	// [size: 4, style: flanked]
	// nb alignment  1
	// [size: 3, style: free]
	// nb alignment  1
	// [size: 4, style: half]
	// nb alignment  2
	// [size: 2, style: free]
	// [size: 2, style: free]
}

func ExampleRules_analyzeAlign_tree() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetTreeP1_1()
	r = New(Player2, 7, 7)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()
	r.PrintThrees()

	b = boards.GetTreeP1_1()
	r = New(Player2, 7, 13)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()
	r.PrintThrees()

	b = boards.GetTreeP1_1()
	r = New(Player2, 8, 9)
	r = prepareAnalyzeAlign(b, r)
	r.PrintAlignments()
	r.PrintThrees()

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

	b = boards.GetTreeP1_1()
	r = New(Player2, 10, 10)
	r = prepareAnalyzeAlign(b, r)
	r.PrintThrees()
	r.analyzeMoveCondition()
	fmt.Println(r.Info)

	// Output:
	// nb threes  2
	// this spot compose one double three
}

func ExampleRules_analyzeAlign_winNoCapturable() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetWinNoCapturableP2()
	r = New(Player1, 9, 11)
	r = prepareAnalyzeAlign(b, r)
	r.analyzeWinCondition(b, 0)
	r.PrintAlignments()
	fmt.Println(r.Info)

	// Output:
	// nb alignment  1
	// [size: 5, style: flanked, no capturable]
	// congratulation, winner by alignment
}

func ExampleRules_analyzeAlign_winCapturable() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetWinCapturableP2()
	r = New(Player1, 9, 11)
	r = prepareAnalyzeAlign(b, r)
	r.analyzeWinCondition(b, 0)
	r.PrintAlignments()
	fmt.Println(r.Info)

	// Output:
	// nb alignment  1
	// [size: 5, style: flanked, capturable]
}

func TestCheckRules(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	// test: 0 > simple invalid move
	b = boards.GetSimpleP2()
	r = New(Player1, 3, 3)
	r.CheckRules(b, 3)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > invalid move by double three action
	b = boards.GetTreeP1_1()
	r = New(Player2, 10, 10)
	r.CheckRules(b, 3)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > win by capture
	b = boards.GetCaptureP1_1()
	r = New(Player2, 0, 18)
	r.CheckRules(b, 4)
	if r.Win == false {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3 > align five token but align is capturable
	b = boards.GetWinCapturableP2()
	r = New(Player1, 9, 11)
	r.CheckRules(b, 3)
	if r.Win == true {
		t.Error(t.Name() + " > test: 3")
	}

	// test: 4 > win by align five token
	b = boards.GetWinNoCapturableP2()
	r = New(Player1, 9, 11)
	r.CheckRules(b, 3)
	if r.Win == false {
		t.Error(t.Name() + " > test: 4")
	}

	// test: 5 > no win but align five token (not consecutive)
	b = boards.GetWinSituationP2()
	r = New(Player1, 9, 14)
	r.CheckRules(b, 3)
	if r.Win == true {
		t.Error(t.Name() + " > test: 5")
	}

	// test: 6 > invalid move by double three action #2
	b = boards.GetTreeP1_2()
	r = New(Player1, 9, 8)
	r.CheckRules(b, 0)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 6 ")
	}

	// test: 7 > invalid move by double three action #2
	b = boards.GetTreeP1_3()
	r = New(Player1, 9, 10)
	r.CheckRules(b, 0)
	if r.Movable == false {
		t.Error(t.Name() + " > test: 7 ")
	}
}

func TestLenAlignment_P1Only(t *testing.T) {
	// Player1 rules neutral instance
	r := New(Player1, 0, 0)

	// test: 1 > 2 token aligned consecutively
	mask1 := [11]uint8{0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask1); a == nil || a.size != 2 {
		t.Error(t.Name()+" > test: 1 / alignment struct : ", a)
	}

	// test: 2 > 2 token aligned not consecutively
	mask2 := [11]uint8{0, 0, 0, 0, 0, 1, 0, 1, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask2); a == nil || a.size != 2 {
		t.Error(t.Name()+" > test: 2 / alignment struct : ", a)
	}

	// test: 3 > 3 token aligned consecutively
	mask3 := [11]uint8{0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask3); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 3 / alignment struct : ", a)
	}

	// test: 4 > 3 token aligned not consecutively
	mask4 := [11]uint8{0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask4); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 4 / alignment struct : ", a)
	}

	// test: 5 > 3 token aligned not consecutively
	mask5 := [11]uint8{0, 0, 1, 0, 0, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask5); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 5 / alignment struct : ", a)
	}

	// test: 6 > 3 token aligned not consecutively
	mask6 := [11]uint8{0, 0, 1, 0, 0, 1, 1, 1, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask6); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 6 / alignment struct : ", a)
	}

	// test: 7 > 4 token aligned consecutively
	mask7 := [11]uint8{0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask7); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 7 / alignment struct : ", a)
	}

	// test: 8 > 4 token aligned not consecutively
	mask8 := [11]uint8{0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask8); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 8 / alignment struct : ", a)
	}

	// test: 9 > 4 token aligned not consecutively
	mask9 := [11]uint8{0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask9); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 9 / alignment struct : ", a)
	}

}

func TestLenAlignmentP1_P2(t *testing.T) {
	// Player1 rules neutral instance
	r := New(Player1, 0, 0)

	// test: 1 > 2 token aligned consecutively
	mask1 := [11]uint8{0, 0, 0, 0, 2, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask1); a == nil || a.size != 2 {
		t.Error(t.Name()+" > test: 1 / alignment struct : ", a)
	}

	// test: 2 > 2 token aligned not consecutively
	mask2 := [11]uint8{0, 0, 2, 0, 2, 1, 0, 1, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask2); a == nil || a.size != 2 {
		t.Error(t.Name()+" > test: 2 / alignment struct : ", a)
	}

	// test: 3 > 3 token aligned consecutively
	mask3 := [11]uint8{0, 0, 0, 0, 2, 1, 1, 1, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask3); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 3 / alignment struct : ", a)
	}

	// test: 4 > 3 token aligned not consecutively
	mask4 := [11]uint8{0, 2, 0, 1, 0, 1, 1, 2, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask4); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 4 / alignment struct : ", a)
	}

	// test: 5 > 3 token aligned not consecutively
	mask5 := [11]uint8{0, 2, 1, 0, 0, 1, 1, 2, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask5); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 5 / alignment struct : ", a)
	}

	// test: 6 > 3 token aligned not consecutively
	mask6 := [11]uint8{0, 0, 1, 0, 0, 1, 1, 1, 2, 0, 0}
	if a := r.analyzeLenAlignment(&mask6); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 6 / alignment struct : ", a)
	}

	// test: 7 > 4 token aligned consecutively
	mask7 := [11]uint8{0, 0, 0, 1, 1, 1, 1, 2, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask7); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 7 / alignment struct : ", a)
	}

	// test: 8 > 4 token aligned not consecutively
	mask8 := [11]uint8{0, 2, 1, 1, 0, 1, 1, 0, 0, 2, 0}
	if a := r.analyzeLenAlignment(&mask8); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 8 / alignment struct : ", a)
	}

	// test: 9 > 4 token aligned not consecutively
	mask9 := [11]uint8{0, 1, 1, 0, 0, 1, 1, 2, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask9); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 9 / alignment struct : ", a)
	}

}

func TestLenAlignment(t *testing.T) {
	// Player1 rules neutral instance
	r := New(Player1, 0, 0)

	// test: 1 > 2 token aligned consecutively
	mask1 := [11]uint8{3, 3, 3, 3, 2, 1, 1, 0, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask1); a == nil || a.size != 2 {
		t.Error(t.Name()+" > test: 1 / alignment struct : ", a)
	}

	// test: 2 > 2 token aligned not consecutively
	mask2 := [11]uint8{3, 3, 2, 0, 2, 1, 0, 1, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask2); a == nil || a.size != 2 {
		t.Error(t.Name()+" > test: 2 / alignment struct : ", a)
	}

	// test: 3 > 3 token aligned consecutively
	mask3 := [11]uint8{0, 0, 0, 0, 1, 1, 1, 0, 3, 3, 3}
	if a := r.analyzeLenAlignment(&mask3); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 3 / alignment struct : ", a)
	}

	// test: 4 > 3 token aligned not consecutively
	mask4 := [11]uint8{0, 2, 0, 1, 0, 1, 1, 2, 3, 3, 3}
	if a := r.analyzeLenAlignment(&mask4); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 4 / alignment struct : ", a)
	}

	// test: 5 > 3 token aligned not consecutively
	mask5 := [11]uint8{0, 2, 1, 0, 0, 1, 1, 2, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask5); a == nil || a.size != 3 {
		t.Error(t.Name()+" > test: 5 / alignment struct : ", a)
	}

	// test: 6 > 3 token aligned not consecutively
	mask6 := [11]uint8{0, 0, 1, 0, 0, 1, 1, 1, 2, 0, 0}
	if a := r.analyzeLenAlignment(&mask6); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 6 / alignment struct : ", a)
	}

	// test: 7 > 4 token aligned consecutively
	mask7 := [11]uint8{3, 3, 0, 1, 1, 1, 1, 2, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask7); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 7 / alignment struct : ", a)
	}

	// test: 8 > 4 token aligned not consecutively
	mask8 := [11]uint8{0, 2, 1, 1, 0, 1, 1, 0, 0, 2, 0}
	if a := r.analyzeLenAlignment(&mask8); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 8 / alignment struct : ", a)
	}

	// test: 9 > 4 token aligned not consecutively
	mask9 := [11]uint8{3, 1, 1, 0, 0, 1, 1, 2, 0, 0, 0}
	if a := r.analyzeLenAlignment(&mask9); a == nil || a.size != 4 {
		t.Error(t.Name()+" > test: 9 / alignment struct : ", a)
	}

}
