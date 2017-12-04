package alphabeta

import (
	"math"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

const (
	levelLow    = 1
	levelMedium = 2
	levelHard   = 3
)

type IA struct {
	level uint8
	moves *Node
}

func NewIA() *IA {
	return &IA{}
}

// InfoPlayer of previous evaluation
type InfoPlayer struct {
	totalCapture uint8
}

// State of evaluation
type State struct {
	maxDepth      uint16
	limitDepth    uint16
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
	if ok, _ := ruler.IsAvailablePosition(s.board, y, x); !ok {
		return nil
	}
	n := new(Node)
	n.rule.Init(s.currentPlayer, y, x)
	n.rule.CheckRules(s.board, s.getTotalCapture(s.currentPlayer))
	if !n.rule.Movable {
		return nil
	}
	return n
}

func (s *State) updateTokenPlayer(nodes *[10]*Node) {
	var n *Node

	for i := 0; i < 10; i++ {
		if n = (*nodes)[i]; n == nil {
			break
		}
		y, x := n.rule.GetPosition()
		(*s.board)[y][x] = rdef.Empty
	}
}

func (s *State) restoreTokenPlayer(nodes *[10]*Node) {
	var n *Node

	for i := 0; i < 10; i++ {
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

func (s *State) alphabetaNegaScout(alpha, beta int16, depth uint16, n *Node) int16 {
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

			if alpha < node.weight && depth == 1 && s.maxDepth < s.limitDepth && !node.rule.Win {
				s.maxDepth += 2
				node.weight = -s.alphabetaNegaScout(-beta, -node.weight, depth+2, node)
				s.maxDepth -= 2
			}

			// restore move and restore data
			s.restoreData(node, n)

			// save best moves for the player
			if depth == s.maxDepth-1 && alpha < node.weight {
				s.addNode(node)
			}
			// save best moves for the IA
			if depth == s.maxDepth && alpha < node.weight {
				s.save = node
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
func (i *IA) Play(b *[19][19]uint8, s *database.Session, c chan uint8) {
	state := New(b, rdef.Player2)
	state.addTotalCapture(rdef.Player1, uint8(s.NbCaptureP1))
	state.addTotalCapture(rdef.Player2, uint8(s.NbCaptureP2))

	switch i.level {
	case levelLow:
		state.limitDepth = 6
	case levelMedium:
		state.limitDepth = 8
	case levelHard:
		state.limitDepth = 10
	default:
		state.limitDepth = 4
	}

	_ = state.alphabetaNegaScout(math.MinInt16+1, math.MaxInt16, state.maxDepth, nil)

	i.moves = state.lst
	y, x := state.save.rule.GetPosition()
	c <- uint8(y)
	c <- uint8(x)
}

func (i *IA) PlayOpponnent(y, x int8) {
	var weightMove int16
	var bestScore int16

	bestScore = math.MinInt16 + 1
	for n := i.moves; n != nil; n = n.next {
		if n.weight > bestScore {
			bestScore = n.weight
		}
		posY, posX := n.rule.GetPosition()
		if posY == y && posX == x {
			weightMove = n.weight
		}
	}

	if weightMove == bestScore && i.level < levelHard {
		i.level++
	}
}
