package sslstatus

import (
	"log"
	"strings"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/util"
)

// SSLStatus service
type SSLStatus struct {
	url               string
	recovered         bool
	lastDownTime      string
	enabled           bool
	verbose           bool
	logger            *log.Logger
	checkInterval     int
	stopChan          chan bool
	message           string
	configs           config.Config
	notificatorConfig config.Config
}

var SEVERITIES = []string{"Warning", "Danger", "Critical"}

// NewSSLStatus SSLStatus's constructor
func NewSSLStatus(verbose bool, logger *log.Logger) *SSLStatus {
	stopChan := make(chan bool, 1)
	return &SSLStatus{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *SSLStatus) Name() string {
	return "sslstatus"
}

// Ping will try to ping the service
func (d *SSLStatus) Ping() []byte {
	domains, ok := d.configs["domains"].([]interface{})
	if !ok {
		if d.verbose {
			d.logger.Println("domains is not in the SSL_status config")
		}
		return []byte("NOT_OK")
	}

	var domianStrs []string

	for _, domain := range domains {
		domainStr, ok := domain.(string)
		if ok {
			domianStrs = append(domianStrs, domainStr)
		}
	}

	sslStatusData := checkSSLExpiryMulti(domianStrs, d.logger)

	d.SetMessage(sslStatusData)

	if containsSeverity(sslStatusData, SEVERITIES) {
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *SSLStatus) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *SSLStatus) Connect() error {
	if d.verbose {
		d.logger.Println("connect SSLStatus")
	}

	return nil
}

// Close will close the service resources if needed
func (d *SSLStatus) Close() error {
	if d.verbose {
		d.logger.Println("close SSLStatus")
	}

	return nil
}

// SetRecover will set recovered status
func (d *SSLStatus) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *SSLStatus) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *SSLStatus) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *SSLStatus) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *SSLStatus) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *SSLStatus) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *SSLStatus) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *SSLStatus) IsEnabled() bool {
	return d.enabled
}

// SetMessage will set additional message
func (d *SSLStatus) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *SSLStatus) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *SSLStatus) SetConfig(configs config.Config) {
	d.configs = configs
}

// SetNotificatorConfig will set config
func (d *SSLStatus) SetNotificatorConfig(configs config.Config) {
	d.notificatorConfig = configs
}

// GetNotificators will return notificators
func (d *SSLStatus) GetNotificators() []tob.Notificator {
	return tob.InitNotificatorFactory(d.notificatorConfig, d.verbose)
}

// Stop will receive stop channel
func (d *SSLStatus) Stop() chan bool {
	return d.stopChan
}

func containsSeverity(sslStatusData string, severities []string) bool {
	for _, sev := range severities {
		if strings.Contains(sslStatusData, sev) {
			return true
		}
	}
	return false
}
