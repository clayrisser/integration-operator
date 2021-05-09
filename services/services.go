package services

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

type Services struct {
	Events   *EventsService
	Intrface *InterfaceService
	Plug     *PlugService
	Util     *UtilService
	ctx      *context.Context
	req      *ctrl.Request
}

func NewServices(ctx *context.Context, req *ctrl.Request) *Services {
	services := &Services{
		ctx: ctx,
		req: req,
	}
	services.Events = NewEventsService(services)
	services.Intrface = NewInterfaceService(services)
	services.Plug = NewPlugService(services)
	services.Util = NewUtilService(services)
	return services
}
