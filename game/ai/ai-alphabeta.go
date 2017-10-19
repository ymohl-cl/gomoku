package ai

import (
	"fmt"
	"math"
	"time"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
)

const (
	ExactBound = 0
	LowerBound = 1
	UpperBound = 2

	maxStape = 3
)

type AI struct {
	s *State
}

func New() *AI {
	return &AI{}
}

type State struct {
	board          *[][]uint8
	nbCapsCurrent  uint8
	nbCapsOther    uint8
	nbThreeCurrent uint8
	nbThreeOther   uint8
	player         uint8
	pq             *Node
}

func (a AI) newState(b *[][]uint8, nbCOldCurrent, nbCOldOther, nbTOldCurrent, nbTOldOther, p uint8) *State {
	s := State{
		board:          b,
		nbCapsCurrent:  nbCOldOther,
		nbCapsOther:    nbCOldCurrent,
		nbThreeCurrent: nbTOldOther,
		nbThreeOther:   nbTOldCurrent,
		player:         p,
	}

	return &s
}

type Node struct {
	r         ruler.Rules
	x         uint8
	y         uint8
	weight    int8
	nextState *State
	next      *Node
	flag      uint8
	value     int8
}

func (a AI) newNode(s *State, x, y uint8) (bool, *Node) {

	node := &Node{
		x: x,
		y: y,
	}

	if n := s.SearchNode(node); n != nil {
		return true, n
	}
	return false, node
}

func (s *State) addNode(n *Node) {
	if s.pq == nil {
		s.pq = n
	} else {
		tmp := s.pq
		s.pq = n
		s.pq.next = tmp
	}
}
func (a *AI) copyBoard(b [][]uint8) [][]uint8 {
	var newBoard [][]uint8
	for _, line := range b {
		newLine := make([]uint8, 19)
		copy(newLine, line)
		newBoard = append(newBoard, newLine)
	}
	return newBoard
}

func (a *AI) genBoard(b *[][]uint8, node *Node, player uint8) *[][]uint8 {
	var newBoard [][]uint8
	for _, line := range *b {
		newLine := make([]uint8, 19)
		copy(newLine, line)
		newBoard = append(newBoard, newLine)
	}
	newBoard[node.y][node.x] = player
	if node.r.IsCaptured == true {
		for _, cap := range node.r.GetCaptures() {
			newBoard[cap.Y][cap.X] = ruler.TokenEmpty
		}
	}
	return &newBoard
}

func (s State) switchPlayer() uint8 {
	if s.player == ruler.TokenP1 {
		return ruler.TokenP2
	}
	return ruler.TokenP1
}

func (a AI) eval(s *State, n *Node, stape int8) int8 {
	var ret int8

	if !n.r.IsWin {
		ret = 30
		ret += 4 * int8(s.nbCapsOther-s.nbCapsCurrent)
		ret += 2 * int8(s.nbThreeOther)
		ret -= (maxStape - stape)
	} else {
		ret = 100
	}

	if s.player == ruler.TokenP1 {
		return ret
	}
	return -ret
}

func (s *State) SearchNode(n *Node) *Node {
	if n == nil {
		return nil
	}
	for test := s.pq; test != nil; test = test.next {
		if test.x == n.x && test.y == n.y {
			return test
		}
	}
	return nil
}

var timeTotalNode time.Duration
var timeTotalRule time.Duration

//var timeTotalAddNode time.Duration
var timeTotalCopyBoard time.Duration

//var timeTotalState time.Duration

func (a *AI) alphabeta(s *State, alpha, beta int8, stape int8, workNode *Node) int8 {
	var best int8
	var val int8
	var node, n *Node
	var find bool
	//	var findNode *Node

	if stape == 0 || (workNode != nil && workNode.r.IsWin) {
		return a.eval(s, workNode, stape)
	}

	alphaOrg := alpha
	if n = s.SearchNode(workNode); n != nil {
		if n.flag == ExactBound {
			return n.value
		} else if n.flag == LowerBound {
			alpha = a.max(alpha, n.value)
		} else if n.flag == UpperBound {
			beta = a.min(beta, n.value)
		}
		if alpha >= beta {
			return n.value
		}
	}

	best = math.MinInt8
	var x uint8
	for y := uint8(0); y < 19; y++ {
		for x = 0; x < 19; x++ {
			find, node = a.newNode(s, x, y)
			if find == false {
				node.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
			}
			if node.r.IsMoved {
				if node.nextState == nil {
					node.nextState = a.newState(a.genBoard(s.board, node, s.player), s.nbCapsCurrent+node.r.NbCaps, s.nbCapsOther, s.nbThreeCurrent+node.r.NbThree, s.nbThreeOther, s.switchPlayer())
				}
				val = -a.alphabeta(node.nextState, -beta, -alpha, stape-1, node)
				node.weight = val
				if val > best {
					best = val
					if best > alpha {
						alpha = best
						if alpha >= beta {
							break
						}
					}
				}
				if find == false {
					s.addNode(node)
				}
			}
		}
	}
	if n = s.SearchNode(workNode); n != nil {
		workNode = n
	}
	// Transposition Table Store; node is the lookup key for ttEntry
	workNode.value = best
	if best <= alphaOrg {
		workNode.flag = UpperBound
	} else if best >= beta {
		workNode.flag = LowerBound
	} else {
		workNode.flag = ExactBound
	}
	//ttEntry.depth := depth
	//TranspositionTableStore( node, ttEntry )

	return best
}

