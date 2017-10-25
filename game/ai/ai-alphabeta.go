package ai

import (
	"fmt"
	"math"
	"time"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
)

const (
	posCenter = 9
	maxDepth  = 4
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
	nbTCurrent    uint8
	nbTOther      uint8
	player        uint8
	pq            *Node
}

func (a AI) newState(nbCOldCurrent, nbCOldOther, nbTOldCurrent, nbTOldOther, p uint8) *State {
	s := State{
		nbCapsCurrent: nbCOldOther,
		nbCapsOther:   nbCOldCurrent,
		nbTCurrent:    nbTOldOther,
		nbTOther:      nbTOldCurrent,
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

func (s *State) Set(nbCOldCurrent, nbCOldOther, nbTOldCurrent, nbTOldOther, p uint8) {
	s.nbCapsCurrent = nbCOldOther
	s.nbCapsOther = nbCOldCurrent
	s.nbTCurrent = nbTOldCurrent
	s.nbTOther = nbTOldOther
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
func (a *AI) copyBoard(b *[][]uint8) *[][]uint8 {
	var newBoard [][]uint8
	for _, line := range *b {
		newLine := make([]uint8, 19)
		copy(newLine, line)
		newBoard = append(newBoard, newLine)
	}
	return &newBoard
}

func (s State) switchPlayer() uint8 {
	if s.player == ruler.TokenP1 {
		return ruler.TokenP2
	}
	return ruler.TokenP1
}

func (a AI) eval(s *State, stape uint8, r *ruler.Rules) int8 {
	var ret int8

	ret = math.MaxInt8 - (maxDepth - int8(stape))
	//ret -= (4 - int8(stape))

	//fmt.Println("Value tree: ", s.nbTOther, " : ", s.nbTCurrent)
	//	if r.NbThree == 0 && len(r.CapturableWin) == 0 {
	ret -= 1
	//} else
	if !r.IsWin {
		//if r.NbThree == 1 {
		//	ret -= 1
		//}
		//	ret = 100
		//ret += int8(s.nbCapsOther - s.nbCapsCurrent)
		ret -= int8(s.nbCapsCurrent - s.nbCapsOther)
		//	if r.NbThree != 1 && len(r.CapturableWin) != 0 {
		//	} else {
		//		ret -= 1
		//	}
	}
	/* else if !r.IsWin {
		//		fmt.Println("current: ", s.nbTCurrent, " other: ", s.nbTOther)
		//		if s.nbTCurrent != 0 && (s.nbTCurrent >= s.nbTOther) {
		//			ret += 10
		//		}
		//ret = 10
		//ret = math.MaxInt8 / 2
		//		ret -= int8(s.nbTCurrent - s.nbTOther)
	}
	*/
	//		valueCaps := 5 - uint8(s.nbCapsCurrent+r.NbCaps)
	//		valueCaps := 5 - uint8(s.nbCapsOther)
	//		valueLine := uint8(0)
	//		if (r.NbThree > 0 && r.NbToken == 3) || r.NbToken == 4 {
	//valueLine = 1
	//			ret -= 1
	//		}
	/* else {
		valueLine = 5 - r.NbToken
	}

	fmt.Print(" - ValueLine: ", valueLine)
	fmt.Print(" - ValueCaps: ", valueCaps)
	if valueLine < valueCaps {
		fmt.Println("ValueLine: ", valueLine)
		ret -= int8(valueLine)
	} else {
		ret -= int8(valueCaps)
	}*/

	//	}

	// if ruler.TokenP1 == s.player, TokenP2 playing.
	//	if s.player == ruler.TokenP2 {
	//		ret *= -1
	//	}
	return ret
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

func getPointsOpti(b *[][]uint8) ([]int8, []int8) {
	var yS []int8
	var xS []int8

	// checkDiagonalTopLeft
	for y, x := int8(0), int8(0); y != posCenter && x != posCenter; y, x = y+1, x+1 {
		r := ruler.New()
		r.CheckRules(b, x, y, ruler.TokenP2, 0)
		if r.IsMoved {
			yS = append(yS, y)
			xS = append(xS, x)
			break
		}
	}

	// checkDiagonalBotRight
	for y, x := int8(18), int8(18); y != posCenter && x != posCenter; y, x = y-1, x-1 {
		r := ruler.New()
		r.CheckRules(b, x, y, ruler.TokenP2, 0)
		if r.IsMoved {
			yS = append(yS, y)
			xS = append(xS, x)
			break
		}
	}

	// checkDiagonalBotLeft
	for y, x := int8(18), int8(0); y != posCenter && x != posCenter; y, x = y-1, x+1 {
		r := ruler.New()
		r.CheckRules(b, x, y, ruler.TokenP2, 0)
		if r.IsMoved {
			yS = append(yS, y)
			xS = append(xS, x)
			break
		}
	}

	// checkDiagonalTopRight
	for y, x := int8(0), int8(18); y != posCenter && x != posCenter; y, x = y+1, x-1 {
		r := ruler.New()
		r.CheckRules(b, x, y, ruler.TokenP2, 0)
		if r.IsMoved {
			yS = append(yS, y)
			xS = append(xS, x)
			break
		}
	}

	return yS, xS
}

func (a *AI) alreadyCheck(x, y int8, xS, yS []int8) bool {
	for index, y := range yS {
		if y == y && x == xS[index] {
			return true
		}
	}
	return false
}

/*
func (a *AI) preAlphaBeta(s *State, b *[][]uint8) {
	var yS []int8
	var xS []int8
	alphaS := []int8{int8(-125), int8(-125), int8(-125), int8(-125)}
	betaS := []int8{int8(125), int8(125), int8(125), int8(125)}
	var depth uint8

	//	var aA, aB, aC, aD int8
	//	var bA, bB, bC, bD int8
	depth = 4
	wg := new(sync.WaitGroup)
	yS, xS = getPointsOpti(b)

	for index, y := range yS {
		cBoard := a.copyBoard(b)
		Mnode := a.newNode(uint8(xS[index]), uint8(y))
		state := a.newState(0, 0, ruler.TokenP1)
		r := ruler.New()
		r.CheckRules(b, int8(Mnode.x), int8(Mnode.y), s.player, s.nbCapsCurrent)
		//node.weight = -a.eval(s, stape, r)
		a.applyMove(cBoard, r, s.player, Mnode.x, Mnode.y)

		Mnode.nextState.Set(s.nbCapsCurrent+r.NbCaps, s.nbCapsOther, s.switchPlayer())
		wg.Add(1)
		go func(i int, N *Node, wait *sync.WaitGroup) {
			defer wait.Done()
			defer fmt.Println("Finish one")
			N.weight = -a.alphabeta(&N.nextState, cBoard, &alphaS[i], &betaS[i], depth-1, state, r)
			a.s.addNode(N)
		}(index, Mnode, wg)
	}

	wg.Wait()

	weight := int8(math.MinInt8)
	var saveIndex int8
	fmt.Println("------------ View value -------------")
	for index, y := range yS {
		fmt.Println("Position Y: ", y, " - X: ", xS[index])
		fmt.Println("Alpha: ", alphaS[index], " - Beta: ", betaS[index])

		for n := a.s.pq; n != nil; n = n.next {
			if int8(n.y) == y && int8(n.x) == xS[index] && weight <= n.weight {
				fmt.Print("node Y: ", n.y, " - node X: ", n.x)
				fmt.Println(" - weight: ", n.weight)
				weight = n.weight
				saveIndex = int8(index)
			}
		}
	}
	fmt.Println("------------- End Values ------------")

	for y := uint8(0); y < 19; y++ {
		for x := uint8(0); x < 19; x++ {
			if a.alreadyCheck(int8(x), int8(y), xS, yS) {
				continue
			}
			Mnode := a.newNode(x, y)
			state := a.newState(0, 0, ruler.TokenP1)
			r := ruler.New()
			r.CheckRules(b, int8(Mnode.x), int8(Mnode.y), s.player, s.nbCapsCurrent)
			if r.IsMoved {
				a.applyMove(b, r, s.player, Mnode.x, Mnode.y)
				Mnode.nextState.Set(s.nbCapsCurrent+r.NbCaps, s.nbCapsOther, s.switchPlayer())
				Mnode.weight = -a.alphabeta(&Mnode.nextState, b, &alphaS[saveIndex], &betaS[saveIndex], depth-1, state, r)
				a.restoreMove(b, r, s.player, x, y)
				a.s.addNode(Mnode)
			}
		}
	}
}
*/
func (a *AI) alphabeta(s *State, b *[][]uint8, alpha, beta int8, stape uint8, oldRule *ruler.Rules) int8 {
	score := int8(math.MinInt8) + 1
	var node *Node

	if stape == 0 || oldRule.IsWin {
		return -a.eval(s, stape+1, oldRule)
	}

	for y := uint8(0); y < 19; y++ {
		for x := uint8(0); x < 19; x++ {
			r := ruler.New()
			r.CheckRules(b, int8(x), int8(y), s.player, s.nbCapsCurrent)
			if r.IsMoved {
				node = a.newNode(x, y)
				//node.weight = -a.eval(s, stape, r)
				a.applyMove(b, r, s.player, x, y)
				//	countToken := uint8(0)
				//	if r.IsCaptured {
				//			countToken = 2
				//		}
				node.nextState.Set(s.nbCapsCurrent+r.NbCaps, s.nbCapsOther, s.nbTCurrent+r.NbThree, s.nbTOther, s.switchPlayer())
				//*beta = -*beta
				//*alpha = -*alpha
				node.weight = -a.alphabeta(&node.nextState, b, -beta, -alpha, stape-1, r)
				a.restoreMove(b, r, s.player, x, y)
				s.addNode(node)

				//				fmt.Println("node-weight: ", node.weight, " score: ", score, " stape: ", stape)
				//				fmt.Println("Test begin -- alpha: ", alpha, " beta: ", beta)
				if score < node.weight {
					score = node.weight
				}

				if alpha < node.weight {
					alpha = node.weight
					if alpha >= beta {
						//						fmt.Println("alpha: ", alpha, " beta: ", beta)
						return score
					}
				}
				//	score = a.max(score, node.weight)
				//*alpha = a.max(*alpha, score)

				//					fmt.Println("stape: ", stape, " - Alpha: ", alpha, " - Beta: ", beta, " - Score: ", score)
				//			if c != nil {
				//				c <- node.weight
				//			}
			}
		}
	}
	//	if c != nil {
	//		c <- score
	//	}
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
		a.s = a.newState(uint8(s.NbCaptureP1), uint8(s.NbCaptureP2), 0, 0, ruler.TokenP2)
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
	//state := a.newState(0, 0, ruler.TokenP1)
	r := ruler.New()
	//a.preAlphaBeta(a.s, b)
	ret := a.alphabeta(a.s, b, math.MinInt8+1, math.MaxInt8, maxDepth, r)
	fmt.Println("Ret: ", ret)
	//	fmt.Println("min int8: ", math.MinInt8)
	//	fmt.Println("max int8: ", math.MaxInt8)

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
*/
func (a *AI) max(x, y int8) int8 {
	if x > y {
		return x
	}
	return y
}
