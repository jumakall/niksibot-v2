package commands

import (
	"../player"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type File struct{}

func (p File) Commands() []string {
	return []string{"f", "file"}
}

func (_ File) Execute(s *discordgo.Session, g *discordgo.Guild, _ *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	parts := strings.SplitN(m.Content, " ", 2)

	if len(parts) < 2 {
		log.WithFields(log.Fields{
			"guild":   m.GuildID,
			"channel": m.ChannelID,
			"user":    m.Author.Username + "#" + m.Author.Discriminator,
			"message": m.Content,
			"reason":  "No file was specified",
		}).Warning("Unable to create a play")
		return
	}

	foundFile := false
	for _, sound := range *p.Sounds {
		if sound.Name == parts[1] {
			foundFile = true

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

			play := player.CreatePlay(sound, m.Author, voiceChannel, g)
			play.Forced = true

			ps := player.CreatePlaySet(play)
			p.Enqueue(ps)
			p.StartPlayback()

			/*err := Discord.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Warning("Failed to remove message")
			}*/
		}
	}

	if !foundFile {
		log.WithFields(log.Fields{
			"guild":   m.GuildID,
			"channel": m.ChannelID,
			"user":    m.Author.Username + "#" + m.Author.Discriminator,
			"message": m.Content,
			"reason":  "No matching file for the request",
		}).Warning("Unable to create a play")
	}
}
