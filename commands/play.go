package commands

import (
	"../player"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
)

type Play struct{}

func (_ *Play) Commands() []string {
	return []string{"p", "play"}
}

func (_ *Play) Execute(s *discordgo.Session, g *discordgo.Guild, _ *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
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
		soundInventory := *p.Sounds
		sound := soundInventory[rand.Intn(len(soundInventory))]

		play := player.CreatePlay(sound, m.Author, voiceChannel, g)
		ps := player.CreatePlaySet([]*player.Play{play})
		p.Enqueue(ps)
		p.StartPlayback()
	}
}
