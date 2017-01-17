package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cetacean/magiism/dominos"
	"github.com/cetacean/magiism/dominos/game"
)

func main() {
	gg, err := game.New([]string{"Xena", "Vic"})
	if err != nil {
		log.Fatal(err)
	}
	g := &wrapper{Game: gg}
	log.Printf("%s is the starting player!", g.GetActivePlayer().ID)
	for {
		err := g.Menu()
		if err != nil {
			switch err {
			case game.ErrEndOfTurn:
				continue
			default:
				log.Println(err)
			}
		}
	}
}

type wrapper struct {
	*game.Game
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

// End of turn sentry error
var (
	ErrEndOfTurn = errors.New("end of turn")
)

func (g *wrapper) Menu() error {
	if g.UnresolvedDouble {
		log.Println("Unresolved double")
	}

	p := g.GetActivePlayer()

	log.Printf("%s IS NOW UP", p.ID)
	log.Printf("CENTER PIECE: (%d, %d)\n", g.Center.Left, g.Center.Right)
	for i, e := range g.Trains {
		log.Printf("%d: %s", i, e.Display())
	}

	log.Println(g.GetActivePlayer().Display())
	log.Printf("Commands: (p)lace | (b)ig turn | (k)nock | (d)raw | (e)ndturn")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	for scanner.Scan() {
		log.Println(g.GetActivePlayer().Display())
		if scanner.Err() != nil {
			log.Println(scanner.Err())
			return scanner.Err()
		}

		t := scanner.Text()
		ev := &game.Event{
			PlayerID: p.ID,
		}

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
			ev.Action = game.PlayDomino
			ev.PathID = pIndexInt
			ev.HandIndex = hIndexInt

		case "b":
			log.Println("big turn not implemented")
		case "k":
			ev.Action = game.Knock
		case "d":
			ev.Action = game.DrawDomino
		case "e":
			ev.Action = game.EndTurn
		}

		resp, err := g.HandleEvent(ev)
		if err == game.ErrEndOfTurn {
		}
		if err != nil {
			switch err {
			case game.ErrEndOfTurn:
				log.Println(resp.GlobalMessage)
				log.Println("user: ", resp.UserMessage)
				return nil

			case dominos.ErrDontOwnPath:
				log.Println("You do not own the path you tried to play on and it is not marked to be playable on")
			case dominos.ErrNotPlayable:
				log.Println("That domino is unplayable on that path.")
			case dominos.ErrDanglingDouble:
				log.Println("There is a dangling double that must be resolved")

			default:
				return err
			}
		}

		fmt.Print("> ")
	}

	return nil
}
