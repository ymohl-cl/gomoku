package alphabeta

import (
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
	scoreWinDetection = scoreWin + 5000
	scoreByCapture    = int16(15)
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

	// get number align to the spot
	if nbr = n.rule.GetNumberAlignment(); nbr == 0 {
		return
	}

	// remove from counter the better alignment
	nbr--
	a = n.rule.GetBetterAlignment()

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
	//	if sc.alignment < score {
	sc.alignment += score
	//	sc.nbAlignements = nbr
	//}
	return
}

// evalAlignment return score to alignment parameter on this evaluation
func (s *State) evalAlignment(n *Node, current, opponent *Score) {
	spots := new([10]*Node)
	index := 0

	flag := false
	// get score current - flag define current player
	for node := n; node != nil; node = node.prev {
		if flag == false && node.rule.IsMyPosition(s.board) {
			// remove previous spots from the board.
			s.updateTokenPlayer(spots)
			// get alignments
			node.rule.UpdateAlignments(s.board)
			// score
			s.scoreAlignment(node, current, true)
			spots[index] = node
			index++
		}
		flag = !flag
	}
	// restore spots deleted
	s.restoreTokenPlayer(spots)

	spots = new([10]*Node)
	index = 0
	// get score opponent - flag define the opponent turn
	flag = true
	for node := n.prev; node != nil; node = node.prev {
		if flag == true && node.rule.IsMyPosition(s.board) {
			// remove previous spots from the board.
			s.updateTokenPlayer(spots)
			// get alignments
			node.rule.UpdateAlignments(s.board)
			// score
			s.scoreAlignment(node, opponent, true)
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

	/*
		var scoreCurrent int16
		var scoreCurrentBis int16
		var scoreOpponent int16
		var scoreOpponentBis int16
	*/
	ret = scoreNeutral
	ret -= (current.alignment + current.capture)
	ret += (opponent.alignment + opponent.capture)
	/*	if current.alignment > current.capture {
			scoreCurrent = current.alignment
			scoreCurrentBis = current.capture
		} else {
			scoreCurrent = current.capture
			scoreCurrentBis = current.alignment
		}

		if opponent.alignment > opponent.capture {
			scoreOpponent = opponent.alignment
			scoreOpponentBis = opponent.capture
		} else {
			scoreOpponent = opponent.capture
			scoreOpponentBis = opponent.alignment
		}

		ret = scoreNeutral

		if scoreCurrent > scoreOpponent && scoreOpponentBis == 0 {
			ret -= (scoreCurrent - scoreOpponent + scoreCurrentBis)
		} else if scoreOpponent > scoreCurrent && scoreCurrentBis == 0 {
			ret += (scoreOpponent - scoreCurrent + scoreOpponentBis)
		} else {
			ret -= (scoreCurrent - scoreOpponent)
		}
	*/
	return ret
}

// eval function define the weight to the evaluation
// see constantes to check the score parameters
func (s *State) eval(n *Node, depth uint16) int16 {
	var current Score
	var opponent Score

	current.idPlayer = n.rule.GetPlayer()
	opponent.idPlayer = rdef.GetOtherPlayer(current.idPlayer)

	s.evalCapture(n, &current, &opponent)
	s.evalAlignment(n, &current, &opponent)

	// wins situations
	if n.rule.Win {
		a := n.rule.GetBetterAlignment()
		if a != nil && !a.GetCaptureStatus() {
			return scoreWin + int16(s.maxDepth-depth)
		}
	}
	return s.analyzeScore(&current, &opponent, n)
}
