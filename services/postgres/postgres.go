package postgres

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/telkomdev/tob/util"
)

// Postgres service
type Postgres struct {
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

// NewPostgres Postgres's constructor
func NewPostgres(verbose bool, logger *log.Logger) *Postgres {
	stopChan := make(chan bool, 1)
	return &Postgres{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *Postgres) Name() string {
	return "postgresql"
}

// Ping will try to ping the service
func (d *Postgres) Ping() []byte {
	if d.db == nil {
		return []byte("NOT_OK")
	}

	if err := d.db.Ping(); err != nil {
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Postgres) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *Postgres) Connect() error {
	if d.verbose {
		d.logger.Println("connecting to Postgres server")
	}

	db, err := sql.Open("postgres", d.url)
	if err != nil {
		return err
	}

	if d.verbose {
		d.logger.Println("connecting to Postgres server succeed")
	}

	// set connected db
	d.db = db

	return nil
}

// Close will close the service resources if needed
func (d *Postgres) Close() error {
	if d.verbose {
		d.logger.Println("closing Postgresql connection")
	}

	if d.db != nil {
		err := d.db.Close()
		if err != nil {
			return err
		}

		if d.verbose {
			d.logger.Println("closing Postgresql connection succeed")
		}
	}

	return nil
}

// SetRecover will set recovered status
func (d *Postgres) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Postgres) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *Postgres) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *Postgres) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *Postgres) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *Postgres) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *Postgres) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *Postgres) IsEnabled() bool {
	return d.enabled
}

// Stop will receive stop channel
func (d *Postgres) Stop() chan bool {
	return d.stopChan
}
