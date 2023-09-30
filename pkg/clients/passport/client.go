package passport

import (
	"encoding/json"
	"fmt"
	"github.com/akxcix/nomadcore/pkg/cache"
	"net/http"
	"strings"
	"time"

	"github.com/akxcix/nomadcore/pkg/commons/httpclient"
)

type Client struct {
	httpClient          *http.Client
	host                string
	getUsersFromIDCache *cache.Cache[GetUsersFromIDsResponse]
}

func NewClient(host string) *Client {
	httpClient := httpclient.NewPooledHttpClient(100, 100, 100, 1000)
	getUsersFromIDCache := cache.New[GetUsersFromIDsResponse]()

	return &Client{
		httpClient:          httpClient,
		host:                host,
		getUsersFromIDCache: getUsersFromIDCache,
	}
}

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

func (c *Client) CachedGetUsersFromUsernames(usernames []string, ttl time.Duration) (GetUsersFromIDsResponse, error) {
	joinedUsernames := strings.Join(usernames, ",")
	cachedValue := c.getUsersFromIDCache.GetValue(joinedUsernames)

	if !cachedValue.None {
		return cachedValue.Some, nil
	}

	result, err := c.GetUsersFromUsernames(usernames)
	if err != nil {
		return GetUsersFromIDsResponse{}, err
	}

	c.getUsersFromIDCache.SetValue(joinedUsernames, result, ttl)
	return result, nil
}
