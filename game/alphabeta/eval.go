package alphabeta

import (
	"fmt"

	"github.com/ymohl-cl/gomoku/game/ruler"
)

const (
	// scoreFirst is a bonus to the first player which play to differ the equal score
	scoreFirst = 1
	// maxScore is return when a win situation is detectable. it's a buff of situation
	scoreMax = 25

	/* Score Captures */
	// score for each captures
	scoreByCapture = 5
	// scoreLimitCapture is set when the player make the first capture of eval and has 4 captures
	scoreLimitCapture = 21

	/* Score Alignments */
	// scoreHalf is the weight by spot on half-free alignment
	scoreHalf = 5
	// scoreFlanked is the weight by spot on flanked alignment
	scoreFlanked = 4
	// scoreFree is the weight by spot on free alignment
	scoreFree = 6
	// scoreAlign is a bonus to the additional alignments
	scoreByAlign = 1
)

func maxWeight(v1, v2 int8) int8 {
	if v1 >= v2 {
		return v1
	}
	return v2
}

func (s *State) printEval(n *Node, depth uint8, score int8) {
	fmt.Println("Score eval: ", score, " - depth: ", depth)

	depthNode := s.maxDepth - depth
	for node := n; node != nil; node = node.prev {
		y, x := node.rule.GetPosition()
		fmt.Println("move ", depthNode, " on y: ", y, " - x: ", x)
		depthNode--
	}

}

func (s *State) scoreAlignment(n *Node) int8 {
	var coef int8
	var size int8
	var number int8
	var a *ruler.Align

	number = n.rule.GetNumberAlignment()
	if number == 0 {
		return 0
	}

	a = n.rule.GetMaxAlignment()
	size = a.GetSize()
	//	fmt.Println("size: ", size)

	if size == 4 {
		return scoreMax
	} else if size == 3 && a.IsStyle(ruler.AlignFree) {
		return scoreMax - 1
	}

	if a.IsStyle(ruler.AlignHalf) {
		coef = scoreHalf
	} else if a.IsStyle(ruler.AlignFree) {
		coef = scoreFree
	} else {
		coef = scoreFlanked
	}
	//	fmt.Println("coef: ", coef)
	//	fmt.Println("number align: ", number)

	ret := size * coef
	ret += (number - 1) * scoreByAlign
	//	fmt.Println("ret: ", ret)
	return ret
}

// evalAlignment return score to alignment parameter on this evaluation
func (s *State) evalAlignment(n *Node) int8 {
	var scoreCurrent int8
	var scoreOpponent int8

	// check first alignment opponent because he it played on first
	if n.prev != nil {
		if !n.prev.rule.IsMyPosition(s.board) {
			scoreOpponent = 0
		} else {
			n.prev.rule.UpdateAlignment(s.board)
			scoreOpponent = s.scoreAlignment(n.prev)
			if scoreOpponent >= (scoreMax - 1) {
				return -scoreOpponent
			}
		}
	}

	// check score alignment for the last move
	scoreCurrent = s.scoreAlignment(n)
	if scoreCurrent >= (scoreMax - 1) {
		return scoreCurrent
	}

	return scoreCurrent - scoreOpponent
}

// evalCapture return the score to the captures parameter on this evaluation
func (s *State) evalCapture(n *Node, current, opponent uint8) int8 {
	var flagCurrent bool
	var flagOpponent bool
	var scoreCurrent int8
	var scoreOpponent int8
	var first uint8

	for node := n; node != nil; node = node.prev {
		if node.rule.NumberCapture > 0 {
			first = node.rule.GetPlayer()
			if first == current {
				flagCurrent = true
			} else {
				flagOpponent = true
			}
		}
	}

	if first == current {
		scoreCurrent += scoreFirst
	} else if first == opponent {
		scoreOpponent += scoreFirst
	}

	if flagCurrent == true {
		scoreCurrent += int8(s.getTotalCapture(current)) * scoreByCapture

		if scoreCurrent == scoreLimitCapture {
			return scoreMax
		}
	}
	if flagOpponent == true {
		scoreOpponent += int8(s.getTotalCapture(opponent)) * scoreByCapture

		if scoreOpponent == scoreLimitCapture {
			return -scoreMax
		}
	}

	return scoreCurrent - scoreOpponent
}

// analyzeScoreCapture return true if win condition detected and adapt the score
func (s *State) analyzeScoreCapture(score *int8, depth uint8) bool {
	if *score == scoreMax {
		*score = -127 + (int8(s.maxDepth-depth) + 2)
		return true
	}

	if *score == -scoreMax {
		*score = 127 - (int8(s.maxDepth-depth) + 2)
		return true
	}

	return false
}

// analyzeScoreAlignment return true if win condition is detected and adapt the score
func (s *State) analyzeScoreAlignment(score *int8, depth uint8) bool {
	// Need to invert sign
	if *score == scoreMax {
		*score = -127 + (int8(s.maxDepth-depth) + 2)
		return true
	} else if *score == scoreMax-1 {
		*score = -127 + (int8(s.maxDepth-depth) + 4)
		return true
	}

	if *score == -scoreMax {
		*score = 127 - (int8(s.maxDepth-depth) + 2)
		return true
	} else if *score == -(scoreMax - 1) {
		*score = 127 - (int8(s.maxDepth-depth) + 4)
		return true
	}

	return false
}

// eval function define the weight to the last node
// keep score range 100 > 127 to wins situations
// init score to 50
// 25 point to capture score. 25 points to alignment score
// ret += (totalScoreCurrent - totalScoreOpponent). Certify positive score
// between 100 and 0 value
func (s *State) eval(n *Node, depth uint8) int8 {
	var score int8
	var ret int8

	//fmt.Println(" ||||||||||||||||||| EVAL ||||||||||||||||||| ")
	current := n.rule.GetPlayer()
	opponent := ruler.GetOtherPlayer(current)

	if n.rule.Win {
		// wins situations
		ret = -127 + int8(s.maxDepth-depth)
	} else {
		// init score
		ret = 50

		score = s.evalCapture(n, current, opponent)
		//fmt.Println("score capture: ", score)
		//if s.analyzeScoreCapture(&score, depth) {
		//fmt.Println("return capture")
		//return score
		//}
		ret += score

		score = s.evalAlignment(n)
		//fmt.Println("score alignement: ", score)
		if s.analyzeScoreAlignment(&score, depth) {
			//fmt.Println("return align")
			return score
		}
		ret += score

		ret = -ret
	}

	// debug
	//s.printEval(n, depth, ret)
	return ret
	/*
		if current == ruler.Player1 {
			return -ret
		}
		return ret
	*/
}
