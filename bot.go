package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jumakall/niksibot-v2/player"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// BotName specifies the name of the bot
	BotName = "NiksiBot"

	// CommandPrefix specifies how all commands should begin
	CommandPrefix = "!"

	// DatabaseLocation specifies the location of SQLite database file
	DatabaseLocation = "./db.sqlite"

	// DatabaseSchemaVersion defines the schema version this build expects
	DatabaseSchemaVersion = 1

	// SoundsDirectory specifies the directory for sounds files
	SoundsDirectory = "./audio"

	// SoundExtension is used to filter sounds in SoundsDirectory
	SoundExtension = ".dca"
)

var (
	// Sounds is a list of all sound files
	Sounds []*player.Sound

	// Commands is a list of all available commands
	Commands []Command

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
	createDB(DatabaseLocation)
}

func DiscoverSounds(path string) []*player.Sound {
	log.WithFields(log.Fields{
		"path": path,
	}).Debug("Discovering sounds")

	var sounds []*player.Sound
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.WithFields(log.Fields{
			"path": path,
			"err":  err,
		}).Fatal("Sound discovery failed")
	}

	for _, f := range files {
		// filter files with wrong extension
		if filepath.Ext(f.Name()) != SoundExtension {
			continue
		}

		log.WithFields(log.Fields{
			"path": path,
			"file": f.Name(),
			"size": f.Size(),
		}).Trace("Discovered sound")

		name := f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))]
		sounds = append(sounds, player.CreateSound(name, f.Name(), path))
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

	firstStart()
	log.WithFields(log.Fields{
		"discordgo": discordgo.VERSION,
	}).Info(fmt.Sprintf("%s is starting", BotName))
	rand.Seed(time.Now().Unix())
	Sounds = DiscoverSounds(SoundsDirectory)
	Commands = DiscoverCommands()
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
