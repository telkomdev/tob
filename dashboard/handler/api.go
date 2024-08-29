package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/dashboard/shared"
)

// WebhookMessage type
type WebhookMessage struct {
	Message string `json:"message"`
}

// DashboardHTTPHandler type
type DashboardHTTPHandler struct {
	tobConfig        config.Config
	serviceData      map[string]map[string]interface{}
	logger           *log.Logger
	webhookTobTokens []string
}

// NewDashboardHTTPHandler DashboardHTTPHandler's constructor
func NewDashboardHTTPHandler(tobConfig config.Config, logger *log.Logger) (*DashboardHTTPHandler, error) {
	notificatorConfigInterface, ok := tobConfig["notificator"]
	if !ok {
		return nil, errors.New("notificator key from tob config is undefined")
	}

	notificators, ok := notificatorConfigInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("cannot convert tob config notificator to map")
	}

	webhookNotificatorInterfaces, ok := notificators["webhook"]
	if !ok {
		return nil, errors.New("webhook notificator key is not in config")
	}

	var webhookTobTokens []string
	webhookConfigList := webhookNotificatorInterfaces.([]interface{})
	for _, webhookConfigInterface := range webhookConfigList {
		webhookConfig := webhookConfigInterface.(map[string]interface{})

		logger.Print("webhook key")
		logger.Print(webhookConfig)

		webhookTobToken, ok := webhookConfig["tobToken"].(string)
		if !ok {
			return nil, errors.New("cannot convert webhookConfig tobToken to string")
		}

		webhookTobTokens = append(webhookTobTokens, strings.Trim(webhookTobToken, " "))
	}

	serviceConfigInterface, ok := tobConfig["service"]
	if !ok {
		return nil, errors.New("service key from tob config is undefined")
	}

	services, ok := serviceConfigInterface.(map[string]interface{})
	if !ok {
		return nil, errors.New("cannot convert tob config services to map")
	}

	serviceData := make(map[string]map[string]interface{})

	for name, serviceInteface := range services {
		services, ok := serviceInteface.(map[string]interface{})
		if !ok {
			return nil, errors.New("cannot convert tob service to map")
		}

		services["status"] = "UP"
		serviceData[name] = services
	}

	return &DashboardHTTPHandler{
		tobConfig:        tobConfig,
		serviceData:      serviceData,
		logger:           logger,
		webhookTobTokens: webhookTobTokens,
	}, nil
}

// GetServices will return tob services
func (h *DashboardHTTPHandler) GetServices() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodGet {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    405,
				Message: "http method not valid",
				Data:    shared.EmptyJSON{},
			}, 405)
			return
		}

		shared.BuildJSONResponse(resp, shared.Response[map[string]map[string]interface{}]{
			Success: true,
			Code:    200,
			Message: "get all services succeed",
			Data:    h.serviceData,
		}, 200)
	}
}

// HandleTobWebhook will handle webhook that send by Tob
func (h *DashboardHTTPHandler) HandleTobWebhook() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		if req.Method != http.MethodPost {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    405,
				Message: "http method not valid",
				Data:    shared.EmptyJSON{},
			}, 405)
			return
		}

		if len(req.Header["X-Tob-Token"]) <= 0 {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: "X-Tob-Token cannot be empty",
				Data:    shared.EmptyJSON{},
			}, 401)
			return
		}

		for _, token := range h.webhookTobTokens {
			if req.Header["X-Tob-Token"][0] != token {
				shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
					Success: false,
					Code:    401,
					Message: "X-Tob-Token is not valid",
					Data:    shared.EmptyJSON{},
				}, 401)
				return
			}
		}

		var message WebhookMessage

		err := json.NewDecoder(req.Body).Decode(&message)
		if err != nil {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    400,
				Message: "webhook payload is not valid",
				Data:    shared.EmptyJSON{},
			}, 400)
			return
		}

		messages := strings.Split(message.Message, " ")
		if len(messages) > 0 {
			h.logger.Print("--------------------")
			h.logger.Print(messages[0])

			serviceName := strings.Trim(messages[0], " ")
			status := strings.Trim(messages[2], " ")

			service, ok := h.serviceData[serviceName]
			if ok {
				service["status"] = status
			}
		}

		shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
			Success: true,
			Code:    200,
			Message: "handle tob webhook succeed",
			Data:    shared.EmptyJSON{},
		}, 200)
	}
}
