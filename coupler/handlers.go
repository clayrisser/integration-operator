package coupler

import (
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/services"
)

type Payload interface{}

type Handlers struct {
	s *services.Services
}

func NewHandlers() *Handlers {
	return &Handlers{s: services.NewServices()}
}

func (h *Handlers) HandlePlugCreated(plug *integrationv1alpha2.Plug) {}

func (h *Handlers) HandlePlugJoined(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug, payload Payload) {
}

func (h *Handlers) HandlePlugChanged(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug, payload Payload) {
}

func (h *Handlers) HandlePlugDeparted() {
}

func (h *Handlers) HandlePlugBroken(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug, payload Payload) {
}

func (h *Handlers) HandleSocketCreated(plug *integrationv1alpha2.Plug) {}

func (h *Handlers) HandleSocketJoined(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
}

func (h *Handlers) HandleSocketChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
}

func (h *Handlers) HandleSocketDeparted() {
}

func (h *Handlers) HandleSocketBroken(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
}
