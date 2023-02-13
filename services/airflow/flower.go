package airflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/telkomdev/tob/httpx"
	"github.com/telkomdev/tob/util"
)

// Airflow flower service
type AirflowFlower struct {
	url           string
	recovered     bool
	serviceName   string
	lastDownTime  string
	workers       []map[string]interface{}
	workerErr     bool
	enabled       bool
	verbose       bool
	logger        *log.Logger
	checkInterval int
	stopChan      chan bool
}

// Airflow flower's constructor
func NewAirflowFlower(verbose bool, logger *log.Logger) *AirflowFlower {
	stopChan := make(chan bool, 1)
	return &AirflowFlower{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (af *AirflowFlower) Name() string {
	return "airflowflower"
}

// checkWorkerStatus will check available worker status in Airflow cluster
func (af *AirflowFlower) checkWorkerStatus(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if af.verbose {
			af.logger.Println(fmt.Sprintf("cannot read response body: %v", err))
		}

		return err
	}
	defer func() { resp.Body.Close() }()

	var data map[string][]map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		if af.verbose {
			af.logger.Println(fmt.Sprintf("cannot read parse JSON body: %v", err))
		}

		return err
	}

	af.workers = af.workers[:0]
	af.workerErr = true
	for _, worker := range data["data"] {
		wStatus := worker["status"].(bool)
		wName := worker["hostname"]
		if !wStatus && af.verbose {
			af.logger.Println(fmt.Sprintf("airflow worker %s is offline", wName))
		} else {
			// if there is any worker alive, means that airflow worker is healthy
			// but need to check manually for offline worker
			af.workerErr = false
		}

		// also append to workers if needed later
		af.workers = append(af.workers, worker)
	}

	if af.workerErr {
		if af.verbose {
			af.logger.Println("no available online worker")
		}

		return errors.New("no available online worker")
	}

	return nil
}

// Ping will try to ping the service
func (af *AirflowFlower) Ping() []byte {
	resp, err := httpx.HTTPGet(af.url+"?json=1", nil, 5)
	if err != nil {
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if af.verbose {
			af.logger.Println(fmt.Sprintf("airflow-flower Ping status: %d", resp.StatusCode))
		}

		return []byte("NOT_OK")
	}

	if err := af.checkWorkerStatus(resp); err != nil {
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (af *AirflowFlower) SetURL(url string) {
	af.url = url
}

// Connect to service if needed
func (af *AirflowFlower) Connect() error {
	if af.verbose {
		af.logger.Println("connecting to Airflow-flower server")
	}

	return nil
}

// Close will close the service resources if needed
func (af *AirflowFlower) Close() error {
	if af.verbose {
		af.logger.Println("close Airflow-flower")
	}

	return nil
}

// SetRecover will set recovered status
func (af *AirflowFlower) SetRecover(recovered bool) {
	af.recovered = recovered
}

// IsRecover will return recovered status
func (af *AirflowFlower) IsRecover() bool {
	return af.recovered
}

// LastDownTime will set last down time of service to current time
func (af *AirflowFlower) SetLastDownTimeNow() {
	if af.recovered {
		af.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (af *AirflowFlower) GetDownTimeDiff() string {
	return util.TimeDifference(af.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (af *AirflowFlower) SetCheckInterval(interval int) {
	af.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (af *AirflowFlower) GetCheckInterval() int {
	return af.checkInterval
}

// Enable will set enabled status to service
func (af *AirflowFlower) Enable(enabled bool) {
	af.enabled = enabled
}

// IsEnabled will return enable status
func (af *AirflowFlower) IsEnabled() bool {
	return af.enabled
}

// Stop will receive stop channel
func (af *AirflowFlower) Stop() chan bool {
	return af.stopChan
}
