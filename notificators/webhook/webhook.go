package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
)

// WebhookMessage represent Webhook request
type WebhookMessage struct {
	Message string `json:"message"`
}

// WebhookResponse represent Webhook response
type WebhookResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

// Webhook represent Webhook notificator
type Webhook struct {
	webhookURL string
	headers    map[string]string
	enabled    bool
}

// NewWebhook Webhook's constructor
func NewWebhook(configs config.Config) (*Webhook, error) {
	notificatorConfigInterface, ok := configs["notificator"]
	if !ok {
		return nil, errors.New("error: cannot find notificator field in the config file")
	}

	notificatorConfig := notificatorConfigInterface.(map[string]interface{})

	webhookConfigInterface, ok := notificatorConfig["webhook"]
	if !ok {
		return nil, errors.New("error: cannot find Webhook field in the config file")
	}

	webhookConfig := webhookConfigInterface.(map[string]interface{})

	webhookURL, ok := webhookConfig["url"].(string)
	if !ok {
		return nil, errors.New("error: cannot find Webhook url field in the config file")
	}

	tobToken, ok := webhookConfig["tobToken"].(string)
	if !ok {
		return nil, errors.New("error: cannot find Webhook tobToken field in the config file")
	}

	enabled, ok := webhookConfig["enable"].(bool)
	if !ok {
		return nil, errors.New("error: cannot find Webhook enable field in the config file")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	// set x-tob-token header with tobToken
	headers["x-tob-token"] = tobToken

	return &Webhook{
		webhookURL: webhookURL,
		enabled:    enabled,
		headers:    headers,
	}, nil
}

// Provider will return Notificator provider
func (d *Webhook) Provider() string {
	return "webhook"
}

// Send will send notification
func (d *Webhook) Send(msg string) error {
	webhookMessage := WebhookMessage{
		Message: msg,
	}

	messageJSON, err := json.Marshal(webhookMessage)
	if err != nil {
		return err
	}

	resp, err := httpx.HTTPPost(d.webhookURL, bytes.NewBuffer(messageJSON), d.headers, 5)
	if err != nil {
		return err
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		return httpx.ErrorStatusNot200
	}

	return nil
}

// IsEnabled will return enable status
func (d *Webhook) IsEnabled() bool {
	return d.enabled
}
