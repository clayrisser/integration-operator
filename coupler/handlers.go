package coupler

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/silicon-hills/integration-operator/util"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/json"
	"github.com/tidwall/gjson"
)

type Config []byte

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HandlePlugCreated(plug gjson.Result) error {
	return h.processEvent(
		&plug,
		nil,
		nil,
		nil,
		plug.Get("spec.integrationEndpoint").String(),
		"created",
	)
}

func (h *Handlers) HandlePlugCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return h.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.integrationEndpoint").String(),
		"coupled",
	)
}

func (h *Handlers) HandlePlugUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return h.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.integrationEndpoint").String(),
		"updated",
	)
}

func (h *Handlers) HandlePlugDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return h.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.integrationEndpoint").String(),
		"decoupled",
	)
}

func (h *Handlers) HandlePlugBroken(
	plug gjson.Result,
) error {
	return h.processEvent(
		&plug,
		nil,
		nil,
		nil,
		plug.Get("spec.integrationEndpoint").String(),
		"broken",
	)
}

func (h *Handlers) HandleSocketCreated(socket gjson.Result) error {
	return h.processEvent(
		nil,
		&socket,
		nil,
		nil,
		socket.Get("spec.integrationEndpoint").String(),
		"created",
	)
}

func (h *Handlers) HandleSocketCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return h.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.integrationEndpoint").String(),
		"coupled",
	)
}

func (h *Handlers) HandleSocketUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return h.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.integrationEndpoint").String(),
		"updated",
	)
}

func (h *Handlers) HandleSocketDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return h.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.integrationEndpoint").String(),
		"decoupled",
	)
}

func (h *Handlers) HandleSocketBroken(
	socket gjson.Result,
) error {
	return h.processEvent(
		nil,
		&socket,
		nil,
		nil,
		socket.Get("spec.integrationEndpoint").String(),
		"broken",
	)
}

func (h *Handlers) processEvent(
	plug *gjson.Result,
	socket *gjson.Result,
	plugConfig *gjson.Result,
	socketConfig *gjson.Result,
	endpoint string,
	eventName string,
) error {
	m := minify.New()
	m.AddFunc("application/json", json.Minify)
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)

	body := `{"version":"1"`
	if plug != nil {
		body += fmt.Sprintf(`,"plug":%s`, plug.String())
	}
	if socket != nil {
		body += fmt.Sprintf(`,"socket":%s`, socket.String())
	}
	if plugConfig != nil {
		body += fmt.Sprintf(`,"plugConfig":%s`, plugConfig.String())
	}
	if socketConfig != nil {
		body += fmt.Sprintf(`,"socketConfig":%s`, socketConfig.String())
	}
	body += "}"
	body, err := m.String("application/json", body)
	if err != nil {
		return err
	}
	go func() {
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(util.GetEndpoint(endpoint) + "/" + eventName)
		if err != nil {
			errCh <- err
		}
		rCh <- r
	}()
	select {
	case _ = <-rCh:
		return nil
	case err := <-errCh:
		return err
	}
}
