package web

import (
	"fmt"
	"log"
	"time"

	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
	"github.com/telkomdev/tob/util"
)

// Web service
type Web struct {
	url           string
	recovered     bool
	lastDownTime  string
	enabled       bool
	verbose       bool
	logger        *log.Logger
	checkInterval int
	stopChan      chan bool
	message       string
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
		d.SetMessage(err.Error())
		if d.verbose {
			d.logger.Printf("error: Ping() %s\n", err.Error())
		}
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		d.SetMessage(fmt.Sprintf("error: web Ping status: %d\n", resp.StatusCode))
		if d.verbose {
			d.logger.Printf("web Ping status: %d\n", resp.StatusCode)
		}

		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Printf("web Ping status: %d\n", resp.StatusCode)
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

// LastDownTime will set last down time of service to current time
func (d *Web) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *Web) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
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

// SetMessage will set additional message
func (d *Web) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *Web) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *Web) SetConfig(configs config.Config) {

}

// Stop will receive stop channel
func (d *Web) Stop() chan bool {
	return d.stopChan
}
