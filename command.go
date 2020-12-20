package main

import (
	"./commands"
	"./player"
	"github.com/bwmarrin/discordgo"
)

func DiscoverCommands() []Command {
	return []Command{&commands.File{}, &commands.Skip{}, &commands.DD{}, &commands.Play{}, &commands.Rng{}}
}

type Command interface {
	Commands() []string
	Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player)
}
