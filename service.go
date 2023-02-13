package tob

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

	// Redis service kind
	Redis ServiceKind = "redis"

	// Airflow service kind
	Airflow ServiceKind = "airflow"

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

	// Stop will receive stop channel
	Stop() chan bool
}
