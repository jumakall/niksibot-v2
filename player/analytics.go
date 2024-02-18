package player

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

const (
	// AnalyticsServiceName is the name of the analytics service
	AnalyticsServiceName = "NiksiOnline"
)

type Analytics struct {
	// Endpoint is where analytics data is sent
	Endpoint string

	// ServerHostname name is included in analytics
	ServerHostname string
}

func InitializeAnalytics(endpoint string) *Analytics {
	hostname, err := os.Hostname()

	if err != nil {
		log.Warning(fmt.Sprintf("Failed to initialize %s", AnalyticsServiceName))
		return nil
	}

	log.Info(fmt.Sprintf("Connected to %s, statistics will be reported", AnalyticsServiceName))

	return &Analytics{
		Endpoint:       endpoint,
		ServerHostname: hostname,
	}
}

func (a *Analytics) Play(p *Play) {
	a.PublishData(p, "play")
}

func (a *Analytics) Skip(p *Play) {
	a.PublishData(p, "skip")
}

func (a *Analytics) PublishData(p *Play, action string) {
	data := map[string]string{
		"server": a.ServerHostname,
		"sound":  p.Sound.Name,
		"action": action,
	}

	log.WithFields(log.Fields{
		"server": a.ServerHostname,
		"sound":  p.Sound.Name,
		"action": action,
	}).Trace(fmt.Sprintf("Sending data to %s", AnalyticsServiceName))

	// convert data to json and send it to endpoint
	jsonEnc, _ := json.Marshal(data)
	response, err := http.Post(a.Endpoint+"/v1", "application/json", bytes.NewBuffer(jsonEnc))

	if err != nil {
		log.Warning("Failed to send analytics data")
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Warning("Failed to send analytics data")
	}
}
