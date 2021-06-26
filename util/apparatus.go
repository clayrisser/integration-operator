package util

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/tdewolff/minify"
	minifyJson "github.com/tdewolff/minify/json"
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

func (u *ApparatusUtil) PlugCreated(plug *integrationv1alpha2.Plug) error {
	return u.processEvent(
		plug,
		nil,
		nil,
		nil,
		plug.Spec.Apparatus.Endpoint,
		"created",
	)
}

func (u *ApparatusUtil) PlugCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.processEvent(
		plug,
		socket,
		plugConfig,
		socketConfig,
		plug.Spec.Apparatus.Endpoint,
		"coupled",
	)
}

func (u *ApparatusUtil) PlugUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.processEvent(
		plug,
		socket,
		plugConfig,
		socketConfig,
		plug.Spec.Apparatus.Endpoint,
		"updated",
	)
}

func (u *ApparatusUtil) PlugDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.processEvent(
		plug,
		socket,
		plugConfig,
		socketConfig,
		plug.Spec.Apparatus.Endpoint,
		"decoupled",
	)
}

func (u *ApparatusUtil) PlugDeleted(
	plug *integrationv1alpha2.Plug,
) error {
	return u.processEvent(
		plug,
		nil,
		nil,
		nil,
		plug.Spec.Apparatus.Endpoint,
		"deleted",
	)
}

func (u *ApparatusUtil) PlugBroken(
	plug *integrationv1alpha2.Plug,
) error {
	return u.processEvent(
		plug,
		nil,
		nil,
		nil,
		plug.Spec.Apparatus.Endpoint,
		"broken",
	)
}

func (u *ApparatusUtil) SocketCreated(socket *integrationv1alpha2.Socket) error {
	return u.processEvent(
		nil,
		socket,
		nil,
		nil,
		socket.Spec.Apparatus.Endpoint,
		"created",
	)
}

func (u *ApparatusUtil) SocketCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.processEvent(
		plug,
		socket,
		plugConfig,
		socketConfig,
		socket.Spec.Apparatus.Endpoint,
		"coupled",
	)
}

func (u *ApparatusUtil) SocketUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,

) error {
	return u.processEvent(
		plug,
		socket,
		plugConfig,
		socketConfig,
		socket.Spec.Apparatus.Endpoint,
		"updated",
	)
}

func (u *ApparatusUtil) SocketDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.processEvent(
		plug,
		socket,
		plugConfig,
		socketConfig,
		socket.Spec.Apparatus.Endpoint,
		"decoupled",
	)
}

func (u *ApparatusUtil) SocketDeleted(
	socket *integrationv1alpha2.Socket,
) error {
	return u.processEvent(
		nil,
		socket,
		nil,
		nil,
		socket.Spec.Apparatus.Endpoint,
		"deleted",
	)
}

func (u *ApparatusUtil) SocketBroken(
	socket *integrationv1alpha2.Socket,
) error {
	return u.processEvent(
		nil,
		socket,
		nil,
		nil,
		socket.Spec.Apparatus.Endpoint,
		"broken",
	)
}

func (u *ApparatusUtil) processEvent(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
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
		body, err = sjson.Set(body, "plug", plug)
		if err != nil {
			return err
		}
	}
	if socket != nil {
		body, err = sjson.Set(body, "socket", socket)
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
