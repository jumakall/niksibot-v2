package player

type PlaySet struct {
	original []*Play
	queue    []*Play
}

func CreatePlaySet(plays []*Play) *PlaySet {
	return &PlaySet{
		original: plays,
		queue:    plays,
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
	ps.queue = ps.original
	// reset must re-create all plays, so skips won't be remembered
}

func (ps *PlaySet) Length() int {
	return len(ps.queue)
}
