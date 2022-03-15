package ifconfigme

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

type Client interface {
	Ip(ctx context.Context) string
}

type client struct {
	timeout time.Duration
}

func NewClient(timeout time.Duration) Client {
	return newClient(timeout)
}

func newClient(timeout time.Duration) *client {
	return &client{
		timeout: timeout,
	}
}

func (c *client) Ip(ctx context.Context) string {
	h := c.createHTTPClient()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ifconfig.me", http.NoBody)
	if err != nil {
		return ""
	}

	response, err := h.Do(request)
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

func (c *client) createHTTPClient() *http.Client {
	return &http.Client{
		Timeout: c.timeout,
	}
}
