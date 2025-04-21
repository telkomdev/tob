package tob

import (
	"github.com/telkomdev/tob/config"
)

// ServiceKind represent a type/kind of service
type ServiceKind string

var (
	// Postgresql service kind
	Postgresql ServiceKind = "postgresql"

	// MySQL service kind
	MySQL ServiceKind = "mysql"

	// Web service kind
	Web ServiceKind = "web"

	// MongoDB service kind
	MongoDB ServiceKind = "mongodb"

	// Oracle service kind
	Oracle ServiceKind = "oracle"

	// Redis service kind
	Redis ServiceKind = "redis"

	// Elasticsearch service kind
	Elasticsearch ServiceKind = "elasticsearch"

	// Airflow service kind
	Airflow ServiceKind = "airflow"

	// AirflowFlower service kind
	AirflowFlower ServiceKind = "airflowflower"

	// DiskStatus service kind
	DiskStatus ServiceKind = "diskstatus"

	// Kafka servie kind
	Kafka ServiceKind = "kafka"

	// Plugin service kind
	Plugin ServiceKind = "plugin"

	// DiskStatus service kind
	SSLStatus ServiceKind = "sslstatus"

	// Dummy service kind
	Dummy ServiceKind = "dummy"
)

// Service represent base of all available services
type Service interface {

	// Name the name of the service
	Name() string

	// Ping will try to ping the service
	Ping() []byte

	// SetURL will set the service URL
	SetURL(url string)

	// Connect to service if needed
	Connect() error

	// Close will close the service resources if needed
	Close() error

	// SetRecover will set recovered status
	SetRecover(recovered bool)

	// IsRecover will return recovered status
	IsRecover() bool

	// LastDownTime will set last down time of service to current time
	SetLastDownTimeNow()

	// GetDownTimeDiff will return down time service difference in minutes
	GetDownTimeDiff() string

	// SetCheckInterval will set check interval to service
	SetCheckInterval(interval int)

	// GetCheckInterval will return check interval to service
	GetCheckInterval() int

	// Enable will set enabled status to service
	Enable(enabled bool)

	// IsEnabled will return enable status
	IsEnabled() bool

	// SetMessage will set additional message
	SetMessage(message string)

	// GetMessage will return additional message
	GetMessage() string

	// SetConfig will set config
	SetConfig(configs config.Config)

	// SetNotificatorConfig will set config
	SetNotificatorConfig(configs config.Config)

	// GetNotificators will return notificators
	GetNotificators() []Notificator

	// Stop will receive stop channel
	Stop() chan bool
}
