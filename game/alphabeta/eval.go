package alphabeta

import "github.com/ymohl-cl/gomoku/game/ruler"

const (
	// maxScoreCapture is return when scoreLimitCapture is found. it's a buff of situation
	maxScoreCapture = 25
	// score for each captures
	scoreByCapture = 5
	// scoreLimitCapture is set when the player make the first capture of eval and has 4 captures
	scoreLimitCapture = 21
)

func maxWeight(v1, v2 int8) int8 {
	if v1 >= v2 {
		return v1
	}
	return v2
}

// evalAlignment return score to alignment parameter on this evaluation
func (s *State) evalAlignment(n *Node) int8 {
	return 0
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
		scoreCurrent++
	} else if first == opponent {
		scoreOpponent++
	}

	if flagCurrent == true {
		scoreCurrent += int8(s.getTotalCapture(current)) * scoreByCapture

		if scoreCurrent == scoreLimitCapture {
			return maxScoreCapture
		}
	}
	if flagOpponent == true {
		scoreOpponent += int8(s.getTotalCapture(opponent)) * scoreByCapture

		if scoreOpponent == scoreLimitCapture {
			return -maxScoreCapture
		}
	}

	return scoreCurrent - scoreOpponent
}

// eval function define the weight to the last node
// keep score range 100 > 127 to wins situations
// init score to 50
// 25 point to capture score. 25 points to alignment score
// ret += (totalScoreCurrent - totalScoreOpponent). Certify positive score
// between 100 and 0 value
func (s *State) eval(n *Node, depth uint8) int8 {
	var ret int8

	current := n.rule.GetPlayer()
	opponent := ruler.GetOtherPlayer(current)

	if n.rule.Win {
		// wins situations
		return -127 + int8(maxDepth-depth)
	}

	// init score
	ret = 50

	ret += s.evalCapture(n, current, opponent)
	ret += s.evalAlignment(n)

	if current == ruler.Player1 {
		return -ret
	}
	return ret
}

/*func eval(s *State, stape uint8, r *ruler.Rules, base *State, maxDepth uint8) int8 {
	ret := int8(63)

	// Var Caps
	var basicCurrent int8
	var basicOther int8
	if s.player == ruler.Player2 {
		basicCurrent = int8(s.nbCapsOther - base.nbCapsOther)
		basicOther = int8(s.nbCapsCurrent - base.nbCapsCurrent)
	} else {
		basicCurrent = int8(s.nbCapsOther - base.nbCapsCurrent)
		basicOther = int8(s.nbCapsCurrent - base.nbCapsOther)
	}

	//Win
	if r.Win {
		if basicCurrent > 0 {
			return -127 + int8(maxDepth-stape)
		} else {
			return -117 + int8(maxDepth-stape)
		}

	}

	//Caps
	ret += (basicCurrent*7 - basicOther*6)

	//Three
	treeCurrent := s.nbTOther
	treeOther := s.nbTCurrent

	//	if len(r.CapturableWin) > 0 {
	//		ret += int8(treeCurrent - treeOther*3)
	//	} else {
	ret += int8(treeCurrent*3 - treeOther)
	//	}

	//Align
	AlignCurrent := s.nbAlignOther
	AlignOther := s.nbAlignCurrent

	ret += int8(AlignCurrent*4 - AlignOther)

	return -ret
}*/
