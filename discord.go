package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	log "github.com/sirupsen/logrus"
	"strings"
)

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Debug("Discord websocket connected")
	log.Info(fmt.Sprintf("%s is ready to serve", BotName))
}

func onMessageReceive(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{
		"guild":   m.GuildID,
		"channel": m.ChannelID,
		"author":  m.Author.Username + "#" + m.Author.Discriminator,
		"message": m.Content,
	}).Trace("Received a message")

	// make sure the message starts with specified prefix
	if len(m.Content) <= 0 || m.Content[:len(CommandPrefix)] != CommandPrefix || m.Author.ID == s.State.Ready.User.ID {
		return
	}

	parts := strings.SplitN(m.Content, " ", 2)

	channel, _ := s.State.Channel(m.ChannelID)
	if channel == nil {
		log.WithFields(log.Fields{
			"channel": m.ChannelID,
			"message": m.ID,
		}).Warning("Failed to grab channel")
		return
	}

	guild, _ := s.State.Guild(channel.GuildID)
	if guild == nil {
		log.WithFields(log.Fields{
			"guild":   channel.GuildID,
			"channel": channel.ID,
			"message": m.ID,
		}).Warning("Failed to grab guild")
		return
	}

	if Players[guild.ID] == nil {
		Players[guild.ID] = player.CreatePlayer(Discord, guild, &Sounds, PlayerHistoryLength)
	}
	player := Players[guild.ID]

	for _, c := range Commands {
		if Contains(parts[0][len(CommandPrefix):], c.Commands()) {
			c.Execute(s, guild, channel, m, player)
		}
	}
}

func Contains(str string, list []string) bool {
	for _, a := range list {
		if a == str {
			return true
		}
	}

	return false
}

func OpenDiscordWebsocket(token string) *discordgo.Session {
	discord, err := discordgo.New(token)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Failed to create Discord session")
	}

	discord.AddHandler(ready)
	discord.AddHandler(onMessageReceive)

	log.Debug("Discord websocket is connecting")
	err = discord.Open()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Discord websocket connection failed")
	}

	return discord
}
