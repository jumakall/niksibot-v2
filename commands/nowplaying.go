package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type NowPlaying struct{}

func (n NowPlaying) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "np",
			Description: "Displays currently playing sound",
		},
	}
}
func (n NowPlaying) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"np": np,
	}
}

func np(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	if p.Playlist.NowPlaying == nil {
		SendResponse(s, i, "Nothing's playing... :woman_facepalming:")
		return
	}

	SendResponse(s, i, ":mag_right: "+p.Playlist.NowPlaying.Sound.Name)
}
