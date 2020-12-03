package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	log "github.com/Sirupsen/logrus"
)

const (
	// BotName specifies the name of the bot
	BotName = "NiksiBot"

	// DatabaseLocation specifies the location of SQLite database file
	DatabaseLocation = "./db.sqlite"

	// DatabaseSchemaVersion defines the schema version this build expects
	DatabaseSchemaVersion = 1
)

func main() {
	var (
		Verbose = flag.Bool("v", false, "Verbose")
	)
	flag.Parse()

	if *Verbose {
		log.SetLevel(log.TraceLevel)
	}

	if _, err := os.Stat(DatabaseLocation); os.IsNotExist(err) {
		log.Info(fmt.Sprintf("Greetings! %s is doing a bit of preparation work since it is started for the first time", BotName))
		createDB(DatabaseLocation)
	}

	log.Info(fmt.Sprintf("%s is starting", BotName))

	// Wait for a signal to quit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	log.Info(fmt.Sprintf("%s is done for this time, see you later", BotName))
}
