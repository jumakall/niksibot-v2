package commands

import (
	"../player"
	"github.com/bwmarrin/discordgo"
)

type Clear struct{}

func (p *Clear) Commands() []string {
	return []string{"clear"}
}

func (_ *Clear) Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	p.Playlist.SetFiller(nil)
	p.Playlist.Clear()
	p.Skip(m.Author)
}
