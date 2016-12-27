package main

import (
	"bufio"
	"log"
	"os"

	"github.com/cetacean/magiism/dominos"
)

func main() {
	g := &game{dominos.NewGame([]string{"Xena", "Vic"})}
	for {
		g.Menu()
	}
}

type game struct {
	*dominos.Game
}

func (g *game) Menu() {
	log.Printf("%s is now up", g.Players[g.ActivePlayer].ID)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Println(scanner.Err())
			return
		}

		t := scanner.Text()

		switch t {
		case "place":
		case "knock":
		case "draw":
		case "endturn":
			g.NextTurn()
			return
		default:
			log.Println("Command not understood, please try again.")
		}
	}
}
