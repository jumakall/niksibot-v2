package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	log "github.com/sirupsen/logrus"
	"math/rand"
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
	fallbackUsed := false
	if sound == nil {

		sounds := *p.Library.SearchFile(file)
		if len(sounds) == 0 {
			SendResponse(s, i, "Couldn't find \""+file+"\" :interrobang:")
			return
		}

		sound = sounds[rand.Intn(len(sounds))]
		fallbackUsed = true
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

	if fallbackUsed {
		SendResponse(s, i, ":game_die: "+sound.Name)
	} else {
		SendResponse(s, i, ":loud_sound: "+sound.Name)
	}
}
