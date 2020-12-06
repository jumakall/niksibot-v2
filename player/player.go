package player

import (
	"container/list"
	log "github.com/Sirupsen/logrus"
	"github.com/bwmarrin/discordgo"
	"sync"
)

type Player struct {
	// Sounds is a list of all available sounds
	Sounds *[]*Sound

	// Discord references to the current Discord session
	Discord *discordgo.Session

	// Guild specifies to which guild this instance of player is bound to
	Guild *discordgo.Guild

	// VC is reference to currently active voice connection
	VC *discordgo.VoiceConnection

	// Queue is the current play queue
	Queue      *list.List
	NowPlaying *Play

	playerServiceLock    *sync.Mutex
	playerServiceRunning bool

	DisconnectPending bool
}

func CreatePlayer(discord *discordgo.Session, guild *discordgo.Guild, sounds *[]*Sound) *Player {
	log.WithFields(log.Fields{
		"guild": guild.Name,
	}).Trace("Instancing a player")

	return &Player{
		Sounds:               sounds,
		Discord:              discord,
		Guild:                guild,
		Queue:                list.New(),
		playerServiceLock:    &sync.Mutex{},
		playerServiceRunning: false,
		DisconnectPending:    false,
	}
}

func (p *Player) StartPlayback() {
	if p.Queue.Len() <= 0 {
		log.WithFields(log.Fields{
			"guild":  p.Guild.Name,
			"reason": "The queue is empty",
		}).Warning("Failed to start playback")
		return
	}

	p.playerServiceLock.Lock()
	if !p.playerServiceRunning {
		p.playerServiceRunning = true
		go p.backgroundPlayer()
	}
	p.playerServiceLock.Unlock()
}

func (p *Player) backgroundPlayer() {
	for p.Queue.Len() > 0 && !p.DisconnectPending {
		front := p.Queue.Front()
		p.Queue.Remove(front)
		p.NowPlaying = front.Value.(*Play)

		vc, err := p.connect(p.NowPlaying.Channel)
		if err != nil {
			break
		}

		log.WithFields(log.Fields{
			"guild":   p.NowPlaying.Guild.Name,
			"channel": p.NowPlaying.Channel.Name,
			"user":    p.NowPlaying.User.Username + "#" + p.NowPlaying.User.Discriminator,
			"sound":   p.NowPlaying.Sound.File,
			"forced":  p.NowPlaying.Forced,
		}).Info("Playing sound")

		err = p.NowPlaying.PlayToVoiceChannel(vc)
		if err != nil {
			log.WithFields(log.Fields{
				"guild":   p.NowPlaying.Guild.Name,
				"channel": p.NowPlaying.Channel.Name,
				"user":    p.NowPlaying.User.Username + "#" + p.NowPlaying.User.Discriminator,
				"sound":   p.NowPlaying.Sound.File,
				"err":     err,
			}).Warning("Error while playing sound")
		}
	}

	if p.DisconnectPending {
		p.disconnect()
	}

	p.playerServiceLock.Lock()
	p.playerServiceRunning = false
	p.playerServiceLock.Unlock()
}

func (p *Player) connect(voiceChannel *discordgo.Channel) (*discordgo.VoiceConnection, error) {
	if p.VC == nil {
		log.WithFields(log.Fields{
			"guild":   p.Guild.Name,
			"channel": voiceChannel.Name,
		}).Debug("Connecting to a voice channel")

		vc, err := p.Discord.ChannelVoiceJoin(p.Guild.ID, voiceChannel.ID, false, true)
		if err != nil {
			log.WithFields(log.Fields{
				"guild":   p.Guild.Name,
				"channel": voiceChannel.Name,
				"err":     err,
			}).Warning("Failed to connect to a voice channel")
			return nil, err
		}
		p.VC = vc
	} else if p.VC.ChannelID != voiceChannel.ID {
		log.WithFields(log.Fields{
			"guild":   p.Guild.Name,
			"channel": voiceChannel.Name,
		}).Debug("Changing voice channel")

		err := p.VC.ChangeChannel(voiceChannel.ID, false, true)
		if err != nil {
			log.WithFields(log.Fields{
				"guild":   p.Guild.Name,
				"channel": voiceChannel.Name,
				"err":     err,
			}).Warning("Failed to change voice channel")
			return nil, err
		}
	}

	return p.VC, nil
}

func (p *Player) Enqueue(play *Play) {
	if play == nil {
		log.WithFields(log.Fields{
			"guild": p.Guild.Name,
		}).Warning("Cannot queue null")
		return
	}

	log.WithFields(log.Fields{
		"guild":   p.Guild.Name,
		"channel": play.Channel.Name,
		"user":    play.User.Username + "#" + play.User.Discriminator,
		"file":    play.Sound.File,
		"forced":  play.Forced,
	}).Info("Queuing play")

	p.Queue.PushBack(play)
}

func (p *Player) Skip() {
	log.WithFields(log.Fields{
		"guild":   p.NowPlaying.Guild.Name,
		"channel": p.NowPlaying.Channel.Name,
		"sound":   p.NowPlaying.Sound.File,
	}).Info("Skipping play")
	p.NowPlaying.Skipped = true
}

func (p *Player) ClearQueue() {
	log.WithFields(log.Fields{
		"guild": p.Guild.Name,
	}).Info("Clearing play queue")
	p.Queue = list.New()
}

func (p *Player) Disconnect() {
	log.WithFields(log.Fields{
		"guild": p.Guild.Name,
	}).Trace("Disconnect from voice channel requested")

	if p.playerServiceRunning {
		p.DisconnectPending = true
		p.Skip()
	} else {
		p.disconnect()
	}
}

func (p *Player) disconnect() {
	log.WithFields(log.Fields{
		"guild": p.Guild.Name,
	}).Info("Disconnecting from voice channel")

	p.DisconnectPending = false
	err := p.VC.Disconnect()
	p.VC = nil

	if err != nil {
		log.WithFields(log.Fields{
			"guild": p.Guild.Name,
			"err":   err,
		}).Warning("Disconnect from voice channel failed")
	}
}

func (p *Player) FindSound(name string) *Sound {
	for _, s := range *p.Sounds {
		if s.Name == name {
			return s
		}
	}

	return nil
}