func (a *AI) getCoord(weight int8) (uint8, uint8) {
	var x, y uint8
	//var refNode *Node

	//	a.s = &a.s.pq.nextState
	//	var tmp int8
	//tmp := int8(math.MinInt8)
	for node := a.s.pq; node != nil; node = node.next {
		fmt.Println("x : ", node.x, " - y : ", node.y, "weight: ", node.weight, " - node value", node.value)
		//weight := node.weight
		//		if weight < 0 {
		//		weight *= -1
		//		}
		if node.weight == weight {
			x = node.x
			y = node.y
			//tmp = node.weight
			//refNode = node
			//tmp = node.weight
			//refNode = node
		}
	}

	/*	if mi != tmp {
		fmt.Println("Not first")
	}*/
	a.s = nil
	//	a.s.pq = nil
	//fmt.Println("----- Getting Coord -----")
	//fmt.Println("a.s alpha: ", a.s.alpha)
	//fmt.Println("a.s beta: ", a.s.beta)
	//refNode.nextState.alpha = a.s.alpha
	//*a.s = refNode.nextState
	//if a.s.alpha != refNode.nextState.alpha || a.s.beta != refNode.nextState.beta {
	//	fmt.Println("Error on herit nextState on getCoord")
	//}
	//fmt.Println("Switch to nextState (opponent) ")
	//fmt.Println("a.s alpha: ", a.s.alpha)
	//fmt.Println("a.s beta: ", a.s.beta)
	return x, y
}

func (a *AI) mtdf(s *State, workNode *Node, f, d int8) int8 {

	g := int8(f)
	upperbound := int8(100)
	lowerbound := int8(-100)
	beta := int8(0)
	for lowerbound < upperbound {
		if g == lowerbound {
			beta = g + 1
		} else {
			beta = g
		}
		g = a.alphabeta(s, beta-1, beta, d, workNode)
		if g < beta {
			upperbound = g
		} else {
			lowerbound = g
		}
	}
	return g
}

func (a *AI) iterative_deepening(s *State, workNode *Node) int8 {

	timeStart := time.Now()

	firstguess := int8(0)
	for d := int8(0); d < maxStape; d++ {
		firstguess = a.mtdf(s, workNode, firstguess, d)
		fmt.Println("Time : ", time.Since(timeStart), "- Stape : ", d)

		if time.Since(timeStart) >= time.Millisecond*100 {
			return firstguess
		}
		//if times_up() then break;
	}
	return firstguess
}

func (a *AI) Play(b *[][]uint8, s *database.Session, c chan uint8) {
	//var workNode *Node

	if a.s == nil || a.s.pq == nil {
		a.s = a.newState(b, uint8(s.NbCaptureP1), uint8(s.NbCaptureP2), 0, 0, ruler.TokenP2)
		//a.s.alpha = math.MinInt8
		//a.s.beta = math.MaxInt8
	}
	_, n := a.newNode(a.s, 0, 0)
	timeTotalNode = 0
	//	timeTotalAddNode = 0
	timeTotalCopyBoard = 0
	timeTotalRule = 0
	//	timeTotalState = 0

	//	fmt.Println("AI a.s alpha: ", a.s.alpha)
	//	fmt.Println("AI a.s beta: ", a.s.beta)
	var test time.Duration
	test += time.Since(time.Now())
	//a.prevalphabeta(a.s, nil, math.MinInt8, math.MaxInt8, 4)

	w := a.iterative_deepening(a.s, n)
	fmt.Println("return iterative : ", w)
	//a.alphabeta(a.s, -100, 100, maxStape, n)

	//	a.getLen(a.s)
	x, y := a.getCoord(w)
	//	fmt.Println("------------- End -------------")
	//	fmt.Println("Time total to clean node: ", timeTotalNode)
	//	fmt.Println("Time total to add node on list: ", timeTotalAddNode)
	//	fmt.Println("Time total to copy board, node and state: ", timeTotalCopyBoard)
	//	fmt.Println("Time total to check Rule: ", timeTotalRule)
	//	fmt.Println("Time total to create new state: ", timeTotalState)

	c <- y
	c <- x
}

func (a *AI) PlayOpposing(y, x uint8) {
	if a.s == nil || a.s.pq == nil {
		return
	}
	for node := a.s.pq; node != nil; node = node.next {
		if y == node.y && x == node.x {
			//			fmt.Println("------------------ START ------------------")
			//node.nextState.alpha = a.s.alpha
			a.s = node.nextState
			//			fmt.Println("alpha opposant: ", a.s.alpha)
			//			fmt.Println("beta opposant: ", a.s.beta)
			return
		}
	}
	a.s = nil
	//	fmt.Println("FUCCCCK")
}

func (a *AI) min(x, y int8) int8 {
	if x < y {
		return x
	}
	return y
}

func (a *AI) max(x, y int8) int8 {
	if x > y {
		return x
	}
	return y
}
