package util

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/tdewolff/minify"
	minifyJson "github.com/tdewolff/minify/json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ApparatusUtil struct {
	dataUtil *DataUtil
	ctx      *context.Context
}

func NewApparatusUtil(
	ctx *context.Context,
) *ApparatusUtil {
	return &ApparatusUtil{
		ctx:      ctx,
		dataUtil: NewDataUtil(ctx),
	}
}

func (u *ApparatusUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
) ([]byte, error) {
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	m := minify.New()
	m.AddFunc("application/json", minifyJson.Minify)
	go func() {
		body := `{"version":1}`
		var err error
		if plug != nil {
			bPlug, err := json.Marshal(plug)
			if err != nil {
				errCh <- err
			}
			plugObj := &unstructured.Unstructured{}
			_, _, err = decUnstructured.Decode(bPlug, nil, plugObj)
			if err != nil {
				errCh <- err
			}
			body, err = sjson.Set(body, "plug", plugObj)
			if err != nil {
				errCh <- err
			}
		}
		body, err = m.String("application/json", string(body))
		if err != nil {
			errCh <- err
		}
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(GetEndpoint(plug.Spec.Apparatus.Endpoint) + "/config")
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
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	m := minify.New()
	m.AddFunc("application/json", minifyJson.Minify)
	go func() {
		body := `{"version":1}`
		var err error
		if socket != nil {
			bSocket, err := json.Marshal(socket)
			if err != nil {
				errCh <- err
			}
			socketObj := &unstructured.Unstructured{}
			_, _, err = decUnstructured.Decode(bSocket, nil, socketObj)
			if err != nil {
				errCh <- err
			}
			body, err = sjson.Set(body, "socket", socketObj)
			if err != nil {
				errCh <- err
			}
		}
		body, err = m.String("application/json", body)
		if err != nil {
			errCh <- err
		}
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(GetEndpoint(socket.Spec.Apparatus.Endpoint) + "/config")
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
		plug.Get("spec.apparatus.endpoint").String(),
		"created",
	)
}

func (u *ApparatusUtil) PlugCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.apparatus.endpoint").String(),
		"coupled",
	)
}

func (u *ApparatusUtil) PlugUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.apparatus.endpoint").String(),
		"updated",
	)
}

func (u *ApparatusUtil) PlugDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		plug.Get("spec.apparatus.endpoint").String(),
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
		plug.Get("spec.apparatus.endpoint").String(),
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
		plug.Get("spec.apparatus.endpoint").String(),
		"broken",
	)
}

func (u *ApparatusUtil) SocketCreated(socket gjson.Result) error {
	return u.processEvent(
		nil,
		&socket,
		nil,
		nil,
		socket.Get("spec.apparatus.endpoint").String(),
		"created",
	)
}

func (u *ApparatusUtil) SocketCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.apparatus.endpoint").String(),
		"coupled",
	)
}

func (u *ApparatusUtil) SocketUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,

) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.apparatus.endpoint").String(),
		"updated",
	)
}

func (u *ApparatusUtil) SocketDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.processEvent(
		&plug,
		&socket,
		&plugConfig,
		&socketConfig,
		socket.Get("spec.apparatus.endpoint").String(),
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
		socket.Get("spec.apparatus.endpoint").String(),
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
		socket.Get("spec.apparatus.endpoint").String(),
		"broken",
	)
}

func (u *ApparatusUtil) processEvent(
	plug *gjson.Result,
	socket *gjson.Result,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
	endpoint string,
	eventName string,
) error {
	m := minify.New()
	m.AddFunc("application/json", minifyJson.Minify)
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	body := `{"version":"1"}`
	var err error
	if plug != nil {
		body, err = sjson.Set(body, "plug", plug.String())
		if err != nil {
			return err
		}
	}
	if socket != nil {
		body, err = sjson.Set(body, "socket", socket.String())
		if err != nil {
			return err
		}
	}
	if plugConfig != nil {
		body, err = sjson.Set(body, "plugConfig", plugConfig)
		if err != nil {
			return err
		}
	}
	if socketConfig != nil {
		body, err = sjson.Set(body, "socketConfig", socketConfig)
		if err != nil {
			return err
		}
	}
	body, err = m.String("application/json", body)
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
