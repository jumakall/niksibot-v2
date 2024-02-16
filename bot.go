package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/commands"
	"github.com/jumakall/niksibot-v2/player"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// BotName specifies the name of the bot
	BotName = "NiksiBot"

	// DatabaseLocation specifies the location of SQLite database file
	DatabaseLocation = "./db.sqlite"

	// DatabaseSchemaVersion defines the schema version this build expects
	DatabaseSchemaVersion = 1

	// SoundsDirectory specifies the directory for sounds files
	SoundsDirectory = "audio"

	// SoundExtension is used to filter sounds in SoundsDirectory
	SoundExtension = ".dca"
)

var (
	// Sounds is a list of all sound files
	Sounds []*player.Sound

	// TagManager manages tag and sound relations
	TagManager *player.TagManager

	// Registrations contains Discord command definitions
	Registrations []*discordgo.ApplicationCommand

	// Commands is a list of all available commands
	Commands *map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, p *player.Player)

	// Discord is currently active session
	Discord *discordgo.Session

	// Players is a list of all players
	Players = make(map[string]*player.Player)
)

func firstStart() {

	if _, err := os.Stat(DatabaseLocation); !os.IsNotExist(err) {
		return
	}

	log.Info(fmt.Sprintf("Greetings! %s is doing a bit of preparation work since it is started for the first time", BotName))
	//createDB(DatabaseLocation)
}

func DiscoverSounds(path string) []*player.Sound {
	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Discovering sounds")

	var sounds []*player.Sound

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		// return in case of error
		if err != nil {
			return err
		}

		// filter folders and files with wrong extension
		if info.IsDir() || filepath.Ext(info.Name()) != SoundExtension {
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
		sound := player.CreateSound(name, info.Name(), trimmedPath)
		sounds = append(sounds, sound)

		autotag := strings.TrimLeft(trimmedPath, SoundsDirectory)
		autotag = autotag[1:]
		TagManager.TagSound(autotag, sound)

		return nil
	})
	if err != nil {
		log.WithFields(log.Fields{
			"path": path,
			"err":  err,
		}).Fatal("Sound discovery failed")
	}

	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Sound discovery completed")

	return sounds
}

func main() {
	var (
		Token         = flag.String("t", "", "Discord Bot Token")
		Verbose       = flag.Bool("v", false, "Verbose")
		DoubleVerbose = flag.Bool("vv", false, "More verbose")
		CStatus       = flag.String("status", "", "Custom status for the bot")
	)
	flag.Parse()

	if *Verbose || *DoubleVerbose {
		log.SetLevel(log.TraceLevel)
	}

	//firstStart()
	log.WithFields(log.Fields{
		"discordgo": discordgo.VERSION,
	}).Info(fmt.Sprintf("%s is starting", BotName))
	rand.Seed(time.Now().Unix())
	TagManager = player.CreateTagManager(&Sounds)
	Sounds = DiscoverSounds(SoundsDirectory)
	Registrations = commands.DiscoverRegistrations()
	Commands = commands.DiscoverCommands()
	Discord = OpenDiscordWebsocket(*Token)

	if *DoubleVerbose {
		Discord.LogLevel = discordgo.LogDebug
	}

	status := CreateStatus(Discord)
	if *CStatus != "" {
		status.Messages = append(status.Messages, *CStatus)
	}

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Info(fmt.Sprintf("%s is done for this time, see you later", BotName))
}
