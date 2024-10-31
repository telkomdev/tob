package redisdb

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/util"
)

// Redis service
type Redis struct {
	url               string
	recovered         bool
	lastDownTime      string
	enabled           bool
	verbose           bool
	logger            *log.Logger
	client            *redis.Client
	checkInterval     int
	stopChan          chan bool
	message           string
	notificatorConfig config.Config
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
		d.SetMessage(reply.Err().Error())
		if d.verbose {
			d.logger.Println("Redis error")
			d.logger.Println(reply.Err())
		}
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

// LastDownTime will set last down time of service to current time
func (d *Redis) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *Redis) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
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

// SetMessage will set additional message
func (d *Redis) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *Redis) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *Redis) SetConfig(configs config.Config) {

}

// SetNotificatorConfig will set config
func (d *Redis) SetNotificatorConfig(configs config.Config) {
	d.notificatorConfig = configs
}

// GetNotificators will return notificators
func (d *Redis) GetNotificators() []tob.Notificator {
	return tob.InitNotificatorFactory(d.notificatorConfig, d.verbose)
}

// Stop will receive stop channel
func (d *Redis) Stop() chan bool {
	return d.stopChan
}
