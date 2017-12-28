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
	scoreWinDetection = scoreWin + 2000
	scoreByCapture    = int16(2000)
	// scoreAlign is a bonus to the additional alignments
	scoreByAlign = int16(1)
	// scoreFree / half and flanked are the multiplier to the alignment type
	scoreFree    = int16(300)
	scoreHalf    = int16(100)
	scoreFlanked = int16(50)

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
	win           bool
	depthWin      int16
}

func maxWeight(v1, v2 int16) int16 {
	if v1 >= v2 {
		return v1
	}
	return v2
}

func scoreAlignment(a *alignment.Alignment, depth int16) int16 {
	var coef int16
	var score int16

	// define multiplier from alignment type
	if a.IsStyle(rdef.AlignHalf) {
		coef = scoreHalf / depth
	} else if a.IsStyle(rdef.AlignFree) {
		coef = scoreFree / depth
	} else {
		coef = scoreFlanked / depth
	}
	if a.Size <= 2 {
		coef = scoreFlanked / depth
	}

	// calcul score
	score = int16(a.Size) * coef

	if a.Size >= 5 {
		score = -2
	} else if a.Size == 4 && a.Style&rdef.AlignFree != 0 {
		score = -1
	}
	return score
}

// calcul score on alignment simulation, if flag == true, the win condition is applies
func (s *State) scoreAlignments(n *Node, sc *Score, player bool, depth int16) {
	for _, a := range n.rule.GetAlignment() {
		ret := scoreAlignment(a, depth)
		if ret == -2 {
			sc.win = true
			sc.depthWin = depth
			return
		} else if ret == -1 {
			sc.isTree = true
			sc.depthWin = depth
		}
		sc.alignment += ret
	}

	// remove from counter the better alignment
	/*
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
		if a.Size == 2 {
			coef = scoreFlanked
		}

		// calcul score
		score = int16(a.Size) * coef
		score += nbr * scoreByAlign

		// check better score
		//	if sc.alignment < score {

		if a.Size >= 5 {
			sc.win = true
			sc.alignment -= score
			sc.capture -= scoreByCapture
		} else {
			sc.alignment += score
		}
		//	sc.nbAlignements = nbr
		//}
	*/
	return
}

// evalAlignment return score to alignment parameter on this evaluation
func (s *State) evalAlignment(n *Node, current, opponent *Score) {
	depth := int16(s.maxDepth)
	spots := new([10]*Node)
	index := 0

	flag := false
	tmpDepth := depth
	nbNode := 0
	// get score current - flag define current player
	for node := n; node != nil; node = node.prev {
		nbNode++
		if flag == false && node.rule.IsMyPosition(s.board) {
			// remove previous spots from the board.
			s.updateTokenPlayer(spots)
			// get alignments
			node.rule.UpdateAlignments(s.board)
			// score
			s.scoreAlignments(node, current, true, tmpDepth)
			spots[index] = node
			index++
		}
		tmpDepth--
		flag = !flag
	}
	// restore spots deleted
	s.restoreTokenPlayer(spots)

	spots = new([10]*Node)
	index = 0
	// get score opponent - flag define the opponent turn
	flag = true
	depth -= 1
	for node := n.prev; node != nil; node = node.prev {
		if flag == true && node.rule.IsMyPosition(s.board) {
			// remove previous spots from the board.
			s.updateTokenPlayer(spots)
			// get alignments
			node.rule.UpdateAlignments(s.board)
			// score
			s.scoreAlignments(node, opponent, false, depth)
			spots[index] = node
			index++
		}
		depth--
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
	depth := int16(s.maxDepth + 1)
	var player bool
	var nbCurrent int16
	var nbOpponent int16

	// browse the evalutation to check the captures.
	for node := n; node != nil; node = node.prev {
		if node.rule.NumberCapture > 0 {
			if player == false {
				current.capture += scoreByCapture / depth
				nbCurrent++
			} else {
				opponent.capture += scoreByCapture / depth
				nbOpponent++
			}
		}
		player = !player
	}

	// set score capture to each player who did a capture during the evaluation
	if current.capture > 0 {
		current.capturable = true
		current.capture += (int16(s.getTotalCapture(current.idPlayer)) - nbCurrent) * scoreByCapture
	}
	if opponent.capturable {
		opponent.capturable = true
		opponent.capture += (int16(s.getTotalCapture(opponent.idPlayer)) - nbOpponent) * scoreByCapture
	}

	return
}

// analyzeScore return the final weight
func (s *State) analyzeScore(current, opponent *Score, lastNode *Node, depth uint16) int16 {
	var ret int16

	/*
		var scoreCurrent int16
		var scoreCurrentBis int16
		var scoreOpponent int16
		var scoreOpponentBis int16
	*/
	if current.win && opponent.win {
		if current.depthWin < opponent.depthWin {
			return scoreWin + current.depthWin + opponent.capture
		}
		return -(scoreWin + opponent.depthWin + current.capture)
	} else if current.win {
		return scoreWin + current.depthWin + opponent.capture
	} else if opponent.win {
		return -(scoreWin + opponent.depthWin + current.capture)
	}

	if current.isTree && opponent.isTree {
		if current.depthWin < opponent.depthWin {
			return scoreWin + current.depthWin + depthOutEvalToFourSpots
		}
		return -(scoreWin + opponent.depthWin + depthOutEvalToFourSpots)
	} else if current.isTree {
		return scoreWin + current.depthWin + depthOutEvalToFourSpots
	} else if opponent.isTree {
		return -(scoreWin + opponent.depthWin + depthOutEvalToFourSpots)
	}
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

	// wins situations
	//	if n.rule.Win || len(n.rule.CapturableWin) > 0 {
	if n.rule.Win {
		return scoreWin + int16(s.maxDepth-depth)
	}

	current.idPlayer = n.rule.GetPlayer()
	opponent.idPlayer = rdef.GetOtherPlayer(current.idPlayer)

	s.evalCapture(n, &current, &opponent)
	s.evalAlignment(n, &current, &opponent)

	//	fmt.Println("c align: ", current.alignment, " - capture: ", current.capture)
	//	fmt.Println("c opponent: ", opponent.alignment, " - capture: ", opponent.capture)
	return s.analyzeScore(&current, &opponent, n, depth)
}
