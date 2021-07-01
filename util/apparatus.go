/*
 * File: /util/apparatus.go
 * Project: integration-operator
 * File Created: 23-06-2021 22:14:06
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 01-07-2021 16:40:56
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
	"github.com/silicon-hills/integration-operator/config"
	"github.com/tdewolff/minify"
	minifyJson "github.com/tdewolff/minify/json"
	"github.com/tidwall/sjson"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

var startedApparatusTimers map[string]*time.Timer = map[string]*time.Timer{}

type ApparatusUtil struct {
	client   *kubernetes.Clientset
	ctx      *context.Context
	dataUtil *DataUtil
	log      logr.Logger
	varUtil  *VarUtil
}

func NewApparatusUtil(
	ctx *context.Context,
) *ApparatusUtil {
	return &ApparatusUtil{
		client:   kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		ctx:      ctx,
		dataUtil: NewDataUtil(ctx),
		log:      ctrl.Log.WithName("util.ApparatusUtil"),
		varUtil:  NewVarUtil(ctx),
	}
}

func (u *ApparatusUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
) ([]byte, error) {
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
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
		url := u.getPlugEndpoint(plug) + "/config"
		u.log.Info("getting plug config", "method", "POST", "url", url)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(url)
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
		url := u.getSocketEndpoint(socket) + "/config"
		u.log.Info("getting socket config", "method", "POST", "url", url)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(url)
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
		plug.Spec.Apparatus,
		plug.Name,
		plug.Namespace,
		string(plug.UID),
		u.getPlugEndpoint(plug),
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
		plug.Spec.Apparatus,
		plug.Name,
		plug.Namespace,
		string(plug.UID),
		u.getPlugEndpoint(plug),
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
		plug.Spec.Apparatus,
		plug.Name,
		plug.Namespace,
		string(plug.UID),
		u.getPlugEndpoint(plug),
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
		plug.Spec.Apparatus,
		plug.Name,
		plug.Namespace,
		string(plug.UID),
		u.getPlugEndpoint(plug),
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
		plug.Spec.Apparatus,
		plug.Name,
		plug.Namespace,
		string(plug.UID),
		u.getPlugEndpoint(plug),
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
		plug.Spec.Apparatus,
		plug.Name,
		plug.Namespace,
		string(plug.UID),
		u.getPlugEndpoint(plug),
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
		socket.Spec.Apparatus,
		socket.Name,
		socket.Namespace,
		string(socket.UID),
		u.getSocketEndpoint(socket),
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
		socket.Spec.Apparatus,
		socket.Name,
		socket.Namespace,
		string(socket.UID),
		u.getSocketEndpoint(socket),
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
		socket.Spec.Apparatus,
		socket.Name,
		socket.Namespace,
		string(socket.UID),
		u.getSocketEndpoint(socket),
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
		socket.Spec.Apparatus,
		socket.Name,
		socket.Namespace,
		string(socket.UID),
		u.getSocketEndpoint(socket),
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
		socket.Spec.Apparatus,
		socket.Name,
		socket.Namespace,
		string(socket.UID),
		u.getSocketEndpoint(socket),
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
		socket.Spec.Apparatus,
		socket.Name,
		socket.Namespace,
		string(socket.UID),
		u.getSocketEndpoint(socket),
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

func (u *ApparatusUtil) RenewIdleTimeout(
	apparatus *integrationv1alpha2.SpecApparatus,
	name string,
	namespace string,
	uid string,
) {
	if apparatus == nil {
		return
	}
	idleTimeout := time.Second * 60
	if apparatus.IdleTimeout != 0 {
		idleTimeout = time.Second * time.Duration(apparatus.IdleTimeout)
	}
	if timer, ok := startedApparatusTimers[uid]; ok {
		timer.Reset(idleTimeout)
	} else {
		if _, err := u.client.CoreV1().Pods(namespace).Get(
			*u.ctx,
			name,
			metav1.GetOptions{},
		); err != nil {
			if errors.IsNotFound(err) {
				return
			}
		}
		startedApparatusTimers[uid] = time.AfterFunc(idleTimeout, func() {
			if err := u.client.CoreV1().Pods(namespace).Delete(
				*u.ctx,
				name,
				metav1.DeleteOptions{},
			); err != nil {
				u.log.Error(
					err,
					"failed to terminate idle apparatus "+namespace+"."+name,
				)
			} else {
				u.log.Info("terminated idle apparatus " + namespace + "." + name)
			}
		})
	}
}

func (u *ApparatusUtil) StartFromPlug(plug *integrationv1alpha2.Plug) (bool, error) {
	return u.start(
		plug.Spec.Apparatus,
		plug.Name,
		plug.Namespace,
		string(plug.UID),
		u.createPlugOwnerReference(plug),
	)
}

func (u *ApparatusUtil) StartFromSocket(socket *integrationv1alpha2.Socket) (bool, error) {
	return u.start(
		socket.Spec.Apparatus,
		socket.Name,
		socket.Namespace,
		string(socket.UID),
		u.createSocketOwnerReference(socket),
	)
}

func (u *ApparatusUtil) start(
	apparatus *integrationv1alpha2.SpecApparatus,
	name string,
	namespace string,
	uid string,
	ownerReference metav1.OwnerReference,
) (bool, error) {
	if apparatus == nil || len(*apparatus.Containers) <= 0 {
		return false, nil
	}
	idleTimeout := time.Second * 60
	if apparatus.IdleTimeout != 0 {
		idleTimeout = time.Second * time.Duration(apparatus.IdleTimeout)
	}
	alreadyExists := false
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			OwnerReferences: []metav1.OwnerReference{
				ownerReference,
			},
			Labels: map[string]string{
				"apparatus": name,
			},
		},
		Spec: v1.PodSpec{
			Affinity: &v1.Affinity{
				NodeAffinity: &v1.NodeAffinity{
					RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
						NodeSelectorTerms: []v1.NodeSelectorTerm{
							{
								MatchExpressions: []v1.NodeSelectorRequirement{
									{
										Key:      "kubernetes.io/arch",
										Operator: v1.NodeSelectorOpIn,
										Values:   []string{"amd64"},
									},
								},
							},
						},
					},
				},
			},
			Containers: *apparatus.Containers,
		},
	}
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			OwnerReferences: []metav1.OwnerReference{
				ownerReference,
			},
		},
		Spec: v1.ServiceSpec{
			Type: v1.ServiceTypeClusterIP,
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					Protocol:   v1.ProtocolTCP,
					TargetPort: intstr.FromString("container"),
				},
			},
			Selector: map[string]string{
				"apparatus": name,
			},
		},
	}
	_, err := u.client.CoreV1().Services(namespace).Create(
		*u.ctx,
		service,
		metav1.CreateOptions{
			FieldManager: "integration-operator",
		},
	)
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return false, err
		}
	}
	_, err = u.client.CoreV1().Pods(namespace).Create(
		*u.ctx,
		pod,
		metav1.CreateOptions{
			FieldManager: "integration-operator",
		},
	)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			alreadyExists = true
		} else {
			return false, err
		}
	}
	if timer, ok := startedApparatusTimers[uid]; ok {
		timer.Reset(idleTimeout)
	} else {
		startedApparatusTimers[uid] = time.AfterFunc(idleTimeout, func() {
			if err := u.client.CoreV1().Pods(namespace).Delete(
				*u.ctx,
				name,
				metav1.DeleteOptions{},
			); err != nil {
				u.log.Error(
					err,
					"failed to terminate idle apparatus "+namespace+"."+name,
				)
			} else {
				u.log.Info("terminated idle apparatus " + namespace + "." + name)
			}
		})
	}
	if !alreadyExists {
		u.log.Info("started apparatus " + namespace + "." + name)
	}
	return true, nil
}

func (u *ApparatusUtil) createPlugOwnerReference(plug *integrationv1alpha2.Plug) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: plug.APIVersion,
		Kind:       plug.Kind,
		Name:       plug.Name,
		UID:        plug.UID,
	}
}

func (u *ApparatusUtil) createSocketOwnerReference(socket *integrationv1alpha2.Socket) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: socket.APIVersion,
		Kind:       socket.Kind,
		Name:       socket.Name,
		UID:        socket.UID,
	}
}

func (u *ApparatusUtil) getPlugEndpoint(plug *integrationv1alpha2.Plug) string {
	return u.getEndpoint(plug.Name, plug.Namespace, plug.Spec.Apparatus, config.DebugPlugEndpoint)
}

func (u *ApparatusUtil) getSocketEndpoint(socket *integrationv1alpha2.Socket) string {
	return u.getEndpoint(socket.Name, socket.Namespace, socket.Spec.Apparatus, config.DebugSocketEndpoint)
}

func (u *ApparatusUtil) getEndpoint(
	name string,
	namespace string,
	apparatus *integrationv1alpha2.SpecApparatus,
	debugEndpoint string,
) string {
	endpoint := debugEndpoint
	if endpoint == "" {
		endpoint = apparatus.Endpoint
		if endpoint == "" || endpoint[0] == '/' {
			if apparatus.Containers != nil &&
				len(*apparatus.Containers) > 0 {
				endpoint = name + "." + namespace + ".svc.cluster.local" + endpoint
			} else {
				return "http://localhost" + endpoint
			}
		}
	}
	if len(endpoint) < 7 || (endpoint[0:8] != "https://" && endpoint[0:7] != "http://") {
		endpoint = "http://" + endpoint
	}
	if endpoint[len(endpoint)-1] == '/' {
		endpoint = string(endpoint[0 : len(endpoint)-2])
	}
	return endpoint
}

func (u *ApparatusUtil) processEvent(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
	apparatus *integrationv1alpha2.SpecApparatus,
	name string,
	namespace string,
	uid string,
	endpoint string,
	eventName string,
) error {
	u.RenewIdleTimeout(apparatus, name, namespace, uid)
	min := minify.New()
	min.AddFunc("application/json", minifyJson.Minify)
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
	body, err = min.String("application/json", body)
	if err != nil {
		return err
	}
	go func() {
		url := endpoint + "/" + eventName
		u.log.Info("triggered event "+eventName, "method", "POST", "url", url)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(url)
		if err != nil {
			errCh <- err
		}
		rCh <- r
	}()
	select {
	case <-rCh:
		return nil
	case err := <-errCh:
		return err
	}
}
