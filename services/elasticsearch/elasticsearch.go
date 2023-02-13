package elasticsearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/telkomdev/tob/httpx"
	"github.com/telkomdev/tob/util"
)

// Elasticsearch service
type Elasticsearch struct {
	url           string
	recovered     bool
	serviceName   string
	lastDownTime  string
	enabled       bool
	verbose       bool
	logger        *log.Logger
	checkInterval int
	stopChan      chan bool
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if e.verbose {
			e.logger.Println(fmt.Sprintf("cannot read response body: %v", err))
		}

		return cStatus, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		if e.verbose {
			e.logger.Println(fmt.Sprintf("cannot read parse JSON body: %v", err))
		}

		return cStatus, err
	}

	cStatus, ok := data["status"].(string)
	if !ok {
		if e.verbose {
			e.logger.Println(fmt.Sprintf("cannot read cluster status: %v", err))
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
		return []byte("NOT_OK")
	}

	statusOK := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOK {
		if e.verbose {
			e.logger.Println(fmt.Sprintf("elasticsearch Ping status: %d", resp.StatusCode))
		}

		return []byte("NOT_OK")
	}

	cStatus, err := e.checkClusterStatus(resp)
	if err != nil {
		return []byte("NOT_OK")
	}

	if e.verbose {
		e.logger.Println(fmt.Sprintf("elasticsearch cluster status: %s", cStatus))
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

// Stop will receive stop channel
func (e *Elasticsearch) Stop() chan bool {
	return e.stopChan
}
