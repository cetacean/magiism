package dominos

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Domino is a single tile with two sides. This is a game piece.
type Domino struct {
	Left, Right int // The values of each "side" of the domino.
}

// IsDouble checks if both of the tile values are the same.
func (d Domino) IsDouble() bool {
	return d.Left == d.Right
}

// IsPlayable returns true if d2 can be played on d1.
func (d Domino) IsPlayable(d2 Domino) bool {
	return d.Left == d2.Left ||
		d.Left == d2.Right ||
		d.Right == d2.Left ||
		d.Right == d2.Right
}

// Value returns how many "points" a tile is worth.
func (d Domino) Value() int {
	return d.Left + d.Right
}

func (d Domino) Display() string {
	return fmt.Sprintf("[%d|%d]", d.Left, d.Right)
}

// Game represents the total state for a single game
type Game struct {
	TilePool []Domino
	Trains   []*Path
	Players  []*Player
	Center   Domino

	UnresolvedDouble bool
	ActivePlayer     int
}

// Path represents a single player's path. If no player is set,
// the path is treated as the Mexican train.
type Path struct {
	Player   string
	Train    bool // If true, other players can play on it
	Elements []Element

	UnresolvedDouble bool
	MexicanTrain     bool
}

func (p *Player) Display() string {
	result := "YOUR HAND:"
	for i, e := range p.Hand {
		result = result + fmt.Sprintf(" %d:[%d|%d]", i, e.Right, e.Left)
	}
	return result
}

// Display a path
func (p *Path) Display() string {
	result := ""

	if p.MexicanTrain {
		result = result + "M>>"
	} else {
		result = result + fmt.Sprintf("%s>>", p.Player)
	}

	for i, e := range p.Elements {
		result = result + fmt.Sprintf(" %d:%s", i, e.Display())
	}
	if p.Train {
		result = result + " *"
	}
	if p.UnresolvedDouble {
		result = result + " <!>"
	}
	return result
}

// Element is a wrapper for Domino that indicates if the Domino
// is flipped or not. This is for later UI implementation.
type Element struct {
	Domino
	Flipped bool
}

func (e Element) Display() string {
	if e.Flipped {
		d := Domino{
			Left:  e.Right,
			Right: e.Left,
		}

		return d.Display()
	}

	return e.Display()
}

// NewGame creates a new game board out of a list of
// players.
func NewGame(players []string) *Game {
	g := &Game{
		Trains: make([]*Path, len(players)+1),
	}

	mexicanTrain := &Path{
		Train:        true,
		Player:       "",
		MexicanTrain: true,
	}
	g.Trains[len(players)] = mexicanTrain

	// Generate the pool of tiles for the game
	var doms []Domino
	for i := 0; i <= dominoCount(len(players)); i++ {
		for j := 0; j <= i; j++ {
			doms = append(doms, Domino{i, j})
		}
	}

	// Randomize the order of the tiles
	for _, i := range rand.Perm(len(doms)) {
		g.TilePool = append(g.TilePool, doms[i])
	}

	// How many times should be pre-populated into a player's hand
	hc := handCount(len(players))

	// Create player structures
	for i, p := range players {
		newPlayer := &Player{
			ID: p,
		}
		g.Players = append(g.Players, newPlayer)

		path := &Path{
			Player: p,
		}
		newPlayer.Path = path

		g.Trains[i] = path

		for i := 0; i <= hc; i++ {
			err := g.Draw(newPlayer)
			if err != nil {
				panic(err)
			}
		}

		// Find largest piece in player hands
		var largest Domino
		var starter int
		for i, player := range g.Players {
			for _, dom := range player.Hand {
				if dom.Left == dom.Right {
					if dom.Left > largest.Left {
						largest = dom
						starter = i
					}
				}
			}
		}
		g.Center = largest
		g.ActivePlayer = starter
	}

	return g
}

// Player is a single player in the game
type Player struct {
	Hand    []Domino
	BigPlay bool
	Knocked bool
	ID      string
	Path    *Path
}

// Draw adds a single tile from the game's tile pool to a player's hand.
func (g *Game) Draw(p *Player) error {
	if len(g.TilePool) == 0 {
		return errors.New("no tiles left")
	}

	t := g.TilePool[0]
	g.TilePool = g.TilePool[1:]
	p.Hand = append(p.Hand, t)

	return nil
}

// Place sets given Domino d from Player pl to the Path target if it fits.
func (g *Game) Place(pl *Player, d Domino, target *Path) bool {
	var last Element
	if len(target.Elements) == 0 {
		if !g.Center.IsPlayable(d) {
			return false
		}
	} else {
		last = target.Elements[len(target.Elements)-1]
		if !last.IsPlayable(d) {
			return false // Given domino d is not playable on the given Path.
		}
	}
	if target.Player != pl.ID && !target.Train && !target.MexicanTrain {
		return false // Cannot play on a train you don't own
	}

	e := Element{
		Domino:  d,
		Flipped: last.Left == d.Left || last.Right == d.Right,
	}

	target.Elements = append(target.Elements, e)

	return true
}

// Knock sets the knocked flag if a player has one tile left in their hand.
func (g *Game) Knock(p *Player) bool {
	if len(p.Hand) == 1 {
		p.Knocked = true
	}

	return p.Knocked
}

// NextTurn marks the next player as "up", adding two tiles to their hand if
// they only have one tile in their hand and haven't explicitly knocked.
func (g *Game) NextTurn() (*Player, string) {
	nextPlayer := (g.ActivePlayer + 1) % len(g.Players)
	p := g.Players[nextPlayer]
	g.ActivePlayer = nextPlayer
	status := ""
	if len(p.Hand) == 1 && !p.Knocked {
		g.Draw(p)
		g.Draw(p)
		p.Knocked = false
		status = "noknock"
	}

	return p, status
}

func handCount(playernum int) int {
	switch playernum {
	case 2:
		return 6
	case 3, 4:
		return 10
	case 5, 6:
		return 9
	case 7, 8:
		return 7
	default:
		return 6
	}
}

func dominoCount(playernum int) int {
	switch playernum {
	case 1, 2:
		return 6
	case 3, 4:
		return 9
	case 5, 6, 7, 8:
		return 12
	case 9, 10, 11, 12:
		return 15
	default:
		return 18
	}
}
