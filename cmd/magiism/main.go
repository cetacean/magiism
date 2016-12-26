package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/facebookgo/flagenv"
)

var (
	username = flag.String("username", "", "discord username")
	password = flag.String("password", "", "discord password")
)

func main() {
	flag.Parse()
	flagenv.Parse()

	d, err := discordgo.New(*username, *password)
	if err != nil {
		log.Fatal(err)
	}

	d.AddHandler(messageCreate)

	err = d.Open()
	if err != nil {
		log.Fatal(err)
	}

	runtime.Goexit()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Print message to stdout.
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)
}
