package coupler

import (
	"fmt"

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

func (h *Handlers) HandlePlugCreated(plug *integrationv1alpha2.Plug) {
	fmt.Println("plug created")
}

func (h *Handlers) HandleJoined(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
	fmt.Println("joined")
}

func (h *Handlers) HandlePlugChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
	fmt.Println("plug changed")
}

func (h *Handlers) HandleSocketCreated(socket *integrationv1alpha2.Socket) {
	fmt.Println("socket created")

}

func (h *Handlers) HandleSocketChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
	fmt.Println("socket changed")
}

func (h *Handlers) HandleDeparted(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, payload Payload) {
	fmt.Println("departed")
}

func (h *Handlers) HandleBroken(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket) {
	fmt.Println("broken")
}
