package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/dashboard/shared"
	"github.com/telkomdev/tob/dashboard/utils"
)

var (
	defaultDashboardTitle = "Tob Monitoring Dashboard"
)

// WebhookMessage type
type WebhookMessage struct {
	Message string `json:"message"`
}

// LoginPayload type
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// DashboardHTTPHandler type
type DashboardHTTPHandler struct {
	serviceData       map[string]map[string]interface{}
	logger            *log.Logger
	webhookTobTokens  []string
	dashboardTitle    string
	dashboardUsername string
	dashboardPassword string
}

// Data type
type Data struct {
	Data           map[string]map[string]interface{} `json:"data"`
	DashboardTitle string                            `json:"dashboardTitle"`
}

// LoginResponse type
type LoginResponse struct {
	Username  string `json:"username"`
	JWTString string `json:"jwtString"`
}

// NewDashboardHTTPHandler DashboardHTTPHandler's constructor
func NewDashboardHTTPHandler(tobConfig config.Config, logger *log.Logger) (*DashboardHTTPHandler, error) {
	if dashboardTitle, ok := tobConfig["dashboardTitle"].(string); ok {
		defaultDashboardTitle = dashboardTitle
	}

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

		serviceEnabled, ok := services["enable"].(bool)
		if !ok {
			return nil, errors.New("cannot convert service: enable to bool")
		}

		if !serviceEnabled {
			continue
		}

		// by default services status is UP
		services["status"] = "UP"
		services["url"] = ""
		serviceData[name] = services
	}

	dashboardUsername, ok := tobConfig["dashboardUsername"].(string)
	if !ok {
		return nil, errors.New("cannot parse dashboardUsername from configs")
	}

	dashboardPassword, ok := tobConfig["dashboardPassword"].(string)
	if !ok {
		return nil, errors.New("cannot parse dashboardPassword from configs")
	}

	return &DashboardHTTPHandler{
		dashboardTitle:    defaultDashboardTitle,
		serviceData:       serviceData,
		logger:            logger,
		webhookTobTokens:  webhookTobTokens,
		dashboardUsername: dashboardUsername,
		dashboardPassword: dashboardPassword,
	}, nil
}

// Login will handle user login
func (h *DashboardHTTPHandler) Login(jwtService utils.JwtService) http.HandlerFunc {
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

		var loginPayload LoginPayload

		err := json.NewDecoder(req.Body).Decode(&loginPayload)
		if err != nil {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    400,
				Message: "login payload is not valid",
				Data:    shared.EmptyJSON{},
			}, 400)
			return
		}

		hashedPassword, err := utils.Sha256Hex([]byte(loginPayload.Password))
		if err != nil {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    400,
				Message: "login payload is not valid",
				Data:    shared.EmptyJSON{},
			}, 400)
			return
		}

		if h.dashboardUsername != loginPayload.Username || h.dashboardPassword != hashedPassword {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: "username or password is not valid",
				Data:    shared.EmptyJSON{},
			}, 401)
			return
		}

		var claim utils.Claim
		claim.Alg = utils.HS256
		claim.Subject = h.dashboardUsername
		claim.User.ID = h.dashboardUsername
		claim.User.FullName = h.dashboardUsername
		claim.User.Email = h.dashboardUsername

		// 1 year
		jwtString, err := jwtService.Generate(&claim, time.Hour*8766)
		if err != nil {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: "error generating jwt",
				Data:    shared.EmptyJSON{},
			}, 401)
			return
		}

		shared.BuildJSONResponse(resp, shared.Response[LoginResponse]{
			Success: true,
			Code:    200,
			Message: "login succeed",
			Data: LoginResponse{
				Username:  h.dashboardUsername,
				JWTString: fmt.Sprintf("Bearer %s", jwtString),
			},
		}, 200)
	}
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

		data := Data{
			Data:           h.serviceData,
			DashboardTitle: h.dashboardTitle,
		}

		shared.BuildJSONResponse(resp, shared.Response[Data]{
			Success: true,
			Code:    200,
			Message: "get all services succeed",
			Data:    data,
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

		xTobTokenInConfig := func() bool {
			for _, token := range h.webhookTobTokens {
				if req.Header["X-Tob-Token"][0] == token {
					return true
				}
			}

			return false
		}

		if !xTobTokenInConfig() {
			shared.BuildJSONResponse(resp, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: "X-Tob-Token is not valid",
				Data:    shared.EmptyJSON{},
			}, 401)
			return
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
			h.logger.Print(messages)

			serviceName := strings.Trim(messages[0], " ")
			status := strings.Trim(regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(messages[2], ""), " ")

			service, ok := h.serviceData[serviceName]
			if ok {
				service["status"] = status
			}

			if len(messages) > 3 {
				messageDetails := strings.Join(messages[4:], " ")
				service["messageDetails"] = messageDetails
			}

			if status == "UP" {
				service["messageDetails"] = ""
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
