// Package httpclient provides utilities for creating and managing HTTP clients.
package httpclient

import (
	"net/http"
	"time"
)

// NewPooledHttpClient initializes a new http.Client with specified connection pooling settings.
// Parameters:
// - maxIdleConns: Maximum number of idle connections.
// - maxConnsPerHost: Maximum number of connections per host.
// - maxIdleConnsPerHost: Maximum number of idle connections per host.
// - timeoutMillis: Client timeout in milliseconds.
func NewPooledHttpClient(maxIdleConns, maxConnsPerHost, maxIdleConnsPerHost, timeoutMillis int) *http.Client {
	// Clone the default transport settings
	t := http.DefaultTransport.(*http.Transport).Clone()

	// Customize the transport settings
	t.MaxIdleConns = maxIdleConns
	t.MaxConnsPerHost = maxConnsPerHost
	t.MaxIdleConnsPerHost = maxIdleConnsPerHost

	// Initialize and return the custom HTTP client
	httpClient := &http.Client{
		Timeout:   time.Duration(timeoutMillis) * time.Millisecond,
		Transport: t,
	}

	return httpClient
}
