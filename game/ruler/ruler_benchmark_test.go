package ruler

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	"github.com/ymohl-cl/gomoku/game/ruler/alignment"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

/*var stop bool



// check if win exist
r.analyzeWinCondition(board, nbCaps)
return*/

func BenchmarkAvailablePosition(b *testing.B) {
	var board *[19][19]uint8
	var r *Rules

	board = boards.GetEmpty()
	r = New(rdef.Player2, 9, 8)
	// test: 0 > already used 1/2
	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.isAvailablePosition(board)
	}
}

func BenchmarkMaskFromBoard(b *testing.B) {
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
	}
}

func BenchmarkAnalyzeAlignments(b *testing.B) {
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
	}
}

func BenchmarkAnalyzeMoveCondition(b *testing.B) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)
	r.Movable = true

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeMoveCondition()
	}
}

func BenchmarkAnalyzeiWinCondition(b *testing.B) {
	var r *Rules

	r = New(rdef.Player1, 9, 9)

	r.aligns = append(r.aligns, &alignment.Alignment{Size: 3})
	r.analyzeWinCondition(nil, 0)

	//Benchmark
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r.analyzeWinCondition(nil, 0)
	}
}
