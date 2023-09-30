// Package httpclient provides utilities for making HTTP requests.
package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// MakePOSTRequest performs a POST HTTP request and returns the response as a byte array.
// Parameters:
// - client: The http.Client to use for the request.
// - baseURL: The base URL for the request.
// - endpoint: The API endpoint for the request.
// - payload: The request payload.
// - customHeaders: Additional HTTP headers.
// - customQueryParams: Additional query parameters.
func MakePOSTRequest[T any](client *http.Client, baseURL string, endpoint string, payload T, customHeaders map[string]string, customQueryParams map[string]string) ([]byte, error) {
	// Prepare request body
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	bodyBuffer := bytes.NewBuffer(jsonPayload)

	// Prepare request
	fullURL := baseURL + endpoint
	req, err := http.NewRequest("POST", fullURL, bodyBuffer)
	if err != nil {
		return nil, err
	}

	// Set default and custom headers
	setHeaders(req, customHeaders)

	// Add query parameters
	addQueryParameters(req, customQueryParams)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read and return response
	return io.ReadAll(resp.Body)
}

// setHeaders sets default and custom headers for an HTTP request.
func setHeaders(req *http.Request, customHeaders map[string]string) {
	req.Header.Set("Accept", "application/json")
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}
}

// addQueryParameters adds custom query parameters to an HTTP request.
func addQueryParameters(req *http.Request, customQueryParams map[string]string) {
	q := req.URL.Query()
	for key, value := range customQueryParams {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
}
