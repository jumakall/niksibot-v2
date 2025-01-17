package player

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Discovery struct {
	Path      string
	Extension string
}

type Library struct {
	Folders []*Discovery
	Sounds  []*Sound
	Tags    map[string][]*Sound
}

func CreateLibrary() *Library {
	return &Library{
		Folders: []*Discovery{},
		Sounds:  []*Sound{},
		Tags:    map[string][]*Sound{},
	}
}

func (l *Library) Discover(path string, extension string) {
	l.Folders = append(l.Folders,
		&Discovery{
			Path:      path,
			Extension: extension,
		})

	l.discover(path, extension)
}

func (l *Library) discover(discoveryPath string, extension string) {
	log.WithFields(log.Fields{
		"path": discoveryPath,
	}).Debug("Discovering sounds")
	count := 0

	err := filepath.Walk(discoveryPath, func(path string, info os.FileInfo, err error) error {
		// return in case of error
		if err != nil {
			return err
		}

		// filter folders and files with wrong extension
		if info.IsDir() || filepath.Ext(info.Name()) != extension {
			return nil
		}

		path = strings.ReplaceAll(path, "\\", "/")
		name := info.Name()[:len(info.Name())-len(filepath.Ext(info.Name()))]
		trimmedPath := strings.TrimSuffix(path, info.Name())
		trimmedPath = trimmedPath[:len(trimmedPath)-1]

		// log found file
		log.WithFields(log.Fields{
			"name": name,
			"file": info.Name(),
			"path": trimmedPath,
			"size": info.Size(),
		}).Trace("Discovered sound")
		sound := CreateSound(name, info.Name(), trimmedPath)
		l.AddSound(sound)

		autotag := strings.TrimLeft(trimmedPath, discoveryPath)
		if autotag != "" {
			l.TagSound(autotag[1:], sound)
		}

		count++
		return nil
	})
	if err != nil {
		log.WithFields(log.Fields{
			"path": discoveryPath,
			"err":  err,
		}).Fatal("Sound discovery failed")
	}

	log.WithFields(log.Fields{
		"path":  discoveryPath,
		"count": count,
	}).Info("Sound discovery completed")
}

func (l *Library) Rediscover() {
	l.Sounds = []*Sound{}
	l.Tags = map[string][]*Sound{}

	for _, discovery := range l.Folders {
		l.discover(discovery.Path, discovery.Extension)
	}
}

func (l *Library) AddSound(s *Sound) {
	l.Sounds = append(l.Sounds, s)
}

func (l *Library) TagSound(tag string, s *Sound) {
	tag = strings.ToLower(tag)
	l.Tags[tag] = append(l.Tags[tag], s)
}

func (l *Library) GetSoundByTag(tag string) []*Sound {
	if strings.ToLower(tag) == "all" {
		return l.Sounds
	}

	return l.Tags[strings.ToLower(tag)]
}

func (l *Library) FindSoundByName(name string) *Sound {
	for _, s := range l.Sounds {
		if strings.ToLower(s.Name) == strings.ToLower(name) {
			return s
		}
	}

	return nil
}

func (l *Library) SearchFile(query string) *[]*Sound {
	// split query to search terms
	var searchTerms = strings.Split(strings.ToLower(query), " ")

	// list for results
	var results []*Sound

	// iterate all sounds
	for _, s := range l.Sounds {
		satisfy := true
		soundName := strings.ToLower(s.Name)

		// check if all search terms are satisfied
		for _, q := range searchTerms {
			if strings.Index(soundName, q) == -1 {
				satisfy = false
				break
			}
		}

		// if all search terms we satisfied, add to results
		if satisfy {
			results = append(results, s)
		}
	}

	return &results
}
