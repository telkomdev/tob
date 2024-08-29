package server

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/dashboard/handler"
	"github.com/telkomdev/tob/dashboard/ui"
)

var (
	defaultDashboardHTTPPort = 9115
)

// HTTPServer struct
type HTTPServer struct {
	port    int
	logger  *log.Logger
	configs config.Config

	dashboardStaticAssets fs.FS
	dashboardHTTPHandler  *handler.DashboardHTTPHandler
}

func NewHTTPServer(configs config.Config, logger *log.Logger) (*HTTPServer, error) {
	// dashboard HTTP Port
	if parsedDashboardHTTPPort, ok := configs["dashboardHttpPort"].(float64); ok {
		defaultDashboardHTTPPort = int(parsedDashboardHTTPPort)
	}

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
	}, nil
}

// Run will Run Dashboard HTTP Server
func (s *HTTPServer) Run() {
	mux := http.NewServeMux()

	dashboardFs := http.FileServer(http.FS(s.dashboardStaticAssets))
	mux.Handle("/", http.StripPrefix("/", dashboardFs))

	mux.HandleFunc("/api/services", s.dashboardHTTPHandler.GetServices())
	mux.HandleFunc("/api/tob/webhook", s.dashboardHTTPHandler.HandleTobWebhook())

	log.Printf("Dashboard HTTP server running on port %d\n", s.port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.port), mux))

}

// Exit will exit and cleanup Dashboard HTTP Server
func (s *HTTPServer) Exit() {
	s.logger.Print("exiting Dashboard HTTP server\n")
}
