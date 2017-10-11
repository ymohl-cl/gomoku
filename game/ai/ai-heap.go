package ai

import (
	"fmt"
	"math"
	"sync"

	"github.com/ymohl-cl/gomoku/game/ruler"
)

type AI struct {
	s *State
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

type Node struct {
	r         *ruler.Rules
	x         uint8
	y         uint8
	weight    int8
	nextState *State
	next      *Node
}

func New() *AI {
	return &AI{}
}

func (s *State) addNode(n *Node) {
	var tmp *Node

	if n == nil {
		return
	}
	tmp = s.pq
	if tmp != nil {
		tmp.next = n
	} else {
		tmp = n
	}
	s.pq = tmp
}

func (s State) switchPlayer() uint8 {
	if s.player == ruler.TokenP1 {
		return ruler.TokenP2
	} else {
		return ruler.TokenP1
	}
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

func (a AI) newNode(x, y uint8) *Node {
	return &Node{
		r: ruler.New(),
		x: x,
		y: y,
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

func (a *AI) getCoord(weight int8) (uint8, uint8) {
	var x, y uint8

	for node := a.s.pq; node != nil; node = node.next {
		if weight == node.weight {
			x = node.x
			y = node.y
			a.s = nil
			break
		}
	}
	return x, y
}

func (a *AI) browseThree(s *State, stape int8) int8 {
	var y, x uint8
	wg := new(sync.WaitGroup)
	c := make(chan *Node, 19*19)

	for y = 0; y < 19; y++ {
		for x = 0; x < 19; x++ {
			wg.Add(1)
			go func(xi, yi uint8) {
				defer wg.Done()
				node := a.newNode(xi, yi)
				node.r.CheckRules(s.board, int8(xi), int8(yi), s.player, s.nbCapsCurrent)
				if node.r.IsMoved {
					node.weight = a.eval(s, node, stape) // - ou + stape ?
					c <- node
				} else {
					c <- nil
				}
			}(x, y)
		}
	}

	wg.Wait()
	close(c)
	for n := range c {
		s.addNode(n)
	}
	a.recastThree(s, stape)

	if s.player == ruler.TokenP1 {
		return s.min()
	}
	return s.max()
}

func (a *AI) recastThree(s *State, stape int8) {
	var tmp int8
	wg := new(sync.WaitGroup)

	if stape >= 10 {
		return
	}

	if s.player == ruler.TokenP1 {
		tmp = s.min()
	} else {
		tmp = s.max()
	}

	for node := s.pq; node != nil; node = node.next {
		if tmp == node.weight {
			wg.Add(1)
			go func(n *Node) {
				defer wg.Done()
				b := a.copyBoard(s.board)
				b[n.y][n.x] = s.player
				n.nextState = a.newState(b, s.nbCapsCurrent+n.r.NbCaps, s.nbCapsOther, s.lenMaxP1, s.lenMaxP2, s.switchPlayer())
				n.weight = a.browseThree(n.nextState, stape+1)
			}(node)
		}
	}
	wg.Wait()
}

func (a AI) eval(s *State, n *Node, stape int8) int8 {

	if s.player == ruler.TokenP1 {
		if !n.r.IsWin {
			return 0 - stape
		}
		return math.MinInt8 + stape
	}

	if !n.r.IsWin {
		return 0 + stape
	}
	return math.MaxInt8 - stape
}

func (s State) min() int8 {
	var tmp int8

	tmp = math.MaxInt8
	for node := s.pq; node != nil; node = node.next {
		if tmp > node.weight {
			tmp = node.weight
		}
	}
	return tmp
}

func (s State) max() int8 {
	var tmp int8

	tmp = math.MinInt8
	for node := s.pq; node != nil; node = node.next {
		if tmp < node.weight {
			tmp = node.weight
		}
	}
	return tmp
}

// Play start the game
func (a *AI) Play(b [][]uint8, c chan uint8) {
	if a.s == nil {
		a.s = a.newState(b, 0, 0, 0, 0, ruler.TokenP2)
	}
	weight := a.browseThree(a.s, int8(0))
	fmt.Println("Weight to play: ", weight)
	//	a.getLen(a.s)
	x, y := a.getCoord(weight)
	c <- y
	c <- x
}

/*
func (a AI) PlayOpponent(x, y uint8) {

}*/
func (a AI) getLen(s *State) {
	fmt.Println("Len of node: ", s.getLen())
	for node := s.pq; node != nil; node = node.next {
		a.getLen(node.nextState)
	}
}

func (s State) getLen() int {
	var i int
	for node := s.pq; node != nil; node = node.next {
		i++
	}
	return i
}
