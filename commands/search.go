package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	"strings"
)

type Search struct{}

func (c *Search) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "search",
			Description: "Search files from the library",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "Search query",
					Required:    true,
				},
			},
		},
	}
}
func (c *Search) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"search": search,
	}
}

func search(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	options := i.ApplicationCommandData().Options
	query := options[0].StringValue()

	results := p.Library.SearchFile(query)

	if len(*results) == 0 {
		SendResponse(s, i, ":mailbox_with_no_mail: What a shame, no mail")
	}

	var resultsString []string
	for _, s := range *results {
		resultsString = append(resultsString, "* "+s.Name)
	}

	SendResponse(s, i, ":mailbox_with_mail: Your magazine, sir...\n"+strings.Join(resultsString, "\n"))
}
