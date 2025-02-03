package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://api.shodan.io"

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// GetAPIInfo gets the info of the API user account
func (c *Client) GetAPIInfo() (*APIInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api-info?key=%s",
		c.apiKey, c.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result APIInfo
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
