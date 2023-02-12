package mysqldb

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/telkomdev/tob/util"
)

// MySQL service
type MySQL struct {
	url           string
	recovered     bool
	lastDownTime  string
	enabled       bool
	verbose       bool
	logger        *log.Logger
	db            *sql.DB
	checkInterval int
	stopChan      chan bool
}

// NewMySQL MySQL's constructor
func NewMySQL(verbose bool, logger *log.Logger) *MySQL {
	stopChan := make(chan bool, 1)
	return &MySQL{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *MySQL) Name() string {
	return "mysql"
}

// Ping will try to ping the service
func (d *MySQL) Ping() []byte {
	if d.db == nil {
		return []byte("NOT_OK")
	}

	if err := d.db.Ping(); err != nil {
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *MySQL) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *MySQL) Connect() error {
	if d.verbose {
		d.logger.Println("connecting to MySQL server")
	}

	db, err := sql.Open("mysql", d.url)
	if err != nil {
		return err
	}

	if d.verbose {
		d.logger.Println("connecting to MySQL server succeed")
	}

	// set connected db
	d.db = db

	return nil
}

// Close will close the service resources if needed
func (d *MySQL) Close() error {
	if d.verbose {
		d.logger.Println("closing MySQL connection")
	}

	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return err
		}

		if d.verbose {
			d.logger.Println("closing MySQL connection succeed")
		}
	}

	return nil
}

// SetRecover will set recovered status
func (d *MySQL) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *MySQL) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *MySQL) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *MySQL) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *MySQL) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *MySQL) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *MySQL) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *MySQL) IsEnabled() bool {
	return d.enabled
}

// Stop will receive stop channel
func (d *MySQL) Stop() chan bool {
	return d.stopChan
}
