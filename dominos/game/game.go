// Package game provides a higher level interface to package dominos.
package game

import (
	"errors"

	"github.com/Xe/uuid"
	"github.com/cetacean/magiism/dominos"
)

// Package errors
var (
	ErrGameCreationFailed = errors.New("game: dominos.NewGame failed, please report as a bug")
	ErrNotYourTurn        = errors.New("game: it is not your turn")
)

// Action is the kind of turn action the player is taking.
type Action int

// Possible actions a player can take in their turn.
const (
	EndTurn Action = iota
	PlayDomino
	DrawDomino
	Knock
)

// Event is a single user command -> game state event.
type Event struct {
	Action   Action
	PlayerID string // This must be populated by the server, never by the user directly.

	// If PlayDomino is chosen, these next two fields are filled.
	PathID int
	Domino dominos.Domino
}

// Response is the result of the event being run against the game state.
// This structure is safe to json-encode for the user directly.
type Response struct {
	Success bool

	State *dominos.Game

	GlobalMessage string
	UserMessage   string
	PlayerID      string
}

// Game is a high-level wrapper around the dominos.Game struct.
type Game struct {
	*dominos.Game

	ID string

	// These variables are for the currently active player's turn.
	Drawn  bool
	Played bool
}

// Store represents a in-memory or on-database storage for many domino games.
type Store interface {
	GetGame(id string) (*Game, error)
	PutGame(id string, g *Game) error
}

// New creates a new game with given players.
func New(players []string) (*Game, error) {
	dg := dominos.NewGame(players)
	if dg == nil {
		return nil, ErrGameCreationFailed
	}

	g := &Game{
		Game: dg,
		ID:   uuid.New(),
	}

	return g, nil
}

// HandleEvent handles a single game event, failing if it failed.
func (g *Game) HandleEvent(e *Event) (*Response, error) {
	r := &Response{
		State: g.Game,
	}

	if g.GetActivePlayer().ID != e.PlayerID {
		return nil, ErrNotYourTurn
	}

	// switch on e.Action and then take the appropriate actions.

	r.Success = true
	return r, nil
}
