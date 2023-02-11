package email

import (
	"errors"
	"fmt"
	"github.com/telkomdev/tob/config"
	"net/smtp"
	"strings"
)

// Email represent email notificator
type Email struct {
	authEmail    string
	authPassword string
	//auth host, eg: smtp.gmail.com
	authHost string
	//address should include smtp provider port eg: "smtp.gmail.com:587" google smtp host
	smtpAddress string
	from        string
	to          []string
	subject     string
	enabled     bool
}

// NewEmail Email's constructor
func NewEmail(configs config.Config) (*Email, error) {
	notificatorConfigInterface, ok := configs["notificator"]
	if !ok {
		return nil, errors.New("error: cannot find notificator field in the config file")
	}

	notificatorConfig := notificatorConfigInterface.(map[string]interface{})

	emailConfigInterface, ok := notificatorConfig["email"]
	if !ok {
		return nil, errors.New("error: cannot find email field in the config file")
	}

	emailConfig := emailConfigInterface.(map[string]interface{})

	authEmail, ok := emailConfig["authEmail"].(string)
	if !ok {
		return nil, errors.New("error: cannot find authEmail field in the config file")
	}

	authPassword, ok := emailConfig["authPassword"].(string)
	if !ok {
		return nil, errors.New("error: cannot find authPassword field in the config file")
	}

	authHost, ok := emailConfig["authHost"].(string)
	if !ok {
		return nil, errors.New("error: cannot find authHost field in the config file")
	}

	smtpAddress, ok := emailConfig["smtpAddress"].(string)
	if !ok {
		return nil, errors.New("error: cannot find smtpAddress field in the config file")
	}

	from, ok := emailConfig["from"].(string)
	if !ok {
		return nil, errors.New("error: cannot find from field in the config file")
	}

	tosInterface, ok := emailConfig["to"].([]interface{})
	if !ok {
		return nil, errors.New("error: cannot find to field in the config file")
	}

	var toS []string
	for _, toInterface := range tosInterface {
		to, ok := toInterface.(string)
		if !ok {
			return nil, errors.New("error: to field is not valid string")
		}

		toS = append(toS, to)
	}

	subject, ok := emailConfig["subject"].(string)
	if !ok {
		return nil, errors.New("error: cannot find subject field in the config file")
	}

	enabled, ok := emailConfig["enable"].(bool)
	if !ok {
		return nil, errors.New("error: cannot find email enable field in the config file")
	}

	return &Email{
		authEmail:    authEmail,
		authPassword: authPassword,
		authHost:     authHost,
		smtpAddress:  smtpAddress,
		from:         from,
		to:           toS,
		subject:      subject,
		enabled:      enabled,
	}, nil
}

// Provider will return Notificator provider
func (*Email) Provider() string {
	return "email"
}

// Send will send notification
func (e *Email) Send(msg string) error {

	var messageBuilder strings.Builder

	messageBuilder.WriteString(fmt.Sprintf("Subject: %s!", e.subject))
	messageBuilder.WriteString("\r\n")
	messageBuilder.WriteString("MIME-version: 1.0;")
	messageBuilder.WriteString("\r\n")
	messageBuilder.WriteString("Content-Type: text/html; charset=\"UTF-8\";")
	messageBuilder.WriteString("\r\n")
	messageBuilder.WriteString("\r\n")

	messageBody := msg
	messageBuilder.WriteString(messageBody)
	messageBuilder.WriteString("\r\n")

	message := messageBuilder.String()

	auth := smtp.PlainAuth("", e.authEmail, e.authPassword, e.authHost)

	err := smtp.SendMail(e.smtpAddress, auth, e.from, e.to, []byte(message))
	if err != nil {
		return err
	}

	return nil
}

// IsEnabled will return enable status
func (e *Email) IsEnabled() bool {
	return e.enabled
}
