/*
 * File: /util/apparatus.go
 * Project: integration-operator
 * File Created: 23-06-2021 22:14:06
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 29-06-2021 08:57:18
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Silicon Hills LLC (c) Copyright 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"context"
	"net"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-resty/resty/v2"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/tdewolff/minify"
	minifyJson "github.com/tdewolff/minify/json"
	"github.com/tidwall/sjson"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

var startedApparatus map[string]*time.Timer = map[string]*time.Timer{}

type ApparatusUtil struct {
	client   *kubernetes.Clientset
	ctx      *context.Context
	dataUtil *DataUtil
	varUtil  *VarUtil
}

func NewApparatusUtil(
	ctx *context.Context,
) *ApparatusUtil {
	return &ApparatusUtil{
		client:   kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		ctx:      ctx,
		dataUtil: NewDataUtil(ctx),
		varUtil:  NewVarUtil(ctx),
	}
}

func (u *ApparatusUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
) ([]byte, error) {
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	defer close(rCh)
	defer close(errCh)
	min := minify.New()
	min.AddFunc("application/json", minifyJson.Minify)
	go func() {
		body := `{"version":1}`
		var err error
		if plug != nil {
			body, err = sjson.Set(body, "plug", plug)
			if err != nil {
				errCh <- err
				return
			}

			data, err := u.dataUtil.GetPlugData(plug)
			if err != nil {
				errCh <- err
				return
			}
			body, err = sjson.Set(body, "data", data)
			if err != nil {
				errCh <- err
				return
			}

			if plug.Spec.Vars != nil {
				vars, err := u.varUtil.GetVars(plug.Namespace, plug.Spec.Vars)
				if err != nil {
					errCh <- err
					return
				}
				body, err = sjson.Set(body, "vars", vars)
				if err != nil {
					errCh <- err
					return
				}
			}
		}
		body, err = min.String("application/json", string(body))
		if err != nil {
			errCh <- err
			return
		}
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(GetEndpoint(plug.Spec.Apparatus.Endpoint) + "/config")
		if err != nil {
			errCh <- err
			return
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
	defer close(rCh)
	defer close(errCh)
	min := minify.New()
	min.AddFunc("application/json", minifyJson.Minify)
	go func() {
		body := `{"version":1}`
		var err error
		if socket != nil {
			body, err = sjson.Set(body, "socket", socket)
			if err != nil {
				errCh <- err
				return
			}

			data, err := u.dataUtil.GetSocketData(socket)
			if err != nil {
				errCh <- err
				return
			}
			body, err = sjson.Set(body, "data", data)
			if err != nil {
				errCh <- err
				return
			}

			if socket.Spec.Vars != nil {
				vars, err := u.varUtil.GetVars(socket.Namespace, socket.Spec.Vars)
				if err != nil {
					errCh <- err
					return
				}
				body, err = sjson.Set(body, "vars", vars)
				if err != nil {
					errCh <- err
					return
				}
			}
		}
		body, err = min.String("application/json", body)
		if err != nil {
			errCh <- err
			return
		}
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(GetEndpoint(socket.Spec.Apparatus.Endpoint) + "/config")
		if err != nil {
			errCh <- err
			return
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
	if plug.Spec.Apparatus == nil {
		return nil
	}
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
	if plug.Spec.Apparatus == nil {
		return nil
	}
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
	if plug.Spec.Apparatus == nil {
		return nil
	}
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
	if plug.Spec.Apparatus == nil {
		return nil
	}
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
	if plug.Spec.Apparatus == nil {
		return nil
	}
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
	if plug.Spec.Apparatus == nil {
		return nil
	}
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
	if socket.Spec.Apparatus == nil {
		return nil
	}
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
	if socket.Spec.Apparatus == nil {
		return nil
	}
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
	if socket.Spec.Apparatus == nil {
		return nil
	}
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
	if socket.Spec.Apparatus == nil {
		return nil
	}
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
	if socket.Spec.Apparatus == nil {
		return nil
	}
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
	if socket.Spec.Apparatus == nil {
		return nil
	}
	return u.processEvent(
		nil,
		socket,
		nil,
		nil,
		socket.Spec.Apparatus.Endpoint,
		"broken",
	)
}

func (u *ApparatusUtil) NotRunning(err error) bool {
	// TODO: detect if apparatus api call
	if nerr, ok := err.(net.Error); ok {
		return !nerr.Timeout() && !nerr.Temporary()
	}
	return false
}

func (u *ApparatusUtil) Start(
	log *logr.Logger,
	apparatus *integrationv1alpha2.SpecApparatus,
	namespace string,
	uid string,
) (bool, error) {
	if apparatus == nil || len(apparatus.Containers) <= 0 {
		return false, nil
	}
	configMap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "abc",
		},
		Data: map[string]string{
			"Hello": "world",
		},
	}
	configMap, err := u.client.CoreV1().ConfigMaps(namespace).Create(
		*u.ctx,
		configMap,
		metav1.CreateOptions{
			FieldManager: "integration-operator",
		},
	)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			u.registerIdleTimeout(log, configMap, namespace, uid)
			return false, nil
		} else {
			return false, err
		}
	}
	u.registerIdleTimeout(log, configMap, namespace, uid)
	(*log).Info("started apparatus " + configMap.Namespace + "." + configMap.Name)
	return true, nil
}

func (u *ApparatusUtil) registerIdleTimeout(
	log *logr.Logger,
	configMap *v1.ConfigMap,
	namespace string,
	uid string,
) {
	startedApparatus[uid] = time.AfterFunc(time.Second*15, func() {
		if err := u.client.CoreV1().ConfigMaps(namespace).Delete(
			*u.ctx,
			configMap.Name,
			metav1.DeleteOptions{},
		); err != nil {
			(*log).Error(
				err,
				"failed to terminate idle apparatus "+configMap.Namespace+"."+configMap.Name,
			)
		} else {
			(*log).Info("terminated idle apparatus " + configMap.Namespace + "." + configMap.Name)
		}
	})
}

func (u *ApparatusUtil) processEvent(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
	endpoint string,
	eventName string,
) error {
	min := minify.New()
	min.AddFunc("application/json", minifyJson.Minify)
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	defer close(rCh)
	defer close(errCh)
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
	body, err = min.String("application/json", body)
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
