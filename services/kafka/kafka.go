package kafka

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	kf "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/util"
)

const (
	// ErrorClosedNetwork is an error indicating the connection is closed
	ErrorClosedNetwork = "use of closed network connection"
)

// Kafka service
type Kafka struct {
	url           string
	brokerSize    int
	recovered     bool
	lastDownTime  string
	enabled       bool
	verbose       bool
	logger        *log.Logger
	client        *kf.Conn
	checkInterval int
	stopChan      chan bool
}

// NewKafka Kafka's constructor
func NewKafka(verbose bool, logger *log.Logger) *Kafka {
	stopChan := make(chan bool, 1)
	return &Kafka{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *Kafka) Name() string {
	return "kafka"
}

// Ping will try to ping the service
func (d *Kafka) Ping() []byte {
	if d.client == nil {
		return []byte("NOT_OK")
	}

	reply, err := d.client.Brokers()
	if err != nil {
		if d.verbose {
			d.logger.Println("Kafka error read available brokers")
			d.logger.Println(err)

			// re dial
			if strings.Contains(err.Error(), ErrorClosedNetwork) {
				d.logger.Println(fmt.Sprintf("Kafka: %s | do re dial\n", err.Error()))
				// re dial ignore error
				err = d.dial()
				if err != nil {
					d.logger.Println(fmt.Sprintf("Kafka: %s | do re dial\n", err.Error()))
				}
			}
		}
		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Print(fmt.Sprintf("Expected Kafka Brokers size: %d", d.brokerSize))
		d.logger.Print("Kafka Brokers: ")
		for _, broker := range reply {
			d.logger.Println(broker.Host)
		}
	}

	// if the reply length is less than brokerSize,
	// then there is an indication that the Kafka Cluster is experiencing problems
	if len(reply) < d.brokerSize {
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Kafka) SetURL(url string) {
	d.url = url
}

func (d *Kafka) dial() error {
	if d.verbose {
		d.logger.Println("connecting to Kafka server")
	}

	u, err := url.Parse(d.url)
	if err != nil {
		return err
	}

	hosts := u.Host
	password := ""
	username := ""
	if u.User != nil {
		if u.User.Username() != "" {
			username = u.User.Username()
		}

		p, ok := u.User.Password()
		if ok {
			password = p
		}
	}

	dialer := &kf.Dialer{
		ClientID:  "tob",
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if username != "" && password != "" {
		dialer.SASLMechanism = plain.Mechanism{
			Username: username,
			Password: password,
		}
	}

	hostSplitted := strings.Split(hosts, ",")
	host := ""
	if len(hostSplitted) > 1 {
		host = hostSplitted[0]
	} else {
		host = hosts
	}

	// set brokerSize with total length of Kafka host broker from config
	d.brokerSize = len(hostSplitted)

	conn, err := dialer.Dial("tcp", host)
	if err != nil {
		return err
	}

	// set connected conn
	d.client = conn

	if d.verbose {
		d.logger.Println("connecting to Kafka server succeed")
	}

	return nil
}

// Connect to service if needed
func (d *Kafka) Connect() error {
	return d.dial()
}

// Close will close the service resources if needed
func (d *Kafka) Close() error {
	if d.verbose {
		d.logger.Println("closing Kafka connection")
	}

	if d.client != nil {
		err := d.client.Close()
		if err != nil {
			return err
		}

		if d.verbose {
			d.logger.Println("closing Kafka connection succeed")
		}
	}

	return nil
}

// SetRecover will set recovered status
func (d *Kafka) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Kafka) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *Kafka) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *Kafka) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *Kafka) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *Kafka) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *Kafka) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *Kafka) IsEnabled() bool {
	return d.enabled
}

// SetMessage will set additional message
func (d *Kafka) SetMessage(message string) {

}

// GetMessage will return additional message
func (d *Kafka) GetMessage() string {
	return ""
}

// SetConfig will set config
func (d *Kafka) SetConfig(configs config.Config) {

}

// Stop will receive stop channel
func (d *Kafka) Stop() chan bool {
	return d.stopChan
}
