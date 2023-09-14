package passport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/akxcix/nomadcore/pkg/commons/httpclient"
)

type Client struct {
	httpClient *http.Client
	host       string
}

func NewClient(host string) *Client {
	httpClient := httpclient.NewPooledHttpClient(100, 100, 100, 1000)

	return &Client{
		httpClient: httpClient,
		host:       host,
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
