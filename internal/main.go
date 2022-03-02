package internal

type Service interface {
	UpdateDNS()
}

type service struct {
	config *Config
}

func NewService() Service {
	return newService()
}

func newService() *service {
	return &service{
		config: NewConfig(),
	}
}

func (s *service) UpdateDNS() {
	//ctx := context.Background()
	//cf := cloudflare.NewClient(s.config.APIToken, s.config.Timeout)
	//ip := ifconfigme.NewClient(s.config.Timeout)
}
