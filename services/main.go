package services

type Services struct {
	Util *UtilService
}

func NewServices() *Services {
	services := &Services{}
	services.Util = NewUtilService(services)
	return services
}
