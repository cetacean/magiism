package dominos

import "testing"

func TestNewGame(t *testing.T) {
	g := NewGame([]string{"Xena"})
	if g == nil {
		t.Fatalf("game didn't initialize somehow :(")
	}
}
