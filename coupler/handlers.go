package coupler

import (
	"fmt"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/services"
)

type Config []byte

type Handlers struct {
	s *services.Services
}

func NewHandlers() *Handlers {
	return &Handlers{s: services.NewServices()}
}

func (h *Handlers) HandlePlugCreated(plug *integrationv1alpha2.Plug) {
	fmt.Println("plug created")
}

// func (h *Handlers) HandleGetConfig() (Config error) {
// 	client := resty.New()
// 	rCh := make(chan *resty.Response)
// 	errCh := make(chan error)
// 	go func() {
// 		r, err := client.R().EnableTrace().SetQueryParams(map[string]string{
// 			"version": "1",
// 		}).Get("http://localhost:3000")
// 		if err != nil {
// 			errCh <- err
// 		}
// 		rCh <- r
// 	}()

// 	// do stuff

// 	select {
// 	case r := <-rCh:
// 		return str(r)
// 	case err := <-errCh:
// 		return err
// 	case <-time.After(3 * time.Second):
// 		return errors.New("timeout")
// 	}
// 	return nil, nil
// }

func (h *Handlers) HandleJoined(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, config Config) {
	fmt.Println("joined")
}

func (h *Handlers) HandlePlugChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, config Config) {
	fmt.Println("plug changed")
}

func (h *Handlers) HandleSocketCreated(socket *integrationv1alpha2.Socket) {
	fmt.Println("socket created")

}

func (h *Handlers) HandleSocketChanged(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, config Config) {
	fmt.Println("socket changed")
}

func (h *Handlers) HandleDeparted(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket, config Config) {
	fmt.Println("departed")
}

func (h *Handlers) HandleBroken(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket) {
	fmt.Println("broken")
}
