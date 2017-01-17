// Package game provides a higher level interface to package dominos.
package game

import (
	"errors"
	"fmt"
	"log"

	"github.com/Xe/uuid"
	"github.com/cetacean/magiism/dominos"
)

// Package errors
var (
	ErrGameCreationFailed = errors.New("game: dominos.NewGame failed, please report as a bug")
	ErrNotYourTurn        = errors.New("game: it is not your turn")
	ErrInvalidHandIndex   = errors.New("game: invalid hand index")
	ErrEndOfTurn          = errors.New("game: your turn is now over")
	ErrUnknownAction      = errors.New("game: unknown action")
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

// Possible messages to the client, TODO: translations?
const (
	KnockSuccessfulMsg   = "$EVENT_PLAYER_NAME has only one tile left!"
	CannotKnockMsg       = "You cannot knock, you have more than one tile in your hand"
	OutOfTilesMsg        = "Out of tiles, can't draw"
	MustResolveDoubleMsg = "You must resolve this double if you can"
	PlaySuccessfulMsg    = "$EVENT_PLAYER_NAME has played $DOMINO on $PATH_ID_OWNER"
	MustTryDrawingMsg    = "You must try to draw a tile and see if that works before ending your turn"
	SettingTrainMsg      = "Setting train on $EVENT_PLAYER_NAME"
)

// Event is a single user command -> game state event.
type Event struct {
	Action   Action
	PlayerID string // This must be populated by the server, never by the user directly.

	// If PlayDomino is chosen, these next two fields are filled.
	PathID    int
	HandIndex int
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
	dg, err := dominos.NewGame(players)
	if err != nil {
		log.Println("game creation failed: ", err)
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
		State:    g.Game,
		PlayerID: e.PlayerID,
	}

	p := g.GetActivePlayer()

	if e.Action == Knock {
		if e.PlayerID != p.ID {
			kp, ok := g.GetPlayerByID(e.PlayerID)
			if !ok {
				return nil, errors.New("impossible state, tried to look up a player that doesn't exist?")
			}

			if g.Knock(kp) {
				r.GlobalMessage = KnockSuccessfulMsg
				return r, nil
			}
			return nil, ErrNotYourTurn
		}
	}

	if g.GetActivePlayer().ID != e.PlayerID {
		return nil, ErrNotYourTurn
	}

	// switch on e.Action and then take the appropriate actions.
	switch e.Action {
	case EndTurn:
		nagged := false
		for i, d := range p.Hand {
			for j, path := range g.Trains {
				_, err := g.CanPlace(p, d, path)
				if err == nil {
					nagged = true
					r.UserMessage += fmt.Sprintf("you can place tile %s (%d) in your hand on path %d\n", d.Display(), i, j)
				}
			}
		}

		if nagged {
			r.Success = false
			return r, nil
		}

		if !g.Played && !g.Drawn {
			r.UserMessage = MustTryDrawingMsg
			r.Success = false
			return r, nil
		}

		if g.Drawn && !g.Played {
			r.GlobalMessage += "\n" + SettingTrainMsg
		}

		g.endOfTurn(r)
		return r, ErrEndOfTurn

	case PlayDomino:
		path := g.Trains[e.PathID]
		d, ok := p.RemoveFromHand(e.HandIndex)
		if !ok {
			return nil, ErrInvalidHandIndex
		}

		err := g.Place(p, d, path)
		if err != nil {
			p.Hand = append(p.Hand, d)

			return nil, err
		}

		r.GlobalMessage = PlaySuccessfulMsg
		r.Success = true
		g.Played = true

		if d.IsDouble() {
			g.UnresolvedDouble = true
			path.UnresolvedDouble = true
			r.UserMessage = MustResolveDoubleMsg
			g.Played = false

			return r, nil
		}

		g.endOfTurn(r)
		return r, ErrEndOfTurn

	case DrawDomino:
		if !g.Drawn {
			g.Drawn = true
			err := g.Draw(p)
			if err != nil {
				r.GlobalMessage = OutOfTilesMsg
				g.endOfTurn(r)
				return r, ErrEndOfTurn
			}
		} else {
			r.Success = false
		}

	case Knock:
		if g.Knock(p) {
			r.GlobalMessage = KnockSuccessfulMsg
			r.Success = true
			p.Knocked = true
			return r, nil
		} else {
			r.UserMessage = CannotKnockMsg
			r.Success = false
		}
	default:
		return nil, ErrUnknownAction
	}

	return r, nil
}

func (g *Game) endOfTurn(r *Response) {
	g.Drawn = false
	g.Played = false

	_, status := g.NextTurn()
	if status != "" {
		switch status {
		case "noknock":
			r.GlobalMessage += fmt.Sprintf("\n$CURRENT_PLAYER has drawn two tiles for not knocking when they had one tile left")
		}
	}
}
