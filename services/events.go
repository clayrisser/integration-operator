package services

import (
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
)

type EventsService struct {
	s *Services
}

func NewEventsService(services *Services) *EventsService {
	return &EventsService{s: services}
}

func (e *EventsService) HandlePlugCreated(plug *integrationv1alpha2.Plug) {}

func (e *EventsService) HandlePlugJoined(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug, payload []byte) {
}

func (e *EventsService) HandlePlugChanged(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug, payload []byte) {
}

func (e *EventsService) HandlePlugDeparted(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug, payload []byte) {
}

func (e *EventsService) HandlePlugBroken(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug, payload []byte) {
}

func (e *EventsService) HandleSocketCreated(plug *integrationv1alpha2.Plug) {}

func (e *EventsService) HandleSocketJoined(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload []byte) {
}

func (e *EventsService) HandleSocketChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload []byte) {
}

func (e *EventsService) HandleSocketDeparted(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload []byte) {
}

func (e *EventsService) HandleSocketBroken(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload []byte) {
}
