package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type Disconnect struct{}

func (_ *Disconnect) Commands() []string {
	return []string{"disconnect"}
}

func (_ *Disconnect) Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	p.Disconnect(m.Author)
}
