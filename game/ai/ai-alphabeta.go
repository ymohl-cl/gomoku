package ai

import (
	"math"
	"time"

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
	board         *[][]uint8
	nbCapsCurrent uint8
	nbCapsOther   uint8
	lenMaxP1      uint8
	lenMaxP2      uint8
	alpha         int8
	beta          int8
	player        uint8
	pq            *Node
	nodes         map[uint8]map[uint8]*Node
	save          *Node
}

func (a AI) newState(b *[][]uint8, nbCOldCurrent, nbCOldOther, lMaxP1, lMaxP2, p uint8) *State {
	s := State{
		board:         b,
		nbCapsCurrent: nbCOldOther,
		nbCapsOther:   nbCOldCurrent,
		lenMaxP1:      lMaxP1,
		lenMaxP2:      lMaxP2,
		player:        p,
		nodes:         make(map[uint8]map[uint8]*Node),
	}

	return &s
}

type Node struct {
	r         ruler.Rules
	x         uint8
	y         uint8
	weight    int8
	minV      int8
	maxV      int8
	nextState State
	next      *Node
}

func (a AI) newNode() *Node {
	return &Node{
		minV: math.MinInt8,
		maxV: math.MaxInt8,
	}
}

func (n Node) newByCopy(b *[][]uint8, s *State) *Node {
	newNode := &Node{
		x:      n.x,
		y:      n.y,
		weight: n.weight,
		next:   n.next,
	}
	newNode.r.Copy(&n.r)
	newNode.nextState.Set(b, s.nbCapsCurrent+n.r.NbCaps, s.nbCapsOther, s.lenMaxP1, s.lenMaxP2, s.switchPlayer())
	return newNode
}

func (n *Node) Clean(x, y uint8) {
	n.r.Clean()
	n.x = x
	n.y = y
	n.weight = 0
	n.nextState.Clean()
	n.next = nil
}

func (s *State) Clean() {
	s.board = nil
	s.nbCapsCurrent = 0
	s.nbCapsOther = 0
	s.lenMaxP1 = 0
	s.lenMaxP2 = 0
	s.player = 0
	s.pq = nil
	s.nodes = nil
}

func (s *State) Set(b *[][]uint8, nbCOldCurrent, nbCOldOther, lMaxP1, lMaxP2, p uint8) {
	s.board = b
	s.nbCapsCurrent = nbCOldOther
	s.nbCapsOther = nbCOldCurrent
	s.lenMaxP1 = lMaxP1
	s.lenMaxP2 = lMaxP2
	s.player = p
	s.alpha = math.MinInt8
	s.beta = math.MaxInt8
	s.nodes = make(map[uint8]map[uint8]*Node)
}

