package services

type PlugService struct {
	s *Services
}

func NewPlugService(services *Services) *PlugService {
	return &PlugService{s: services}
}
