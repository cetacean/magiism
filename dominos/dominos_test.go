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
	prev := g.Players[g.ActivePlayer]
	p, _ := g.NextTurn()
	if p == prev {
		t.Fatal("Turn progression did not work")
	}

	if g.Players[g.ActivePlayer] != p {
		t.Fatal("g.ActivePlayer has not been updated")
	}
}
