package commands

import (
	"../player"
	"github.com/bwmarrin/discordgo"
)

type DD struct{}

func (p DD) Commands() []string {
	return []string{"dd"}
}

func (_ DD) Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	p.ClearQueue()
	p.Disconnect()
}
