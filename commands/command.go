package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

var (
	classes = []Command{&Play{}, &DD{}, &Skip{}, &File{}, &Clear{}, &Rng{}}
)

func DiscoverRegistrations() []*discordgo.ApplicationCommand {
	var registrations []*discordgo.ApplicationCommand

	for _, v := range classes {
		registrations = append(registrations, v.Register()...)
	}

	return registrations
}

func DiscoverCommands() *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	// create main command map
	commands := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){}

	// iterate all available classes
	for _, v := range classes {
		// get commands
		c := v.Commands()

		// combine map to main map
		for kk, vv := range c {
			commands[kk] = vv
		}
	}

	return &commands
}

type Command interface {
	Register() []*discordgo.ApplicationCommand
	Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player)
}

func SendResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
	if err != nil {
		return
	}
}
