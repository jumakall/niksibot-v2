package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/getsentry/sentry-go"
	"github.com/jumakall/niksibot-v2/commands"
	"github.com/jumakall/niksibot-v2/player"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
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
	// SentryDSN defines where to send error details
	SentryDSN = ""

	// Environment information is included in error reports
	Environment = "development"

	// CommitHash is the commit from which the app was built from
	CommitHash = "dev"

	// Release is included in error reports
	Release = "niksibot-v2@" + CommitHash

	// AnalyticsEndpoint is where analytics data is sent
	AnalyticsEndpoint = ""

	// Library manages sound and tag information
	Library *player.Library

	// Analytics sends statistics to remote endpoint
	Analytics *player.Analytics

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

func main() {
	// enable error reporting to Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         SentryDSN,
		Environment: Environment,
		Release:     Release,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()

	log.WithFields(log.Fields{
		"sentry-sdk": sentry.SDKVersion,
	}).Info("Error reporting enabled, error and context data will be uploaded in case of an error.")

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
		"discordgo":   discordgo.VERSION,
		"environment": Environment,
		"commit":      CommitHash,
		"release":     Release,
	}).Info(fmt.Sprintf("%s is starting", BotName))
	rand.Seed(time.Now().Unix())
	Library = player.CreateLibrary()
	Library.Discover(SoundsDirectory, SoundExtension)
	Registrations = commands.DiscoverRegistrations()
	Commands = commands.DiscoverCommands()
	Analytics = player.InitializeAnalytics(AnalyticsEndpoint)
	Discord = OpenDiscordWebsocket(*Token)

	if *DoubleVerbose {
		Discord.LogLevel = discordgo.LogDebug
	}

	status := CreateStatus(Discord)
	status.Messages = append(status.Messages, "with "+strconv.Itoa(len(Library.Sounds))+" sounds")
	if *CStatus != "" {
		status.Messages = append(status.Messages, *CStatus)
	}

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Info(fmt.Sprintf("%s is done for this time, see you later", BotName))
}
