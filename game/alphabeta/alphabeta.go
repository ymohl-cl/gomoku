package alphabeta

import (
	"math"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
	"github.com/ymohl-cl/gomoku/game/ruler/alignment"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
)

const (
	levelLow    = 1
	levelMedium = 2
	levelHard   = 3
)

type IA struct {
	level uint8
	nbOk  uint8
	moves *Node
	y     int8
	x     int8
	minY  int8
	minX  int8
	maxY  int8
	maxX  int8
}

func NewIA() *IA {
	return &IA{y: 9, x: 9, minX: 9, minY: 9, maxX: 9, maxY: 9}
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
	y             int8
	x             int8
	minY          int8
	minX          int8
	maxY          int8
	maxX          int8
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

func (i *IA) updatewindow(y, x int8) {
	if y < i.minY {
		i.minY = y
	}
	if x < i.minX {
		i.minX = x
	}
	if x > i.maxX {
		i.maxX = x
	}
	if y > i.maxY {
		i.maxY = y
	}
}

func (s *State) updatewindow(y, x int8) {
	if y < s.minY {
		s.minY = y
	}
	if x < s.minX {
		s.minX = x
	}
	if x > s.maxX {
		s.maxX = x
	}
	if y > s.maxY {
		s.maxY = y
	}
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
	var saveMinY int8
	var saveMinX int8
	var saveMaxY int8
	var saveMaxX int8
	//if n != nil && (n.rule.Win || len(n.rule.CapturableWin) > 0) || depth == 0 {
	if n != nil && (n.rule.Win || depth == 0) {
		return s.eval(n, depth)
	}

	var first bool

	lenY := int8(math.Ceil(float64(s.maxY-s.minY)/2.0)) + 2
	lenX := int8(math.Ceil(float64(s.maxX-s.minX)/2.0)) + 2
	posX := (s.maxX + s.minX) / 2
	posY := (s.maxY + s.minY) / 2
	for tmpCordY, tmpCordX := int8(0), int8(0); true; {
		for yi := -tmpCordY; yi <= tmpCordY; yi++ {
			for xi := -tmpCordX; xi <= tmpCordX; xi++ {
				x := posX + xi
				y := posY + yi

				if yi != tmpCordY && yi != -tmpCordY && xi == -tmpCordX {
					xi = tmpCordX - 1
				}

				node := s.newNode(y, x)
				if node == nil {
					continue
				}
				// apply move and update data
				saveMaxY = s.maxY
				saveMinY = s.minY
				saveMaxX = s.maxX
				saveMinX = s.minX
				s.updateData(node, n)
				s.updatewindow(y, x)
				//				saveY = s.y
				//				saveX = s.x
				//				s.y = y
				//				s.x = x

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
					node.weight = -s.alphabetaNegaScout(-beta, -node.weight, depth+1, node)
					s.maxDepth -= 2
				}

				//				s.y = saveY
				//				s.x = saveX

				// restore move and restore data
				s.restoreData(node, n)

				s.maxY = saveMaxY
				s.minY = saveMinY
				s.maxX = saveMaxX
				s.minX = saveMinX

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
		if tmpCordX == lenX && tmpCordY == lenY {
			break
		}
		if tmpCordX < lenX {
			tmpCordX++
		}
		if tmpCordY < lenY {
			tmpCordY++
		}
	}

	return alpha
}

// Play start the alphabeta algorythm
func (i *IA) Play(b *[19][19]uint8, s *database.Session, c chan uint8, forceSpot []*alignment.Spot) {
	state := New(b, rdef.Player2)
	state.addTotalCapture(rdef.Player1, uint8(s.NbCaptureP1))
	state.addTotalCapture(rdef.Player2, uint8(s.NbCaptureP2))
	state.minY = i.minY
	state.minX = i.minX
	state.maxX = i.maxX
	state.maxY = i.maxY
	//state.y = 9
	//state.x = 9

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

	if len(forceSpot) > 0 {
		c <- uint8(forceSpot[0].Y)
		c <- uint8(forceSpot[0].X)
		return
	}

	_ = state.alphabetaNegaScout(math.MinInt16+1, math.MaxInt16, state.maxDepth, nil)

	i.moves = state.lst
	y, x := state.save.rule.GetPosition()
	i.y = y
	i.x = x
	i.updatewindow(y, x)
	c <- uint8(y)
	c <- uint8(x)
}

func (i *IA) PlayOpponnent(y, x int8) {
	var weightMove int16
	var bestScore int16

	i.updatewindow(y, x)
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
		if i.nbOk >= 1 {
			i.level++
			i.nbOk = 0
		} else {
			i.nbOk++
		}
	}
}
