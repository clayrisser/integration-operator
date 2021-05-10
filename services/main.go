package services

type Services struct {
	Intrface *InterfaceService
	Plug     *PlugService
	Util     *UtilService
}

func NewServices() *Services {
	services := &Services{}
	services.Intrface = NewInterfaceService(services)
	services.Plug = NewPlugService(services)
	services.Util = NewUtilService(services)
	return services
}
