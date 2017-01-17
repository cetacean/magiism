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

// Emoji returns the emoji-fied version of the domino for Discord or slack.
func (d Domino) Emoji() string {
	return fmt.Sprintf("[:d%d:|:d%d:]", d.Left, d.Right)
}

// Display gives a human-readable version of this struct for debugging purposes.
func (d Domino) Display() string {
	if d.IsDouble() {
		return fmt.Sprintf("[%d||%d]", d.Left, d.Right)
	}
	return fmt.Sprintf("[%d|%d]", d.Left, d.Right)
}

// Game represents the total state for a single game
type Game struct {
	TilePool []Domino `json:"-"`
	Trains   []*Path
	Players  []*Player `json:"-"`
	Center   Domino

	UnresolvedDouble bool
	ActivePlayer     int
}

// Path represents a single player's path. If no player is set,
// the path is treated as the Mexican train.
type Path struct {
	Player   string
	Train    bool // If true, other players can play on it
	Elements []*Element

	UnresolvedDouble bool
	MexicanTrain     bool
}

// Display shows the player's hand for debugging purposes.
func (p *Player) Display() string {
	result := "YOUR HAND:"
	for i, e := range p.Hand {
		result = result + fmt.Sprintf(" %d:%s", i, e.Display())
	}
	return result
}

// EmojiHand returns the player's hand emoji-formatted for Discord or Slack.
func (p *Player) EmojiHand() string {
	result := "Your hand: "
	for i, e := range p.Hand {
		result += fmt.Sprintf(", %d: %s", i, e.Emoji())
	}
	return result
}

// Display a path for debugging purposes.
func (p *Path) Display() string {
	result := ""

	if p.MexicanTrain {
		result = result + "       M >>"
	} else {
		result = result + fmt.Sprintf("%8s >>", p.Player)
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

// IsPlayable checks if the given element and previous node are playable against
// a given domino.
func (e *Element) IsPlayable(prev *Element, next Domino) bool {
	// Okay so basically we're getting into graph theory for this so hold on to
	// your orange juice.
	lintersect := false // left intersection for e->prev
	rintersect := false // right intersection for e->prev

	// So from here we want to figure out where the previous domino intersects
	// this element.
	switch {
	case prev.Left == e.Left:
		lintersect = true
	case prev.Right == e.Left:
		lintersect = true
	case prev.Left == e.Right:
		rintersect = true
	case prev.Right == e.Right:
		rintersect = true
	}

	// Now we have to determine if the next domino matches the free side.
	switch {
	// The left side of the domino is physically taken, only the right side is
	// available for matching.
	case lintersect && (next.Right == e.Right):
		fallthrough
	case lintersect && (next.Left == e.Right):
		return true

		// The right side of the domino is physically taken, only the left side is
		// available for matching.
	case rintersect && (next.Right == e.Left):
		fallthrough
	case rintersect && (next.Left == e.Left):
		return true
	}

	return false
}

// Display gives a human-readable version of this struct for debugging purposes.
func (e *Element) Display() string {
	if e.Flipped {
		d := Domino{
			Left:  e.Right,
			Right: e.Left,
		}

		return d.Display()
	}

	return e.Domino.Display()
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

		for i := 0; i < hc; i++ {
			err := g.Draw(newPlayer)
			if err != nil {
				panic(err)
			}
		}
	}

	// Find largest piece in player hands
	// XXX TODO find a better way to do this?
	var largest Domino
	var starter int
	var handIndex int

	for i, player := range g.Players {
		for j, dom := range player.Hand {
			if dom.Left == dom.Right { // XXX maybe Domino.IsDouble :: Domino -> bool
				if dom.Left > largest.Left {
					largest = dom
					starter = i
					handIndex = j
				}
			}
		}
	}
	g.Center = largest
	g.ActivePlayer = starter
	g.GetActivePlayer().RemoveFromHand(handIndex)

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

func removeHandAtIndex(hand []Domino, i int) []Domino {
	hand[len(hand)-1], hand[i] = hand[i], hand[len(hand)-1]
	return hand[:len(hand)-1]
}

// RemoveFromHand when given index `at` will remove that element from the player's
// hand, returning it for future use.
func (p *Player) RemoveFromHand(at int) Domino {
	result := p.Hand[at]
	p.Hand = removeHandAtIndex(p.Hand, at)
	return result
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

// Placement errors
var (
	ErrNotPlayable    = errors.New("domino: domino is not playable on that path")
	ErrDontOwnPath    = errors.New("domino: path is not playable on by this player")
	ErrDanglingDouble = errors.New("domino: there is a dangling double that must be resolved")
)

// CanPlace returns an error if the given tile cannot be placed correctly and
// returns the element of the target path if it is playable
func (g *Game) CanPlace(pl *Player, d Domino, target *Path) (*Element, error) {
	// Ownership checks. Players can play on the target path if they own it, it
	// has a train on it or it is the mexican train.
	if target.Player != pl.ID && !target.Train && !target.MexicanTrain {
		return nil, ErrDontOwnPath
	}

	if g.UnresolvedDouble {
		if !target.UnresolvedDouble {
			return nil, ErrDanglingDouble
		}
	}

	// If the target path is empty, compare simply against the center tile
	// instead of checking for side matching.
	if len(target.Elements) == 0 {
		if g.Center.Left == d.Left || g.Center.Left == d.Right {
			return &Element{Domino: d}, nil
		}

		if g.Center.IsPlayable(d) {
			return &Element{Domino: d}, nil
		}

		return nil, ErrNotPlayable
	}

	var last *Element
	if len(target.Elements) >= 2 {
		last = target.Elements[len(target.Elements)-2]
	} else {
		last = &Element{Domino: g.Center}
	}
	curr := target.Elements[len(target.Elements)-1]

	if curr.IsPlayable(last, d) {
		return &Element{Domino: d}, nil
	}

	return nil, ErrNotPlayable
}

// Place sets given Domino d from Player pl to the Path target if it fits.
func (g *Game) Place(pl *Player, d Domino, target *Path) error {
	last, err := g.CanPlace(pl, d, target)
	if err != nil {
		return err
	}

	e := &Element{
		Domino: d,
	}

	target.Elements = append(target.Elements, e)

	// If the user has their train up and is playing on their own path, remove
	// the train from the player.
	if pl.ID == target.Player && pl.Path.Train {
		pl.Path.Train = false
	}

	// Check if the domino needs to be flipped or not
	if last.IsDouble() && d.Right == last.Right {
		e.Flipped = true
	}
	if last.IsDouble() && d.Left == last.Left {
		e.Flipped = true
	}

	if d.IsDouble() {
		g.UnresolvedDouble = true
		target.UnresolvedDouble = true
	}

	if last.IsDouble() && target.UnresolvedDouble && g.UnresolvedDouble {
		g.UnresolvedDouble = false
		target.UnresolvedDouble = false
	}

	return nil
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

// GetActivePlayer returns the currently active Player structure.
func (g *Game) GetActivePlayer() *Player {
	return g.Players[g.ActivePlayer]
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
