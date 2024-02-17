package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

type Status struct {
	Session  *discordgo.Session
	Messages []string
	Active   int
}

func CreateStatus(discord *discordgo.Session) *Status {
	status := Status{
		Session: discord,
		Messages: []string{
			fmt.Sprintf("with the new %s", BotName),
			fmt.Sprintf("with version %s", CommitHash),
		},
		Active: 0,
	}

	discord.UpdateGameStatus(0, status.Messages[status.Active])

	go func() {
		for range time.Tick(60 * time.Second) {
			status.Active += 1
			if status.Active >= len(status.Messages) {
				status.Active = 0
			}

			discord.UpdateGameStatus(0, status.Messages[status.Active])
		}
	}()

	return &status
}
