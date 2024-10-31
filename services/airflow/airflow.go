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
	url                string
	recovered          bool
	lastDownTime       string
	schedulerStatus    string
	metadatabaseStatus string
	enabled            bool
	verbose            bool
	logger             *log.Logger
	checkInterval      int
	stopChan           chan bool
	message            string
	notificatorConfig  config.Config
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

	schedulerStatus, ok := data["scheduler"].(map[string]interface{})["status"].(string)
	if !ok {
		if a.verbose {
			a.logger.Printf("cannot read scheduler status: %v\n", err)
		}

		return err
	}
	a.schedulerStatus = schedulerStatus

	metadatabaseStatus, ok := data["metadatabase"].(map[string]interface{})["status"].(string)
	if !ok {
		if a.verbose {
			a.logger.Printf("cannot read metadatabase status: %v\n", err)
		}

		return err
	}
	a.metadatabaseStatus = metadatabaseStatus

	if a.schedulerStatus != "healthy" || a.metadatabaseStatus != "healthy" {
		errMsg := fmt.Sprintf("airflow is unhealthy: scheduler (%s), metadatabase (%s)", a.schedulerStatus, a.metadatabaseStatus)
		if a.verbose {
			a.logger.Println(errMsg)
		}

		return errors.New(errMsg)
	}

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
	notificators, err := tob.InitNotificatorFactory(a.notificatorConfig, a.verbose)
	if err != nil {
		a.logger.Printf("Warning: %s service does not activate Notifications, GetNotificators() will be nil\n", a.Name())
		return nil
	}
	return notificators
}

// Stop will receive stop channel
func (a *Airflow) Stop() chan bool {
	return a.stopChan
}
