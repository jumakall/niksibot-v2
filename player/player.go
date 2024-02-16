package player

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Player struct {
	// Sounds is a list of all available sounds
	Sounds *[]*Sound

	// TagManager manages tag and sound relations
	TagManager *TagManager

	// Discord references to the current Discord session
	Discord *discordgo.Session

	// Guild specifies to which guild this instance of player is bound to
	Guild *discordgo.Guild

	// VC is reference to currently active voice connection
	VC *discordgo.VoiceConnection

	// Playlist manages PlaySets that Player receives
	Playlist *Playlist

	playerServiceLock    *sync.Mutex
	playerServiceRunning bool

	DisconnectPending bool
}

func CreatePlayer(discord *discordgo.Session, guild *discordgo.Guild, sounds *[]*Sound, tagManager *TagManager) *Player {
	log.WithFields(log.Fields{
		"guild": guild.Name,
	}).Trace("Instancing a player")

	return &Player{
		Sounds:               sounds,
		TagManager:           tagManager,
		Discord:              discord,
		Guild:                guild,
		Playlist:             CreatePlaylist(),
		playerServiceLock:    &sync.Mutex{},
		playerServiceRunning: false,
		DisconnectPending:    false,
	}
}

func (p *Player) StartPlayback() {
	if p.Playlist.IsPlaylistEmpty() {
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
	for !p.Playlist.IsPlaylistEmpty() && !p.DisconnectPending {
		play := p.Playlist.Advance()
		p.playSound(play)
	}

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

	err = play.PlayToVoiceChannel(vc)
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

func (p *Player) Skip(actor *discordgo.User) {
	np := p.Playlist.NowPlaying
	if np == nil {
		return
	}

	logger := log.WithFields(log.Fields{})
	if actor != nil {
		logger = log.WithFields(log.Fields{
			"actor": actor.Username + "#" + actor.Discriminator,
		})
	}

	logger.WithFields(log.Fields{
		"guild":   np.Guild.Name,
		"channel": np.Channel.Name,
		"sound":   np.Sound.File,
	}).Info("Skipping play")
	np.Skipped = true
}

func (p *Player) Disconnect(actor *discordgo.User) {
	log.WithFields(log.Fields{
		"guild": p.Guild.Name,
	}).Trace("Disconnect from voice channel requested")

	if p.playerServiceRunning {
		p.DisconnectPending = true
		p.Skip(actor)
	} else {
		p.disconnect()
	}
}

func (p *Player) disconnect() {
	if p.VC == nil {
		log.WithFields(log.Fields{
			"guild": p.Guild.Name,
		}).Trace("Trying to disconnect when not connected to any voice channel")
		return
	}

	log.WithFields(log.Fields{
		"guild": p.Guild.Name,
	}).Info("Disconnecting from voice channel")

	p.DisconnectPending = false
	err := p.VC.Disconnect()
	p.VC = nil
	p.Playlist.Stop()

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
