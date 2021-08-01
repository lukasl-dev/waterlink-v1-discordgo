package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
)

var (
	token     = os.Getenv("TOKEN")
	sessionID string

	session *discordgo.Session
)

// initDiscordGo initializes discordgo.
func initDiscordGo() {
	createSession()
	registerHandlers()
	openSession()
}

// createSession creates a new discordgo session.
func createSession() {
	var err error
	session, err = discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		log.Fatalln("Opening discordgo session failed:", err)
	}
	log.Println("Discordgo session created.")
}

// openSession opens the discordgo session.
func openSession() {
	if err := session.Open(); err != nil {
		log.Fatalln("Opening discordgo session failed:", err)
	}
	log.Println("Discordgo session opened.")
}
