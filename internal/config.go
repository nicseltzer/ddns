package internal

type Config interface {
}

type config struct {
	APIKey   string
	ClientID string
	Token    string
}

func NewConfig() Config {

	return &config{}
}
