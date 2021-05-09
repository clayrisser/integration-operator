package services

type Services struct {
	Events   *EventsService
	Intrface *InterfaceService
	Plug     *PlugService
	Util     *UtilService
}

func NewServices() *Services {
	services := &Services{}
	services.Events = NewEventsService(services)
	services.Intrface = NewInterfaceService(services)
	services.Plug = NewPlugService(services)
	services.Util = NewUtilService(services)
	return services
}
