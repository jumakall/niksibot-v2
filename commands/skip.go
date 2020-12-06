package commands

import (
	"../player"
	"github.com/bwmarrin/discordgo"
)

type Skip struct{}

func (p Skip) Commands() []string {
	return []string{"skip"}
}

func (_ Skip) Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	p.Skip()
}
