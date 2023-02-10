package dummy

import (
	"log"
	"math/rand"
)

// Dummy service
type Dummy struct {
	url       string
	recovered bool
	verbose   bool
	logger    *log.Logger
}

// NewDummy Dummy's constructor
func NewDummy(verbose bool, logger *log.Logger) *Dummy {
	return &Dummy{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered: true,
	}
}

// Name the name of the service
func (d *Dummy) Name() string {
	return "dummy"
}

// Ping will try to ping the service
func (d *Dummy) Ping() []byte {
	n := rand.Intn(100)
	if n < 50 {
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Dummy) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *Dummy) Connect() error {
	if d.verbose {
		d.logger.Println("connect dummy")
	}

	return nil
}

// Close will close the service resources if needed
func (d *Dummy) Close() error {
	if d.verbose {
		d.logger.Println("close dummy")
	}

	return nil
}

// SetRecover will set recovered status
func (d *Dummy) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Dummy) IsRecover() bool {
	return d.recovered
}
