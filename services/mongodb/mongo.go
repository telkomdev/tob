package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/telkomdev/tob"
	"github.com/telkomdev/tob/config"
	"github.com/telkomdev/tob/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo service
type Mongo struct {
	url               string
	recovered         bool
	lastDownTime      string
	enabled           bool
	verbose           bool
	logger            *log.Logger
	client            *mongo.Client
	checkInterval     int
	stopChan          chan bool
	message           string
	notificatorConfig config.Config
}

// NewMongo Mongo's constructor
func NewMongo(verbose bool, logger *log.Logger) *Mongo {
	stopChan := make(chan bool, 1)
	return &Mongo{
		logger:  logger,
		verbose: verbose,

		// by default service is recovered
		recovered:     true,
		checkInterval: 0,
		stopChan:      stopChan,
	}
}

// Name the name of the service
func (d *Mongo) Name() string {
	return "mongodb"
}

// Ping will try to ping the service
func (d *Mongo) Ping() []byte {
	if d.client == nil {
		return []byte("NOT_OK")
	}

	if err := d.client.Ping(context.Background(), nil); err != nil {
		d.SetMessage(err.Error())
		if d.verbose {
			d.logger.Println("MongoDB error")
			d.logger.Println(err)
		}
		return []byte("NOT_OK")
	}

	return []byte("OK")
}

// SetURL will set the service URL
func (d *Mongo) SetURL(url string) {
	d.url = url
}

// Connect to service if needed
func (d *Mongo) Connect() error {
	if d.verbose {
		d.logger.Println("connecting to MongoDB server")
	}

	client, err := mongo.NewClient(
		options.Client().ApplyURI(d.url),
		options.Client().SetConnectTimeout(time.Second*4),
		options.Client().SetServerSelectionTimeout(time.Second*4),
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer func() { cancel() }()

	if err := client.Connect(ctx); err != nil {
		return err
	}

	if d.verbose {
		d.logger.Println("connecting to MongoDB server succeed")
	}

	// set connected client
	d.client = client

	return nil
}

// Close will close the service resources if needed
func (d *Mongo) Close() error {
	if d.verbose {
		d.logger.Println("closing MongoDB connection")
	}

	if d.client != nil {
		err := d.client.Disconnect(context.Background())
		if err != nil {
			return err
		}

		if d.verbose {
			d.logger.Println("closing MongoDB connection succeed")
		}
	}

	return nil
}

// SetRecover will set recovered status
func (d *Mongo) SetRecover(recovered bool) {
	d.recovered = recovered
}

// IsRecover will return recovered status
func (d *Mongo) IsRecover() bool {
	return d.recovered
}

// LastDownTime will set last down time of service to current time
func (d *Mongo) SetLastDownTimeNow() {
	if d.recovered {
		d.lastDownTime = time.Now().Format(util.YYMMDD)
	}
}

// GetDownTimeDiff will return down time service difference in minutes
func (d *Mongo) GetDownTimeDiff() string {
	return util.TimeDifference(d.lastDownTime, time.Now().Format(util.YYMMDD))
}

// SetCheckInterval will set check interval to service
func (d *Mongo) SetCheckInterval(interval int) {
	d.checkInterval = interval
}

// GetCheckInterval will return check interval to service
func (d *Mongo) GetCheckInterval() int {
	return d.checkInterval
}

// Enable will set enabled status to service
func (d *Mongo) Enable(enabled bool) {
	d.enabled = enabled
}

// IsEnabled will return enable status
func (d *Mongo) IsEnabled() bool {
	return d.enabled
}

// SetMessage will set additional message
func (d *Mongo) SetMessage(message string) {
	d.message = message
}

// GetMessage will return additional message
func (d *Mongo) GetMessage() string {
	return d.message
}

// SetConfig will set config
func (d *Mongo) SetConfig(configs config.Config) {

}

// SetNotificatorConfig will set config
func (d *Mongo) SetNotificatorConfig(configs config.Config) {
	d.notificatorConfig = configs
}

// GetNotificators will return notificators
func (d *Mongo) GetNotificators() []tob.Notificator {
	notificators, err := tob.InitNotificatorFactory(d.notificatorConfig, d.verbose)
	if err != nil {
		d.logger.Printf("Warning: %s service does not activate Notifications, GetNotificators() will be nil\n", d.Name())
		return nil
	}
	return notificators
}

// Stop will receive stop channel
func (d *Mongo) Stop() chan bool {
	return d.stopChan
}
