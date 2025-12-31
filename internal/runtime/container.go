// Package runtime provides application runtime dependency injection container
// This package uses uber's dig library to manage dependency injection
package runtime

import (
	"go.uber.org/dig"
)

// container is the application's global dependency injection container
// All services and components are registered and resolved through it
var container *dig.Container

// init initializes the dependency injection container
// Automatically called when the program starts
func init() {
	container = dig.New()
}

// GetContainer returns a reference to the global dependency injection container
// For use by other packages to register or get services
func GetContainer() *dig.Container {
	return container
}
