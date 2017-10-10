package ai

import (
	"sync"

	"github.com/ymohl-cl/gomoku/game/ruler"
)

type AI struct {
	node *StateMove
}

func New() *AI {
	return &AI{node: nil}
}

type StateMove struct {
	b             *[][]uint8
	r             *ruler.Rules
	nodes         []*StateMove
	posX          uint8
	posY          uint8
	currentPlayer uint8
	nbCapsP1      int32
	nbCapsP2      int32
}

func newNode(board *[][]uint8, rule *ruler.Rules, x, y uint8, player uint8, nbCapsP1, nbCapsP2 int32) *StateMove {
	st := &StateMove{
		b:             board,
		r:             rule.Copy(),
		currentPlayer: player,
		nbCapsP1:      nbCapsP1,
		nbCapsP2:      nbCapsP2,
	}
	return st
}

func (s StateMove) switchPlayer() uint8 {
	if s.currentPlayer == ruler.TokenP1 {
		return ruler.TokenP2
	} else {
		return ruler.TokenP1
	}
}

func (a AI) Play(b [][]uint8, c chan uint8) {
	var x, y uint8
	wg := new(sync.WaitGroup)

	if a.node == nil {
		node := newNode(b, nil, x, y, ruler.TokenP2, 0, 0)
	}
	a.makeThree(a.node, 0, wg)

	wg.Wait()

	x, y = a.getCoord()
	c <- y
	c <- x
}

func (a AI) getCoord() (uint8, uint8) {
	return 10, 9
	// HAHA
}

func (a *AI) copyBoard(b *[][]uint8, x, y uint8, player uint8) *[][]uint8 {
	var newBoard [][]uint8
	for _, line := range *b {
		newLine := make([]uint8, 19)
		copy(newLine, line)
		newBoard = append(newBoard, newLine)
	}
	newBoard[y][x] = player
	return &newBoard
}

func (a *AI) makeThree(node *StateMove, state int, wg *sync.WaitGroup) {
	var x, y uint8
	var nbCapsCurrent *int32
	var nbCapsTmp int32

	wg.Add(1)
	defer wg.Done()

	if state == 10 {
		return
	}

	if node.currentPlayer == ruler.TokenP1 {
		nbCapsCurrent = &node.nbCapsP1
	} else {
		nbCapsCurrent = &node.nbCapsP2
	}

	rule := ruler.New()
	for y = 0; y < 19; y++ {
		for x = 0; x < 19; x++ {
			nbCapsTmp = *nbCapsCurrent
			rule.CheckRules(node.b, int8(x), int8(y), node.currentPlayer, &nbCapsTmp)
			if rule.IsMoved {
				node := newNode(b, rule, x, y, currentPlayer, nbCaps)
				nodes = append(nodes, node)
				if !rule.IsWin {
					go a.makeThree(a.copyBoard(b, x, y, currentPlayer), state+1, node.switchPlayer(), node.nodes, node.nbCaps, wg)
				}
			}
			rule.Clean()
		}
	}

	return
}
