package slack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
	"strings"
)

// https://api.slack.com/messaging/sending

// SlackMessage represent Slack request
type SlackMessage struct {
	Text string `json:"text"`
}

// SlackResponse represent Slack response
type SlackResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

// Slack represent Slack notificator
type Slack struct {
	webhookURL string
	headers    map[string]string
	mentions   []string
	enabled    bool
}

// NewSlack Slack's constructor
func NewSlack(configs config.Config) (*Slack, error) {
	notificatorConfigInterface, ok := configs["notificator"]
	if !ok {
		return nil, errors.New("error: cannot find notificator field in the config file")
	}

	notificatorConfig := notificatorConfigInterface.(map[string]interface{})

	slackConfigInterface, ok := notificatorConfig["slack"]
	if !ok {
		return nil, errors.New("error: cannot find Slack field in the config file")
	}

	slackConfig := slackConfigInterface.(map[string]interface{})

	webhookURL, ok := slackConfig["webhookUrl"].(string)
	if !ok {
		return nil, errors.New("error: cannot find Slack webhookUrl field in the config file")
	}

	mentionsInterface, ok := slackConfig["mentions"].([]interface{})
	if !ok {
		return nil, errors.New("error: cannot find Slack mentions field in the config file")
	}

	var mentions []string
	for _, mentionInterface := range mentionsInterface {
		mention, ok := mentionInterface.(string)
		if !ok {
			return nil, errors.New("error: mention field is not valid string")
		}

		mentions = append(mentions, mention)
	}

	enabled, ok := slackConfig["enable"].(bool)
	if !ok {
		return nil, errors.New("error: cannot find Slack enable field in the config file")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	return &Slack{
		webhookURL: webhookURL,
		enabled:    enabled,
		headers:    headers,
		mentions:   mentions,
	}, nil
}

// Provider will return Notificator provider
func (d *Slack) Provider() string {
	return "slack"
}

// Send will send notification
func (d *Slack) Send(msg string) error {
	var messageBuilder strings.Builder

	messageBuilder.WriteString("Hey ")
	for _, mention := range d.mentions {
		messageBuilder.WriteString(fmt.Sprintf("<%s>", mention))
		messageBuilder.WriteString(", ")
	}

	messageBuilder.WriteString(" ")
	messageBuilder.WriteString(msg)

	msg = messageBuilder.String()

	slackMessage := SlackMessage{
		Text: msg,
	}

	messageJSON, err := json.Marshal(slackMessage)
	if err != nil {
		return err
	}

	_, err = httpx.HTTPPost(d.webhookURL, bytes.NewBuffer(messageJSON), d.headers, 5)
	if err != nil {
		return err
	}

	return nil
}

// IsEnabled will return enable status
func (d *Slack) IsEnabled() bool {
	return d.enabled
}
