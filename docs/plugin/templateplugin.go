package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/util"
)

func main() {}

// TemplatePlugin service
type TemplatePlugin struct {
	url               string
	recovered         bool
	lastDownTime      string
	enabled           bool
	verbose           bool
	logger            *log.Logger
	checkInterval     int
	stopChan          chan bool
	message           string
	notificatorConfig config.Config
}

// NewDummy Dummy's constructor
func NewTemplatePlugin(verbose bool, logger *log.Logger) *TemplatePlugin {
	stopChan := make(chan bool, 1)
	return &TemplatePlugin{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *TemplatePlugin) Name() string {
	return "plugin"
}

// Ping will try to ping the service
func (d *TemplatePlugin) Ping() []byte {
	n := rand.Intn(100)
	if n < 50 {
		d.SetMessage("dummy plugin has an error")
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *TemplatePlugin) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *TemplatePlugin) Connect() error {
	if d.verbose {
		d.logger.Println("connect dummy")
	}

	return nil
}

// Close will close the service resources if needed
func (d *TemplatePlugin) Close() error {
	if d.verbose {
		d.logger.Println("close dummy plugin")
	}

	return nil
}

// SetRecover will set recovered status
func (d *TemplatePlugin) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *TemplatePlugin) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *TemplatePlugin) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *TemplatePlugin) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *TemplatePlugin) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *TemplatePlugin) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *TemplatePlugin) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *TemplatePlugin) IsEnabled() bool {
	return d.enabled
}

// SetMessage will set additional message
func (d *TemplatePlugin) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *TemplatePlugin) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *TemplatePlugin) SetConfig(configs config.Config) {

}

// SetNotificatorConfig will set config
func (d *TemplatePlugin) SetNotificatorConfig(configs config.Config) {
	d.notificatorConfig = configs
}

// GetNotificators will return notificators
func (d *TemplatePlugin) GetNotificators() []tob.Notificator {
	return tob.InitNotificatorFactory(d.notificatorConfig, d.verbose)
}

// Stop will receive stop channel
func (d *TemplatePlugin) Stop() chan bool {
	if d.stopChan == nil {
		d.stopChan = make(chan bool, 1)
	}

	return d.stopChan
}

// Exported
var Service TemplatePlugin
