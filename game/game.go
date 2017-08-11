package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/stats"
)

const (
	tokenEmpty = 0
	tokenP1    = 1
	tokenP2    = 2
	sizeY      = 19
	sizeX      = 19

	errorData   = "data provided to game is empty"
	errorPlayer = "missing player to run the game"
	errorMove   = "this move is not possible"
)

type Game struct {
	players playersInfo
	board   [][]uint8
	timer   time.Time
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

func (g Game) CheckMove(x, y uint8) bool {
	if g.board[y][x] == tokenEmpty {
		return true
	}
	return false
}

func (g *Game) Move(x, y uint8) error {
	// veriff move
	if !g.CheckMove(x, y) {
		return errors.New(errorMove)
	}

	// create stat move
	s := stats.New(x, y, time.Since(g.timer), false)

	// applies move and change current player
	if g.players.currentPlayer == g.players.p1 {
		g.board[y][x] = tokenP1
		g.players.statsP1 = append(g.players.statsP1, s)
		g.players.currentPlayer = g.players.p2
	} else {
		g.board[y][x] = tokenP2
		g.players.statsP2 = append(g.players.statsP2, s)
		g.players.currentPlayer = g.players.p1
	}

	return nil
}

func (g *Game) Playing() {
	g.timer = time.Now()
}
