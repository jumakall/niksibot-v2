package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type Clear struct{}

func (p *Clear) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "clear",
			Description: "Clear queue",
		},
	}
}

func (p *Clear) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"clear": clearCommand,
	}
}

func clearCommand(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	p.Playlist.SetFiller(nil)
	p.Playlist.Clear()
	p.Skip(i.Member.User)

	SendResponse(s, i, ":put_litter_in_its_place: Queue cleared")
}
