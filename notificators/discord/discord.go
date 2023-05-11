package discord

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
	"io"
	"log"
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

// DiscordConfig represent discord config
type DiscordConfig struct {
	threadURL string
	name      string
	avatarURL string
	headers   map[string]string
	mentions  []string
	enabled   bool
}

// Discord represent discord notificator
type Discord struct {
	configs []DiscordConfig
	logger  *log.Logger
	verbose bool
}

// NewDiscord Discord's constructor
func NewDiscord(configs config.Config, verbose bool, logger *log.Logger) (*Discord, error) {
	notificatorConfigInterface, ok := configs["notificator"]
	if !ok {
		return nil, errors.New("error: cannot find notificator field in the config file")
	}

	notificatorConfig := notificatorConfigInterface.(map[string]interface{})

	discordConfigInterfaces, ok := notificatorConfig["discord"]
	if !ok {
		return nil, errors.New("error: cannot find discord field in the config file")
	}

	var discordConfigs []DiscordConfig
	discordConfigList := discordConfigInterfaces.([]interface{})
	for _, discordConfigInterface := range discordConfigList {
		discordConfig, ok := discordConfigInterface.(map[string]interface{})
		if !ok {
			return nil, errors.New("error: discord config is not valid")
		}

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

		conf := DiscordConfig{
			name:      name,
			threadURL: threadURL,
			avatarURL: avatarURL,
			enabled:   enabled,
			headers:   headers,
			mentions:  mentions,
		}

		discordConfigs = append(discordConfigs, conf)
	}

	return &Discord{
		configs: discordConfigs,
		logger:  logger,
		verbose: verbose,
	}, nil
}

// Provider will return Notificator provider
func (d *Discord) Provider() string {
	return "discord"
}

// Send will send notification
func (d *Discord) Send(msg string) error {
	for _, conf := range d.configs {
		if conf.enabled {
			var messageBuilder strings.Builder

			messageBuilder.WriteString("Hey ")
			for _, mention := range conf.mentions {
				if strings.Contains(mention, "here") {
					messageBuilder.WriteString(fmt.Sprintf("%s", mention))
				} else {
					messageBuilder.WriteString(fmt.Sprintf("<%s>", mention))
				}

				messageBuilder.WriteString(", ")
			}

			messageBuilder.WriteString(" ")
			messageBuilder.WriteString(msg)

			content := messageBuilder.String()

			discordMessage := DiscordMessage{
				Username:  conf.name,
				AvatarURL: conf.avatarURL,
				Content:   content,
			}

			messageJSON, err := json.Marshal(discordMessage)
			if err != nil {
				return err
			}

			go func(threadURL string, body io.Reader, headers map[string]string, timeout int) {
				_, err = httpx.HTTPPost(threadURL, body, headers, timeout)
				if err != nil {
					if d.verbose {
						d.logger.Println(err)
					}
				}
			}(conf.threadURL, bytes.NewBuffer(messageJSON), conf.headers, 5)
		}
	}

	return nil
}

// IsEnabled will return enable status
func (d *Discord) IsEnabled() bool {
	if len(d.configs) > 0 {
		temp := d.configs[0].enabled
		for _, conf := range d.configs {
			temp = temp || conf.enabled
		}
		return temp
	}
	return false
}
