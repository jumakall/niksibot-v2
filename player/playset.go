package player

import "container/list"

type PlaySet struct {
	queue *list.List
}

func CreatePlaySet(sound *Play) *PlaySet {
	ps := &PlaySet{
		queue: list.New(),
	}

	if sound != nil {
		ps.Push(sound)
	}

	return ps
}

func (ps PlaySet) Push(play *Play) {
	ps.queue.PushBack(play)
}

func (ps PlaySet) Peek() *Play {
	return ps.queue.Front().Value.(*Play)
}

func (ps PlaySet) Take() *Play {
	front := ps.queue.Front()
	ps.queue.Remove(front)
	return front.Value.(*Play)
}

func (ps PlaySet) IsExhausted() bool {
	return ps.queue.Len() <= 0
}

func (ps PlaySet) Length() int {
	return ps.queue.Len()
}

func (ps PlaySet) Clear() {
	ps.queue = list.New()
}
