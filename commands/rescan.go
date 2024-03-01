package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type Rescan struct{}

func (c *Rescan) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "rescan",
			Description: "Rescan library to detect changes",
		},
	}
}
func (c *Rescan) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"rescan": rescan,
	}
}

func rescan(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	p.Library.Rediscover()

	SendResponse(s, i, ":recycle:")
}
