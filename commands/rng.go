package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Rng struct{}

func (p *Rng) Commands() []string {
	return []string{"rng", "rng4ever"}
}

func (_ *Rng) Execute(s *discordgo.Session, g *discordgo.Guild, _ *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	parts := strings.SplitN(m.Content, " ", 2)

	if len(parts) < 2 {
		log.WithFields(log.Fields{
			"guild":   m.GuildID,
			"channel": m.ChannelID,
			"user":    m.Author.Username + "#" + m.Author.Discriminator,
			"message": m.Content,
			"reason":  "No tag was specified",
		}).Warning("Unable to create a play")
		return
	}

	voiceChannel := player.FindUsersVoiceChannel(s.State, g, m.Author)
	if voiceChannel == nil {
		log.WithFields(log.Fields{
			"guild":   m.GuildID,
			"channel": m.ChannelID,
			"user":    m.Author.Username + "#" + m.Author.Discriminator,
			"message": m.Content,
			"reason":  "The user hasn't connected to a voice channel",
		}).Warning("Unable to create a play")
		return
	}

	if parts[1] == "all" {
		var plays []*player.Play
		for _, sound := range *p.Sounds {
			plays = append(plays, player.CreatePlay(sound, m.Author, voiceChannel, g))
		}

		ps := player.CreatePlaySet(plays)
		ps.Shuffle()
		ps.ShuffleOnReset = true

		p.Playlist.SetFiller(ps)
		p.Playlist.Enqueue(ps)
		p.StartPlayback()
	}
}
