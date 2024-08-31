package server

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/dashboard/handler"
	"github.com/telkomdev/tob/dashboard/middleware"
	"github.com/telkomdev/tob/dashboard/ui"
	"github.com/telkomdev/tob/dashboard/utils"
)

var (
	defaultDashboardHTTPPort = 9115
)

// HTTPServer struct
type HTTPServer struct {
	port    int
	logger  *log.Logger
	configs config.Config

	jwtService utils.JwtService

	dashboardStaticAssets fs.FS
	dashboardHTTPHandler  *handler.DashboardHTTPHandler
}

func NewHTTPServer(configs config.Config, logger *log.Logger) (*HTTPServer, error) {
	// dashboard HTTP Port
	if parsedDashboardHTTPPort, ok := configs["dashboardHttpPort"].(float64); ok {
		defaultDashboardHTTPPort = int(parsedDashboardHTTPPort)
	}

	dashboardJwtKey, ok := configs["dashboardJwtKey"].(string)
	if !ok {
		return nil, errors.New("cannot parse dashboardJwtKey from configs")
	}

	jwtService := utils.NewJWT(dashboardJwtKey)

	dashboardStaticAssets, err := ui.Assets()
	if err != nil {
		return nil, err
	}

	dashboardHTTPHandler, err := handler.NewDashboardHTTPHandler(configs, logger)
	if err != nil {
		return nil, err
	}

	return &HTTPServer{
		configs:               configs,
		port:                  defaultDashboardHTTPPort,
		logger:                logger,
		dashboardStaticAssets: dashboardStaticAssets,
		dashboardHTTPHandler:  dashboardHTTPHandler,
		jwtService:            jwtService,
	}, nil
}

// Run will Run Dashboard HTTP Server
func (s *HTTPServer) Run() {
	mux := http.NewServeMux()

	dashboardFs := http.FileServer(http.FS(s.dashboardStaticAssets))

	var index = func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				dashboardFs.ServeHTTP(w, r)
				return
			}

			f, err := s.dashboardStaticAssets.Open(strings.TrimPrefix(path.Clean(r.URL.Path), "/"))
			if err == nil {
				defer f.Close()
			}
			if os.IsNotExist(err) {
				r.URL.Path = "/"
			}

			dashboardFs.ServeHTTP(w, r)
		})
	}

	// mux.Handle("/", http.StripPrefix("/", dashboardFs))
	mux.Handle("/", index())

	mux.HandleFunc("/api/login", s.dashboardHTTPHandler.Login(s.jwtService))
	mux.Handle("/api/services", middleware.JWTMiddleware(s.jwtService, s.dashboardHTTPHandler.GetServices()))
	mux.HandleFunc("/api/tob/webhook", s.dashboardHTTPHandler.HandleTobWebhook())

	log.Printf("Dashboard HTTP server running on port %d\n", s.port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux))

}

// Exit will exit and cleanup Dashboard HTTP Server
func (s *HTTPServer) Exit() {
	s.logger.Print("exiting Dashboard HTTP server\n")
}
