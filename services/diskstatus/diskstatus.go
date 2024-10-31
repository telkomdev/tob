package diskstatus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/url"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/data"
	"github.com/telkomdev/tob/httpx"
	"github.com/telkomdev/tob/util"
)

// DiskStatus service
type DiskStatus struct {
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
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// NewDiskStatus DiskStatus's constructor
func NewDiskStatus(verbose bool, logger *log.Logger) *DiskStatus {
	stopChan := make(chan bool, 1)
	return &DiskStatus{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *DiskStatus) Name() string {
	return "diskstatus"
}

// Resolve hostname to IPv4 address
func (d *DiskStatus) resolveIPv4() string {
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
func (d *DiskStatus) Ping() []byte {
	fileSystemPathStr, ok := d.configs["fileSystem"].(string)
	if !ok {
		if d.verbose {
			d.logger.Println("fileSystemPathStr is not valid")
		}
		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Printf("tob-http-agent check %s file system\n", fileSystemPathStr)
	}

	fileSystemPayload := data.FileSystem{
		Path: fileSystemPathStr,
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

	resp, err := httpx.HTTPPost(fmt.Sprintf("%s/check-disk", d.url), bytes.NewBuffer(payloadJSON), headers, 5)
	if err != nil {
		if d.verbose {
			d.logger.Println(err)
		}
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if d.verbose {
			d.logger.Printf("DiskStatus Ping status: %d\n", resp.StatusCode)
		}

		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Printf("DiskStatus Ping status: %d\n", resp.StatusCode)
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

	thresholdDiskUsage := util.InterfaceToFloat64(d.configs["thresholdDiskUsage"])

	diskUsed := util.InterfaceToFloat64(target.Data["diskUsed"])
	filesystem := target.Data["filesystem"]

	ipv4 := d.resolveIPv4()

	if d.verbose {
		d.logger.Println("IP :", ipv4)
		d.logger.Println("threshold disk usage: ", thresholdDiskUsage)
		d.logger.Println("disk used: ", diskUsed)
		d.logger.Println("file system: ", filesystem)
	}

	if diskUsed >= thresholdDiskUsage {
		d.SetMessage(fmt.Sprintf("disk used exceeds the threshold\nIP: %s\nthreshold: %d%s\ndisk used: %d%s\nfile system: %s\n%s",
			ipv4, int(thresholdDiskUsage), "%", int(diskUsed), "%", filesystem, "-------------------------------------"))
		return []byte("NOT_OK")
	}

	d.SetMessage(fmt.Sprintf("disk storage has been increased\nIP: %s\nthreshold: %d%s\ndisk used: %d%s\nfile system: %s\n%s",
		ipv4, int(thresholdDiskUsage), "%", int(diskUsed), "%", filesystem, "-------------------------------------"))
	return []byte("OK")
}

// SetURL will set the service URL
func (d *DiskStatus) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *DiskStatus) Connect() error {
	if d.verbose {
		d.logger.Println("connect DiskStatus")
	}

	return nil
}

// Close will close the service resources if needed
func (d *DiskStatus) Close() error {
	if d.verbose {
		d.logger.Println("close DiskStatus")
	}

	return nil
}

// SetRecover will set recovered status
func (d *DiskStatus) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *DiskStatus) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *DiskStatus) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *DiskStatus) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *DiskStatus) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *DiskStatus) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *DiskStatus) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *DiskStatus) IsEnabled() bool {
	return d.enabled
}

// SetMessage will set additional message
func (d *DiskStatus) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *DiskStatus) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *DiskStatus) SetConfig(configs config.Config) {
	d.configs = configs
}

// SetNotificatorConfig will set config
func (d *DiskStatus) SetNotificatorConfig(configs config.Config) {
	d.notificatorConfig = configs
}

// GetNotificators will return notificators
func (d *DiskStatus) GetNotificators() []tob.Notificator {
	notificators, err := tob.InitNotificatorFactory(d.notificatorConfig, d.verbose)
	if err != nil {
		d.logger.Printf("Warning: %s service does not activate Notifications, GetNotificators() will be nil\n", d.Name())
		return nil
	}
	return notificators
}

// Stop will receive stop channel
func (d *DiskStatus) Stop() chan bool {
	return d.stopChan
}
