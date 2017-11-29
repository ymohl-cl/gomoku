package ruler

import (
	"fmt"
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	"github.com/ymohl-cl/gomoku/game/ruler/alignment"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
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
	r.analyzeAlignments(&mask, -1, -1)

	r.getMaskFromBoard(b, -1, 0, &mask)
	r.analyzeAlignments(&mask, -1, 0)

	r.getMaskFromBoard(b, -1, 1, &mask)
	r.analyzeAlignments(&mask, -1, 1)

	r.getMaskFromBoard(b, 0, -1, &mask)
	r.analyzeAlignments(&mask, 0, -1)

	return r
}

func ExampleNew() {
	r := New(rdef.Player1, 9, 9)
	r.PrintMove()
	// Output:
	// player:  1  (y:  9  - x:  9 )
}

func ExampleRules_Init() {
	r := New(rdef.Player1, 9, 9)
	r.Init(rdef.Player2, 8, 8)
	r.PrintMove()
	// Output:
	// player:  2  (y:  8  - x:  8 )
}

func TestApplyMove(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetEmpty()
	r = New(rdef.Player2, 9, 9)

	// test: 0
	r.ApplyMove(b)
	if (*b)[9][9] != rdef.Player2 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	b = boards.GetEmpty()
	r = New(rdef.Player2, 9, 9)
	(*b)[9][10] = rdef.GetOtherPlayer(rdef.Player2)
	(*b)[9][11] = rdef.GetOtherPlayer(rdef.Player2)
	r.captures = append(r.captures, &Spot{Y: 9, X: 10})
	r.captures = append(r.captures, &Spot{Y: 9, X: 11})
	r.ApplyMove(b)
	if (*b)[9][10] != rdef.Empty && (*b)[9][11] != rdef.Empty {
		t.Error(t.Name() + " > test: 2")
	}
}

func TestRestoreMove(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetEmpty()
	r = New(rdef.Player2, 9, 9)
	(*b)[9][10] = rdef.Player2

	// test: 0
	r.RestoreMove(b)
	if (*b)[9][9] != rdef.Empty {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	b = boards.GetEmpty()
	r = New(rdef.Player2, 9, 9)
	(*b)[9][9] = rdef.Player2
	(*b)[9][10] = rdef.Empty
	(*b)[9][11] = rdef.Empty
	r.captures = append(r.captures, &Spot{Y: 9, X: 10})
	r.captures = append(r.captures, &Spot{Y: 9, X: 11})
	r.RestoreMove(b)
	if (*b)[9][9] != rdef.Empty && (*b)[9][10] != rdef.Player1 || (*b)[9][11] != rdef.Player1 {
		t.Error(t.Name() + " > test: 1")
	}
}

func TestGetCapture(t *testing.T) {

	// test: 0
	r := New(rdef.Player2, 9, 9)
	captures := r.GetCaptures()
	if len(captures) != 0 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	r.captures = append(r.captures, &Spot{Y: 9, X: 10})
	r.captures = append(r.captures, &Spot{Y: 9, X: 11})
	captures = r.GetCaptures()
	if len(captures) != 2 {
		t.Error(t.Name() + " > test: 2")
	}
}

func TestGetPlayer(t *testing.T) {
	var r *Rules

	// test: 0
	r = New(rdef.Player1, 9, 9)
	if r.GetPlayer() != rdef.Player1 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	r = New(rdef.Player2, 8, 8)
	if r.GetPlayer() != rdef.Player2 {
		t.Error(t.Name() + " > test: 1")
	}
}

func TestGetPosition(t *testing.T) {
	// test: 0
	r := New(rdef.Player1, 9, 9)
	if y, x := r.GetPosition(); y != 9 && x != 9 {
		t.Error(t.Name() + " > test: 0")
	}
}

func TestGetNumberAlignment(t *testing.T) {
	// test: 0
	r := New(rdef.Player1, 9, 9)
	if r.GetNumberAlignment() != 0 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	r = New(rdef.Player1, 9, 9)
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 2})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 4})
	if r.GetNumberAlignment() != 3 {
		t.Error(t.Name() + " > test: 3")
	}
}

