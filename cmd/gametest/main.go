package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cetacean/magiism/dominos"
)

func main() {
	g := &game{dominos.NewGame([]string{"Xena", "Vic"})}
	log.Printf("%s is the starting player!", g.GetActivePlayer().ID)
	for {
		g.Menu()

		p, status := g.NextTurn()
		if status == "noknock" {
			log.Printf("%s did not knock, drawn two tiles", p.ID)
		}
	}
}

type game struct {
	*dominos.Game
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

var (
	EndOfTurn = errors.New("end of turn")
)

func (g *game) Menu() error {
	log.Printf("%s IS NOW UP", g.GetActivePlayer().ID)
	log.Printf("CENTER PIECE: (%d, %d)\n", g.Center.Left, g.Center.Right)
	for i, e := range g.Trains {
		log.Printf("%d: %s", i, e.Display())
	}

	log.Printf("%s", g.GetActivePlayer().Display())
	log.Printf("Commands: (p)lace | (b)ig turn | (k)nock | (d)raw | (e)ndturn")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Println(scanner.Err())
			return scanner.Err()
		}

		t := scanner.Text()

		switch t {
		case "p":
			fmt.Printf("hand index to place> ")
			scanner.Scan()
			hIndex := scanner.Text()

			fmt.Printf("path index to play on> ")
			scanner.Scan()
			pIndex := scanner.Text()

			hIndexInt := atoi(hIndex)
			pIndexInt := atoi(pIndex)
			path := g.Trains[pIndexInt]
			p := g.GetActivePlayer()

			err := g.Place(p, p.RemoveFromHand(hIndexInt), path)
			if err != nil {
				log.Printf("could not place tile: %v", err)
			}

			return EndOfTurn
		case "b":
			log.Println("big turn not implemented")
		case "k":
			p := g.GetActivePlayer()
			if g.Knock(p) {
				log.Printf("%s has knocked, they only have one domino left!", p.ID)
			} else {
				log.Println("cannot knock, you have more than one tile in your hand")
			}
		case "d":
			log.Println("draw not implemented")
		case "e":
			return EndOfTurn
		default:
			log.Println("Command not understood, please try again.")
		}
	}

	return nil
}
