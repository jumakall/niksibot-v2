package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	"math/rand"
)

type DD struct{}

func (p *DD) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "dd",
			Description: "Stop playback, clear queue and disconnect",
		},
	}
}
func (p *DD) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"dd": dd,
	}
}

func dd(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	p.Playlist.SetFiller(nil)
	p.Playlist.Clear()
	p.Disconnect(i.Member.User)

	if rand.Intn(10) == 0 {
		SendResponse(s, i, ":middle_finger:")
	} else {
		SendResponse(s, i, ":japanese_goblin:")
	}

}