func TestIsMyPosition(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	r = New(rdef.Player1, 9, 8)
	b = boards.GetEmpty()

	// test: 0 > position not set
	if r.IsMyPosition(b) {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > position set
	(*b)[9][8] = rdef.Player1
	if !r.IsMyPosition(b) {
		t.Error(t.Name() + " > test: 1")
	}
}

func TestIsAvailablePosition(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetEmpty()
	r = New(rdef.Player2, 9, 8)
	// test: 0 > already used 1/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 0")
	}

	b = boards.GetSimpleP2()
	r = New(rdef.Player1, 9, 6)
	// test: 1 > allowed position
	if !r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 1")
	}

	r = New(rdef.Player1, 9, 5)
	// test: 2 > no neighbour 1/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 2")
	}

	r = New(rdef.Player1, 9, 0)
	// test: 3 > neighbour 2/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 3")
	}

	r = New(rdef.Player1, 9, 9)
	// test: 4 > already used 2/2
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 4")
	}

	r = New(rdef.Player1, -1, 6)
	// test: 5 > out of board
	if r.isAvailablePosition(b) {
		t.Error(t.Name() + " > test: 5")
	}
}

func TestAnalyzeWin(t *testing.T) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)
	// test: 0
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3})
	r.analyzeWinCondition(nil, 0)
	if r.Win == true {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	r.NumberCapture = 3
	r.analyzeWinCondition(nil, 2)
	if r.Win == false {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 2
	r.NumberCapture = 4
	r.analyzeWinCondition(nil, 3)
	if r.Win == false {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 5})
	r.analyzeWinCondition(nil, 0)
	if r.Win == false {
		t.Error(t.Name() + " > test: 3")
	}
}

