package dummy

import (
	"log"
	"math/rand"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/util"
)

// Dummy service
type Dummy struct {
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
func NewDummy(verbose bool, logger *log.Logger) *Dummy {
	stopChan := make(chan bool, 1)
	return &Dummy{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *Dummy) Name() string {
	return "dummy"
}

// Ping will try to ping the service
func (d *Dummy) Ping() []byte {
	n := rand.Intn(100)
	if n < 50 {
		d.SetMessage("dummy has an error")
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Dummy) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *Dummy) Connect() error {
	if d.verbose {
		d.logger.Println("connect dummy")
	}

	return nil
}

// Close will close the service resources if needed
func (d *Dummy) Close() error {
	if d.verbose {
		d.logger.Println("close dummy")
	}

	return nil
}

// SetRecover will set recovered status
func (d *Dummy) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Dummy) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *Dummy) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *Dummy) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *Dummy) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *Dummy) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *Dummy) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *Dummy) IsEnabled() bool {
	return d.enabled
}

// SetMessage will set additional message
func (d *Dummy) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *Dummy) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *Dummy) SetConfig(configs config.Config) {

}

// SetNotificatorConfig will set config
func (d *Dummy) SetNotificatorConfig(configs config.Config) {
	d.notificatorConfig = configs
}

// GetNotificators will return notificators
func (d *Dummy) GetNotificators() []tob.Notificator {
	notificators, err := tob.InitNotificatorFactory(d.notificatorConfig, d.verbose)
	if err != nil {
		d.logger.Printf("Warning: %s service does not activate Notifications, GetNotificators() will be nil\n", d.Name())
		return nil
	}
	return notificators
}

// Stop will receive stop channel
func (d *Dummy) Stop() chan bool {
	return d.stopChan
}
