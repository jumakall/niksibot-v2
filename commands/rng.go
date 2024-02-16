package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	log "github.com/sirupsen/logrus"
)

type Rng struct{}

func (p *Rng) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "rng4ever",
			Description: "Loop all sounds with given tag forever",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tag",
					Description: "A tag that should be played and played and played....",
					Required:    true,
				},
			},
		},
	}
}
func (p *Rng) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"rng4ever": rngCommand,
	}
}

func rngCommand(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	user := i.Member.User
	options := i.ApplicationCommandData().Options
	tag := options[0].StringValue()

	guild, _ := s.State.Guild(i.GuildID)
	if guild == nil {
		return
	}

	voiceChannel := player.FindUsersVoiceChannel(s.State, guild, user)
	if voiceChannel == nil {
		log.WithFields(log.Fields{
			"guild":  guild.Name,
			"user":   user.Username,
			"reason": "The user hasn't connected to a voice channel",
		}).Warning("Unable to create a play")
		SendResponse(s, i, ":telephone: Connect to a voice channel and try again")
		return
	}

	if tag == "all" {
		var plays []*player.Play
		for _, sound := range *p.Sounds {
			plays = append(plays, player.CreatePlay(sound, user, voiceChannel, guild))
		}

		ps := player.CreatePlaySet(plays)
		ps.Shuffle()
		ps.ShuffleOnReset = true

		p.Playlist.SetFiller(ps)
		p.Playlist.Enqueue(ps)
		p.StartPlayback()

		SendResponse(s, i, ":loudspeaker: "+tag)
	} else {
		SendResponse(s, i, "Sorry, currently only \"all\" tag is supported.")
	}
}