func TestAnalyzeMoveCondition(t *testing.T) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)
	r.Movable = true
	// test: 0
	r.analyzeMoveCondition()
	if r.Movable == false {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	r.NumberThree = 2
	r.NumberCapture = 1
	r.analyzeMoveCondition()
	if r.Movable == false {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2
	r.NumberCapture = 0
	r.analyzeMoveCondition()
	if r.Movable == true {
		t.Error(t.Name() + " > test: 2")
	}
}

func TestAnalyzeAlignments(t *testing.T) {
	var r *Rules
	var mask [11]uint8

	// test: 0 > not enought size to align five spot
	r = New(rdef.Player1, 9, 9)
	mask = [11]uint8{3, 3, 2, 0, 1, 1, 1, 2, 0, 0, 0}
	r.analyzeAlignments(&mask, -1, 1)
	if r.NumberThree != 0 && len(r.aligns) != 0 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > basic alignment
	r = New(rdef.Player1, 9, 9)
	mask = [11]uint8{3, 3, 2, 0, 1, 1, 1, 0, 2, 0, 0}
	r.analyzeAlignments(&mask, -1, 1)
	if r.NumberThree != 0 && len(r.aligns) != 0 {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > basic alignment
	r = New(rdef.Player1, 9, 9)
	mask = [11]uint8{3, 3, 1, 0, 1, 1, 1, 0, 2, 0, 0}
	r.analyzeAlignments(&mask, -1, 1)
	if r.NumberThree != 0 && len(r.aligns) != 0 {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3 > three alignment
	r = New(rdef.Player1, 9, 9)
	mask = [11]uint8{3, 3, 0, 0, 1, 1, 1, 0, 0, 2, 0}
	r.analyzeAlignments(&mask, -1, 1)
	if r.NumberThree != 1 && len(r.aligns) != 0 {
		t.Error(t.Name() + " > test: 3")
	}

	// test: 4 > four alignment
	r = New(rdef.Player1, 9, 9)
	mask = [11]uint8{3, 3, 1, 1, 1, 1, 0, 1, 0, 2, 0}
	r.analyzeAlignments(&mask, -1, 1)
	if r.NumberThree != 0 && len(r.aligns) != 0 {
		t.Error(t.Name() + " > test: 4")
	}

	// test: 5 > three alignment
	r = New(rdef.Player1, 9, 9)
	mask = [11]uint8{3, 3, 1, 1, 1, 1, 1, 0, 0, 2, 0}
	r.analyzeAlignments(&mask, -1, 1)
	if r.NumberThree != 0 && len(r.aligns) != 1 {
		t.Error(t.Name() + " > test: 5")
	}
	for _, a := range r.aligns {
		if a.Size != 5 && a.GetInfosWin() == nil {
			t.Error(t.Name() + " > test: 5 > verification size alignment")
		}
	}
}

func TestAddCapture(t *testing.T) {
	r := New(rdef.Player1, 9, 9)

	// test: 0 > 0 capture
	if len(r.captures) != 0 {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > 1 capture
	r.addCapture(9, 8)
	if len(r.captures) != 1 {
		t.Error(t.Name() + " > test: 1")
	}
	// test: 2 > 1 capture
	r.addCapture(9, 10)
	if len(r.captures) != 2 {
		t.Error(t.Name() + " > test: 2")
	}
}

func ExampleRules_analyzeCapture() {
	var b *[19][19]uint8
	var r *Rules

	b = boards.GetCaptureP1_1()
	r = New(rdef.Player2, 0, 18)
	r = prepareAnalyzeCapture(b, r)
	r.PrintCaptures()

	b = boards.GetCaptureP2_1()
	r = New(rdef.Player1, 7, 9)
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

func ExampleRules_getMaskFrom() {
	var b *[19][19]uint8
	var r *Rules
	var mask [11]uint8

	b = boards.GetSimpleP2()

	// simple test_part1
	r = New(rdef.Player1, 9, 6)
	r.getMaskFromBoard(b, -1, -1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 0, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, -1, 1, &mask)
	fmt.Println(mask)

	r.getMaskFromBoard(b, 0, -1, &mask)
	fmt.Println(mask)

	// simple test_part2
	r = New(rdef.Player1, 8, 9)
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
	r = New(rdef.Player1, 10, 18)
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

func TestCheckRules(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	// test: 0 > simple invalid move
	b = boards.GetSimpleP2()
	r = New(rdef.Player1, 3, 3)
	r.CheckRules(b, 3)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > invalid move by double three action
	b = boards.GetTreeP1_1()
	r = New(rdef.Player2, 10, 10)
	r.CheckRules(b, 3)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > win by capture
	b = boards.GetCaptureP1_1()
	r = New(rdef.Player2, 0, 18)
	r.CheckRules(b, 4)
	if r.Movable == false || r.Win == false {
		t.Error(t.Name() + " > test: 2")
	}

	/*
		// test: 3 > align five token but align is capturable
		b = boards.GetWinCapturableP2()
		r = New(Player1, 9, 11)
		r.CheckRules(b, 3)
		if r.Win == true {
			t.Error(t.Name() + " > test: 3")
		}
	*/

	// test: 4 > win by align five token
	b = boards.GetWinNoCapturableP2()
	r = New(rdef.Player1, 9, 11)
	r.CheckRules(b, 3)
	if r.Movable == false || r.Win == false {
		t.Error(t.Name() + " > test: 4")
	}

	// test: 5 > no win but align five token (not consecutive)
	b = boards.GetWinSituationP2()
	r = New(rdef.Player1, 9, 14)
	r.CheckRules(b, 3)
	if r.Movable == false || r.Win == true {
		t.Error(t.Name() + " > test: 5")
	}

	// test: 6 > invalid move by double three action #1
	b = boards.GetThreeP1_2()
	r = New(rdef.Player1, 9, 8)
	r.CheckRules(b, 0)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 6 ")
	}

	// test: 7 > invalid move by double three action #2
	b = boards.GetTreeP1_3()
	r = New(rdef.Player1, 9, 10)
	r.CheckRules(b, 0)
	if r.Movable == false {
		t.Error(t.Name() + " > test: 7 ")
	}
}

func TestUpdateAlignments(t *testing.T) {
	var b *[19][19]uint8
	var r *Rules

	// test: 0 > simple invalid move
	b = boards.GetSimpleP2()
	r = New(rdef.Player1, 3, 3)
	r.CheckRules(b, 3)
	r.UpdateAlignments(b)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > invalid move by double three action
	b = boards.GetTreeP1_1()
	r = New(rdef.Player2, 10, 10)
	r.CheckRules(b, 3)
	r.UpdateAlignments(b)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2 > win by capture
	b = boards.GetCaptureP1_1()
	r = New(rdef.Player2, 0, 18)
	r.CheckRules(b, 4)
	r.UpdateAlignments(b)
	if r.Movable == false || r.Win == false {
		t.Error(t.Name() + " > test: 2")
	}

	/*
		// test: 3 > align five token but align is capturable
		b = boards.GetWinCapturableP2()
		r = New(Player1, 9, 11)
		r.CheckRules(b, 3)
		r.UpdateAlignments(b)
		if r.Win == true {
			t.Error(t.Name() + " > test: 3")
		}
	*/

	// test: 4 > win by align five token
	b = boards.GetWinNoCapturableP2()
	r = New(rdef.Player1, 9, 11)
	r.CheckRules(b, 3)
	r.UpdateAlignments(b)
	if r.Movable == false || r.Win == false {
		t.Error(t.Name() + " > test: 4")
	}

	// test: 5 > no win but align five token (not consecutive)
	b = boards.GetWinSituationP2()
	r = New(rdef.Player1, 9, 14)
	r.CheckRules(b, 3)
	r.UpdateAlignments(b)
	if r.Movable == false || r.Win == true {
		t.Error(t.Name() + " > test: 5")
	}

	// test: 6 > invalid move by double three action #1
	b = boards.GetThreeP1_2()
	r = New(rdef.Player1, 9, 8)
	r.CheckRules(b, 0)
	r.UpdateAlignments(b)
	if r.Movable == true {
		t.Error(t.Name() + " > test: 6 ")
	}

	// test: 7 > invalid move by double three action #2
	b = boards.GetTreeP1_3()
	r = New(rdef.Player1, 9, 10)
	r.CheckRules(b, 0)
	r.UpdateAlignments(b)
	if r.Movable == false {
		t.Error(t.Name() + " > test: 7 ")
	}
}

/*
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
*/

func TestGetBetterAlignment(t *testing.T) {
	// test: 0
	r := New(rdef.Player1, 9, 9)
	if r.GetBetterAlignment() != nil {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1
	r = New(rdef.Player1, 9, 9)
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3, Style: rdef.AlignFlanked})
	a := r.GetBetterAlignment()
	if a.Style&rdef.AlignFlanked == 0 || a.Style&rdef.AlignFree != 0 || a.Style&rdef.AlignHalf != 0 || a.Size != 3 {
		t.Error(t.Name() + " > test: 1")
	}

	// test: 2
	r = New(rdef.Player1, 9, 9)
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3, Style: rdef.AlignFlanked})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 4, Style: rdef.AlignFlanked})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 2, Style: rdef.AlignFlanked})
	a = r.GetBetterAlignment()
	if a.Style&rdef.AlignFlanked == 0 || a.Style&rdef.AlignFree != 0 || a.Style&rdef.AlignHalf != 0 || a.Size != 4 {
		t.Error(t.Name() + " > test: 2")
	}

	// test: 3
	r = New(rdef.Player1, 9, 9)
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3, Style: rdef.AlignFlanked})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3, Style: rdef.AlignFree})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3, Style: rdef.AlignFlanked})
	a = r.GetBetterAlignment()
	if a.Style&rdef.AlignFlanked != 0 || a.Style&rdef.AlignFree == 0 || a.Style&rdef.AlignHalf != 0 || a.Size != 3 {
		t.Error(t.Name() + " > test: 3")
	}

	// test: 4
	r = New(rdef.Player1, 9, 9)
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 4, Style: rdef.AlignHalf})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 2, Style: rdef.AlignFree})
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3, Style: rdef.AlignFlanked})
	a = r.GetBetterAlignment()
	if a.Style&rdef.AlignFlanked != 0 || a.Style&rdef.AlignFree != 0 || a.Style&rdef.AlignHalf == 0 || a.Size != 4 {
		t.Error(t.Name() + " > test: 4")
	}
}
