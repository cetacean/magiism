package dominos

import "testing"

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
			Elements: []Element{{
				Domino: Domino{
					Left:  6,
					Right: 1,
				},
			}},
			Player: "A",
		},
	}
	p := g.GetActivePlayer()
	p.ID = "A"
	did := g.Place(p, Domino{1, 4}, g.Trains[0])
	if !did {
		t.Fatal("could not place domino")
	}

	t.Logf("%s", g.Trains[0].Display())
}
