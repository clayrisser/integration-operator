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

func (h *Handlers) HandleJoined(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
}

func (h *Handlers) HandlePlugChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
}

func (h *Handlers) HandleSocketCreated(socket *integrationv1alpha2.Socket) {}

func (h *Handlers) HandleSocketChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
}

func (h *Handlers) HandleDeparted(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
}

func (h *Handlers) HandleBroken(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket) {
}
