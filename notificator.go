package tob

import (
	"github.com/telkomdev/tob/config"

	"github.com/telkomdev/tob/notificators/discord"
	"github.com/telkomdev/tob/notificators/email"
	"github.com/telkomdev/tob/notificators/slack"
	"github.com/telkomdev/tob/notificators/telegram"
)

// Notificator the notificator base
type Notificator interface {
	// Provider will return Notificator provider
	Provider() string

	// Send will send message to Notificator
	Send(msg string) error

	// IsEnabled will return enable status
	IsEnabled() bool
}

// InitNotificatorFactory will init all notificator
func InitNotificatorFactory(configs config.Config) ([]Notificator, error) {
	// discord notificator
	discordNotificator, err := discord.NewDiscord(configs)
	if err != nil {
		return nil, err
	}

	// email notificator
	emailNotificator, err := email.NewEmail(configs)
	if err != nil {
		return nil, err
	}

	// slack notificator
	slackNotificator, err := slack.NewSlack(configs)
	if err != nil {
		return nil, err
	}

	// telegram notificator
	telegramNotificator, err := telegram.NewTelegram(configs)
	if err != nil {
		return nil, err
	}

	notificators := []Notificator{
		emailNotificator,
		discordNotificator,
		slackNotificator,
		telegramNotificator,
	}

	return notificators, nil
}
