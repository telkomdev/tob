package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
	"strings"
)

// DiscordMessage represent discord request
type DiscordMessage struct {
	Content   string `json:"content"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

// DiscordResponse represent discord response
type DiscordResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

// Discord represent discord notificator
type Discord struct {
	threadURL string
	name      string
	avatarURL string
	headers   map[string]string
	mentions  []string
	enabled   bool
}

// NewDiscord Discord's constructor
func NewDiscord(configs config.Config) (*Discord, error) {
	notificatorConfigInterface, ok := configs["notificator"]
	if !ok {
		return nil, errors.New("error: cannot find notificator field in the config file")
	}

	notificatorConfig := notificatorConfigInterface.(map[string]interface{})

	discordConfigInterface, ok := notificatorConfig["discord"]
	if !ok {
		return nil, errors.New("error: cannot find discord field in the config file")
	}

	discordConfig := discordConfigInterface.(map[string]interface{})

	name, ok := discordConfig["name"].(string)
	if !ok {
		return nil, errors.New("error: cannot find discord name field in the config file")
	}

	threadURL, ok := discordConfig["url"].(string)
	if !ok {
		return nil, errors.New("error: cannot find discord url field in the config file")
	}

	avatarURL, ok := discordConfig["avatarUrl"].(string)
	if !ok {
		return nil, errors.New("error: cannot find discord avatarUrl field in the config file")
	}

	mentionsInterface, ok := discordConfig["mentions"].([]interface{})
	if !ok {
		return nil, errors.New("error: cannot find discord mentions field in the config file")
	}

	var mentions []string
	for _, mentionInterface := range mentionsInterface {
		mention, ok := mentionInterface.(string)
		if !ok {
			return nil, errors.New("error: mention field is not valid string")
		}

		mentions = append(mentions, mention)
	}

	enabled, ok := discordConfig["enable"].(bool)
	if !ok {
		return nil, errors.New("error: cannot find discord enable field in the config file")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	return &Discord{
		name:      name,
		threadURL: threadURL,
		avatarURL: avatarURL,
		enabled:   enabled,
		headers:   headers,
		mentions:  mentions,
	}, nil
}

// Provider will return Notificator provider
func (d *Discord) Provider() string {
	return "discord"
}

// Send will send notification
func (d *Discord) Send(msg string) error {
	var messageBuilder strings.Builder

	messageBuilder.WriteString("Hey ")
	for _, mention := range d.mentions {
		if strings.Contains(mention, "here") {
			messageBuilder.WriteString(fmt.Sprintf("%s", mention))
		} else {
			messageBuilder.WriteString(fmt.Sprintf("<%s>", mention))
		}

		messageBuilder.WriteString(", ")
	}

	messageBuilder.WriteString(" ")
	messageBuilder.WriteString(msg)

	msg = messageBuilder.String()

	discordMessage := DiscordMessage{
		Username:  d.name,
		AvatarURL: d.avatarURL,
		Content:   msg,
	}

	messageJSON, err := json.Marshal(discordMessage)
	if err != nil {
		return err
	}

	_, err = httpx.HTTPPost(d.threadURL, bytes.NewBuffer(messageJSON), d.headers, 5)
	if err != nil {
		return err
	}

	return nil
}

// IsEnabled will return enable status
func (d *Discord) IsEnabled() bool {
	return d.enabled
}
