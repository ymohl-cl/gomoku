package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
	"github.com/ymohl-cl/gomoku/game/stats"
)

const (
	errorData   = "data provided to game is empty"
	errorPlayer = "missing player to run the game"
	errorMove   = "this move is not possible"
)

type Game struct {
	players   playersInfo
	board     [][]uint8
	timerPlay time.Time
	timerGame time.Time
	rules     *ruler.Rules
}

type playersInfo struct {
	p1            *database.Player
	p2            *database.Player
	statsP1       []stats.Stat
	statsP2       []stats.Stat
	currentPlayer *database.Player
}

func New(d *database.Data) (*Game, error) {
	g := Game{}

	// set Player
	if d == nil {
		return nil, errors.New(errorData)
	}
	if d.Current.P1 == nil || d.Current.P2 == nil {
		return nil, errors.New(errorPlayer)
	}
	g.players = playersInfo{p1: d.Current.P1, p2: d.Current.P2}

	// define first player randomly
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	if r.Int31()%2 == 0 {
		g.players.currentPlayer = g.players.p1
	} else {
		g.players.currentPlayer = g.players.p2
	}

	// init board
	for y := 0; y < 19; y++ {
		line := []uint8{}
		for x := 0; x < 19; x++ {
			line = append(line, uint8(0))
		}
		g.board = append(g.board, line)
	}
	fmt.Println("len de board: ", len(g.board))
	return &g, nil
}

func (g Game) GetCurrentPlayer() *database.Player {
	return g.players.currentPlayer
}

func (g *Game) AppliesMove(x, y uint8) {
	// create stat move
	s := stats.New(x, y, time.Since(g.timerPlay), false)

	if g.players.currentPlayer == g.players.p1 {
		g.board[y][x] = ruler.TokenP1
		g.players.statsP1 = append(g.players.statsP1, s)
	} else {
		g.board[y][x] = ruler.TokenP2
		g.players.statsP2 = append(g.players.statsP2, s)
	}
}

func (g *Game) SwitchPlayer() {
	if g.players.currentPlayer == g.players.p1 {
		g.players.currentPlayer = g.players.p2
	} else {
		g.players.currentPlayer = g.players.p1
	}
}

func (g *Game) Move(x, y uint8) error {
	g.rules = ruler.New()

	// veriff move
	if b, str := g.rules.CheckMove(g.board, int8(x), int8(y)); b == false {
		return errors.New(errorMove + str)
	}

	g.AppliesMove(x, y)
	// check and save capture token
	var valueToken uint8
	if g.players.currentPlayer == g.players.p1 {
		valueToken = ruler.TokenP1
	} else {
		valueToken = ruler.TokenP2
	}
	g.rules.CheckCapture(g.board, int8(x), int8(y), valueToken)
	for _, cap := range g.rules.GetCaptures() {
		g.board[cap.Y][cap.X] = ruler.TokenEmpty
	}

	g.SwitchPlayer()
	return nil
}

func (g Game) GetCaptures() []ruler.Capture {
	return g.rules.GetCaptures()
}

func (g *Game) Playing() {
	g.timerPlay = time.Now()
	g.rules = nil
}

func (g *Game) Run() {
	g.Playing()
	g.timerGame = time.Now()
}

func (g *Game) GetTimeToPlay() time.Duration {
	return time.Since(g.timerPlay)
}

func (g *Game) GetTimeGame() time.Duration {
	return time.Since(g.timerGame)
}
