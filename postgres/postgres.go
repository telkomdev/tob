package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

// Postgres service
type Postgres struct {
	url       string
	recovered bool
	verbose   bool
	logger    *log.Logger
	db        *sql.DB
}

// NewPostgres Postgres's constructor
func NewPostgres(verbose bool, logger *log.Logger) *Postgres {
	return &Postgres{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered: true,
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
