// Package passport contains client logic for interacting with the user service.
package passport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/wanderfusion/nomadcore/pkg/cache"
	"github.com/wanderfusion/nomadcore/pkg/commons/httpclient"
)

// Client struct encapsulates HTTP client and host information for Passport service.
type Client struct {
	httpClient          *http.Client
	host                string
	getUsersFromIDCache *cache.Cache[GetUsersFromIDsResponse] // Cache for GetUser calls
}

// NewClient initializes a new Passport client.
func NewClient(host string) *Client {
	httpClient := httpclient.NewPooledHttpClient(100, 100, 100, 1000)
	getUsersFromIDCache := cache.New[GetUsersFromIDsResponse]()

	return &Client{
		httpClient:          httpClient,
		host:                host,
		getUsersFromIDCache: getUsersFromIDCache,
	}
}

// GetUsersFromUsernames fetches user information based on usernames.
func (c *Client) GetUsersFromUsernames(usernames []string) (GetUsersFromIDsResponse, error) {
	joinedIDs := strings.Join(usernames, ",")
	url := fmt.Sprintf("%s/users/username/%s", c.host, joinedIDs)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return GetUsersFromIDsResponse{}, err
	}
	defer resp.Body.Close()

	var result GetUsersFromIDsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return GetUsersFromIDsResponse{}, err
	}

	return result, nil
}

// CachedGetUsersFromUsernames fetches user information from cache or falls back to API call.
func (c *Client) CachedGetUsersFromUsernames(usernames []string, ttl time.Duration) (GetUsersFromIDsResponse, error) {
	joinedUsernames := strings.Join(usernames, ",")
	cachedValue := c.getUsersFromIDCache.GetValue(joinedUsernames)

	// Check if value exists in cache
	if !cachedValue.None {
		return cachedValue.Some, nil
	}

	// Fetch from API and update cache
	result, err := c.GetUsersFromUsernames(usernames)
	if err != nil {
		return GetUsersFromIDsResponse{}, err
	}
	c.getUsersFromIDCache.SetValue(joinedUsernames, result, ttl)

	return result, nil
}
