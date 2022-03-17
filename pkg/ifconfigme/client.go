package ifconfigme

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	ExternalIP(ctx context.Context) string
}

type IP struct {
	string
}

type client struct {
	httpClient *http.Client
}

func NewClient(timeout time.Duration) Client {
	return newClient(timeout)
}

func newClient(timeout time.Duration) *client {
	return &client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *client) ExternalIP(ctx context.Context) string {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ifconfig.me", http.NoBody)
	if err != nil {
		return ""
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return ""
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return ""
	}

	return string(body)
}