func (s *State) addNode(n *Node) {
	//	if _, ok := s.nodes[n.y]; !ok {
	//		s.nodes[n.y] = make(map[uint8]*Node)
	//	}
	//	s.nodes[n.y][n.x] = n
	//	s.save = n

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

func (a *AI) genBoard(b *[][]uint8, node *Node, player uint8) *[][]uint8 {
	var newBoard [][]uint8
	for _, line := range *b {
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
	return &newBoard
}

func (s State) switchPlayer() uint8 {
	if s.player == ruler.TokenP1 {
		return ruler.TokenP2
	}
	return ruler.TokenP1
}

func (a AI) eval(s *State, n *Node, stape int8) int8 {
	var ret int8

	if !n.r.IsWin {
		ret = 30
		ret += int8(s.nbCapsCurrent) + int8(s.nbCapsOther)
	} else {
		ret = math.MaxInt8
	}

	ret -= stape

	if s.player == ruler.TokenP1 {
		return ret
	}
	return -ret
}

/*func (a *AI) prevalphabeta(s *State, prevnode *Node, alpha, beta int8, stape int8) {
var y, x uint8
var Val int8
wg := new(sync.WaitGroup)

Val = math.MinInt8
for y = 0; y < 19; y++ {
	for x = 0; x < 19; x++ { /*pour tout enfant si de s faire */
/*			wg.Add(1)
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
	}
	wg.Wait()
}*/

var timeTotalNode time.Duration
var timeTotalRule time.Duration

//var timeTotalAddNode time.Duration
var timeTotalCopyBoard time.Duration

//var timeTotalState time.Duration

func (s *State) SetAlpha(a int8) {
	s.alpha = a
}

func (s *State) SetBeta(b int8) {
	s.beta = b
}

func (s *State) SearchByPosition(y, x uint8) *Node {
	if _, ok := s.nodes[y]; !ok {
		return nil
	}
	if node, ok := s.nodes[y][x]; ok {
		return node
	}
	return nil
}

func (s *State) Search(n *Node) *Node {
	if _, ok := s.nodes[n.y]; !ok {
		return nil
	}
	if node, ok := s.nodes[n.y][n.x]; ok {
		return node
	}
	return nil

	/*	for test := s.pq; test != nil; test = test.next {
			if test.x == n.x && test.y == n.y {
				return test
			}
		}
		return nil*/
}

func (a *AI) alphabeta(s *State, alpha, beta int8, stape int8, workNode *Node) int8 {
	var score int8
	var node *Node
	var test bool
	//	var findNode *Node

	/*	if node = s.Search(workNode); node != nil {
		test = true
		if node.minV >= beta {
			return node.minV
		}
		if node.maxV <= alpha {
			return node.maxV
		}
		alpha = a.max(alpha, node.minV)
		beta = a.min(beta, node.maxV)
	}*/

	if stape == 0 || workNode.r.IsWin {
		score = a.eval(s, workNode, stape)
		if node == nil {
			node = workNode
		}
	} else if s.player == ruler.TokenP2 {
		score = math.MinInt8
		aPrime := alpha
		for y := uint8(0); y < 19; y++ {
			for x := uint8(0); x < 19; x++ {
				workNode.Clean(x, y)
				workNode.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
				if workNode.r.IsMoved {
					node = workNode.newByCopy(a.genBoard(s.board, workNode, s.player), s)
					score = a.max(score, a.alphabeta(&node.nextState, aPrime, beta, stape-1, workNode))
					aPrime = a.max(aPrime, score)
					node.weight = score

					if score >= beta {
						break
					}
					if !test {
						s.addNode(node)
					}
				}
			}
		}
	} else {
		score = math.MaxInt8
		bPrime := beta
		for y := uint8(0); y < 19; y++ {
			for x := uint8(0); x < 19; x++ {
				workNode.Clean(x, y)
				workNode.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
				if workNode.r.IsMoved {
					node = workNode.newByCopy(a.genBoard(s.board, workNode, s.player), s)
					score = a.min(score, a.alphabeta(&node.nextState, alpha, bPrime, stape-1, workNode))
					node.weight = score
					bPrime = a.min(bPrime, score)
					if score <= alpha {
						break
					}

					if !test {
						s.addNode(node)
					}
				}
			}
		}
	}

	node.weight = score
	if score <= alpha {
		node.maxV = score
		s.beta = node.maxV
	}
	if score > alpha && score < beta {
		node.minV = score
		node.maxV = score
		s.beta = node.maxV
		s.alpha = node.minV
	}
	if score >= beta {
		node.minV = score
		s.alpha = node.minV
	}

	return score
}

/*func (a *AI) alphabeta(s *State, alpha, beta int8, stape int8, workNode *Node) int8 {
	if stape <= 0 || workNode.r.IsWin {
		fmt.Println("WorkNode is win")
		return a.eval(s, workNode, stape)
	}

	for y := uint8(0); y < 19; y++ {
		for x := uint8(0); x < 19; x++ {
			workNode.Clean(x, y)
			workNode.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
			if workNode.r.IsMoved {
				node := workNode.newByCopy(a.genBoard(s.board, workNode, s.player), s)
				score := a.alphabeta(&node.nextState, -beta, -alpha, stape-1, workNode)
				fmt.Println("score: ", score)
				if score >= alpha {
					alpha = score
					s.addNode(node)
					if alpha >= beta {
						break
					}
				}
			}
		}
	}
	return alpha
}*/

/* alpha est toujours inférieur à beta */
/*func (a *AI) alphabeta(s *State, prevnode *Node, alpha, beta int8, stape int8, workNode *Node) int8 {
var y, x uint8
var Val int8

if stape <= 0 || (prevnode != nil && prevnode.r.IsWin) {
	return a.eval(s, prevnode, stape)
}
if s.player == ruler.TokenP2 {
	Val = math.MinInt8
	for y = 0; y < 19; y++ {
		for x = 0; x < 19; x++ { /*pour tout enfant si de s faire */
/*				timeNode := time.Now()
			workNode.Clean(x, y)
			//	node := a.newNode(x, y)
			timeTotalNode += time.Since(timeNode)
			timeRule := time.Now()
			workNode.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
			//node.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
			timeTotalRule += time.Since(timeRule)
			if workNode.r.IsMoved {
				timeBoard := time.Now()
				node := workNode.newByCopy(a.genBoard(s.board, workNode, s.player), s)
				timeTotalCopyBoard += time.Since(timeBoard)
				//				wg := new(sync.WaitGroup)
				//				wg.Add(1)
				//				go func() {
				//					timeAdd := time.Now()
				//					defer wg.Done()
				//						s.addNode(node)
				//					timeTotalAddNode += time.Since(timeAdd)
				//				}()
				//					wg.Add(1)
				//					go func() {
				//						defer wg.Done()
				//					}()
				//					wg.Wait()
				//					timeState := time.Now()
				//					timeTotalState += time.Since(timeState)
				Val = a.max(Val, a.alphabeta(&node.nextState, node, alpha, beta, stape-1, workNode))
				node.weight = Val
				alpha = a.max(alpha, Val)
				fmt.Println("t2 alpha: ", alpha)
				fmt.Println("t2 beta: ", beta)
				fmt.Println("t2 Value:", Val)
				//				node.nextState.SetBeta(beta)
				//	node.nextState.SetAlpha(alpha)
				//node.nextState.SetAlpha(alpha)
				if beta <= Val {
					return Val
				}
				//node.nextState.SetBeta(beta)
				s.addNode(node)
			}
		}
	}
} else {
	Val = math.MaxInt8
	for y = 0; y < 19; y++ {
		for x = 0; x < 19; x++ { /*pour tout enfant si de s faire */
/*				timeNode := time.Now()
				workNode.Clean(x, y)
				//				node := a.newNode(x, y)
				timeTotalNode += time.Since(timeNode)
				timeRule := time.Now()
				workNode.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
				//				node.r.CheckRules(s.board, int8(x), int8(y), s.player, s.nbCapsCurrent)
				timeTotalRule += time.Since(timeRule)
				if workNode.r.IsMoved {
					timeBoard := time.Now()
					node := workNode.newByCopy(a.genBoard(s.board, workNode, s.player), s)
					timeTotalCopyBoard += time.Since(timeBoard)
					//					wg := new(sync.WaitGroup)
					//					wg.Add(1)
					//					go func() {
					//						timeAdd := time.Now()
					//						defer wg.Done()
					//						s.addNode(node)
					//						timeTotalAddNode += time.Since(timeAdd)
					//					}()
					//					wg.Add(1)
					//					go func() {
					//						defer wg.Done()
					//					}()
					//					wg.Wait()
					//					timeState := time.Now()
					//					timeTotalState += time.Since(timeState)
					Val = a.min(Val, a.alphabeta(&node.nextState, node, alpha, beta, stape-1, workNode))
					node.weight = Val
					beta = a.min(beta, Val)

					fmt.Println("t1 alpha: ", alpha)
					fmt.Println("t1 beta: ", beta)
					fmt.Println("t1 Value:", Val)
					//		node.nextState.SetBeta(beta)
					//					node.nextState.SetAlpha(alpha)
					if Val <= alpha {
						return Val
					}
					//node.nextState.SetAlpha(alpha)
					s.addNode(node)
				}
			}
		}
	}
	return Val
}*/

func (a *AI) getCoord(weight int8) (uint8, uint8) {
	//	return a.s.save.y, a.s.save.x
	//var x, y uint8
	var refNode *Node

	//	fmt.Println("max and min of state: ", a.s.alpha, " - ", a.s.beta)
	//	fmt.Println("save node: ", a.s.save)
	//	a.s = &a.s.pq.nextState
	//	min := int8(math.MaxInt8)
	/*	max := int8(math.MinInt8)
		for _, m := range a.s.nodes {
			for _, node := range m {
				fmt.Print("Y: ", node.y, " - X: ", node.x, " | Weight Node: ", node.weight, " | ")
				fmt.Print(" min: ", node.minV, " - max: ", node.maxV, " | ")
				if max <= node.weight {
					max = node.weight
					fmt.Println("Choice this move")
					//				min = node.minV
					//				max = node.maxV
					refNode = node
				} else {
					fmt.Println("No choice")
				}*/
	/*			weight := node.weight
				fmt.Print("Y: ", node.y, " - X: ", node.x, " | Weight Node: ", weight, " | ")
				fmt.Print(" min: ", node.minV, " - max: ", node.maxV, " | ")
				if tmp >= weight {
					fmt.Println("Choice this move")
					//x = node.x
					//y = node.y
					tmp = node.weight
					refNode = node
				} else {
					fmt.Println("No choice")
				}*/
	//	}
	//	}
	//a.s = &refNode.nextState
	//if a.s.save != nil {
	//	refNode = a.s.save
	//}
	//a.s = nil
	//return refNode.y, refNode.x
	tmp := int8(math.MinInt8)
	for node := a.s.pq; node != nil; node = node.next {
		//fmt.Println("lower: ", node.minV, " - upper: ", node.maxV, " - weight: ", node.weight)
		weight := node.weight
		//		if weight < 0 {
		//		weight *= -1
		//		}
		if tmp <= weight {
			//		x = node.x
			//		y = node.y
			tmp = node.weight
			refNode = node
			//tmp = node.weight
			//refNode = node
		}
	}
	return refNode.y, refNode.x
}

/*	if mi != tmp {
	fmt.Println("Not first")
}*/
//	a.s = &refNode.nextState
//	a.s.pq = nil
//fmt.Println("----- Getting Coord -----")
//fmt.Println("a.s alpha: ", a.s.alpha)
//fmt.Println("a.s beta: ", a.s.beta)
//refNode.nextState.alpha = a.s.alpha
//*a.s = refNode.nextState
//if a.s.alpha != refNode.nextState.alpha || a.s.beta != refNode.nextState.beta {
//	fmt.Println("Error on herit nextState on getCoord")
//}
//fmt.Println("Switch to nextState (opponent) ")
//fmt.Println("a.s alpha: ", a.s.alpha)
//fmt.Println("a.s beta: ", a.s.beta)
//	return x, y

func (a *AI) Play(b *[][]uint8, s *database.Session, c chan uint8) {
	var workNode *Node

	if a.s == nil { //|| a.s.pq == nil {
		a.s = a.newState(b, uint8(s.NbCaptureP1), uint8(s.NbCaptureP2), 0, 0, ruler.TokenP2)
		a.s.alpha = math.MinInt8
		a.s.beta = math.MaxInt8
	}
	workNode = a.newNode()
	timeTotalNode = 0
	//	timeTotalAddNode = 0
	timeTotalCopyBoard = 0
	timeTotalRule = 0
	//	timeTotalState = 0

	//	fmt.Println("AI a.s alpha: ", a.s.alpha)
	//	fmt.Println("AI a.s beta: ", a.s.beta)
	var test time.Duration
	test += time.Since(time.Now())
	//a.prevalphabeta(a.s, nil, math.MinInt8, math.MaxInt8, 4)
	a.alphabeta(a.s, math.MinInt8, math.MaxInt8, 4, workNode)

	//	a.getLen(a.s)
	y, x := a.getCoord(0)
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
	/*	if a.s == nil || a.s.pq == nil {
		return
	}*/

	if a.s == nil {
		return
	}
	node := a.s.SearchByPosition(y, x)
	if node != nil {
		*a.s = node.nextState
	} else {
		a.s = nil
	}
	return
	/*
		for node := a.s.pq; node != nil; node = node.next {
			if y == node.y && x == node.x {
				//			fmt.Println("------------------ START ------------------")
				//node.nextState.alpha = a.s.alpha
				*a.s = node.nextState
				//			fmt.Println("alpha opposant: ", a.s.alpha)
				//			fmt.Println("beta opposant: ", a.s.beta)
				return
			}
		}*/
	//	a.s = nil
	//	fmt.Println("FUCCCCK")
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
