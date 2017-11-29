package alphabeta

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

/* functions to help creation tests */
func createNodesB(b *testing.B, s *State, spots []int8) *Node {
	var y, x int8
	var node, prev *Node

	for index, spot := range spots {
		if index%2 == 0 {
			y = spot
		} else {
			x = spot
			if node = s.newNode(y, x); node == nil {
				b.Fatal(b.Name(), " can't create a new node on (y: ", y, " - x: ", x, ")")
			}
			s.updateData(node, prev)
			prev = node
		}
	}
	return prev
}

func BenchmarkScoreAlignment_noAlignment(b *testing.B) {
	var current *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)
	node := createNodesB(b, state, []int8{9, 7})
	current, _ = getNewScore(node)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.scoreAlignment(node, current, state.maxDepth, true)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkScoreAlignment_freeThree(b *testing.B) {
	var current *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)
	_ = createNodesB(b, state, []int8{9, 7})
	node := createNodesB(b, state, []int8{8, 8, 10, 8, 9, 10, 8, 11, 8, 6, 8, 10})
	current, _ = getNewScore(node)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.scoreAlignment(node, current, 0, true)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkEvalAlignment_noAlignment(b *testing.B) {
	var current, opponent *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)
	node := createNodesB(b, state, []int8{9, 7})
	current, opponent = getNewScore(node)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.evalAlignment(node, current, opponent)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkEvalAlignment_freeThree(b *testing.B) {
	var current, opponent *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)
	node := createNodesB(b, state, []int8{9, 7, 8, 8, 10, 8, 9, 10})
	current, opponent = getNewScore(node)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.evalAlignment(node, current, opponent)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkEvalAlignment_fourSpotAligned(b *testing.B) {
	var current, opponent *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)
	node := createNodesB(b, state, []int8{9, 7, 8, 8, 10, 8, 9, 10})
	current, opponent = getNewScore(node)

	//StartBenchmark alignment
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.evalAlignment(node, current, opponent)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkEvalCapture(b *testing.B) {
	var current, opponent *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)

	node := createNodesB(b, state, []int8{9, 10, 9, 7, 8, 8, 8, 9, 8, 10, 9, 11, 9, 9, 9, 12, 9, 6, 9, 5, 9, 13, 9, 4})
	current, opponent = getNewScore(node)

	//StartBenchmark capture
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.evalCapture(node, current, opponent)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyseScore_winByCapture(b *testing.B) {
	var current, opponent *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)

	current = &Score{capturable: true, capture: 6, alignment: 17}
	opponent = &Score{capturable: false, capture: 0, alignment: 14}
	//StartBenchmark analyzeScore
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.analyzeScore(current, opponent)
	}
	b.StopTimer()
	b.ReportAllocs()
}

func BenchmarkAnalyseScore_winByAlignment(b *testing.B) {
	var current, opponent *Score

	board := boards.GetStartP1_1()
	state := New(board, rdef.Player2)

	current = &Score{capturable: true, capture: 10, alignment: -31765}
	opponent = &Score{capturable: true, capture: 11, alignment: 19}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.analyzeScore(current, opponent)
	}
	b.StopTimer()
	b.ReportAllocs()
}
