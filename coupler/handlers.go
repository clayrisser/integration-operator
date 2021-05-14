package coupler

import (
	"fmt"

	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"
)

type Config []byte

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
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

func (h *Handlers) HandleDeparted() error {
	fmt.Println("departed")
	return nil
}

func (h *Handlers) HandleBroken() error {
	fmt.Println("broken")
	return nil
}
