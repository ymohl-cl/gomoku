package alphabeta

import (
	"fmt"
	"math"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
)

type InfoPlayer struct {
	totalCapture uint8
}

type State struct {
	maxDepth      uint8
	board         *[19][19]uint8
	currentPlayer uint8
	infoP1        InfoPlayer
	infoP2        InfoPlayer
	save          *Node
	lst           *Node
}

type Node struct {
	rule   ruler.Rules
	weight int8
	next   *Node
	prev   *Node
}

func New(b *[19][19]uint8, player uint8) *State {
	return &State{maxDepth: 4, board: b, currentPlayer: player}
}

func (s *State) getTotalCapture(player uint8) uint8 {
	if player == ruler.Player1 {
		return s.infoP1.totalCapture
	}
	return s.infoP2.totalCapture
}

func (s *State) addTotalCapture(player uint8, number uint8) {
	if player == ruler.Player1 {
		s.infoP1.totalCapture += number
	} else {
		s.infoP2.totalCapture += number
	}
}

func (s *State) subTotalCapture(player uint8, number uint8) {
	if player == ruler.Player1 {
		s.infoP1.totalCapture -= number
	} else {
		s.infoP2.totalCapture -= number
	}
}

func (s *State) newNode(y, x int8) *Node {
	n := new(Node)
	n.rule.Init(s.currentPlayer, y, x)
	n.rule.CheckRules(s.board, s.getTotalCapture(s.currentPlayer))
	if !n.rule.Movable {
		n = nil
	}
	return n
}

func (s *State) updateData(n *Node, prev *Node) {
	n.rule.ApplyMove(s.board)

	s.addTotalCapture(n.rule.GetPlayer(), n.rule.NumberCapture)
	s.currentPlayer = ruler.GetOtherPlayer(n.rule.GetPlayer())

	n.prev = prev
}

func (s *State) restoreData(n *Node, prev *Node) {
	n.rule.RestoreMove(s.board)

	s.subTotalCapture(n.rule.GetPlayer(), n.rule.NumberCapture)
	s.currentPlayer = ruler.GetOtherPlayer(s.currentPlayer)

	n.prev = nil
}

func (s *State) addNode(n *Node) {
	if s.lst == nil {
		s.lst = n
		return
	}
	n.next = s.lst
	s.lst = n
}

func (s *State) alphabetaNegaScout(alpha, beta int8, depth uint8, n *Node) int8 {
	if (n != nil && n.rule.Win) || depth == 0 {
		return s.eval(n, depth)
	}

	first := true
	for y := int8(0); y < 19; y++ {
		for x := int8(0); x < 19; x++ {
			node := s.newNode(y, x)
			if node == nil {
				continue
			}
			// apply move and update data
			s.updateData(node, n)

			if first == true {
				first = false
				node.weight = -s.alphabetaNegaScout(-beta, -alpha, depth-1, node)
			} else {
				node.weight = -s.alphabetaNegaScout(-alpha-1, -alpha, depth-1, node)
				if alpha < node.weight && node.weight < beta {
					node.weight = -s.alphabetaNegaScout(-beta, -node.weight, depth-1, node)
				}
			}

			//fmt.Print("Depth: ", depth, " - node y && x: ", y, " - ", x)
			//fmt.Print(" - weight: ", node.weight)
			//fmt.Println(" | Alpha: ", alpha, " Beta: ", beta)

			// restore move and restore data
			s.restoreData(node, n)

			if depth == s.maxDepth {
				s.addNode(node)
			}

			alpha = maxWeight(alpha, node.weight)

			if alpha >= beta {
				return alpha
			}
		}
	}
	return alpha
}

// Play start the alphabeta algorythm
func Play(b *[19][19]uint8, s *database.Session, c chan uint8) {
	state := New(b, ruler.Player2)
	state.addTotalCapture(ruler.Player1, uint8(s.NbCaptureP1))
	state.addTotalCapture(ruler.Player2, uint8(s.NbCaptureP2))

	//ret := state.alphabetaNegaScout(math.MinInt8+1, math.MaxInt8, state.maxDepth, nil)
	state.alphabetaNegaScout(math.MinInt8+1, math.MaxInt8, state.maxDepth, nil)
	//fmt.Println("ret: ", ret)

	tmp := int8(math.MinInt8)

	for n := state.lst; n != nil; n = n.next {
		y, x := n.rule.GetPosition()
		fmt.Println("Node weight: ", n.weight, " y et x: ", y, " - ", x)
		if tmp <= n.weight {
			tmp = n.weight
			state.save = n
		}
	}

	y, x := state.save.rule.GetPosition()
	c <- uint8(y)
	c <- uint8(x)
}
