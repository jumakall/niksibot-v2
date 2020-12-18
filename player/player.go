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

	// Queue is the current PlaySet queue
	Queue *list.List

	// CurrentPlaySet is the PlaySet currently being played
	CurrentPlaySet *PlaySet

	// NowPlaying is the currently playing Play
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
		p.CurrentPlaySet = front.Value.(*PlaySet)

		for !p.CurrentPlaySet.IsExhausted() && !p.DisconnectPending {
			p.NowPlaying = p.CurrentPlaySet.Take()
			p.playSound(p.NowPlaying)
		}
	}

	p.NowPlaying = nil
	if p.DisconnectPending {
		p.disconnect()
	}

	p.playerServiceLock.Lock()
	p.playerServiceRunning = false
	p.playerServiceLock.Unlock()
}

func (p *Player) playSound(play *Play) {
	vc, err := p.connect(play.Channel)
	if err != nil {
		return
	}

	log.WithFields(log.Fields{
		"guild":   play.Guild.Name,
		"channel": play.Channel.Name,
		"user":    play.User.Username + "#" + play.User.Discriminator,
		"sound":   play.Sound.File,
		"forced":  play.Forced,
	}).Info("Playing sound")

	err = p.NowPlaying.PlayToVoiceChannel(vc)
	if err != nil {
		log.WithFields(log.Fields{
			"guild":   play.Guild.Name,
			"channel": play.Channel.Name,
			"user":    play.User.Username + "#" + play.User.Discriminator,
			"sound":   play.Sound.File,
			"err":     err,
		}).Warning("Error while playing sound")
	}
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

func (p *Player) Enqueue(playSet *PlaySet) {
	if playSet == nil {
		log.WithFields(log.Fields{
			"guild": p.Guild.Name,
		}).Warning("Cannot queue null")
		return
	}

	logger := log.WithFields(log.Fields{
		"guild": p.Guild.Name,
	})

	if playSet.Length() == 1 {
		play := playSet.Peek()
		logger.WithFields(log.Fields{
			"channel": play.Channel.Name,
			"user":    play.User.Username + "#" + play.User.Discriminator,
			"file":    play.Sound.File,
			"forced":  play.Forced,
		}).Info("Queuing play")
	} else {
		logger.WithFields(log.Fields{
			"count": playSet.Length(),
		}).Info("Queuing multiple plays")
	}

	p.Queue.PushBack(playSet)
}

func (p *Player) Skip() {
	if p.NowPlaying == nil {
		return
	}

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
	}).Info("Clearing queue")
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
