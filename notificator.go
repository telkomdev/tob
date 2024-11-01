package tob

import (
	"github.com/telkomdev/tob/config"

	"github.com/telkomdev/tob/notificators/discord"
	"github.com/telkomdev/tob/notificators/email"
	"github.com/telkomdev/tob/notificators/slack"
	"github.com/telkomdev/tob/notificators/telegram"
	"github.com/telkomdev/tob/notificators/webhook"
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
func InitNotificatorFactory(configs config.Config, verbose bool) []Notificator {
	// discord notificator
	discordNotificator, err := discord.NewDiscord(configs, verbose, Logger)
	if err != nil {
		discordNotificator = nil
	}

	// email notificator
	emailNotificator, err := email.NewEmail(configs)
	if err != nil {
		emailNotificator = nil
	}

	// slack notificator
	slackNotificator, err := slack.NewSlack(configs)
	if err != nil {
		slackNotificator = nil
	}

	// telegram notificator
	telegramNotificator, err := telegram.NewTelegram(configs)
	if err != nil {
		telegramNotificator = nil
	}

	// webhook notificator
	webhookNotificator, err := webhook.NewWebhook(configs, verbose, Logger)
	if err != nil {
		webhookNotificator = nil
	}

	notificators := []Notificator{
		emailNotificator,
		discordNotificator,
		slackNotificator,
		telegramNotificator,
		webhookNotificator,
	}

	return notificators
}
