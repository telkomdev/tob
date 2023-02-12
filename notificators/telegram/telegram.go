package telegram

import (
	"errors"
	"fmt"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
	"net/url"
)

const (
	// TelegramAPIURL Telegram base URL API
	TelegramAPIURL = "https://api.telegram.org/bot%s/sendMessage?%s"
)

// Telegram represent Telegram notificator
type Telegram struct {
	botToken string
	groupID  string
	headers  map[string]string
	enabled  bool
}

// NewTelegram Telegram's constructor
func NewTelegram(configs config.Config) (*Telegram, error) {
	notificatorConfigInterface, ok := configs["notificator"]
	if !ok {
		return nil, errors.New("error: cannot find notificator field in the config file")
	}

	notificatorConfig := notificatorConfigInterface.(map[string]interface{})

	TelegramConfigInterface, ok := notificatorConfig["telegram"]
	if !ok {
		return nil, errors.New("error: cannot find Telegram field in the config file")
	}

	TelegramConfig := TelegramConfigInterface.(map[string]interface{})

	botToken, ok := TelegramConfig["botToken"].(string)
	if !ok {
		return nil, errors.New("error: cannot find Telegram botToken field in the config file")
	}

	groupID, ok := TelegramConfig["groupId"].(string)
	if !ok {
		return nil, errors.New("error: cannot find Telegram groupId field in the config file")
	}

	enabled, ok := TelegramConfig["enable"].(bool)
	if !ok {
		return nil, errors.New("error: cannot find Telegram enable field in the config file")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	return &Telegram{
		botToken: botToken,
		groupID:  groupID,
		enabled:  enabled,
		headers:  headers,
	}, nil
}

// Provider will return Notificator provider
func (d *Telegram) Provider() string {
	return "telegram"
}

// Send will send notification
func (d *Telegram) Send(msg string) error {
	encodedMessageVal := url.Values{}
	encodedMessageVal.Add("text", msg)
	encodedMessageVal.Add("chat_id", d.groupID)
	encodedMessageVal.Add("disable_web_page_preview", "true")

	telegramURL := fmt.Sprintf(TelegramAPIURL, d.botToken, encodedMessageVal.Encode())

	_, err := httpx.HTTPGet(telegramURL, d.headers, 5)
	if err != nil {
		return err
	}

	return nil
}

// IsEnabled will return enable status
func (d *Telegram) IsEnabled() bool {
	return d.enabled
}
