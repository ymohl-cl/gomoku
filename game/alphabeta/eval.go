package alphabeta

import (
	"math"

	"github.com/ymohl-cl/gomoku/game/ruler"
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
	scoreByCapture    = int16(5)
	// scoreAlign is a bonus to the additional alignments
	scoreByAlign = int16(1)
	// scoreFree / half and flanked are the multiplier to the alignment type
	scoreFree    = int16(6)
	scoreHalf    = int16(5)
	scoreFlanked = int16(4)

	// depthOutEvalToFreeThree is a the additional depths to make a win from a free tree.
	depthOutEvalToFreeThree = uint8(4)
	// depthOutEvalToFourSpots is a the additional depths to make a win from a four spot alignment.
	depthOutEvalToFourSpots = uint8(2)
)

type Score struct {
	idPlayer       uint8
	capturable     bool
	capture        int16
	alignment      int16
	depthAlignment uint8
}

func maxWeight(v1, v2 int16) int16 {
	if v1 >= v2 {
		return v1
	}
	return v2
}

// calcul score on alignment simulation
func (s *State) scoreAlignment(n *Node, sc *Score, depth uint8) {
	var coef int16
	var size int8
	var a *ruler.Align

	// get number align to the spot
	if nbr := n.rule.GetNumberAlignment(); nbr == 0 {
		return
	}

	// get max alignment creted by this spot
	a = n.rule.GetMaxAlignment()
	size = a.GetSize()

	// check wins situations
	if size >= 3 && a.IsStyle(ruler.AlignFree) {
		if sc.alignment < scoreWinDetection {
			sc.alignment = scoreWinDetection + int16(n.rule.GetNumberAlignment())
			sc.depthAlignment = depthOutEvalToFreeThree + depth
		}
		return
	} else if size == 4 {
		if sc.alignment < scoreWinDetection {
			sc.alignment = scoreWinDetection + int16(n.rule.GetNumberAlignment())
			sc.depthAlignment = depthOutEvalToFourSpots + depth
		}
		return
	}

	// define multiplier from alignment type
	if a.IsStyle(ruler.AlignHalf) {
		coef = scoreHalf
	} else if a.IsStyle(ruler.AlignFree) {
		coef = scoreFree
	} else {
		coef = scoreFlanked
	}

	// calcul score
	score := int16(size) * coef
	score += int16(n.rule.GetNumberAlignment()-1) * scoreByAlign

	// check better score
	if sc.alignment < score {
		sc.alignment = score
		sc.depthAlignment = depth
	}
	return
}

// evalAlignment return score to alignment parameter on this evaluation
func (s *State) evalAlignment(n *Node, current, opponent *Score) {
	var spots []*Node

	// get score current
	s.scoreAlignment(n, current, s.maxDepth)

	// get score opponent - flag define the opponent turn
	flag := true
	depth := s.maxDepth - 1
	for node := n.prev; node != nil; node = node.prev {
		if flag == true && node.rule.IsMyPosition(s.board) {
			// remove previous spots from the board.
			s.updateTokenPlayer(&spots)
			// check if this spot don't be captured
			node.rule.UpdateAlignment(s.board)
			// get score opponent on this move
			s.scoreAlignment(node, opponent, depth)
			// save the current spot
			spots = append(spots, node)
		}
		depth--
		flag = !flag
	}
	// restore spots deleted
	s.restoreTokenPlayer(&spots)

	// set advantage to the first action
	if opponent.alignment > 0 {
		opponent.alignment += scoreFirst
	} else if current.alignment > 0 {
		current.alignment += scoreFirst
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
func (s *State) analyzeScore(current, opponent *Score) int16 {
	return 0 // score
}

// eval function define the weight to the evaluation
// see constantes to check the score parameters
func (s *State) eval(n *Node, depth uint8) int16 {
	var current Score
	var opponent Score

	current.idPlayer = n.rule.GetPlayer()
	opponent.idPlayer = ruler.GetOtherPlayer(current.idPlayer)

	// wins situations
	if n.rule.Win {
		return scoreWin + int16(s.maxDepth-depth)
	}

	s.evalCapture(n, &current, &opponent)
	s.evalAlignment(n, &current, &opponent)

	return s.analyzeScore(&current, &opponent)
}
