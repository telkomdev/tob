package webhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"

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

// WebhookConfig represent webhook config
type WebhookConfig struct {
	webhookURL string
	headers    map[string]string
	enabled    bool
}

// Webhook represent Webhook notificator
type Webhook struct {
	configs []WebhookConfig
	logger  *log.Logger
	verbose bool
}

// NewWebhook Webhook's constructor
func NewWebhook(configs config.Config, verbose bool, logger *log.Logger) (*Webhook, error) {
	notificatorConfigInterface, ok := configs["notificator"]
	if !ok {
		return nil, errors.New("error: cannot find notificator field in the config file")
	}

	notificatorConfig := notificatorConfigInterface.(map[string]interface{})

	webhookConfigInterfaces, ok := notificatorConfig["webhook"]
	if !ok {
		return nil, errors.New("error: cannot find Webhook field in the config file")
	}

	var webhookConfigs []WebhookConfig
	webhookConfigList := webhookConfigInterfaces.([]interface{})
	for _, webhookConfigInterface := range webhookConfigList {
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

		conf := WebhookConfig{
			webhookURL: webhookURL,
			headers:    headers,
			enabled:    enabled,
		}

		webhookConfigs = append(webhookConfigs, conf)
	}

	return &Webhook{
		configs: webhookConfigs,
		logger:  logger,
		verbose: verbose,
	}, nil
}

// Provider will return Notificator provider
func (d *Webhook) Provider() string {
	return "webhook"
}

// Send will send notification
func (d *Webhook) Send(msg string) error {
	for _, conf := range d.configs {
		if conf.enabled {
			webhookMessage := WebhookMessage{
				Message: msg,
			}

			messageJSON, err := json.Marshal(webhookMessage)
			if err != nil {
				return err
			}

			go func(webhookURL string, body io.Reader, headers map[string]string, timeout int, d *Webhook) {
				resp, err := httpx.HTTPPost(webhookURL, body, headers, timeout)
				if err != nil {
					if d.verbose {
						d.logger.Println(err)
					}
				} else {
					statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
					if !statusOK {
						if d.verbose {
							d.logger.Printf("notificator %s error: %s\n", d.Provider(), httpx.ErrorStatusNot200)
						}
					}
				}
			}(conf.webhookURL, bytes.NewBuffer(messageJSON), conf.headers, 5, d)
		}
	}

	return nil
}

// IsEnabled will return enable status
func (d *Webhook) IsEnabled() bool {
	if len(d.configs) > 0 {
		temp := d.configs[0].enabled
		for _, conf := range d.configs {
			temp = temp || conf.enabled
		}
		return temp
	}
	return false
}
