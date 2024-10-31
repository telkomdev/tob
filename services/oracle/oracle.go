package oracle

import (
	"context"
	"log"
	"net/url"
	"time"

	ora "github.com/sijms/go-ora/v2"
	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/util"
)

// Oracle service
type Oracle struct {
	url               string
	recovered         bool
	lastDownTime      string
	enabled           bool
	verbose           bool
	logger            *log.Logger
	db                *ora.Connection
	checkInterval     int
	stopChan          chan bool
	message           string
	notificatorConfig config.Config
}

// NewOracle Oracle's constructor
func NewOracle(verbose bool, logger *log.Logger) *Oracle {
	stopChan := make(chan bool, 1)
	return &Oracle{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *Oracle) Name() string {
	return "oracle"
}

// Ping will try to ping the service
func (d *Oracle) Ping() []byte {
	if d.db == nil {
		return []byte("NOT_OK")
	}

	if err := d.db.Ping(context.Background()); err != nil {
		d.SetMessage(err.Error())
		if d.verbose {
			d.logger.Println("Oracle ping error")
			d.logger.Println(err)
		}
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Oracle) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *Oracle) Connect() error {
	if d.verbose {
		d.logger.Println("connecting to Oracle server")
	}

	// url = "oracle://username:pass%23123@127.0.0.1:1521/servicename"
	parsedURL, err := url.Parse(d.url)
	if err != nil {
		if d.verbose {
			d.logger.Printf("error parsing oracle url: %s\n", err.Error())
		}
		return err
	}

	parsedPassword := ""
	parsedUsername := ""
	if parsedURL.User != nil {
		if parsedURL.User.Username() != "" {
			parsedUsername = parsedURL.User.Username()
		}

		p, ok := parsedURL.User.Password()
		if ok {
			parsedPassword = p
		}
	}
	u := &url.URL{
		Scheme: parsedURL.Scheme,
		Host:   parsedURL.Host,
		User:   url.UserPassword(parsedUsername, parsedPassword),
		Path:   parsedURL.Path,
	}

	// Convert the URL object to a string
	connStr := u.String()

	conn, err := ora.NewConnection(connStr, nil)
	if err != nil {
		return err
	}

	err = conn.Open()
	if err != nil {
		return err
	}

	if d.verbose {
		d.logger.Println("connecting to Oracle server succeed")
	}

	// set connected db
	d.db = conn

	return nil
}

// Close will close the service resources if needed
func (d *Oracle) Close() error {
	if d.verbose {
		d.logger.Println("closing Oracle connection")
	}

	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return err
		}

		if d.verbose {
			d.logger.Println("closing Oracle connection succeed")
		}
	}

	return nil
}

// SetRecover will set recovered status
func (d *Oracle) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Oracle) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *Oracle) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *Oracle) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *Oracle) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *Oracle) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *Oracle) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *Oracle) IsEnabled() bool {
	return d.enabled
}

// SetMessage will set additional message
func (d *Oracle) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *Oracle) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *Oracle) SetConfig(configs config.Config) {

}

// SetNotificatorConfig will set config
func (d *Oracle) SetNotificatorConfig(configs config.Config) {
	d.notificatorConfig = configs
}

// GetNotificators will return notificators
func (d *Oracle) GetNotificators() []tob.Notificator {
	return tob.InitNotificatorFactory(d.notificatorConfig, d.verbose)
}

// Stop will receive stop channel
func (d *Oracle) Stop() chan bool {
	return d.stopChan
}
