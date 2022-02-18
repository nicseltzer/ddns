package cloudflare

type Client interface {
}

type client struct {
}

func NewClient() Client {
	return newClient()
}

func newClient() *client {
	return &client{}
}