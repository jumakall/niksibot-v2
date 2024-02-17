package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	"math/rand"
)

type Play struct{}

func (play *Play) Register() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Play a song",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tag",
					Description: "A random sound with given tag will be queued",
					Required:    true,
				},
			},
		},
	}
}
func (play *Play) Commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player){
		"play": PlayTag,
	}
}

func PlayTag(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player) {
	options := i.ApplicationCommandData().Options
	tag := options[0].StringValue()

	guild, _ := s.State.Guild(i.GuildID)
	if guild == nil {
		return
	}

	voiceChannel := player.FindUsersVoiceChannel(s.State, guild, i.Member.User)
	if voiceChannel == nil {
		return
	}

	soundInventory := p.TagManager.GetTag(tag)
	if soundInventory == nil {
		SendResponse(s, i, "Sorry, couldn't find anything")
		return
	}
	sound := soundInventory[rand.Intn(len(soundInventory))]

	play := player.CreatePlay(sound, i.Member.User, voiceChannel, guild)
	ps := player.CreatePlaySet([]*player.Play{play})
	p.Playlist.Enqueue(ps)
	p.StartPlayback()

	SendResponse(s, i, ":loud_sound: "+sound.Name)
}
