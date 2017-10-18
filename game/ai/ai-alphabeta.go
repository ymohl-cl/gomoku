package ai

import (
	"fmt"
	"math"
	"time"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
)

var timeTotalNode time.Duration
var timeTotalRule time.Duration
var timeTotalAddNode time.Duration
var timeTotalCopyBoard time.Duration
var timeTotalState time.Duration

type AI struct {
	s *State
}

func New() *AI {
	return &AI{}
}

type State struct {
	nbCapsCurrent uint8
	nbCapsOther   uint8
	player        uint8
	pq            *Node
}

func (a AI) newState(nbCOldCurrent, nbCOldOther, p uint8) *State {
	s := State{
		nbCapsCurrent: nbCOldOther,
		nbCapsOther:   nbCOldCurrent,
		player:        p,
	}

	return &s
}

type Node struct {
	x      uint8
	y      uint8
	weight int8
	//	minV      int8
	//	maxV      int8
	nextState State
	next      *Node
}

func (a AI) newNode(x, y uint8) *Node {
	return &Node{
		//		minV: math.MinInt8,
		//		maxV: math.MaxInt8,
		x: x,
		y: y,
	}
}

func (n *Node) Clean(x, y uint8) {
	n.x = x
	n.y = y
	n.nextState.Clean()
	n.next = nil
}

func (s *State) Clean() {
	s.nbCapsCurrent = 0
	s.nbCapsOther = 0
	s.player = 0
	s.pq = nil
}

func (s *State) Set(nbCOldCurrent, nbCOldOther, p uint8) {
	s.nbCapsCurrent = nbCOldOther
	s.nbCapsOther = nbCOldCurrent
	s.player = p
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

func (s State) switchPlayer() uint8 {
	if s.player == ruler.TokenP1 {
		return ruler.TokenP2
	}
	return ruler.TokenP1
}

func (a AI) eval(s *State, stape uint8, r *ruler.Rules) int8 {
	var salt int8

	if s.player == ruler.TokenP1 {
		salt = 10
	} else {
		salt = 9
	}

	if r.IsWin {
		salt += 100
	} else {
		salt += int8((s.nbCapsCurrent + r.NbCaps) * (s.nbCapsCurrent + r.NbCaps))
	}

	salt += (int8(stape) - 4)

	if s.player == ruler.TokenP2 {
		return salt
	}
	return -salt
}

func (s *State) Search(n *Node) *Node {
	for test := s.pq; test != nil; test = test.next {
		if test.x == n.x && test.y == n.y {
			return test
		}
	}
	return nil
}

func (a *AI) applyMove(b *[][]uint8, r *ruler.Rules, player uint8, x, y uint8) {
	(*b)[y][x] = player
	if r.IsCaptured == true {
		for _, cap := range r.GetCaptures() {
			(*b)[cap.Y][cap.X] = ruler.TokenEmpty
		}
	}
}

func (a *AI) restoreMove(b *[][]uint8, r *ruler.Rules, player uint8, x, y uint8) {
	opponent := uint8(ruler.TokenP1)

	if player == ruler.TokenP1 {
		opponent = ruler.TokenP2
	}

	(*b)[y][x] = ruler.TokenEmpty
	if r.IsCaptured == true {
		for _, cap := range r.GetCaptures() {
			(*b)[cap.Y][cap.X] = opponent
		}
	}
}

func (a *AI) alphabeta(s *State, b *[][]uint8, alpha, beta int8, stape uint8) int8 {
	score := int8(-125)
	var node *Node

	if stape == 0 {
		return 0
	}

	for y := uint8(0); y < 19; y++ {
		for x := uint8(0); x < 19; x++ {
			r := ruler.New()
			r.CheckRules(b, int8(x), int8(y), s.player, s.nbCapsCurrent)
			if r.IsMoved {
				node = a.newNode(x, y)
				node.weight = -a.eval(s, stape, r)
				weight := node.weight

				if !r.IsWin || score > node.weight {
					// applymove and capture
					a.applyMove(b, r, s.player, x, y)
					node.nextState.Set(s.nbCapsCurrent+r.NbCaps, s.nbCapsOther, s.switchPlayer())
					node.weight = -a.alphabeta(&node.nextState, b, -beta, -alpha, stape-1)
					// restore move and capture
					a.restoreMove(b, r, s.player, x, y)

					if node.weight == 0 {
						node.weight = weight
					}
				}

				if score < node.weight {
					score = node.weight
				}

				s.addNode(node)
				//				score = node.weight
				if node.weight > alpha {
					alpha = node.weight
					if alpha >= beta {
						//					fmt.Println("stape: ", stape, " - Alpha: ", alpha, " - Beta: ", beta, " - Score: ", score)
						return score
					}
				}
			}
		}
	}
	return score
}

func (a *AI) getCoord(weight int8) (uint8, uint8) {
	var x, y uint8
	//	var refNode *Node

	tmp := int8(math.MinInt8)
	for node := a.s.pq; node != nil; node = node.next {
		weight := node.weight
		fmt.Print("Y: ", node.y, " - X: ", node.x, " | Weight Node: ", node.weight, " | ")
		if tmp <= weight {
			fmt.Println("Choiced !")
			x = node.x
			y = node.y
			tmp = node.weight
			//			refNode = node
		} else {
			fmt.Println("Not")
		}
	}

	//	a.s = &refNode.nextState
	a.s = nil
	return x, y
}

func (a *AI) Play(b *[][]uint8, s *database.Session, c chan uint8) {
	//	var workNode *Node

	if a.s == nil || a.s.pq == nil {
		a.s = a.newState(uint8(s.NbCaptureP1), uint8(s.NbCaptureP2), ruler.TokenP2)
		//		a.s.alpha = math.MinInt8
		//		a.s.beta = math.MaxInt8
	}
	//	workNode = a.newNode()
	//timeTotalNode = 0
	//	timeTotalAddNode = 0
	//timeTotalCopyBoard = 0
	//timeTotalRule = 0
	//	timeTotalState = 0

	//	fmt.Println("AI a.s alpha: ", a.s.alpha)
	//	fmt.Println("AI a.s beta: ", a.s.beta)
	//var test time.Duration
	//test += time.Since(time.Now())
	//a.prevalphabeta(a.s, nil, math.MinInt8, math.MaxInt8, 4)
	//	r := ruler.New()
	a.alphabeta(a.s, b, int8(-125), int8(125), 4)

	//	a.getLen(a.s)
	x, y := a.getCoord(0)
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
			*a.s = node.nextState
			//			fmt.Println("alpha opposant: ", a.s.alpha)
			//			fmt.Println("beta opposant: ", a.s.beta)
			return
		}
	}
	a.s = nil
	//	fmt.Println("FUCCCCK")
}

/*func (a *AI) min(x, y int8) int8 {
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
}*/
