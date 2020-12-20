package commands

import (
	"../player"
	"github.com/bwmarrin/discordgo"
)

type Disconnect struct{}

func (_ *Disconnect) Commands() []string {
	return []string{"disconnect"}
}

func (_ *Disconnect) Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	p.Disconnect()
}
