// Package main serves as the entry point of the application.
package main

import (
	"fmt"

	"github.com/akxcix/nomadcore/pkg/app"
)

// main is the entry function for the application.
func main() {
	// Print a message indicating the application is starting.
	fmt.Println("Starting nomad-core...")

	// Run the application logic from the app package.
	app.Run()
}
