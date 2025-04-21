package sslstatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/data"
	"github.com/telkomdev/tob/httpx"
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

type target struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
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

// Resolve hostname to IPv4 address
func (d *SSLStatus) resolveIPv4() string {
	ipv4 := "N/A"

	hostUrl, err := url.Parse(d.url)
	if err != nil {
		if d.verbose {
			d.logger.Println(err)
		}
		return ipv4
	}
	hostname := hostUrl.Hostname()

	ips, err := net.LookupIP(hostname)
	if err != nil {
		if d.verbose {
			d.logger.Println(err)
		}
		return ipv4
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			ipv4 = ip.String()
			break
		}
	}

	return ipv4
}

// Ping will try to ping the service
func (d *SSLStatus) Ping() []byte {
	shellFilePathStr, ok := d.configs["shellFile"].(string)
	if !ok {
		if d.verbose {
			d.logger.Println("shellFilePathStr is not valid")
		}
		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Printf("tob-http-agent check %s SSL Status\n", shellFilePathStr)
	}

	fileSystemPayload := data.FileSystem{
		Path: shellFilePathStr,
	}

	payloadJSON, err := json.Marshal(fileSystemPayload)
	if err != nil {
		if d.verbose {
			d.logger.Println(err)
		}
		return []byte("NOT_OK")
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"

	resp, err := httpx.HTTPPost(fmt.Sprintf("%s/check-ssl", d.url), bytes.NewBuffer(payloadJSON), headers, 120)
	if err != nil {
		if d.verbose {
			d.logger.Println(err)
		}
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if d.verbose {
			d.logger.Printf("SSLStatus Ping status: %d\n", resp.StatusCode)
		}

		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Printf("SSLStatus Ping status: %d\n", resp.StatusCode)
	}

	defer func() { resp.Body.Close() }()

	var target target

	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		if d.verbose {
			d.logger.Println(err)
		}

		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Println(target)
	}

	sslStatusData := target.Data

	ipv4 := d.resolveIPv4()

	if d.verbose {
		d.logger.Println("IP :", ipv4)
		d.logger.Println("sslStatusData: ", sslStatusData)
	}

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
