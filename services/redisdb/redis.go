package redisdb

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"net/url"
)

// Redis service
type Redis struct {
	url           string
	recovered     bool
	enabled       bool
	verbose       bool
	logger        *log.Logger
	client        *redis.Client
	checkInterval int
	stopChan      chan bool
}

// NewRedis Redis's constructor
func NewRedis(verbose bool, logger *log.Logger) *Redis {
	stopChan := make(chan bool, 1)
	return &Redis{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *Redis) Name() string {
	return "redis"
}

// Ping will try to ping the service
func (d *Redis) Ping() []byte {
	if d.client == nil {
		return []byte("NOT_OK")
	}

	reply := d.client.Ping(context.Background())
	if reply.Err() != nil {
		return []byte("NOT_OK")
	}

	if d.verbose {
		d.logger.Print("redis reply: ")
		d.logger.Println(reply.String())
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Redis) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *Redis) Connect() error {
	if d.verbose {
		d.logger.Println("connecting to Redis server")
	}

	u, err := url.Parse(d.url)
	if err != nil {
		return err
	}

	host := u.Host
	password := ""
	if u.User != nil {
		p, ok := u.User.Password()
		if ok {
			password = p
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       0, // use default DB
	})

	if d.verbose {
		d.logger.Println("connecting to Redis server succeed")
	}

	// set connected conn
	d.client = client

	return nil
}

// Close will close the service resources if needed
func (d *Redis) Close() error {
	if d.verbose {
		d.logger.Println("closing Redis connection")
	}

	if d.client != nil {
		err := d.client.Close()
		if err != nil {
			return err
		}

		if d.verbose {
			d.logger.Println("closing Redis connection succeed")
		}
	}

	return nil
}

// SetRecover will set recovered status
func (d *Redis) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Redis) IsRecover() bool {
	return d.recovered
}

// SetCheckInterval will set check interval to service
func (d *Redis) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *Redis) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *Redis) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *Redis) IsEnabled() bool {
	return d.enabled
}

// Stop will receive stop channel
func (d *Redis) Stop() chan bool {
	return d.stopChan
}
