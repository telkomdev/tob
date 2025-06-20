package airflow

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/httpx"
	"github.com/telkomdev/tob/util"
)

// Airflow service
type Airflow struct {
	url                      string
	recovered                bool
	lastDownTime             string
	schedulerStatus          string
	latestSchedulerHeartbeat string
	metadatabaseStatus       string
	enabled                  bool
	verbose                  bool
	logger                   *log.Logger
	checkInterval            int
	stopChan                 chan bool
	message                  string
	notificatorConfig        config.Config
}

// Airflow's constructor
func NewAirflow(verbose bool, logger *log.Logger) *Airflow {
	stopChan := make(chan bool, 1)
	return &Airflow{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (a *Airflow) Name() string {
	return "airflow"
}

// checkClusterStatus will check status of Airflow's metadatabase & scheduler
func (a *Airflow) checkClusterStatus(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if a.verbose {
			a.logger.Printf("cannot read response body: %v\n", err)
		}

		return err
	}
	defer func() { resp.Body.Close() }()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		if a.verbose {
			a.logger.Printf("cannot read parse JSON body: %v\n", err)
		}

		return err
	}

	schedulerRaw, ok := data["scheduler"].(map[string]interface{})
	if !ok {
		if a.verbose {
			a.logger.Println("cannot read scheduler block (not a map)")
		}
		return errors.New("invalid scheduler block")
	}

	schedulerStatus, ok := schedulerRaw["status"].(string)
	if !ok {
		if a.verbose {
			a.logger.Println("cannot read scheduler status (not a string)")
		}

		return errors.New("invalid scheduler status")
	}

	latestSchedulerHeartbeat, ok := schedulerRaw["latest_scheduler_heartbeat"].(string)
	if !ok {
		if a.verbose {
			a.logger.Println("cannot read scheduler latest_scheduler_heartbeat (not a string)")
		}

		return errors.New("invalid latest_scheduler_heartbeat")
	}

	utcTime, err := time.Parse(time.RFC3339Nano, latestSchedulerHeartbeat)
	if err != nil {
		if a.verbose {
			a.logger.Printf("failed to parse time: %v\n", err)
		}

		return errors.New("invalid latest_scheduler_heartbeat")
	}

	timezoneJakarta, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		if a.verbose {
			a.logger.Printf("failed to load WIB location: %v\n", err)
		}

		return errors.New("invalid timezone")
	}

	wibTime := utcTime.In(timezoneJakarta)
	a.latestSchedulerHeartbeat = wibTime.String()

	a.schedulerStatus = schedulerStatus

	metadatabaseRaw, ok := data["metadatabase"].(map[string]interface{})
	if !ok {
		if a.verbose {
			a.logger.Println("cannot read metadatabase block (not a map)")
		}
		return errors.New("invalid metadatabase block")
	}

	metadatabaseStatus, ok := metadatabaseRaw["status"].(string)
	if !ok {
		if a.verbose {
			a.logger.Println("cannot read metadatabase status")
		}

		return errors.New("invalid metadatabase status")
	}

	a.metadatabaseStatus = metadatabaseStatus

	message := fmt.Sprintf("Airflow Scheduler Status: %s\nLatest Scheduler Heartbeat: %s\n\n\n Airflow Metadatabase: %s",
		a.schedulerStatus, a.latestSchedulerHeartbeat, a.metadatabaseStatus)
	if a.schedulerStatus != "healthy" || a.metadatabaseStatus != "healthy" {
		if a.verbose {
			a.logger.Println(message)
		}

		return errors.New(message)
	}

	a.SetMessage(message)

	return nil
}

// Ping will try to ping the service
func (a *Airflow) Ping() []byte {
	resp, err := httpx.HTTPGet(a.url, nil, 5)
	if err != nil {
		a.SetMessage(err.Error())
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		a.SetMessage(fmt.Sprintf("airflow Ping status: %d", resp.StatusCode))
		if a.verbose {
			a.logger.Printf("airflow Ping status: %d\n", resp.StatusCode)
		}

		return []byte("NOT_OK")
	}

	if err := a.checkClusterStatus(resp); err != nil {
		a.SetMessage(err.Error())
		return []byte("NOT_OK")
	}

	if a.verbose {
		a.logger.Printf("airflow: scheduler (%s), metadatabase (%s)\n", a.schedulerStatus, a.metadatabaseStatus)
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (a *Airflow) SetURL(url string) {
	a.url = url
}

// Connect to service if needed
func (a *Airflow) Connect() error {
	if a.verbose {
		a.logger.Println("connecting to Airflow server")
	}

	return nil
}

// Close will close the service resources if needed
func (a *Airflow) Close() error {
	if a.verbose {
		a.logger.Println("close Airflow")
	}

	return nil
}

// SetRecover will set recovered status
func (a *Airflow) SetRecover(recovered bool) {
	a.recovered = recovered
}

// IsRecover will return recovered status
func (a *Airflow) IsRecover() bool {
	return a.recovered
}

// LastDownTime will set last down time of service to current time
func (a *Airflow) SetLastDownTimeNow() {
	if a.recovered {
		a.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (a *Airflow) GetDownTimeDiff() string {
	return util.TimeDifference(a.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (a *Airflow) SetCheckInterval(interval int) {
	a.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (a *Airflow) GetCheckInterval() int {
	return a.checkInterval
}

// Enable will set enabled status to service
func (a *Airflow) Enable(enabled bool) {
	a.enabled = enabled
}

// IsEnabled will return enable status
func (a *Airflow) IsEnabled() bool {
	return a.enabled
}

// SetMessage will set additional message
func (a *Airflow) SetMessage(message string) {
	a.message = message
}

// GetMessage will return additional message
func (a *Airflow) GetMessage() string {
	return a.message
}

// SetConfig will set config
func (a *Airflow) SetConfig(configs config.Config) {

}

// SetNotificatorConfig will set config
func (a *Airflow) SetNotificatorConfig(configs config.Config) {
	a.notificatorConfig = configs
}

// GetNotificators will return notificators
func (a *Airflow) GetNotificators() []tob.Notificator {
	return tob.InitNotificatorFactory(a.notificatorConfig, a.verbose)
}

// Stop will receive stop channel
func (a *Airflow) Stop() chan bool {
	return a.stopChan
}
