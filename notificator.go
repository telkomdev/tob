package tob

// Notificator type
type Notificator interface {
	// Provider will return Notificator provider
	Provider() string

	// Send will send message to Notificator
	Send(msg string) error

	// IsEnabled will return enable status
	IsEnabled() bool
}
