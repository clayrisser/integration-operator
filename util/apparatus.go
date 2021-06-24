package util

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/tdewolff/minify"
	"github.com/tidwall/gjson"
	"sigs.k8s.io/yaml"

	minifyJson "github.com/tdewolff/minify/json"
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
