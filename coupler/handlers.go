package coupler

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/json"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"
)

type Config []byte

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HandlePlugCreated(plug gjson.Result) error {
	m := minify.New()
	m.AddFunc("application/json", json.Minify)
	endpoint := plug.Get("spec.configEndpoint").String()
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	body, err := m.String("application/json", fmt.Sprintf(`{
  "version": "1",
  "plug": %s
}`, plug.String()))
	if err != nil {
		return err
	}
	go func() {
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(endpoint)
		if err != nil {
			errCh <- err
		}
		rCh <- r
	}()
	select {
	case r := <-rCh:
		fmt.Printf(string(r.Body()))
		return nil
	case err := <-errCh:
		return err
	case <-time.After(3 * time.Second):
		return errors.New("timeout")
	}
}

func (h *Handlers) HandlePlugJoined(plug gjson.Result, socket gjson.Result, config gjson.Result) error {
	fmt.Printf("plug joined\n")
	y, err := yaml.JSONToYAML([]byte(config.String()))
	if err != nil {
		return err
	}
	fmt.Println(string(y))
	return nil
}

func (h *Handlers) HandleSocketJoined(plug gjson.Result, socket gjson.Result, config gjson.Result) error {
	fmt.Printf("socket joined\n")
	y, err := yaml.JSONToYAML([]byte(config.String()))
	if err != nil {
		return err
	}
	fmt.Println(string(y))
	return nil
}

func (h *Handlers) HandlePlugChanged(plug gjson.Result, socket gjson.Result, config gjson.Result) error {
	fmt.Println("plug changed")
	return nil
}

func (h *Handlers) HandleSocketCreated(socket gjson.Result) error {
	fmt.Println("socket created")
	return nil

}

func (h *Handlers) HandleSocketChanged(plug gjson.Result, socket gjson.Result, config gjson.Result) error {
	fmt.Println("socket changed")
	return nil
}

func (h *Handlers) HandleDeparted(plug gjson.Result, socket gjson.Result) error {
	fmt.Println("departed")
	return nil
}

func (h *Handlers) HandleBroken() error {
	fmt.Println("broken")
	return nil
}

type CouplerBody struct {
	Config  gjson.Result `json:"config,omitempty"`
	Plug    gjson.Result `json:"plug,omitempty"`
	Socket  gjson.Result `json:"socket,omitempty"`
	Version string       `json:"version,omitempty"`
}
