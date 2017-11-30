package alphabeta

import (
	"fmt"
	"math"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

// InfoPlayer of previous evaluation
type InfoPlayer struct {
	totalCapture uint8
}

// State of evaluation
type State struct {
	maxDepth      uint8
	board         *[19][19]uint8
	currentPlayer uint8
	infoP1        InfoPlayer
	infoP2        InfoPlayer
	save          *Node
	lst           *Node
}

// Node represent one move
type Node struct {
	rule   ruler.Rules
	weight int16
	next   *Node
	prev   *Node
}

// New provide a State
func New(b *[19][19]uint8, player uint8) *State {
	return &State{maxDepth: 4, board: b, currentPlayer: player}
}

func (s *State) getTotalCapture(player uint8) uint8 {
	if player == rdef.Player1 {
		return s.infoP1.totalCapture
	}
	return s.infoP2.totalCapture
}

func (s *State) addTotalCapture(player uint8, number uint8) {
	if player == rdef.Player1 {
		s.infoP1.totalCapture += number
	} else {
		s.infoP2.totalCapture += number
	}
}

func (s *State) subTotalCapture(player uint8, number uint8) {
	if player == rdef.Player1 {
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

func (s *State) updateTokenPlayer(nodes *[5]*Node) {
	var n *Node

	for i := 0; i < 5; i++ {
		if n = (*nodes)[i]; n == nil {
			break
		}
		y, x := n.rule.GetPosition()
		(*s.board)[y][x] = rdef.Empty
	}
	/*
		for _, n := range *nodes {
			y, x := n.rule.GetPosition()
			(*s.board)[y][x] = rdef.Empty
		}
	*/
}

func (s *State) restoreTokenPlayer(nodes *[5]*Node) {
	var n *Node

	for i := 0; i < 5; i++ {
		if n = (*nodes)[i]; n == nil {
			break
		}
		y, x := n.rule.GetPosition()
		(*s.board)[y][x] = n.rule.GetPlayer()
	}
}

func (s *State) updateData(n *Node, prev *Node) {
	n.rule.ApplyMove(s.board)

	s.addTotalCapture(n.rule.GetPlayer(), n.rule.NumberCapture)
	s.currentPlayer = rdef.GetOtherPlayer(n.rule.GetPlayer())

	n.prev = prev
}

func (s *State) restoreData(n *Node, prev *Node) {
	n.rule.RestoreMove(s.board)

	s.subTotalCapture(n.rule.GetPlayer(), n.rule.NumberCapture)
	s.currentPlayer = rdef.GetOtherPlayer(s.currentPlayer)

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

func (s *State) alphabetaNegaScout(alpha, beta int16, depth uint8, n *Node) int16 {
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

			//node.weight = -s.alphabetaNegaScout(-beta, -alpha, depth-1, node)

			if first == true {
				first = false
				node.weight = -s.alphabetaNegaScout(-beta, -alpha, depth-1, node)
			} else {
				node.weight = -s.alphabetaNegaScout(-alpha-1, -alpha, depth-1, node)
				if alpha < node.weight && node.weight < beta {
					node.weight = -s.alphabetaNegaScout(-beta, -node.weight, depth-1, node)
				}
			}

			// restore move and restore data
			s.restoreData(node, n)

			if depth == s.maxDepth {
				s.addNode(node)

				if alpha < node.weight {
					s.save = node
				}
			}

			alpha = maxWeight(alpha, node.weight)

			if alpha >= beta {
				if depth == 4 {
					fmt.Println("alpha: ", alpha, " - beta: ", beta)
				}
				return alpha
			}
		}
	}
	return alpha
}

// Play start the alphabeta algorythm
func Play(b *[19][19]uint8, s *database.Session, c chan uint8) {
	state := New(b, rdef.Player2)
	state.addTotalCapture(rdef.Player1, uint8(s.NbCaptureP1))
	state.addTotalCapture(rdef.Player2, uint8(s.NbCaptureP2))

	//ret := state.alphabetaNegaScout(math.MinInt8+1, math.MaxInt8, state.maxDepth, nil)
	state.alphabetaNegaScout(math.MinInt16+1, math.MaxInt16, state.maxDepth, nil)
	//fmt.Println("ret: ", ret)

	tmp := int16(math.MinInt16)

	yi, xi := state.save.rule.GetPosition()
	fmt.Println("save weight: ", state.save.weight, " y et x: ", yi, " - ", xi)

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
