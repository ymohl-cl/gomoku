package alphabeta

import (
	"fmt"
	"math"

	"github.com/ymohl-cl/gomoku/game/ruler/alignment"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

const (
	/* general evaluation score */
	// scoreMin and scoreMax define the range of possibilities to the evaluation
	scoreMin = math.MinInt16 + 1
	scoreMax = math.MaxInt16
	scoreWin = scoreMin
	// scoreNeutral is the score without advantage for the players
	scoreNeutral = int16(-16355)

	/* specific evaluation score */
	// scoreFirst is a bonus to the first player which play to differ the equal score
	scoreFirst = int16(1)
	// scoreWinDetection is set when win situation is detected on out simulation
	scoreWinDetection = scoreWin + 1000
	scoreByCapture    = int16(6)
	// scoreAlign is a bonus to the additional alignments
	scoreByAlign = int16(1)
	// scoreFree / half and flanked are the multiplier to the alignment type
	scoreFree    = int16(6)
	scoreHalf    = int16(5)
	scoreFlanked = int16(4)

	// depthOutEvalToFreeThree is a the additional depths to make a win from a free tree.
	depthOutEvalToFreeThree = int16(4)
	// depthOutEvalToFourSpots is a the additional depths to make a win from a four spot alignment.
	depthOutEvalToFourSpots = int16(2)
)

// Score is the data to each player pending the evaluation
type Score struct {
	idPlayer      uint8
	capturable    bool
	capture       int16
	alignment     int16
	nbAlignements int16
	isTree        bool
}

func maxWeight(v1, v2 int16) int16 {
	if v1 >= v2 {
		return v1
	}
	return v2
}

// calcul score on alignment simulation, if flag == true, the win condition is applies
func (s *State) scoreAlignment(n *Node, sc *Score, flag bool) {
	var coef int16
	var a *alignment.Alignment
	var nbr int16
	var score int16
	var tree bool

	// get number align to the spot
	if nbr = n.rule.GetNumberAlignment(); nbr == 0 {
		return
	}

	// remove from counter the better alignment
	nbr--
	a = n.rule.GetBetterAlignment()

	// check wins situations
	if flag == true && a.Size == 4 {
		score = scoreWinDetection - nbr + depthOutEvalToFourSpots
	} else if flag == true && a.Size == 3 && a.IsThree {
		score = scoreWinDetection - nbr + depthOutEvalToFreeThree
		tree = true
	}
	if score != 0 && sc.alignment > score {
		sc.alignment = score
		sc.nbAlignements = nbr
		sc.isTree = tree
		return
	}

	// don't continue, because the next step can't be better if the score has already win detected
	if sc.alignment < 0 {
		return
	}

	// define multiplier from alignment type
	if a.IsStyle(rdef.AlignHalf) {
		coef = scoreHalf
	} else if a.IsStyle(rdef.AlignFree) {
		coef = scoreFree
	} else {
		coef = scoreFlanked
	}

	// calcul score
	score = int16(a.Size) * coef
	score += nbr * scoreByAlign

	// check better score
	if sc.alignment < score {
		sc.alignment = score
		sc.nbAlignements = nbr
	}
	return
}

// evalAlignment return score to alignment parameter on this evaluation
func (s *State) evalAlignment(n *Node, current, opponent *Score) {
	spots := new([5]*Node)
	index := 0

	// get score current
	n.rule.UpdateAlignments(s.board)
	s.scoreAlignment(n, current, true)

	// get score opponent - flag define the opponent turn
	flag := true
	for node := n.prev; node != nil; node = node.prev {
		if flag == true && node.rule.IsMyPosition(s.board) {
			// remove previous spots from the board.
			s.updateTokenPlayer(spots)
			// get alignments
			node.rule.UpdateAlignments(s.board)
			// score
			s.scoreAlignment(node, opponent, false)
			spots[index] = node
			index++
		}
		flag = !flag
	}
	// restore spots deleted
	s.restoreTokenPlayer(spots)

	// if there are not winneable situation and equality score.
	// Give advantage to the first player which played
	if opponent.alignment == current.alignment && opponent.alignment > 0 {
		opponent.alignment += scoreFirst
	}

	return
}

// evalCapture return the score to the captures parameter on this evaluation
func (s *State) evalCapture(n *Node, current, opponent *Score) {
	var first uint8

	// browse the evalutation to check the captures.
	for node := n; node != nil; node = node.prev {
		if node.rule.NumberCapture > 0 {
			first = node.rule.GetPlayer()
			if first == current.idPlayer {
				current.capturable = true
			} else {
				opponent.capturable = true
			}
		}
	}

	// set advantage to the first action
	if first == current.idPlayer {
		current.capture += scoreFirst
	} else if first == opponent.idPlayer {
		opponent.capture += scoreFirst
	}

	// set score capture to each player who did a capture during the evaluation
	if current.capturable {
		current.capture += int16(s.getTotalCapture(current.idPlayer)) * scoreByCapture
	}
	if opponent.capturable {
		opponent.capture += int16(s.getTotalCapture(opponent.idPlayer)) * scoreByCapture
	}

	return
}

// analyzeScore return the final weight
func (s *State) analyzeScore(current, opponent *Score, lastNode *Node) int16 {
	var ret int16

	// win condition
	if current.alignment < 0 {
		ret = current.alignment
		ret -= current.capture
		ret += opponent.alignment
		ret += opponent.capture

		if current.isTree && opponent.capturable && s.maxDepth < 10 {
			s.maxDepth += 2
			ret2 := s.alphabetaNegaScout(ret, math.MaxInt16, 2, nil)
			fmt.Println("call with maxDepth: ", s.maxDepth, " - oldeRet: ", ret, " - newRet: ", ret2)
			s.maxDepth -= 2
			ret = ret2
		}

		return ret
	}

	/*else if opponent.alignment < 0 {
		ret = opponent.alignment
		ret -= opponent.capture
		ret += current.alignment
		ret += current.capture
		ret *= -1

		if current.capturable && s.maxDepth < 10 {
			newState := New(s.board, rdef.GetOtherPlayer(current.idPlayer))
			newState.addTotalCapture(current.idPlayer, s.getTotalCapture(current.idPlayer))
			newState.addTotalCapture(opponent.idPlayer, s.getTotalCapture(opponent.idPlayer))
			newState.maxDepth = s.maxDepth + 2
			save := lastNode.prev
			lastNode.prev = nil
			ret = newState.alphabetaNegaScout(math.MinInt16+1, math.MaxInt16, 2, lastNode)
			lastNode.prev = save
			fmt.Println("depth: ", newState.maxDepth)
			return ret
		}

		return ret
	}*/
	/*score = opponent.alignment
	score -= opponent.capture
	score += current.alignment
	score += current.capture
	score *= -1
	// check deap blue
	if opponent.alignment < scoreWinDetection && s.maxDepth < 10 {
		newState := New(s.board, rdef.GetOtherPlayer(current.idPlayer))
		newState.addTotalCapture(current.idPlayer, s.getTotalCapture(current.idPlayer))
		newState.addTotalCapture(opponent.idPlayer, s.getTotalCapture(opponent.idPlayer))
		newState.maxDepth = s.maxDepth + 2
		save := lastNode.prev
		lastNode.prev = nil
		ret := newState.alphabetaNegaScout(math.MinInt16+1, math.MaxInt16, 2, nil)
		lastNode.prev = save
		return ret
	}*/

	//	return score
	//}

	// no win

	ret = scoreNeutral
	ret -= current.alignment - opponent.alignment
	ret -= current.capture - opponent.capture

	// check deap blue
	/*	if s.maxDepth < 10 {
			newState := New(s.board, rdef.GetOtherPlayer(current.idPlayer))
			newState.addTotalCapture(current.idPlayer, s.getTotalCapture(current.idPlayer))
			newState.addTotalCapture(opponent.idPlayer, s.getTotalCapture(opponent.idPlayer))
			newState.maxDepth = s.maxDepth + 2
			save := lastNode.prev
			lastNode.prev = nil
			ret = newState.alphabetaNegaScout(ret, math.MaxInt16, 2, lastNode)
			lastNode.prev = save
			fmt.Println("depth: ", newState.maxDepth)
			return ret
		}
	*/
	return ret
}

// eval function define the weight to the evaluation
// see constantes to check the score parameters
func (s *State) eval(n *Node, depth uint8) int16 {
	var current Score
	var opponent Score

	current.idPlayer = n.rule.GetPlayer()
	opponent.idPlayer = rdef.GetOtherPlayer(current.idPlayer)

	// wins situations
	if n.rule.Win {
		return scoreWin + int16(s.maxDepth-depth)
	}

	s.evalCapture(n, &current, &opponent)
	s.evalAlignment(n, &current, &opponent)

	//	fmt.Println("current: ", current)
	//	fmt.Println("opponent: ", opponent)
	return s.analyzeScore(&current, &opponent, n)
}
