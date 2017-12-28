package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/ymohl-cl/gomoku/database"
	"github.com/ymohl-cl/gomoku/game/alphabeta"
	"github.com/ymohl-cl/gomoku/game/ruler"
	"github.com/ymohl-cl/gomoku/game/ruler/alignment"
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
	Bot       *alphabeta.IA
	data      *database.Data
	End       bool
	Start     bool
	RulerWin  *ruler.Rules
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
	g := Game{}

	// set Player
	if d == nil {
		return nil, errors.New(errorData)
	}
	if d.Current.P1 == nil {
		return nil, errors.New(errorPlayer)
	}
	if d.Current.P2 == nil {
		d.Current.P2 = new(database.Player)
		d.Current.P2.Name = database.Bot
		g.Bot = alphabeta.NewIA()
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

	g.data = d
	return &g, nil
}

func (g Game) GetOtherName() string {
	if g.players.currentPlayer == g.players.p1 {
		return g.players.p2.Name
	} else if g.players.currentPlayer == g.players.p2 {
		return g.players.p1.Name
	}
	return "Moutchatcho :D"
}

// IsBot define is the player it's a bot
func (g Game) IsBot(p *database.Player) bool {
	if p == g.players.p2 && g.Bot != nil {
		return true
	}
	return false
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
func (g *Game) Move(x, y uint8) (bool, string) {
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

	if g.Start == true {
		g.Start = false
		// init first move
		g.rules = ruler.New(valueToken, int8(y), int8(x))
		g.rules.Movable = true
	} else {
		if ok, message := ruler.IsAvailablePosition(&g.board, int8(y), int8(x)); !ok {
			return ok, message
		}

		// Create the ruler
		g.rules = ruler.New(valueToken, int8(y), int8(x))
		g.rules.Movable = true

		//CheckAllRules
		g.rules.CheckRules(&g.board, uint8(*nbCaps))
		if !g.rules.Movable {
			return false, g.rules.Info
		}
	}

	//add Capture nb
	*nbCaps += int32(g.rules.NumberCapture)

	// applies the move on the board
	g.rules.ApplyMove(&g.board)
	// applies the move on the data game
	g.AppliesMove(x, y)

	g.SwitchPlayer()

	return true, ""
}

// IsWin return status win created by the ruler
func (g *Game) IsWin() (bool, string) {
	return g.rules.Win, g.rules.Info
}

// GetCaptures : return a reference of slice of actual player captures
func (g Game) GetCaptures() []*alignment.Spot {
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

func (g Game) GetRules() *ruler.Rules {
	return g.rules
}

func (g Game) IsCapturable() bool {
	a := g.rules.GetBetterAlignment()
	return a.GetCaptureStatus()
}
