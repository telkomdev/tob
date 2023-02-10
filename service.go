package tob

// Service represent available services
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
}
