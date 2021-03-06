package player

import "math/rand"

type PlaySet struct {
	original       []*Play
	queue          []*Play
	ShuffleOnReset bool
}

func CreatePlaySet(plays []*Play) *PlaySet {
	return &PlaySet{
		original:       plays,
		queue:          plays,
		ShuffleOnReset: false,
	}
}

func (ps *PlaySet) Peek() *Play {
	return ps.queue[0]
}

func (ps *PlaySet) Take() *Play {
	play := ps.queue[0]
	ps.queue = ps.queue[1:]
	return play
}

func (ps *PlaySet) IsExhausted() bool {
	return len(ps.queue) <= 0
}

func (ps *PlaySet) Reset() {
	newQueue := make([]*Play, len(ps.original))
	for i := 0; i < len(ps.original); i++ {
		play := ps.original[i]
		newQueue[i] = CreatePlay(play.Sound, play.User, play.Channel, play.Guild)
	}
	ps.queue = newQueue

	if ps.ShuffleOnReset {
		ps.Shuffle()
	}
}

func (ps *PlaySet) Length() int {
	return len(ps.queue)
}

func (ps *PlaySet) Shuffle() {
	a := ps.queue
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
}
