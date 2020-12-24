package player

import (
	"container/list"
	log "github.com/sirupsen/logrus"
)

type Playlist struct {
	// Queue is the current PlaySet queue
	Queue *list.List

	// CurrentPlaySet is the PlaySet currently being played
	CurrentPlaySet *PlaySet

	// NowPlaying is the currently playing Play
	NowPlaying *Play

	// Filler is being played on loop when queue is empty
	filler *PlaySet
}

func CreatePlaylist() *Playlist {
	return &Playlist{
		Queue:  list.New(),
		filler: nil,
	}
}

func (pl *Playlist) IsPlaylistEmpty() bool {
	return pl.IsQueueEmpty() && pl.filler == nil
}

func (pl *Playlist) IsQueueEmpty() bool {
	return pl.Queue.Len() <= 0
}

func (pl *Playlist) AdvancePlaySet() *PlaySet {
	if pl.Queue.Len() > 0 {
		front := pl.Queue.Front()
		pl.Queue.Remove(front)
		pl.CurrentPlaySet = front.Value.(*PlaySet)
	} else if pl.filler != nil {
		if pl.filler.IsExhausted() {
			pl.filler.Reset()
		}

		pl.CurrentPlaySet = pl.filler
	} else {
		log.Warning("Trying to advance to the next playset when playlist is empty")
	}

	return pl.CurrentPlaySet
}

func (pl *Playlist) Advance() *Play {
	if pl.CurrentPlaySet == nil || pl.CurrentPlaySet.IsExhausted() || (pl.CurrentPlaySet == pl.filler && !pl.IsQueueEmpty()) {
		pl.AdvancePlaySet()
	}

	pl.NowPlaying = pl.CurrentPlaySet.Take()
	return pl.NowPlaying
}

func (pl *Playlist) Enqueue(ps *PlaySet) {
	if ps == nil {
		log.Warning("Trying to enqueue empty playset")
		return
	}

	if ps.Length() == 1 {
		play := ps.Peek()
		log.WithFields(log.Fields{
			"guild":   play.Guild.Name,
			"channel": play.Channel.Name,
			"user":    play.User.Username + "#" + play.User.Discriminator,
			"file":    play.Sound.File,
			"forced":  play.Forced,
		}).Info("Queuing play")
	} else {
		play := ps.Peek()
		log.WithFields(log.Fields{
			"guild": play.Guild.Name,
			"size":  ps.Length(),
		}).Info("Queuing PlaySet")
	}

	pl.Queue.PushBack(ps)
}

func (pl *Playlist) SetFiller(ps *PlaySet) {
	pl.filler = ps
}

func (pl *Playlist) Clear() {
	log.Info("Clearing queue")
	pl.Queue = list.New()
	pl.CurrentPlaySet = nil
}

func (pl *Playlist) Stop() {
	pl.NowPlaying = nil
}
