package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/commands"
	"github.com/jumakall/niksibot-v2/player"
)

func DiscoverCommands() []Command {
	return []Command{&commands.File{}, &commands.Skip{}, &commands.DD{}, &commands.Play{}, &commands.Rng{}, &commands.Clear{}}
}

type Command interface {
	Commands() []string
	Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player)
}
