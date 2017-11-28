package ruler

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	"github.com/ymohl-cl/gomoku/game/ruler/alignment"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

func BenchmarkIsAvailablePosition_true(b *testing.B) {
	board := boards.GetSimpleP2()
	r := New(rdef.Player1, 8, 9)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.isAvailablePosition(board)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkIsAvailablePosition_false(b *testing.B) {
	board := boards.GetSimpleP2()
	r := New(rdef.Player1, 11, 11)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.isAvailablePosition(board)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeWinCondition_byCapture(b *testing.B) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)
	r.NumberCapture = 2

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeWinCondition(nil, 3)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeWinCondition_byAlignment(b *testing.B) {
	mask := [11]uint8{0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0}
	r := New(rdef.Player1, 9, 9)
	a := &alignment.Alignment{Size: 5}
	a.NewInfosWin(&mask, -1, -1, false)
	r.aligns = append(r.aligns, &alignment.Alignment{Size: 5})

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeWinCondition(nil, 0)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeWinCondition_false(b *testing.B) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeWinCondition(nil, 0)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeMoveCondition_true(b *testing.B) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)
	r.Movable = true

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeMoveCondition()
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeMoveCondition_false(b *testing.B) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)
	r.Movable = true
	r.NumberThree = 3

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeMoveCondition()
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeAlignments_win(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules
	var mask [11]uint8

	board = boards.GetWinNoCapturableP2()
	r = New(rdef.Player1, 9, 11)

	r.getMaskFromBoard(board, -1, -1, &mask)

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeAlignments(&mask, -1, -1)
		r.aligns = nil
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeAlignments_noWin(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules
	var mask [11]uint8

	board = boards.GetAlignFlankedP2()
	r = New(rdef.Player1, 9, 14)

	r.getMaskFromBoard(board, -1, -1, &mask)

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeAlignments(&mask, -1, -1)
		r.aligns = nil
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyzeCapture(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules
	var mask [11]uint8

	board = boards.GetCaptureP1_1()
	r = New(rdef.Player2, 0, 18)

	r.getMaskFromBoard(board, -1, -1, &mask)

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeCapture(&mask, -1, -1)
		r.captures = nil
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkGetMaskFromBoard(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules
	var mask [11]uint8

	board = boards.GetFilledRightP2()
	r = New(rdef.Player1, 10, 18)

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.getMaskFromBoard(board, -1, -1, &mask)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkCheckRules_invalidPosition(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules

	board = boards.GetTreeP1_1()

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = New(rdef.Player2, 13, 13)
		r.CheckRules(board, 0)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkCheckRules_three(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules

	board = boards.GetTreeP1_1()

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = New(rdef.Player2, 9, 10)
		r.CheckRules(board, 0)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkCheckRules_win(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules

	board = boards.GetWinNoCapturableP2()

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r = New(rdef.Player1, 9, 11)
		r.CheckRules(board, 0)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkGetAlignment(b *testing.B) {
	var mask [11]uint8
	board := boards.GetAlignFlankedP2()
	r := New(rdef.Player1, 9, 14)

	r.getMaskFromBoard(board, -1, -1, &mask)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.getAlignment(&mask, -1, -1)
		r.aligns = nil
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkUpdateAlignments(b *testing.B) {
	board := boards.GetTreeP1_1()
	r := New(rdef.Player2, 9, 10)

	r.CheckRules(board, 0)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.UpdateAlignments(board)
	}
	b.StopTimer()
	b.ReportAllocs()
}
