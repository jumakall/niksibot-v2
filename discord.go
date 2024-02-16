package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	log "github.com/sirupsen/logrus"
)

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Debug("Discord websocket connected")

	// register commands
	for _, v := range Registrations {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", v)

		if err != nil {
			log.WithFields(log.Fields{
				"command": v.Name,
			}).Error("Failed to register command")
			continue
		}

		log.WithFields(log.Fields{
			"command": v.Name,
		}).Trace("Command registered")
	}

	log.Info(fmt.Sprintf("%s is ready to serve", BotName))
}

func onBotInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData().Name
	user := i.Member.User.Username

	// get guild
	guild, _ := s.State.Guild(i.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"Command": command,
			"User":    user,
		}).Warning("Failed to grab guild")
		return
	}

	// log interaction
	log.WithFields(log.Fields{
		"Command": command,
		"User":    user,
		"Guild":   guild.Name,
	}).Info("Interaction received")

	// find the guild's player or create a new one
	if Players[guild.ID] == nil {
		Players[guild.ID] = player.CreatePlayer(Discord, guild, &Sounds)
	}
	p := Players[guild.ID]

	// execute command
	f := (*Commands)[command]
	f(s, i, p)
}

func OpenDiscordWebsocket(token string) *discordgo.Session {
	discord, err := discordgo.New(token)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Failed to create Discord session")
	}

	discord.AddHandler(ready)
	discord.AddHandler(onBotInteraction)

	log.Debug("Discord websocket is connecting")
	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Discord websocket connection failed")
	}

	return discord
}
