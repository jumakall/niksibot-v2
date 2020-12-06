package player

import (
	"encoding/binary"
	"github.com/bwmarrin/discordgo"
	"io"
	"os"
)

type Play struct {
	Guild   *discordgo.Guild
	Channel *discordgo.Channel
	User    *discordgo.User
	Sound   *Sound

	// Forced indicates whether this is the exact sound requested (queued with !file command)
	Forced bool

	// Skipped indicates if this play was skipped
	Skipped bool
}

func CreatePlay(s *Sound, u *discordgo.User, c *discordgo.Channel, g *discordgo.Guild) *Play {
	return &Play{
		Guild:   g,
		Channel: c,
		User:    u,
		Sound:   s,
		Forced:  false,
		Skipped: false,
	}
}

func FindUsersVoiceChannel(discord *discordgo.State, guild *discordgo.Guild, user *discordgo.User) *discordgo.Channel {
	// try to find user from one of the guild's voice channels
	for _, vs := range guild.VoiceStates {
		if vs.UserID == user.ID {
			channel, _ := discord.Channel(vs.ChannelID)
			return channel
		}
	}

	return nil
}

func (p *Play) PlayToVoiceChannel(vc *discordgo.VoiceConnection) error {
	// open file
	file, err := os.Open(p.Sound.PathToFile())
	if err != nil {
		return err
	}

	// set speaking status
	vc.Speaking(true)
	defer vc.Speaking(false)

	var opusFrameLength int16
	for {
		// read opus frame length from dca file
		err = binary.Read(file, binary.LittleEndian, &opusFrameLength)

		// if this is the end of the file, just return
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}

		// fill buffer with data
		buf := make([]byte, opusFrameLength)
		err = binary.Read(file, binary.LittleEndian, &buf)
		if err != nil {
			return err
		}

		// send data to voice channel
		vc.OpusSend <- buf

		// break here if skip is pending
		if p.Skipped {
			break
		}
	}

	return nil
}
