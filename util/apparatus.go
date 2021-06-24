package util

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/tdewolff/minify"
	minifyJson "github.com/tdewolff/minify/json"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"
)

type ApparatusUtil struct{}

func NewApparatusUtil() *ApparatusUtil {
	return &ApparatusUtil{}
}

func (u *ApparatusUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
) ([]byte, error) {
	bPlug, err := json.Marshal(plug)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	m := minify.New()
	m.AddFunc("application/json", minifyJson.Minify)
	go func() {
		body := `{"version":"1"`
		if plug != nil {
			jsonPlug := gjson.Parse(string(bPlug))
			body += fmt.Sprintf(`,"plug":%s`, jsonPlug)
			meta, _ := yaml.YAMLToJSON([]byte(jsonPlug.Get("spec").Get("meta").String()))
			if meta == nil {
				meta = []byte("{}")
			}
			body += fmt.Sprintf(`,"plugMeta":%s`, meta)
		}
		body += "}"
		body, err := m.String("application/json", body)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(GetEndpoint(plug.Spec.IntegrationEndpoint) + "/config")
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
	}
}

func (u *ApparatusUtil) GetSocketConfig(
	socket *integrationv1alpha2.Socket,
) ([]byte, error) {
	bSocket, err := json.Marshal(socket)
	if err != nil {
		return nil, err
	}
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	m := minify.New()
	m.AddFunc("application/json", minifyJson.Minify)
	go func() {
		body := `{"version":"1"`
		if socket != nil {
			jsonSocket := gjson.Parse(string(bSocket))
			body += fmt.Sprintf(`,"socket":%s`, jsonSocket)
			meta, _ := yaml.YAMLToJSON([]byte(jsonSocket.Get("spec").Get("meta").String()))
			if meta == nil {
				meta = []byte("{}")
			}
			body += fmt.Sprintf(`,"socketMeta":%s`, meta)
		}
		body += "}"
		body, err := m.String("application/json", body)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(GetEndpoint(socket.Spec.IntegrationEndpoint) + "/config")
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
	}
}

func (u *ApparatusUtil) PlugCreated(plug gjson.Result) error {
	return u.processEvent(
		&plug,
		nil,
		nil,
		nil,
		plug.Get("spec.integrationEndpoint").String(),
		"created",
	)
}

func (u *ApparatusUtil) PlugCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.integrationEndpoint").String(),
		"coupled",
	)
}

func (u *ApparatusUtil) PlugUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.integrationEndpoint").String(),
		"updated",
	)
}

func (u *ApparatusUtil) PlugDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.integrationEndpoint").String(),
		"decoupled",
	)
}

func (u *ApparatusUtil) PlugDeleted(
	plug gjson.Result,
) error {
	return u.processEvent(
		&plug,
		nil,
		nil,
		nil,
		plug.Get("spec.integrationEndpoint").String(),
		"deleted",
	)
}

func (u *ApparatusUtil) PlugBroken(
	plug gjson.Result,
) error {
	return u.processEvent(
		&plug,
		nil,
		nil,
		nil,
		plug.Get("spec.integrationEndpoint").String(),
		"broken",
	)
}

func (u *ApparatusUtil) SocketCreated(socket gjson.Result) error {
	return u.processEvent(
		nil,
		&socket,
		nil,
		nil,
		socket.Get("spec.integrationEndpoint").String(),
		"created",
	)
}

func (u *ApparatusUtil) SocketCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.integrationEndpoint").String(),
		"coupled",
	)
}

func (u *ApparatusUtil) SocketUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.integrationEndpoint").String(),
		"updated",
	)
}

func (u *ApparatusUtil) SocketDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig gjson.Result,
	socketConfig gjson.Result,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.integrationEndpoint").String(),
		"decoupled",
	)
}

func (u *ApparatusUtil) SocketDeleted(
	socket gjson.Result,
) error {
	return u.processEvent(
		nil,
		&socket,
		nil,
		nil,
		socket.Get("spec.integrationEndpoint").String(),
		"deleted",
	)
}

func (u *ApparatusUtil) SocketBroken(
	socket gjson.Result,
) error {
	return u.processEvent(
		nil,
		&socket,
		nil,
		nil,
		socket.Get("spec.integrationEndpoint").String(),
		"broken",
	)
}

func (u *ApparatusUtil) processEvent(
	plug *gjson.Result,
	socket *gjson.Result,
	plugConfig *gjson.Result,
	socketConfig *gjson.Result,
	endpoint string,
	eventName string,
) error {
	m := minify.New()
	m.AddFunc("application/json", minifyJson.Minify)
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
		}).SetBody([]byte(body)).Post(GetEndpoint(endpoint) + "/" + eventName)
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
