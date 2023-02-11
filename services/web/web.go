package web

import (
	"fmt"
	"github.com/telkomdev/tob/httpx"
	"log"
)

// Web service
type Web struct {
	url           string
	recovered     bool
	enabled       bool
	verbose       bool
	logger        *log.Logger
	checkInterval int
	stopChan      chan bool
}

// NewWeb Web's constructor
func NewWeb(verbose bool, logger *log.Logger) *Web {
	stopChan := make(chan bool, 1)
	return &Web{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *Web) Name() string {
	return "web"
}

// Ping will try to ping the service
func (d *Web) Ping() []byte {
	resp, err := httpx.HTTPGet(d.url, nil, 5)
	if err != nil {
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if d.verbose {
			d.logger.Println(fmt.Sprintf("web Ping status: %d", resp.StatusCode))
		}

		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Println(fmt.Sprintf("web Ping status: %d", resp.StatusCode))
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Web) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *Web) Connect() error {
	if d.verbose {
		d.logger.Println("connect Web")
	}

	return nil
}

// Close will close the service resources if needed
func (d *Web) Close() error {
	if d.verbose {
		d.logger.Println("close Web")
	}

	return nil
}

// SetRecover will set recovered status
func (d *Web) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Web) IsRecover() bool {
	return d.recovered
}

// SetCheckInterval will set check interval to service
func (d *Web) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *Web) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *Web) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *Web) IsEnabled() bool {
	return d.enabled
}

// Stop will receive stop channel
func (d *Web) Stop() chan bool {
	return d.stopChan
}
