package main

import (
	"bufio"
	"log"
	"os"

	"github.com/cetacean/magiism/dominos"
)

func main() {
	g := &game{dominos.NewGame([]string{"Xena", "Vic"})}
	log.Printf("%s is the starting player!", g.Players[g.ActivePlayer].ID)
	for {
		g.Menu()
	}
}

type game struct {
	*dominos.Game
}

func (g *game) Menu() {
	log.Printf("%s IS NOW UP", g.Players[g.ActivePlayer].ID)
	log.Printf("CENTER PIECE: (%d, %d)\n", g.Center.Left, g.Center.Right)
	for _, e := range g.Trains {
		log.Println(e.Display())
	}

	log.Printf("%s", g.Players[g.ActivePlayer].Display())
	log.Printf("Commands: (p)lace | (b)ig turn | (k)nock | (d)raw | (e)ndturn")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Println(scanner.Err())
			return
		}

		t := scanner.Text()

		switch t {
		case "p":
			log.Println("place not implemented")
		case "b":
			log.Println("big turn not implemented")
		case "k":
			log.Println("knock not implemented")
		case "d":
			log.Println("draw not implemented")
		case "e":
			p, status := g.NextTurn()
			if status == "noknock" {
				log.Printf("%s did not knock, drawn two tiles", p.ID)
			}
			return
		default:
			log.Println("Command not understood, please try again.")
		}
	}
}
