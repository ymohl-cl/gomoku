package alphabeta

import (
	"math"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

type IA struct {
	s     *State
	turn  uint16
	level uint16
}

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
	save   *Node
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

func (s *State) setTotalCapture(player uint8, number uint8) {
	if player == rdef.Player1 {
		s.infoP1.totalCapture = number
	} else {
		s.infoP2.totalCapture = number
	}
}

func (s *State) subTotalCapture(player uint8, number uint8) {
	if player == rdef.Player1 {
		s.infoP1.totalCapture -= number
	} else {
		s.infoP2.totalCapture -= number
	}
}

func (s *State) refresh(b *[19][19]uint8, player uint8) {
	s.board = b
	s.currentPlayer = player
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

func (n *Node) addNext(node *Node, alpha int16) {
	lst := n.next
	n.next = node

	if lst != nil && alpha == node.weight {
		node.next = lst
	}
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

			//	if depth == s.maxDepth {
			//		s.addNode(node)

			//if alpha < node.weight {
			//	s.save = node
			//}
			//		}

			if depth == s.maxDepth && alpha < node.weight {
				s.save = node
			} else if depth%2 != 0 && alpha <= node.weight {
				n.addNext(node, alpha)
			} else if depth%2 == 0 && alpha < node.weight {
				n.save = node
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
	var alpha int16
	if i.s == nil {
		alpha = math.MaxInt16
		i.s = New(b, rdef.Player2)
		i.s.addTotalCapture(rdef.Player1, uint8(s.NbCaptureP1))
		i.s.addTotalCapture(rdef.Player2, uint8(s.NbCaptureP2))
	} else {
		if i.s.save == nil {
			alpha = math.MaxInt16
		} else {
			alpha = i.s.save.weight * -1
		}
		//		fmt.Println("alpha: ", alpha)
		i.s.refresh(b, rdef.Player2)
		i.s.setTotalCapture(rdef.Player1, uint8(s.NbCaptureP1))
		i.s.setTotalCapture(rdef.Player2, uint8(s.NbCaptureP2))
		i.s.save = nil
		i.s.lst = nil
	}

	//ret := state.alphabetaNegaScout(math.MinInt8+1, math.MaxInt8, state.maxDepth, nil)
	i.s.alphabetaNegaScout(math.MinInt16+1, alpha, i.s.maxDepth, nil)
	//fmt.Println("ret: ", ret)

	//	tmp := int16(math.MinInt16)

	//	yi, xi := state.save.rule.GetPosition()
	/*
		fmt.Println("save weight: ", state.save.weight, " y et x: ", yi, " - ", xi)

		for n := state.lst; n != nil; n = n.next {
			y, x := n.rule.GetPosition()
			fmt.Println("Node weight: ", n.weight, " y et x: ", y, " - ", x)
			if tmp <= n.weight {
				tmp = n.weight
				state.save = n
			}
		}
	*/
	y, x := i.s.save.rule.GetPosition()
	i.s.save = i.s.save.next
	c <- uint8(y)
	c <- uint8(x)
}

func (i *IA) PlayOpposing(y, x uint8) {
	var nb int8
	var flag bool

	if i.s == nil || i.s.save == nil {
		return
	}
	for n := i.s.save; n != nil; n = n.next {
		nb++
		ny, nx := n.rule.GetPosition()
		//	fmt.Print("weight: ", n.weight)
		if uint8(ny) == y && uint8(nx) == x {
			//fmt.Println(" choiced")
			i.s.save = n.save
			flag = true
		} else {
			//fmt.Println(" not choice")
		}
	}
	if flag == true {
		//fmt.Println("GREAT !!!! nb: ", nb)
		return
	}

	i.level++
	//fmt.Println("NOOB, nb: ", nb)
}
