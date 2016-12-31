package dominos

import (
	"fmt"
	"testing"

	"github.com/kr/pretty"
)

func TestNewGame(t *testing.T) {
	g := NewGame([]string{"Xena"})
	if g == nil {
		t.Fatalf("game didn't initialize somehow :(")
	}
}

func TestEndTurn(t *testing.T) {
	g := NewGame([]string{"Xena", "Vic"})
	prev := g.GetActivePlayer()
	p, _ := g.NextTurn()
	if p == prev {
		t.Fatal("Turn progression did not work")
	}

	if g.GetActivePlayer() != p {
		t.Fatalf("g.ActivePlayer has not been updated, %s", p.Display())
	}
}

func TestRemoveFromHand(t *testing.T) {
	g := NewGame([]string{"A", "B"})
	p := g.GetActivePlayer()
	d := p.RemoveFromHand(0)
	t.Logf("Removed %s from %s's hand", d.Display(), p.ID)
}

func TestCantDraw(t *testing.T) {
	g := NewGame([]string{"A", "B"})
	g.TilePool = nil
	err := g.Draw(g.GetActivePlayer())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestPlace(t *testing.T) {
	g := NewGame([]string{"A", "B"})
	g.Trains = []*Path{
		&Path{
			Elements: []*Element{{
				Domino: Domino{
					Left:  6,
					Right: 1,
				},
			}},
			Player: "A",
		},
	}
	g.Center = Domino{6, 6}
	p := g.GetActivePlayer()
	p.ID = "A"
	err := g.Place(p, Domino{1, 4}, g.Trains[0])
	if err != nil {
		pretty.Println(g)
		t.Fatalf("could not place %v", err)
	}

	t.Logf("%s", g.Trains[0].Display())
}

func TestIsPlayable(t *testing.T) {
	cases := []struct {
		d1, d2     Domino
		shouldwork bool
	}{
		{
			d1: Domino{
				Left:  0,
				Right: 1,
			},
			d2: Domino{
				Left:  1,
				Right: 5,
			},
			shouldwork: true,
		},
		{
			d1: Domino{
				Left:  5,
				Right: 5,
			},
			d2: Domino{
				Left:  5,
				Right: 2,
			},
			shouldwork: true,
		},
	}

	for _, tcase := range cases {
		t.Run(fmt.Sprintf("%s %s %v", tcase.d1.Display(), tcase.d2.Display(), tcase.shouldwork), func(t *testing.T) {
			if tcase.d1.IsPlayable(tcase.d2) != tcase.shouldwork {
				t.Fatalf("%s %s %v but expected %v", tcase.d1.Display(), tcase.d2.Display(), !tcase.shouldwork, tcase.shouldwork)
			}
		})
	}
}
