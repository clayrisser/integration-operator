package coupler

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/silicon-hills/integration-operator/services"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"
)

type Config []byte

type Handlers struct {
	s *services.Services
}

func NewHandlers() *Handlers {
	return &Handlers{s: services.NewServices()}
}

func (h *Handlers) HandlePlugCreated(plug gjson.Result) error {
	fmt.Println("plug created")
	y, err := yaml.JSONToYAML([]byte(plug.String()))
	if err != nil {
		return err
	}
	fmt.Println(string(y))
	return nil
}

func (h *Handlers) HandleGetConfig() (Config, error) {
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	go func() {
		r, err := client.R().EnableTrace().SetQueryParams(map[string]string{
			"version": "1",
		}).Get("http://localhost:3000/config")
		if err != nil {
			errCh <- err
		}
		rCh <- r
	}()
	select {
	case r := <-rCh:
		return r.Body(), nil
	case err := <-errCh:
		return nil, err
	case <-time.After(3 * time.Second):
		return nil, errors.New("timeout")
	}
}

func (h *Handlers) HandleJoined(plug gjson.Result, socket gjson.Result, config gjson.Result) {
	fmt.Println("joined")
}

func (h *Handlers) HandlePlugChanged(plug gjson.Result, socket gjson.Result, config gjson.Result) {
	fmt.Println("plug changed")
}

func (h *Handlers) HandleSocketCreated(socket gjson.Result) {
	fmt.Println("socket created")

}

func (h *Handlers) HandleSocketChanged(plug gjson.Result, socket gjson.Result, config gjson.Result) {
	fmt.Println("socket changed")
}

func (h *Handlers) HandleDeparted() {
	fmt.Println("departed")
}

func (h *Handlers) HandleBroken() {
	fmt.Println("broken")
}
