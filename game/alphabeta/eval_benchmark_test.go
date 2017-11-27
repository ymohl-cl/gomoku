package alphabeta

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

/* functions to help creation tests */
func createNodesB(t *testing.B, s *State, spots []int8) *Node {
	var y, x int8
	var node, prev *Node

	for index, spot := range spots {
		if index%2 == 0 {
			y = spot
		} else {
			x = spot
			if node = s.newNode(y, x); node == nil {
				t.Fatal(t.Name(), " can't create a new node on (y: ", y, " - x: ", x, ")")
			}
			s.updateData(node, prev)
			prev = node
		}
	}
	return prev
}

func BenchmarkEvalCapture(b *testing.B) {

	var board *[19][19]uint8
	var state *State
	var node *Node
	var current, opponent *Score

	board = boards.GetStartP1_1()
	state = New(board, rdef.Player2)

	/* test: P1: 0 capture | P2: 1 capture */
	// State board
	//                     |
	//   . . . . . . . . o x o . . . . . . . .
	// - . . . . . . . x x o . x x . . . . . .
	// createSimulation [P2: 9-6 | P1: 9-5 | P2: 9-13 | P1: 9-4]
	node = createNodesB(b, state, []int8{9, 10, 9, 7, 8, 8, 8, 9, 8, 10, 9, 11, 9, 9, 9, 12, 9, 6, 9, 5, 9, 13, 9, 4})
	current, opponent = getNewScore(node)

	//StartBenchmark capture
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.evalCapture(node, current, opponent)
	}
}

func BenchmarkEvalAlignment(b *testing.B) {
	var board *[19][19]uint8
	var state *State
	var node *Node
	var current, opponent *Score

	board = boards.GetStartP1_1()
	state = New(board, rdef.Player2)

	/* test: 0 > last spot is not winneable situation */
	// State board
	//                     |
	//   . . . . . . x . x o o o o . . . . . .
	// - . . . . . . x o x o x . . . . . . . .
	//   . . . . . . x . o . . . . . . . . . .
	// createSimulation [P1: 9-7 | P2: 8-8 | P1: 10-9 | P2: 9-10]
	node = createNodesB(b, state, []int8{9, 7, 8, 8, 10, 8, 9, 10})
	current, opponent = getNewScore(node)

	//StartBenchmark alignment
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.evalAlignment(node, current, opponent)
	}
}

func BenchmarkAnalyseScore(b *testing.B) {
	var board *[19][19]uint8
	var state *State
	var current, opponent *Score

	board = boards.GetStartP1_1()
	state = New(board, rdef.Player2)

	// test: 0
	current = &Score{capturable: true, capture: 6, alignment: 17}
	opponent = &Score{capturable: false, capture: 0, alignment: 14}
	//StartBenchmark analyzeScore
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		state.analyzeScore(current, opponent)
	}
}
