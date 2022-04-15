package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
)

type Queue struct{}

func (_ *Queue) Commands() []string {
	return []string{"q", "queue"}
}

func (_ *Queue) Execute(s *discordgo.Session, g *discordgo.Guild, c *discordgo.Channel, m *discordgo.MessageCreate, p *player.Player) {
	if p.Playlist.IsQueueEmpty() {
		s.ChannelMessageSend(m.ChannelID, "Nothing's queued right now!")
		return
	}

	msg := ">>> :fast_forward: Upcoming plays:\n"
	i := 0
	for e := p.Playlist.Queue.Front(); e != nil; e = e.Next() {
		ps := e.Value.(*player.PlaySet)

		i++
		msg += fmt.Sprintf("%d. %s _(requested by %s)_\n", i, ps.Name, ps.User.Username+"#"+ps.User.Discriminator)
	}

	s.ChannelMessageSend(m.ChannelID, msg)
}
