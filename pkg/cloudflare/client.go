package cloudflare

type Client interface {
}

type client struct {
	APIKey   string
	ClientID string
	Token    string
}

func NewClient(apiKey, clientID, token string) Client {
	return newClient(apiKey, clientID, token)
}

func newClient(apiKey, clientID, token string) *client {
	return &client{
		APIKey:   apiKey,
		ClientID: clientID,
		Token:    token,
	}
}
