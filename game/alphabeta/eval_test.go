package alphabeta

import (
	"testing"

	"github.com/ymohl-cl/gomoku/game/boards"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

/* functions to help creation tests */
func createNodes(t *testing.T, s *State, spots []int8) *Node {
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

// getNewScore return new structs Score on order current and opponent
func getNewScore(node *Node) (*Score, *Score) {
	current := new(Score)
	opponent := new(Score)

	current.idPlayer = node.rule.GetPlayer()
	opponent.idPlayer = rdef.GetOtherPlayer(current.idPlayer)
	return current, opponent
}

/* ************************************************************************** */

func TestMaxWeight(t *testing.T) {
	if maxWeight(10, 15) != 15 {
		t.Error(t.Name() + " > test: 0")
	}
	if maxWeight(-10, 15) != 15 {
		t.Error(t.Name() + " > test: 1")
	}
	if maxWeight(10, -15) != 10 {
		t.Error(t.Name() + " > test: 2")
	}
	if maxWeight(-10, -15) != -10 {
		t.Error(t.Name() + " > test: 3")
	}
	if maxWeight(0, 0) != 0 {
		t.Error(t.Name() + " > test: 4")
	}
	if maxWeight(scoreMax, scoreMin) != scoreMax {
		t.Error(t.Name() + " > test: 5")
	}
}

func TestScoreAlignment(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var node *Node
	var current *Score

	b = boards.GetStartP1_1()
	state = New(b, rdef.Player2)

	/* test: 0 > No alignment */
	// State board
	//                     |
	// - . . . . . . . . x o . . . . . . . . .
	// createSimulation [P2: 9-7]
	node = createNodes(t, state, []int8{9, 7})
	current, _ = getNewScore(node)
	// call scoreAlignment to test
	state.scoreAlignment(node, current, state.maxDepth)
	if current.alignment > 0 { // || current.depthAlignment > 0 {
		t.Error(t.Name()+" > test: 0 > score current: ", current)
	}

	/* test: 1 > Align free three + one minor align */
	// State board
	//                     |
	//   . . . . . . x . x . o o . . . . . . .
	// - . . . . . . . o x o . . . . . . . . .
	//   . . . . . . . . o . . . . . . . . . .
	// createSimulation [P1: 8-8 | P2: 10-8 | p1: 9-10 | p2: 8-11 | p1: 8-6 | p2: 8-10]
	node = createNodes(t, state, []int8{8, 8, 10, 8, 9, 10, 8, 11, 8, 6, 8, 10})
	current, _ = getNewScore(node)
	// call scoreAlignment to test
	state.scoreAlignment(node, current, 0)
	if current.alignment != scoreWinDetection-(scoreByAlign)+(depthOutEvalToFreeThree-0) {
		t.Error(t.Name()+" > test: 1 > score current: ", current, " rule: ", node.rule)
	}

	/* test: 2 > Align four three, score already set by free three */
	// State board
	//                     |
	//   . . . . . . x . x . o o . . . . . . .
	// - . . . . . . . o x o x . . . . . . . .
	//   . . . . . . . . o . . . . . . . . . .
	// createSimulation [P1: 9-6 | P2: 8-9 | p1: 10-6 | p2: 8-12]
	node = createNodes(t, state, []int8{9, 6, 8, 9, 10, 6, 8, 12})
	// call scoreAlignment to test
	state.scoreAlignment(node, current, state.maxDepth-2)
	if current.alignment != scoreWinDetection+(depthOutEvalToFourSpots-2) {
		t.Error(t.Name()+" > test: 2 > score current: ", current)
	}

	/* test: 3 > score is already win but the new spot is not winneable situation */
	// State board
	//                     |
	//   . . . . . . x . x o o o o . . . . . .
	// - . . . . . . x o x o x . . . . . . . .
	//   . . . . . . x . o . . . . . . . . . .
	// createSimulation [P1: 8-7 | P2: 9-11]
	node = createNodes(t, state, []int8{8, 7, 9, 11})
	// call scoreAlignment to test
	state.scoreAlignment(node, current, 0)
	if current.alignment != scoreWinDetection+(depthOutEvalToFourSpots-2) {
		t.Error(t.Name()+" > test: 3 > score current: ", current)
	}

	/* test: 4 > two align half three */
	// State board
	//                     |
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . . x o x o x o . . . . . . .
	//   . . . . . . x . o . . . . . . . . . .
	// createSimulation [P1: 9-5 | P2: 11-8]
	node = createNodes(t, state, []int8{9, 5, 11, 8})
	current, _ = getNewScore(node)
	// call scoreAlignment to test
	state.scoreAlignment(node, current, 0)
	if current.alignment != scoreHalf*2 {
		t.Error(t.Name()+" > test: 4 > score current: ", current)
	}

	/* test: 5 > tree spots aligned and flanked  */
	// State board
	//                     |
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . x x o x o x o . . . . . . .
	//   . . . . . . x . o . . . . . . . . . .
	//   . . . . . . . . o . . . . . . . . . .
	// createSimulation [P1: 12-8 | P2: 13-7 | P1: 14-7 | P2: 12-7]
	node = createNodes(t, state, []int8{12, 8, 13, 7, 14, 7, 12, 7})
	current, _ = getNewScore(node)
	// call scoreAlignment to test
	state.scoreAlignment(node, current, 0)
	if current.alignment != scoreFlanked*3 {
		t.Error(t.Name()+" > test: 5 > score current: ", current)
	}

	/* test: 6 > double align of two spots type free  */
	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . x x o x o x o . . . . . . .
	//   . . . . . . x . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . o x . . . . . . . . . .
	//   . . . . . . . o . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .
	// createSimulation [P1: 12-9 | P2: 13-6]
	node = createNodes(t, state, []int8{12, 9, 13, 6})
	current, _ = getNewScore(node)
	// call scoreAlignment to test
	state.scoreAlignment(node, current, 0)
	if current.alignment != (scoreFree*2)+(scoreByAlign*1) {
		t.Error(t.Name()+" > test: 6 > score current: ", current)
	}
	/* final board */
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . x x o x o x o . . . . . . .
	//   . . . . . . x . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . o x x . . . . . . . . .
	//   . . . . . . o o . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .

	b = boards.GetThreeP1_2()
	state = New(b, rdef.Player2)

	/* test: 7 > free three with 4 spots align */
	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . o . . . . . . . . . . . .
	//   . . . . . . . x . . . o . . . . . . .
	//   . . . . . . . o x . x . . . . . . . .
	// - . . . . . . . . . x x . . . . . . . .
	//   . . . . . . x x x o . . . . . . . . .
	//   . . . . . . o . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// createSimulation [P2: 11-9 | P1: 9-5 | P2 11-10 | P1: 6-7 | P2: 9-7 | P1: 9-8]
	node = createNodes(t, state, []int8{11, 9, 9, 5, 11, 10, 6, 7, 9, 7, 9, 8})
	current, _ = getNewScore(node)
	// call scoreAlignment to test
	state.scoreAlignment(node, current, 0)
	if current.alignment != scoreWinDetection-(scoreByAlign*2)+(depthOutEvalToFreeThree-0) {
		t.Error(t.Name()+" > test: 7 > score current: ", current, " - want: ", scoreWinDetection-(scoreByAlign*2)+(depthOutEvalToFreeThree-0))
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . o x . . . . . . . . . . .
	//   . . . . . . . x . . . o . . . . . . .
	//   . . . . . . . o x . x . . . . . . . .
	// - . . . . . x . o x x x . . . . . . . .
	//   . . . . . . x x x o . . . . . . . . .
	//   . . . . . . o . . o o . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
}

func TestEvalAlignment(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var node *Node
	var current, opponent *Score

	b = boards.GetStartP1_1()
	state = New(b, rdef.Player2)

	/* test: 0 > No alignment */
	// State board
	//                     |
	// - . . . . . . . . x o . . . . . . . . .
	// createSimulation [P2: 9-7]
	node = createNodes(t, state, []int8{9, 7})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment > 0 {
		t.Error(t.Name()+" > test: 0 > current: ", current)
	}
	if opponent.alignment > 0 {
		t.Error(t.Name()+" > test: 0 > opponent: ", opponent)
	}

	/* test: 1 > Align free three + one minor align */
	// State board
	//                     |
	//   . . . . . . x . x . o o . . . . . . .
	// - . . . . . . . o x o . . . . . . . . .
	//   . . . . . . . . o . . . . . . . . . .
	// createSimulation [P1: 8-8 | P2: 10-8 | p1: 9-10 | p2: 8-11 | p1: 8-6 | p2: 8-10]
	node = createNodes(t, state, []int8{8, 8, 10, 8, 9, 10, 8, 11, 8, 6, 8, 10})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != scoreWinDetection-(scoreByAlign)+(depthOutEvalToFreeThree-0) {
		t.Error(t.Name()+" > test: 1 > current: ", current, ", score want: ", scoreWinDetection-(scoreByAlign)+(depthOutEvalToFreeThree-0))
	}
	if opponent.alignment != scoreFree*2 {
		t.Error(t.Name()+" > test: 1 > opponent: ", opponent, ", score want: ", scoreFree*2)
	}

	/* test: 2 > Align four three */
	// State board
	//                     |
	//   . . . . . . x . x . o o . . . . . . .
	// - . . . . . . . o x o x . . . . . . . .
	//   . . . . . . . . o . . . . . . . . . .
	// createSimulation [P1: 9-6 | P2: 8-9 | p1: 10-6 | p2: 8-12]
	node = createNodes(t, state, []int8{9, 6, 8, 9, 10, 6, 8, 12})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != scoreWinDetection+(depthOutEvalToFourSpots-0) {
		t.Error(t.Name()+" > test: 2 > current: ", current, ", score want: ", scoreWinDetection+(depthOutEvalToFourSpots-0))
	}
	if opponent.alignment != scoreWinDetection+(depthOutEvalToFreeThree-1) {
		t.Error(t.Name()+" > test: 2 > opponent: ", opponent, ", score want: ", scoreWinDetection+(depthOutEvalToFreeThree-1))
	}

	/* test: 3 > last spot is not winneable situation */
	// State board
	//                     |
	//   . . . . . . x . x o o o o . . . . . .
	// - . . . . . . x o x o x . . . . . . . .
	//   . . . . . . x . o . . . . . . . . . .
	// createSimulation [P1: 8-7 | P2: 9-11]
	node = createNodes(t, state, []int8{8, 7, 9, 11})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != (scoreFree*2)+(scoreByAlign*2) {
		t.Error(t.Name()+" > test: 3 > current: ", current, ", score want: ", (scoreFree*2)+(scoreByAlign*2))
	}
	if opponent.alignment != (scoreHalf*3)+(scoreByAlign*2) {
		t.Error(t.Name()+" > test: 3 > opponent: ", opponent, ", score want: ", (scoreHalf*3)+(scoreByAlign))
	}

	/* test: 4 > two align half three */
	// State board
	//                     |
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . . x o x o x o . . . . . . .
	//   . . . . . . x . o . . . . . . . . . .
	// createSimulation [P1: 9-5 | P2: 11-8]
	node = createNodes(t, state, []int8{9, 5, 11, 8})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != scoreHalf*2 {
		t.Error(t.Name()+" > test: 4 > current: ", current, ", score want: ", scoreHalf*2)
	}
	if opponent.alignment != (scoreFree*2)+(scoreByAlign*2) {
		t.Error(t.Name()+" > test: 4 > opponent: ", opponent, ", score want: ", (scoreFree*2)+(scoreByAlign*2))
	}

	/* test: 5 > tree spots aligned and flanked  */
	// State board
	//                     |
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . x x o x o x o . . . . . . .
	//   . . . . . . x . o . . . . . . . . . .
	//   . . . . . . . . o . . . . . . . . . .
	// createSimulation [P1: 12-8 | P2: 13-7 | P1: 14-7 | P2: 12-7]
	node = createNodes(t, state, []int8{12, 8, 13, 7, 14, 7, 12, 7})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != scoreFlanked*3 {
		t.Error(t.Name()+" > test: 5 > current: ", current, ", score want: ", scoreFlanked*3)
	}
	if opponent.alignment != scoreWinDetection-1+(depthOutEvalToFreeThree-3) {
		t.Error(t.Name()+" > test: 5 > opponent: ", opponent, ", score want: ", scoreWinDetection-1+(depthOutEvalToFreeThree-3))
	}

	/* test: 6 > double align of two spots type free  */
	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . x x o x o x o . . . . . . .
	//   . . . . . . x . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . o x . . . . . . . . . .
	//   . . . . . . . o . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .
	// createSimulation [P1: 12-9 | P2: 13-6]
	node = createNodes(t, state, []int8{12, 9, 13, 6})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != (scoreFree*2)+(scoreByAlign*1) {
		t.Error(t.Name()+" > test: 6 > current: ", current, ", score want: ", (scoreFree*2)+(scoreByAlign))
	}
	if opponent.alignment != (scoreFree*2)+(scoreByAlign*2) {
		t.Error(t.Name()+" > test: 6 > opponent: ", opponent, ", score want: ", (scoreFree*2)+(scoreByAlign*1))
	}

	/* final board */
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . x x x o o o o . . . . . .
	// - . . . . . x x o x o x o . . . . . . .
	//   . . . . . . x . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . o x x . . . . . . . . .
	//   . . . . . . o o . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .

	b = boards.GetThreeP1_2()
	state = New(b, rdef.Player2)

	/* test: 7 > free three with 4 spots align */
	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . o . . . . . . . . . . . .
	//   . . . . . . . x . . . o . . . . . . .
	//   . . . . . . . o x . x . . . . . . . .
	// - . . . . . . . . . x x . . . . . . . .
	//   . . . . . . x x x o . . . . . . . . .
	//   . . . . . . o . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// createSimulation [P2: 11-9 | P1: 9-5 | P2 11-10 | P1: 6-7 | P2: 9-7 | P1: 9-8]
	node = createNodes(t, state, []int8{11, 9, 9, 5, 11, 10, 6, 7, 9, 7, 9, 8})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != scoreWinDetection-(scoreByAlign*2)+(depthOutEvalToFreeThree-0) {
		t.Error(t.Name()+" > test: 7 > current: ", current, ", score want: ", scoreWinDetection-(scoreByAlign*2)+(depthOutEvalToFreeThree-0))
	}
	if opponent.alignment != (scoreFree*3)+(scoreByAlign*1) {
		t.Error(t.Name()+" > test: 7 > opponent: ", opponent, ", score want: ", (scoreFree*3)+(scoreByAlign*1))
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . o x . . . . . . . . . . .
	//   . . . . . . . x . . . o . . . . . . .
	//   . . . . . . . o x . x . . . . . . . .
	// - . . . . . x . o x x x . . . . . . . .
	//   . . . . . . x x x o . . . . . . . . .
	//   . . . . . . o . . o o . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	b = boards.GetAlignFreeP2()
	state = New(b, rdef.Player1)

	/* test: 8 > score no win but equality */
	// State board
	//                     |
	//   . . . . . . . . x . . . . . . . . . .
	// - . . . . . . . . o x . . . . . . . . .
	//   . . . . . . . . o . . . . . . . . . .
	// createSimulation [P1: 8-9 | P2: 10-9]
	node = createNodes(t, state, []int8{8, 9, 10, 9})
	current, opponent = getNewScore(node)
	// call evalAlignment to test
	state.evalAlignment(node, current, opponent)
	if current.alignment != (scoreFree*2)+scoreByAlign {
		t.Error(t.Name()+" > test: 7 > current: ", current, ", score want: ", (scoreFree*2)+scoreByAlign)
	}
	if opponent.alignment != (scoreFree*2)+scoreByAlign+scoreFirst {
		t.Error(t.Name()+" > test: 7 > opponent: ", opponent, ", score want: ", (scoreFree*2)+scoreByAlign+scoreFirst)
	}

	// State board
	//                     |
	//   . . . . . . . . x x . . . . . . . . .
	// - . . . . . . . . o x . . . . . . . . .
	//   . . . . . . . . o o . . . . . . . . .

}

func TestEvalCapture(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var node *Node
	var current, opponent *Score

	b = boards.GetStartP1_1()
	state = New(b, rdef.Player2)

	/* test: 0 > No capture */
	// State board
	//                     |
	// - . . . . . . . . x o . . . . . . . . .
	// createSimulation [P2: 9-10 | P1: 9-7 | P2: 8-8 | P1: 8-9]
	node = createNodes(t, state, []int8{9, 10, 9, 7, 8, 8, 8, 9})
	current, opponent = getNewScore(node)
	// call eval to test
	state.evalCapture(node, current, opponent)
	if current.capturable == true || current.capture != 0 {
		t.Error(t.Name()+" > test: 0 > score current: ", current)
	}
	if opponent.capturable == true || opponent.capture != 0 {
		t.Error(t.Name()+" > test: 0 > score opponent: ", opponent)
	}

	/* test: 1 > P1: 1 capture | P2: 0 capture */
	// State board
	//                     |
	//   . . . . . . . . o x . . . . . . . . .
	// - . . . . . . . x x o o . . . . . . . .
	// createSimulation [P2: 8-10 | P1: 9-11 | P2: 9-9 | P1: 9-12]
	node = createNodes(t, state, []int8{8, 10, 9, 11, 9, 9, 9, 12})
	current, opponent = getNewScore(node)
	// call eval to test
	state.evalCapture(node, current, opponent)
	if current.capturable == false || current.capture != (scoreByCapture+scoreFirst) {
		t.Error(t.Name()+" > test: 1 > score current: ", current)
	}
	if opponent.capturable == true || opponent.capture != 0 {
		t.Error(t.Name()+" > test: 1 > score opponent: ", opponent)
	}

	/* test: 2 > P1: 0 capture | P2: 1 capture */
	// State board
	//                     |
	//   . . . . . . . . o x o . . . . . . . .
	// - . . . . . . . x x o . x x . . . . . .
	// createSimulation [P2: 9-6 | P1: 9-5 | P2: 9-13 | P1: 9-4]
	node = createNodes(t, state, []int8{9, 6, 9, 5, 9, 13, 9, 4})
	current, opponent = getNewScore(node)
	// call eval to test
	state.evalCapture(node, current, opponent)
	if current.capturable == true || current.capture != 0 {
		t.Error(t.Name()+" > test: 2 > score current: ", current)
	}
	if opponent.capturable == false || opponent.capture != (scoreByCapture+scoreFirst) {
		t.Error(t.Name()+" > test: 2 > score opponent: ", opponent)
	}

	/* test: 3 > P1: 1 capture, P2: 1 capture */
	// State board
	//                     |
	//   . . . . . . . . o x o . . . . . . . .
	// - . . . . x x o . . o . x x o . . . . .
	// createSimulation [P2: 9-3 | P1: 8-7 | P2: 9-10 | P1: 9-8]
	node = createNodes(t, state, []int8{9, 3, 8, 7, 9, 14, 9, 15})
	current, opponent = getNewScore(node)
	// call eval to test
	state.evalCapture(node, current, opponent)
	if current.capturable == false || current.capture != (scoreByCapture*2) {
		t.Error(t.Name()+" > test: 3 > score current: ", current)
	}
	if opponent.capturable == false || opponent.capture != ((scoreByCapture*2)+scoreFirst) {
		t.Error(t.Name()+" > test: 3 > score opponent: ", opponent)
	}

	/* test: 4 > P1: 1 capture, P2: 2 capture */
	// State board
	//                     |
	//   . . . . . . . x o x o . . . . . . . .
	// - . . . o . . o . . o . x x . . x . . .
	// createSimulation [P2: 9-13 | P1: 7-7 | P2: 6-7 | P1: 10-10 | P2: 9-10 | P1: 7-8 | P2: 9-7 | P1: 6-8]
	node = createNodes(t, state, []int8{9, 13, 7, 7, 6, 7, 10, 10, 9, 10, 7, 8, 9, 7, 6, 8})
	current, opponent = getNewScore(node)
	// call eval to test
	state.evalCapture(node, current, opponent)
	if current.capturable == false || current.capture != (scoreByCapture*3)+scoreFirst {
		t.Error(t.Name()+" > test: 4 > score current: ", current)
	}
	if opponent.capturable == false || opponent.capture != (scoreByCapture*4) {
		t.Error(t.Name()+" > test: 4 > score opponent: ", opponent)
	}

	// State board
	//                     |
	//   . . . . . . . o x . . . . . . . . . .
	//   . . . . . . . . x . . . . . . . . . .
	//   . . . . . . . . . x o . . . . . . . .
	// - . . . o . . o o . . o . . o . x . . .
	//   . . . . . . . . . . x . . . . . . . .
}

func TestAnalyzeScore(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var current, opponent *Score

	b = boards.GetStartP1_1()
	state = New(b, rdef.Player2)

	// test: 0
	current = &Score{capturable: true, capture: 6, alignment: 17}
	opponent = &Score{capturable: false, capture: 0, alignment: 14}
	if ret := state.analyzeScore(current, opponent); ret != -16364 {
		t.Error(t.Name()+" > test: 0 > want: ", -16364, " got: ", ret)
	}

	// test: 1 > invert est 0
	current = &Score{capturable: false, capture: 0, alignment: 14}
	opponent = &Score{capturable: true, capture: 6, alignment: 17}
	if ret := state.analyzeScore(current, opponent); ret != -16346 {
		t.Error(t.Name()+" > test: 1 > want: ", -16346, " got: ", ret)
	}

	// test: 2 > current make a win out of eval
	current = &Score{capturable: true, capture: 10, alignment: -31765}
	opponent = &Score{capturable: true, capture: 11, alignment: 19}
	if ret := state.analyzeScore(current, opponent); ret != -31745 {
		t.Error(t.Name()+" > test: 0 > want: ", -31745, " got: ", ret)
	}
}

// TestEval

/*
func TestUpdateScoreAlignment(t *testing.T) {
	var b *[19][19]uint8
	var state *State
	var scoreCurrent int16
	var scoreOpponent int16

	b = boards.GetEmpty()
	state = New(b, ruler.Player2)

	flag := false
	scoreCurrent = 12
	scoreOpponent = 10

	saveCurrent := scoreCurrent
	saveOpponent := scoreOpponent

	// score current
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 7)
	if saveCurrent != scoreCurrent {
		t.Error(t.Name() + " > test: 1")
	} else if saveOpponent != scoreOpponent {
		t.Error(t.Name() + " > test: 2")
	} else if flag == false {
		t.Error(t.Name() + " > test: 3")
	}

	// score opponent
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 15)
	if saveCurrent != scoreCurrent {
		t.Error(t.Name() + " > test: 4")
	} else if scoreOpponent != 15 {
		t.Error(t.Name()+" > test: 5 > ret: ", scoreOpponent)
	} else if flag == true {
		t.Error(t.Name() + " > test: 6")
	}

	// score current
	saveOpponent = scoreOpponent
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 0)
	if saveCurrent != scoreCurrent {
		t.Error(t.Name() + " > test: 7")
	} else if saveOpponent != scoreOpponent {
		t.Error(t.Name() + " > test: 8")
	} else if flag == false {
		t.Error(t.Name() + " > test: 9")
	}

	// score opponent
	flag = !flag
	state.updateScoreAlignment(&scoreCurrent, &scoreOpponent, &flag, 22)
	if scoreCurrent != 22 {
		t.Error(t.Name()+" > test: 10 > ret: ", scoreCurrent)
	} else if saveOpponent != scoreOpponent {
		t.Error(t.Name() + " > test: 11")
	} else if flag == false {
		t.Error(t.Name() + " > test: 12")
	}
}

func TestEvalAlignment(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	// State board
	// - . . . . . . . . x o . . . . . . . . .
	b = boards.GetStartP1_1()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 11)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(8, 10)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 12)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// p1 capture: 1
	// p2 capture: 0

	ret := state.evalAlignment(n3, 0)
	// test: 1 > check value of eval
	if ret != 25 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . o . . . . . . . .
	// - . . . . . . . . x . . x x . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(9, 13)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(7, 10)
	state.updateData(n3, n2)

	// test: 2 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 2")
	}

	// p1 capture: 2
	// p2 capture: 1

	ret = state.evalAlignment(n3, 0)
	// test: 3 > check value of eval
	if ret != 13 {
		t.Error(t.Name()+" > test: 3 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . o . . . . .
	//   . . . . . . . . . . x . . . . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(7, 13)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 13)
	state.updateData(n3, n2)

	// test: 4 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 4")
	}

	// p1 capture: 3
	// p2 capture: 1

	ret = state.evalAlignment(n3, 0)
	// test: 5 > check value of eval
	if ret != 13 {
		t.Error(t.Name()+" > test: 5 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . . . . . . . x o . x . . . . .

	// move ai
	n0 = state.newNode(10, 12)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 9)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 8)
	state.updateData(n2, n1)
	// move player
	//	boards.Print(state.board)
	n3 = state.newNode(9, 13)
	state.updateData(n3, n2)

	// test: 6 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 6")
	}

	// p1 capture: 3
	// p2 capture: 2

	ret = state.evalAlignment(n3, 0)
	// test: 7 > check value of eval
	if ret != 25 {
		t.Error(t.Name()+" > test: 7 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . . . o . . o o x . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 7)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 6)
	state.updateData(n3, n2)

	// test: 8 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 8")
	}

	// p1 capture: 4
	// p2 capture: 2

	ret = state.evalAlignment(n3, 0)
	// test: 9 > check value of eval
	if ret != 14 {
		t.Error(t.Name()+" > test: 9 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . x o o . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(11, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 5)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 4)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 11)
	state.updateData(n3, n2)

	// test: 10 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 10")
	}

	// p1 capture: 4
	// p2 capture: 4

	ret = state.evalAlignment(n3, 0)
	// test: 11 > check value of eval
	if ret != -26 {
		t.Error(t.Name()+" > test: 11 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . o . . o o . x x . . . . . . .
	//   . . . . . . . . . . . . . o . . . . .
}

func TestEval_noWin(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	b = boards.GetStartP1_1()

	// State board
	// - . . . . . . . . x o . . . . . . . . .

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 11)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(8, 10)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 12)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// p1 capture: 1
	// p1 alignment: 1 | max: 3 free
	// p2 capture: 0
	// p2 alignment: 0 | max: 0

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval
	if ret != -100 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . o . . . . . . . .
	// - . . . . . . . . x . . x x . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(9, 10)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(9, 13)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(7, 10)
	state.updateData(n3, n2)

	// test: 2 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 2")
	}

	// p1 capture: 2
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 1
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 3 > check value of eval to advantage P1 Caps and align
	if ret != -67 {
		t.Error(t.Name()+" > test: 3 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . o . . . . .
	//   . . . . . . . . . . x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(7, 13)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 13)
	state.updateData(n3, n2)

	// test: 4 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 4")
	}

	// p1 capture: 3
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 1
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 5 > check value of eval to advantage P1 Caps and align
	if ret != -79 {
		t.Error(t.Name()+" > test: 5 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . . . . . . . x o . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(10, 12)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 9)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 8)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(9, 13)
	state.updateData(n3, n2)

	// test: 6 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 6")
	}

	// p1 capture: 3
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 2
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 7 > check value of eval
	if ret != -83 {
		t.Error(t.Name()+" > test: 7 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . . . o . . o o x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(8, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 10)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 7)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 6)
	state.updateData(n3, n2)

	// test: 8 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 8")
	}

	// p1 capture: 4
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 2
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 9 > check value of eval
	if ret != -75 {
		t.Error(t.Name()+" > test: 9 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . x . . . . .
	//   . . . . . . x o o . x . . x . . . . .
	//   . . . . . . . . . . . . . . . . . . .

	// move ai
	n0 = state.newNode(11, 13)
	state.updateData(n0, nil)
	// move player
	n1 = state.newNode(10, 5)
	state.updateData(n1, n0)
	// move ai
	n2 = state.newNode(10, 4)
	state.updateData(n2, n1)
	// move player
	n3 = state.newNode(10, 11)
	state.updateData(n3, n2)

	// test: 10 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 10")
	}

	// p1 capture: 4
	// p1 alignment: 2 | max: 2 free
	// p2 capture: 4
	// p2 alignment: 0 | max: 0

	ret = state.eval(n3, 0)
	// test: 11 > check value of eval to score capture P2 == -25 and align == -24
	if ret != 120 {
		t.Error(t.Name()+" > test: 11 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . x . . x . . . . .
	//   . . . . . . . . . . . . . o . . . . .
	// - . . . . . . . . x . . . . . . . . . .
	//   . . . . o . . o o . x x . . . . . . .
	//   . . . . . . . . . . . . . o . . . . .
}

func TestEval_win(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	b = boards.GetWinNoCapturableP2()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(8, 7)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 15)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(7, 6)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 11)
	state.updateData(n3, n1)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// test: 1 > check rule
	if !n3.rule.Win {
		t.Error(t.Name() + " > test: 1")
	}

	ret := state.eval(n3, 0)
	// test: 2 > check value of eval to align P1
	if ret != -123 {
		t.Error(t.Name()+" > test: 2 > resultat: ", ret)
	}

	// State board
	//   . . . . . o . . . . . . . . . . . . .
	//   . . . . . . o . . . . . . . . . . . .
	// - . . . . . . . o x x x x x o x . . . .
	//   . . . . . . . . . . o . o . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
}

func TestEval_situation1(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . o o x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	b = boards.GetStartP2_1()

	// move ai
	state = New(b, ruler.Player1)
	n0 := state.newNode(9, 7)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(8, 9)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player

	n3 := state.newNode(11, 11)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// p1 capture: 1
	// p2 capture: 0

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval to align P1
	if ret != -32 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . o . . . . . . . . .
	// - . . . . . . . x . . x . . . . . . . .
	//   . . . . . . . . . . . x . . . . . . .
	//   . . . . . . . . . . . o . . . . . . .
}

func TestEval_situation2(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	// - . . . . . . . . o o x . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	b = boards.GetStartP2_1()

	// move ai
	state = New(b, ruler.Player1)
	n0 := state.newNode(9, 7)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(8, 9)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(10, 11)
	state.updateData(n2, n1)
	// move player

	n3 := state.newNode(11, 12)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// p1 capture: 1
	// p2 capture: 1

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval to align P1
	if ret != -61 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . o . . . . . . . . .
	// - . . . . . . . x . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . . . o . . . . . .
}

func TestEval_situation3(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .
	// - . . . . . . . . o x . . . . . . . . .
	//   . . . . . . . x . o . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	b = boards.GetStartP2_2()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(11, 10)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(9, 7)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(7, 7)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(10, 10)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// p1 capture: 0
	// p2 capture: 0

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval to align P1
	if ret != -50 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . o . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .
	// - . . . . . . . x o x . . . . . . . . .
	//   . . . . . . . x . o x . . . . . . . .
	//   . . . . . . . . . . o . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
}

func TestEval_situation4(t *testing.T) {
	var b *[19][19]uint8
	var state *State

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .
	// - . . . . . . . . o x . . . . . . . . .
	//   . . . . . . . x . o . . . . . . . . .
	//   . . . . . . . . . . . . . . . . . . .
	b = boards.GetStartP2_2()

	// move ai
	state = New(b, ruler.Player2)
	n0 := state.newNode(7, 6)
	state.updateData(n0, nil)
	// move player
	n1 := state.newNode(11, 10)
	state.updateData(n1, n0)
	// move ai
	n2 := state.newNode(7, 7)
	state.updateData(n2, n1)
	// move player
	n3 := state.newNode(9, 7)
	state.updateData(n3, n2)

	// test: 0 > check all nodes exists
	if n0 == nil || n1 == nil || n2 == nil || n3 == nil {
		t.Error(t.Name() + " > test: 0")
	}

	// p1 capture: 1
	// p2 capture: 0

	ret := state.eval(n3, 0)
	// test: 1 > check value of eval to align P1
	if ret != -60 {
		t.Error(t.Name()+" > test: 1 > resultat: ", ret)
	}

	// State board
	//                     |
	//   . . . . . . . . . . . . . . . . . . .
	//   . . . . . . o o . . . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .
	// - . . . . . . . x . x . . . . . . . . .
	//   . . . . . . . x . . . . . . . . . . .
	//   . . . . . . . . . . x . . . . . . . .
}
*/
