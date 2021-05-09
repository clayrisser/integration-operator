package services

type InterfaceService struct {
	s *Services
}

func NewInterfaceService(services *Services) *InterfaceService {
	return &InterfaceService{s: services}
}
