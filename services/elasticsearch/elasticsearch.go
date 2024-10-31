package elasticsearch

import (
	"encoding/json"
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

// Elasticsearch service
type Elasticsearch struct {
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

// Elasticsearch's constructor
func NewElasticsearch(verbose bool, logger *log.Logger) *Elasticsearch {
	stopChan := make(chan bool, 1)
	return &Elasticsearch{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (e *Elasticsearch) Name() string {
	return "elasticsearch"
}

// checkClusterStatus checks the status of an Elasticsearch cluster by checking cluster index status.
func (e *Elasticsearch) checkClusterStatus(resp *http.Response) (string, error) {
	cStatus := ""
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		if e.verbose {
			e.logger.Printf("cannot read response body: %v", err)
		}

		return cStatus, err
	}
	defer func() { resp.Body.Close() }()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		if e.verbose {
			e.logger.Printf("cannot parse JSON body: %v", err)
		}

		return cStatus, err
	}

	cStatus, ok := data["status"].(string)
	if !ok {
		if e.verbose {
			e.logger.Printf("cannot read cluster status: %v", err)
		}

		return cStatus, err
	}

	if cStatus != "green" && cStatus != "yellow" {
		if e.verbose {
			e.logger.Println("elasticsearch cluster is unhealthy: ", cStatus)
		}

		return cStatus, err
	}

	return cStatus, nil
}

// Ping will try to ping the service
func (e *Elasticsearch) Ping() []byte {
	resp, err := httpx.HTTPGet(e.url, nil, 5)
	if err != nil {
		e.SetMessage(err.Error())
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		e.SetMessage(fmt.Sprintf("error: elasticsearch Ping status: %d", resp.StatusCode))
		if e.verbose {
			e.logger.Printf("elasticsearch Ping status: %d", resp.StatusCode)
		}

		return []byte("NOT_OK")
	}

	cStatus, err := e.checkClusterStatus(resp)
	if err != nil {
		e.SetMessage(err.Error())
		return []byte("NOT_OK")
	}

	if e.verbose {
		e.logger.Printf("elasticsearch cluster status: %s", cStatus)
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (e *Elasticsearch) SetURL(url string) {
	e.url = url
}

// Connect to service if needed
func (e *Elasticsearch) Connect() error {
	if e.verbose {
		e.logger.Println("connecting to Elasticsearch server")
	}

	return nil
}

// Close will close the service resources if needed
func (e *Elasticsearch) Close() error {
	if e.verbose {
		e.logger.Println("close Elasticsearch")
	}

	return nil
}

// SetRecover will set recovered status
func (e *Elasticsearch) SetRecover(recovered bool) {
	e.recovered = recovered
}

// IsRecover will return recovered status
func (e *Elasticsearch) IsRecover() bool {
	return e.recovered
}

// LastDownTime will set last down time of service to current time
func (e *Elasticsearch) SetLastDownTimeNow() {
	if e.recovered {
		e.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (e *Elasticsearch) GetDownTimeDiff() string {
	return util.TimeDifference(e.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (e *Elasticsearch) SetCheckInterval(interval int) {
	e.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (e *Elasticsearch) GetCheckInterval() int {
	return e.checkInterval
}

// Enable will set enabled status to service
func (e *Elasticsearch) Enable(enabled bool) {
	e.enabled = enabled
}

// IsEnabled will return enable status
func (e *Elasticsearch) IsEnabled() bool {
	return e.enabled
}

// SetMessage will set additional message
func (e *Elasticsearch) SetMessage(message string) {
	e.message = message
}

// GetMessage will return additional message
func (e *Elasticsearch) GetMessage() string {
	return e.message
}

// SetConfig will set config
func (e *Elasticsearch) SetConfig(configs config.Config) {

}

// SetNotificatorConfig will set config
func (e *Elasticsearch) SetNotificatorConfig(configs config.Config) {
	e.notificatorConfig = configs
}

// GetNotificators will return notificators
func (e *Elasticsearch) GetNotificators() []tob.Notificator {
	notificators, err := tob.InitNotificatorFactory(e.notificatorConfig, e.verbose)
	if err != nil {
		e.logger.Printf("Warning: %s service does not activate Notifications, GetNotificators() will be nil\n", e.Name())
		return nil
	}
	return notificators
}

// Stop will receive stop channel
func (e *Elasticsearch) Stop() chan bool {
	return e.stopChan
}
