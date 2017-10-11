package ai

import (
	"fmt"
	"math"
	"sync"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
)

type AI struct {
	s *State
}

func New() *AI {
	return &AI{}
}

type State struct {
	board         [][]uint8
	nbCapsCurrent uint8
	nbCapsOther   uint8
	lenMaxP1      uint8
	lenMaxP2      uint8
	player        uint8
	pq            *Node
}

func (a AI) newState(b [][]uint8, nbCOldCurrent, nbCOldOther, lMaxP1, lMaxP2, p uint8) *State {
	s := State{
		board:         a.copyBoard(b),
		nbCapsCurrent: nbCOldOther,
		nbCapsOther:   nbCOldCurrent,
		lenMaxP1:      lMaxP1,
		lenMaxP2:      lMaxP2,
		player:        p,
	}

	return &s
}

type Node struct {
	r         *ruler.Rules
	x         uint8
	y         uint8
	weight    int8
	nextState *State
	next      *Node
}

func (a AI) newNode(x, y uint8) *Node {
	return &Node{
		r: ruler.New(),
		x: x,
		y: y,
	}
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

func (a *AI) genBoard(b [][]uint8, node *Node, player uint8) [][]uint8 {
	var newBoard [][]uint8
	for _, line := range b {
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
	return newBoard
}

func (s State) switchPlayer() uint8 {
	if s.player == ruler.TokenP1 {
		return ruler.TokenP2
	}
	return ruler.TokenP1
}

func (a AI) eval(s *State, n *Node, stape int8) int8 {
	var ret int8

	fmt.Println("NbThree: ", n.r.NbThree)
	if !n.r.IsWin {
		ret = 30
		ret += int8(s.nbCapsOther) - int8(s.nbCapsCurrent)
	} else {
		ret = math.MaxInt8 - stape
	}

	if s.player == ruler.TokenP1 {
		return ret
	}
	return -ret
}

func (a *AI) prevalphabeta(s *State, prevnode *Node, alpha, beta int8, stape int8) {
	var y, x uint8
	var Val int8
	wg := new(sync.WaitGroup)

	Val = math.MinInt8
	for y = 0; y < 19; y++ {
		for x = 0; x < 19; x++ { /*pour tout enfant si de s faire */
			wg.Add(1)
			go func(xi, yi uint8) {
				defer wg.Done()
				node := a.newNode(xi, yi)
				node.r.CheckRules(s.board, int8(xi), int8(yi), s.player, s.nbCapsCurrent)
				if node.r.IsMoved {
					s.addNode(node)
					b := a.genBoard(s.board, node, s.player)
					node.nextState = a.newState(b, s.nbCapsCurrent+node.r.NbCaps, s.nbCapsOther, s.lenMaxP1, s.lenMaxP2, s.switchPlayer())
					node.weight = a.max(Val, a.alphabeta(node.nextState, node, alpha, beta, stape-1))
				}
			}(x, y)
		}
	}
	wg.Wait()
}

/* alpha est toujours inférieur à beta */
func (a *AI) alphabeta(s *State, prevnode *Node, alpha, beta int8, stape int8) int8 {
	var y, x uint8
	var Val int8

	if stape <= 0 || (prevnode != nil && prevnode.r.IsWin) {
		return a.eval(s, prevnode, stape)
	}
	if s.player == ruler.TokenP2 {
		Val = math.MinInt8
		for y = 0; y < 19; y++ {
			for x = 0; x < 19; x++ { /*pour tout enfant si de s faire */
				node := a.newNode(x, y)
				node.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
				if node.r.IsMoved {
					s.addNode(node)
					b := a.genBoard(s.board, node, s.player)
					node.nextState = a.newState(b, s.nbCapsCurrent+node.r.NbCaps, s.nbCapsOther, s.lenMaxP1, s.lenMaxP2, s.switchPlayer())
					Val = a.max(Val, a.alphabeta(node.nextState, node, alpha, beta, stape-1))
					node.weight = Val
					alpha = a.max(alpha, Val)
					if beta <= Val {
						return Val
					}
				}
			}
		}
	} else {
		Val = math.MaxInt8
		for y = 0; y < 19; y++ {
			for x = 0; x < 19; x++ { /*pour tout enfant si de s faire */
				node := a.newNode(x, y)
				node.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
				if node.r.IsMoved {
					s.addNode(node)
					b := a.genBoard(s.board, node, s.player)
					node.nextState = a.newState(b, s.nbCapsCurrent+node.r.NbCaps, s.nbCapsOther, s.lenMaxP1, s.lenMaxP2, s.switchPlayer())
					Val = a.min(Val, a.alphabeta(node.nextState, node, alpha, beta, stape-1))
					node.weight = Val
					beta = a.min(beta, Val)
					if Val <= alpha {
						return Val
					}
				}
			}
		}
	}
	return Val
}

func (a *AI) getCoord(weight int8) (uint8, uint8) {
	var x, y uint8

	println("weight : ", weight)

	var tmp int8
	tmp = math.MinInt8
	for node := a.s.pq; node != nil; node = node.next {
		println("node.weight : ", node.weight)
		if tmp <= node.weight {
			x = node.x
			y = node.y
			tmp = node.weight
		}
	}
	a.s = nil
	return x, y
}

func (a *AI) Play(b [][]uint8, s *database.Session, c chan uint8) {
	if a.s == nil {
		a.s = a.newState(b, uint8(s.NbCaptureP1), uint8(s.NbCaptureP2), 0, 0, ruler.TokenP2)
	}
	fmt.Println("LETS GO")
	a.prevalphabeta(a.s, nil, math.MinInt8, math.MaxInt8, 4)
	//	a.getLen(a.s)
	x, y := a.getCoord(0)
	c <- y
	c <- x
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
