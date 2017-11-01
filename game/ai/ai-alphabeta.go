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
	s        *State
	NbPlayed uint8
	alpha    int8
	beta     int8
}

func New() *AI {
	return &AI{NbPlayed: 1, alpha: math.MinInt8 + 1, beta: math.MaxInt8}
}

type State struct {
	nbCapsCurrent  int8
	nbCapsOther    int8
	nbTCurrent     int8
	nbTOther       int8
	nbAlignCurrent int8
	nbAlignOther   int8
	player         uint8
	alpha          int8
	beta           int8
	pq             *Node
}

func (a AI) newState(nbCOldCurrent, nbCOldOther, nbTOldCurrent, nbTOldOther, nbAlignOldCurrent, nbAlignOldOther int8, p uint8) *State {
	s := State{
		nbCapsCurrent:  nbCOldOther,
		nbCapsOther:    nbCOldCurrent,
		nbTCurrent:     nbTOldOther,
		nbTOther:       nbTOldCurrent,
		nbAlignCurrent: nbAlignOldOther,
		nbAlignOther:   nbAlignOldCurrent,
		alpha:          math.MinInt8 + 1,
		beta:           math.MaxInt8,
		player:         p,
	}

	return &s
}

type Node struct {
	x         uint8
	y         uint8
	weight    int8
	rule      ruler.Rules
	alpha     int8
	beta      int8
	nextState State
	next      *Node
}

