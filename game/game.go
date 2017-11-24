package game

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/ruler"
	rdef "github.com/ymohl-cl/gomoku/game/ruler/defines"
	"github.com/ymohl-cl/gomoku/game/stats"
)

const (
	errorData        = "data provided to game is empty"
	errorPlayer      = "missing player to run the game"
	errorMove        = "this move is not possible"
	errorDoubleThree = "There are a double three"
)

// Game struct contain playerInfo struct, the actual board, timers
// and a Rules instance
type Game struct {
	players   playersInfo
	board     [19][19]uint8
	timerPlay time.Time
	timerGame time.Time
	rules     *ruler.Rules
	//Bot       *alphabeta.AI
	data *database.Data
}

// playersInfo struct contain db and stats about the 2 players actual players
// and a reference of the current player
type playersInfo struct {
	p1            *database.Player
	p2            *database.Player
	statsP1       []stats.Stat
	statsP2       []stats.Stat
	currentPlayer *database.Player
}

// New : Return a new instance and set default values of Game struct
func New(d *database.Data) (*Game, error) {
	var err error
	g := Game{}

	// set Player
	if d == nil {
		return nil, errors.New(errorData)
	}
	if d.Current.P1 == nil {
		return nil, errors.New(errorPlayer)
	}
	if d.Current.P2 == nil {
		if d.Current.P2, err = d.GetPlayerByName(database.Bot); err != nil {
			return nil, err
		}
		//g.Bot = ai.New()
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
	/*for y := 0; y < 19; y++ {
		line := []uint8{}
		for x := 0; x < 19; x++ {
			line = append(line, uint8(0))
		}
		g.board = append(g.board, line)
	}*/
	g.data = d
	return &g, nil
}

// GetCurrentPlayer : return the reference of the actual player
func (g Game) GetCurrentPlayer() *database.Player {
	return g.players.currentPlayer
}

// AppliesMove : save the move in stats struct and put it in the board
func (g *Game) AppliesMove(x, y uint8) {
	// create stat move
	s := stats.New(x, y, time.Since(g.timerPlay), false)

	if g.players.currentPlayer == g.players.p1 {
		g.board[y][x] = rdef.Player1
		g.players.statsP1 = append(g.players.statsP1, s)
	} else {
		g.board[y][x] = rdef.Player2
		g.players.statsP2 = append(g.players.statsP2, s)
	}
}

// SwitchPlayer : _
func (g *Game) SwitchPlayer() {
	if g.players.currentPlayer == g.players.p1 {
		g.players.currentPlayer = g.players.p2
	} else {
		g.players.currentPlayer = g.players.p1
	}
}

// Move : Call the rules checker
// Do the move according to the circumstances
// return if current player win
func (g *Game) Move(x, y uint8) (bool, string, error) {
	var valueToken uint8
	var nbCaps *int32

	// get current player with token value
	if g.players.currentPlayer == g.players.p1 {
		valueToken = rdef.Player1
		nbCaps = &g.data.Current.NbCaptureP1
	} else {
		valueToken = rdef.Player2
		nbCaps = &g.data.Current.NbCaptureP2
	}

	g.rules = ruler.New(valueToken, int8(y), int8(x))

	//CheckAllRules
	g.rules.CheckRules(&g.board, uint8(*nbCaps))
	//g.rules.Print()
	//add Capture nb
	*nbCaps += int32(g.rules.NumberCapture)
	//Verify Check
	fmt.Println(g.rules)
	if g.rules.Movable == false {
		return false, "", errors.New(g.rules.Info)
	}
	g.rules.ApplyMove(&g.board)
	/*
		if g.rules.IsCaptured == true {
			for _, cap := range g.rules.GetCaptures() {
				g.board[cap.Y][cap.X] = ruler.TokenEmpty
			}
		}
	*/
	g.AppliesMove(x, y)
	g.SwitchPlayer()

	return g.rules.Win, g.rules.Info, nil
}

// GetCaptures : return a reference of slice of actual player captures
func (g Game) GetCaptures() []*ruler.Spot {
	return g.rules.GetCaptures()
}

// Playing : reset timerPlayer and rules
func (g *Game) Playing() {
	g.timerPlay = time.Now()
	g.rules = nil
}

// Run : reset all timer and rules
func (g *Game) Run() {
	g.Playing()
	g.timerGame = time.Now()
}

// GetTimeToPlay : return time of player action duration
func (g Game) GetTimeToPlay() time.Duration {
	return time.Since(g.timerPlay)
}

// GetTimeGame : return time of all the game
func (g Game) GetTimeGame() time.Duration {
	return time.Since(g.timerGame)
}

// GetBoard : _
func (g Game) GetBoard() *[19][19]uint8 {
	return &g.board
}
