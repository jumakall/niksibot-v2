package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type Skip struct{}

func (p *Skip) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "skip",
			Description: "Skip currently playing sound",
		},
	}
}
func (p *Skip) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"skip": skip,
	}
}

func skip(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	p.Skip(i.Member.User)

	SendResponse(s, i, ":no_entry_sign: "+p.Playlist.NowPlaying.Sound.Name)
}
