package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	log "github.com/sirupsen/logrus"
)

type File struct{}

func (_ *File) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "file",
			Description: "Queue a specific file",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "file",
					Description: "File to queue",
					Required:    true,
				},
			},
		},
	}
}
func (_ *File) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"file": fileCommand,
	}
}

func fileCommand(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	user := i.Member.User
	options := i.ApplicationCommandData().Options
	file := options[0].StringValue()

	guild, _ := s.State.Guild(i.GuildID)
	if guild == nil {
		return
	}

	sound := p.Library.FindSoundByName(file)
	if sound == nil {
		log.WithFields(log.Fields{
			"guild":  guild.Name,
			"user":   user.Username,
			"reason": "No matching file for the request",
		}).Warning("Unable to create a play")
		SendResponse(s, i, "Couldn't find \""+file+"\" :interrobang:")
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

	play := player.CreatePlay(sound, user, voiceChannel, guild)
	play.Forced = true
	ps := player.CreatePlaySet([]*player.Play{play})
	p.Playlist.Enqueue(ps)
	p.StartPlayback()

	SendResponse(s, i, ":loud_sound: "+sound.Name)
}
