package frontend

// Frontend defines the interface which must be used by any analyzer.
type Frontend interface {
	// TODO: Pass config

	// Listener should launch a service and stay listening for connections.
	// Should only return (nil or error) when something happened.
	Listener() error
}
