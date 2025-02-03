package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) HostSearch(query string) (*SearchResult, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/search?key=%s&query=%s",
		baseURL, c.apiKey, query))

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := &SearchResult{Matches: []Host{}}
	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