func (a AI) newNode(x, y uint8) *Node {
	return &Node{
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

func (s *State) Set(nbCOldCurrent, nbCOldOther, nbTOldCurrent, nbTOldOther, nbAlignOldCurrent, nbAlignOldOther int8, p uint8) {
	s.nbCapsCurrent = nbCOldOther
	s.nbCapsOther = nbCOldCurrent
	s.nbTCurrent = nbTOldOther
	s.nbTOther = nbTOldCurrent
	tmpNbAlignCurrent := s.nbAlignCurrent
	s.nbAlignCurrent = s.max(s.nbAlignOther, nbAlignOldOther)
	s.nbAlignOther = s.max(tmpNbAlignCurrent, nbAlignOldCurrent)
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

func (a AI) eval(s *State, stape uint8, r *ruler.Rules, base *State, prevRule *ruler.Rules) int8 {
	ret := int8(63)

	// Var Caps
	var basicCurrent int8
	var basicOther int8
	if s.player == ruler.TokenP2 {
		//basicCurrent = int8(s.nbCapsCurrent - base.nbCapsOther)
		basicCurrent = int8(s.nbCapsOther - base.nbCapsOther)
		//basicOther = int8(s.nbCapsOther - base.nbCapsCurrent)
		basicOther = int8(s.nbCapsCurrent - base.nbCapsCurrent)
	} else {
		basicCurrent = int8(s.nbCapsOther - base.nbCapsCurrent)
		basicOther = int8(s.nbCapsCurrent - base.nbCapsOther)
	}

	//Win
	if r.IsWin {
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

	if len(r.CapturableWin) > 0 {
		ret += int8(treeCurrent - treeOther*3)
	} else {
		ret += int8(treeCurrent*3 - treeOther)
	}

	//Align
	AlignCurrent := s.nbAlignOther
	AlignOther := s.nbAlignCurrent

	ret += int8(AlignCurrent*4 - AlignOther)

	return -ret
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
func (s *State) getNode(x, y uint8) *Node {
	for n := s.pq; n != nil; n = n.next {
		if n.x == x && n.y == y {
			return n
		}
	}
	return nil
}

func (a *AI) alphabeta(s *State, b *[][]uint8, alpha, beta int8, stape uint8, oldRule *ruler.Rules, base *State, prevRule *ruler.Rules) int8 {
	score := int8(math.MinInt8) + 1

	if stape == 0 || oldRule.IsWin { // || (oldRule.NbThree > 0 && len(oldRule.CapturableWin) == 0) {
		ret := a.eval(s, stape+1, oldRule, base, prevRule)
		return ret
	}

	for y := uint8(0); y < 19; y++ {
		for x := uint8(0); x < 19; x++ {
			found := true
			node := s.getNode(x, y)
			if node == nil {
				node = a.newNode(x, y)
				found = false
				node.rule.CheckRules(b, int8(x), int8(y), s.player, uint8(s.nbCapsCurrent))
			}
			if node.rule.IsMoved {
				a.applyMove(b, &node.rule, s.player, x, y)
				if found == false {
					node.nextState.Set(s.nbCapsCurrent+int8(node.rule.NbCaps), s.nbCapsOther, s.nbTCurrent+int8(node.rule.NbThree), s.nbTOther, int8(node.rule.NbToken), s.nbAlignOther, s.switchPlayer())
				} else {
					node.nextState.nbTCurrent = s.nbTOther
					node.nextState.nbTOther = s.nbTCurrent + int8(node.rule.NbThree)
				}
				node.weight = -a.alphabeta(&node.nextState, b, -beta, -alpha, stape-1, &node.rule, base, oldRule)
				a.restoreMove(b, &node.rule, s.player, x, y)
				if found == false {
					s.addNode(node)
				}

				if score < node.weight {
					score = node.weight
				}

				s.alpha = alpha
				//				node.beta = beta
				if alpha < node.weight {
					alpha = node.weight
					s.alpha = alpha
					if alpha >= beta {
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
	var refNode *Node

	tmp := int8(math.MinInt8)
	for node := a.s.pq; node != nil; node = node.next {
		weight := node.weight
		if tmp <= weight {
			x = node.x
			y = node.y
			tmp = node.weight
			refNode = node
		}
	}

	a.s = &refNode.nextState
	return x, y
}

func (a *AI) FirstPass(s *State, b *[][]uint8) {
	for y := 0; y < 19; y++ {
		for x := 0; x < 19; x++ {
			r := ruler.New()
			r.CheckRules(b, int8(x), int8(y), s.player, uint8(s.nbCapsCurrent))
			if r.IsMoved {
				node := a.newNode(uint8(x), uint8(y))
				if r.IsWin {
					node.weight = 127
					a.s.addNode(node)
					return
				} else {
					r := ruler.New()
					r.CheckRules(b, int8(x), int8(y), s.switchPlayer(), uint8(s.nbCapsCurrent))
					if r.IsWin {
						node.weight = 126
						a.s.addNode(node)
					}
				}
			}
		}
	}
}

func (a *AI) updateFirstPass(x, y uint8, save *Node) {
	for n := save; n != nil; n = n.next {
		if n.x == x && n.y == y {
			a.s = &n.nextState
			return
		}
	}
	a.s.pq = nil
}

func (a *AI) Play(b *[][]uint8, s *database.Session, c chan uint8) {
	var x, y uint8

	a.NbPlayed++
	if a.s == nil || a.s.pq == nil {
		a.s = a.newState(int8(s.NbCaptureP1), int8(s.NbCaptureP2), 0, 0, 0, 0, ruler.TokenP2)
	}

	//	save := a.s.pq
	//	a.s.pq = nil
	//	a.FirstPass(a.s, b)
	//	if a.s.pq != nil {
	//		x, y = a.getCoord(0)
	//		a.updateFirstPass(x, y, save)
	//	} else {
	a.s.nbTCurrent = 0
	a.s.nbTOther = 0
	//	a.s.pq = save
	r := ruler.New()
	prevRule := ruler.New()

	if a.NbPlayed > maxDepth {
		a.NbPlayed = maxDepth
	}
	fmt.Println("Depth: ", a.NbPlayed, " - alpha: ", -a.beta, " - beta: ", a.beta)
	ret := a.alphabeta(a.s, b, math.MinInt8+1, a.beta, maxDepth, r, a.s, prevRule)
	fmt.Println("Ret: ", ret)
	x, y = a.getCoord(0)
	//	}

	c <- y
	c <- x
}

func (a *AI) PlayOpposing(y, x uint8) {
	if a.s == nil || a.s.pq == nil {
		return
	}
	for node := a.s.pq; node != nil; node = node.next {
		if y == node.y && x == node.x {
			a.s = &node.nextState
			a.alpha = a.s.alpha
			a.beta = node.weight + 2
			return
		}
	}
	a.s = nil
	a.alpha = math.MinInt8 + 1
	a.beta = math.MaxInt8
}

func (s *State) max(x, y int8) int8 {
	if x > y {
		return x
	}
	return y
}
