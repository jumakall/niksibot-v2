package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type History struct{}

func (_ *History) Commands() []string {
	return []string{"h", "history"}
}

func (_ *History) Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	if p.History[0] == nil {
		s.ChannelMessageSend(m.ChannelID, "Nothing's played for a while!")
		return
	}

	msg := ">>> :rewind: Recently played:\n"
	for i, el := range p.History {
		if el == nil {
			break
		}

		decorator := ""
		if el.Forced {
			decorator += "**"
		}
		if el.Skipped {
			decorator += "~~"
		}

		msg += fmt.Sprintf("%s%d. %s%s\n", decorator, i+1, el.Sound.Name, Reverse(decorator))
	}

	s.ChannelMessageSend(m.ChannelID, msg)
}

// Reverse the string
// Source: https://stackoverflow.com/questions/1752414/how-to-reverse-a-string-in-go
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
