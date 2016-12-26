package dominos

import (
	"testing"

	"github.com/kr/pretty"
)

func TestNewGame(t *testing.T) {
	g := NewGame([]string{"Xena"})

	pretty.Println(g)
	pretty.Println(len(g.TilePool))
}
