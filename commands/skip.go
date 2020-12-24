package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type Skip struct{}

func (p *Skip) Commands() []string {
	return []string{"s", "skip"}
}

func (_ *Skip) Execute(_ *discordgo.Session, _ *discordgo.Guild, _ *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	p.Skip(m.Author)
}
