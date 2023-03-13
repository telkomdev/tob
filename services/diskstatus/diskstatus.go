package diskstatus

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
	"github.com/telkomdev/tob/util"
)

// DiskStatus service
type DiskStatus struct {
	url           string
	recovered     bool
	serviceName   string
	lastDownTime  string
	enabled       bool
	verbose       bool
	logger        *log.Logger
	checkInterval int
	stopChan      chan bool
	message       string
	configs       config.Config
}

type target struct {
	Status  bool                   `json:"status"`
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

// Ping will try to ping the service
func (d *DiskStatus) Ping() []byte {
	resp, err := httpx.HTTPGet(fmt.Sprintf("%s/check-disk", d.url), nil, 5)
	if err != nil {
		if d.verbose {
			d.logger.Println(err)
		}
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if d.verbose {
			d.logger.Println(fmt.Sprintf("DiskStatus Ping status: %d", resp.StatusCode))
		}

		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Println(fmt.Sprintf("DiskStatus Ping status: %d", resp.StatusCode))
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

	if d.verbose {
		d.logger.Println("threshold disk usage: ", thresholdDiskUsage)
		d.logger.Println("disk used: ", diskUsed)
		d.logger.Println("file system: ", filesystem)
	}

	if diskUsed >= thresholdDiskUsage {
		d.SetMessage(fmt.Sprintf("disk used exceeds the threshold\nthreshold: %d%s\ndisk used: %d%s\nfile system: %s\n%s",
			int(thresholdDiskUsage), "%", int(diskUsed), "%", filesystem, "-------------------------------------"))
		return []byte("NOT_OK")
	}

	d.SetMessage(fmt.Sprintf("disk storage has been increased\nthreshold: %d%s\ndisk used: %d%s\nfile system: %s\n%s",
		int(thresholdDiskUsage), "%", int(diskUsed), "%", filesystem, "-------------------------------------"))
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

// Stop will receive stop channel
func (d *DiskStatus) Stop() chan bool {
	return d.stopChan
}
