package internal

type Service interface {
}

type service struct {
}

func NewService() Service {
	return newService()
}

func newService() *service {
	return &service{}
}

func (s *service) UpdateDNS() {
	client := cloudflare.NewClient()
}
