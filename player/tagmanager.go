package player

type TagManager struct {
	Tags   map[string][]*Sound
	Sounds *[]*Sound
}

func CreateTagManager(sounds *[]*Sound) *TagManager {
	return &TagManager{
		Tags:   map[string][]*Sound{},
		Sounds: sounds,
	}
}

func (tm *TagManager) TagSound(tag string, s *Sound) {
	tm.Tags[tag] = append(tm.Tags[tag], s)
}

func (tm *TagManager) GetTag(tag string) []*Sound {
	if tag == "all" {
		return *(tm.Sounds)
	}

	return tm.Tags[tag]
}
