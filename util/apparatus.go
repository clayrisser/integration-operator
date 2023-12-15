/**
 * File: /util/apparatus.go
 * Project: integration-operator
 * File Created: 17-10-2023 13:49:54
 * Author: Clay Risser
 * -----
 * BitSpur (c) Copyright 2021 - 2023
 *
 * Licensed under the GNU Affero General Public License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.gnu.org/licenses/agpl-3.0.en.html
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving this software without disclosing
 * the source code of your own applications.
 */

package util

import (
	"context"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-resty/resty/v2"
	"github.com/tdewolff/minify"
	minifyJson "github.com/tdewolff/minify/json"
	"github.com/tidwall/sjson"
	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"gitlab.com/bitspur/rock8s/integration-operator/config"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

var startedApparatusTimers map[string]*time.Timer = map[string]*time.Timer{}

type ApparatusUtil struct {
	client   *kubernetes.Clientset
	ctx      context.Context
	dataUtil *DataUtil
	log      logr.Logger
	varUtil  *VarUtil
}

func NewApparatusUtil(
	ctx context.Context,
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
) ([]byte, error) {
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	min := minify.New()
	min.AddFunc("application/json", minifyJson.Minify)
	url := u.getPlugEndpoint(plug) + "/config"
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	go func() {
		body := `{"version":1}`
		var err error
		if plug != nil {
			body, err = sjson.Set(body, "plug", plug)
			if err != nil {
				errCh <- err
				return
			}

			body, err = sjson.Set(body, "socket", socket)
			if err != nil {
				errCh <- err
				return
			}

			plugData, err := u.dataUtil.GetPlugData(plug)
			if err != nil {
				errCh <- err
				return
			}
			body, err = sjson.Set(body, "plugData", plugData)
			if err != nil {
				errCh <- err
				return
			}

			socketData, err := u.dataUtil.GetSocketData(socket)
			if err != nil {
				errCh <- err
				return
			}
			body, err = sjson.Set(body, "socketData", socketData)
			if err != nil {
				errCh <- err
				return
			}

			if plug.Spec.Vars != nil {
				vars, err := u.varUtil.GetVars(plug.Namespace, plug.Spec.Vars, kubectlUtil, plug, socket)
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
		u.log.Info("getting plug config", "method", "POST", "url", url)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(url)
		if err != nil {
			errCh <- NewApparatusNetError(err, r)
			return
		}
		rCh <- r
	}()
	select {
	case r := <-rCh:
		if r.IsError() {
			return r.Body(), NewApparatusNetError(
				errors.New("config failed with "+strconv.Itoa(r.StatusCode())+" status from POST "+url),
				r,
			)
		}
		return r.Body(), nil
	case err := <-errCh:
		return nil, err
	}
}

func (u *ApparatusUtil) GetSocketConfig(
	socket *integrationv1beta1.Socket,
	plug *integrationv1beta1.Plug,
) ([]byte, error) {
	client := resty.New()
	rCh := make(chan *resty.Response)
	errCh := make(chan error)
	min := minify.New()
	min.AddFunc("application/json", minifyJson.Minify)
	url := u.getSocketEndpoint(socket) + "/config"
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	go func() {
		body := `{"version":1}`
		var err error
		if socket != nil {
			body, err = sjson.Set(body, "socket", socket)
			if err != nil {
				errCh <- err
				return
			}

			body, err = sjson.Set(body, "plug", plug)
			if err != nil {
				errCh <- err
				return
			}

			socketData, err := u.dataUtil.GetSocketData(socket)
			if err != nil {
				errCh <- err
				return
			}
			body, err = sjson.Set(body, "socketData", socketData)
			if err != nil {
				errCh <- err
				return
			}

			plugData, err := u.dataUtil.GetPlugData(plug)
			if err != nil {
				errCh <- err
				return
			}
			body, err = sjson.Set(body, "plugData", plugData)
			if err != nil {
				errCh <- err
				return
			}

			if socket.Spec.Vars != nil {
				vars, err := u.varUtil.GetVars(socket.Namespace, socket.Spec.Vars, kubectlUtil, plug, socket)
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
		u.log.Info("getting socket config", "method", "POST", "url", url)
		r, err := client.R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(url)
		if err != nil {
			errCh <- NewApparatusNetError(err, r)
			return
		}
		rCh <- r
	}()
	select {
	case r := <-rCh:
		if r.IsError() {
			return r.Body(), NewApparatusNetError(
				errors.New("config failed with "+strconv.Itoa(r.StatusCode())+" status from POST "+url),
				r,
			)
		}
		return r.Body(), nil
	case err := <-errCh:
		return nil, err
	}
}

func (u *ApparatusUtil) PlugCreated(plug *integrationv1beta1.Plug) error {
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
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
	plug *integrationv1beta1.Plug,
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

func (u *ApparatusUtil) SocketCreated(socket *integrationv1beta1.Socket) error {
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
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
	socket *integrationv1beta1.Socket,
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

func (u *ApparatusUtil) NotRunning(err error) bool {
	if netErr, ok := err.(ApparatusNetError); ok {
		return netErr.NotRunning()
	}
	return false
}

func (u *ApparatusUtil) RenewIdleTimeout(
	apparatus *integrationv1beta1.SpecApparatus,
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
			u.ctx,
			name,
			metav1.GetOptions{},
		); err != nil {
			if k8serrors.IsNotFound(err) {
				return
			}
		}
		startedApparatusTimers[uid] = time.AfterFunc(idleTimeout, func() {
			if err := u.client.CoreV1().Pods(namespace).Delete(
				u.ctx,
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

func (u *ApparatusUtil) StartFromPlug(plug *integrationv1beta1.Plug) (bool, error) {
	return u.start(
		plug.Spec.Apparatus,
		plug.Name+"-apparatus",
		plug.Namespace,
		string(plug.UID),
		u.createPlugOwnerReference(plug),
		EnsureServiceAccount(plug.Spec.ServiceAccountName),
	)
}

func (u *ApparatusUtil) StartFromSocket(socket *integrationv1beta1.Socket) (bool, error) {
	return u.start(
		socket.Spec.Apparatus,
		socket.Name+"-apparatus",
		socket.Namespace,
		string(socket.UID),
		u.createSocketOwnerReference(socket),
		EnsureServiceAccount(socket.Spec.ServiceAccountName),
	)
}

func (u *ApparatusUtil) start(
	apparatus *integrationv1beta1.SpecApparatus,
	name string,
	namespace string,
	uid string,
	ownerReference metav1.OwnerReference,
	serviceAccountName string,
) (bool, error) {
	if apparatus == nil || len(*apparatus.Containers) <= 0 {
		return false, nil
	}
	idleTimeout := time.Second * 60
	if apparatus.IdleTimeout != 0 {
		idleTimeout = time.Second * time.Duration(apparatus.IdleTimeout)
	}
	alreadyExists := false
	automountServiceAccountToken := true
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			OwnerReferences: []metav1.OwnerReference{
				ownerReference,
			},
			Labels: map[string]string{
				"apparatus": name,
			},
			Annotations: map[string]string{
				"sidecar.istio.io/inject": "false",
			},
		},
		Spec: v1.PodSpec{
			AutomountServiceAccountToken: &automountServiceAccountToken,
			ServiceAccountName:           serviceAccountName,
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
		u.ctx,
		service,
		metav1.CreateOptions{
			FieldManager: "integration-operator",
		},
	)
	if err != nil {
		if !k8serrors.IsAlreadyExists(err) {
			return false, err
		}
	}
	_, err = u.client.CoreV1().Pods(namespace).Create(
		u.ctx,
		pod,
		metav1.CreateOptions{
			FieldManager: "integration-operator",
		},
	)
	if err != nil {
		if k8serrors.IsAlreadyExists(err) {
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
				u.ctx,
				name,
				metav1.DeleteOptions{},
			); err != nil {
				u.log.Error(
					err,
					"failed to terminate idle apparatus "+namespace+"/"+name,
				)
			} else {
				u.log.Info("terminated idle apparatus " + namespace + "/" + name)
			}
		})
	}
	if !alreadyExists {
		u.log.Info("started apparatus " + namespace + "/" + name)
	}
	return true, nil
}

func (u *ApparatusUtil) createPlugOwnerReference(plug *integrationv1beta1.Plug) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: plug.APIVersion,
		Kind:       plug.Kind,
		Name:       plug.Name,
		UID:        plug.UID,
	}
}

func (u *ApparatusUtil) createSocketOwnerReference(socket *integrationv1beta1.Socket) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion: socket.APIVersion,
		Kind:       socket.Kind,
		Name:       socket.Name,
		UID:        socket.UID,
	}
}

func (u *ApparatusUtil) getPlugEndpoint(plug *integrationv1beta1.Plug) string {
	return u.getEndpoint(plug.Name+"-apparatus", plug.Namespace, plug.Spec.Apparatus, config.DebugPlugEndpoint)
}

func (u *ApparatusUtil) getSocketEndpoint(socket *integrationv1beta1.Socket) string {
	return u.getEndpoint(socket.Name+"-apparatus", socket.Namespace, socket.Spec.Apparatus, config.DebugSocketEndpoint)
}

func (u *ApparatusUtil) getEndpoint(
	name string,
	namespace string,
	apparatus *integrationv1beta1.SpecApparatus,
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	apparatus *integrationv1beta1.SpecApparatus,
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
	url := endpoint + "/" + eventName
	go func() {
		u.log.Info("triggered event "+eventName, "method", "POST", "url", url)
		r, err := client.SetRetryCount(3).R().EnableTrace().SetHeaders(map[string]string{
			"Content-Type": "application/json",
		}).SetBody([]byte(body)).Post(url)
		if err != nil {
			errCh <- NewApparatusNetError(err, r)
		}
		rCh <- r
	}()
	select {
	case r := <-rCh:
		if r.IsError() {
			return NewApparatusNetError(
				errors.New("event "+eventName+" failed with "+strconv.Itoa(r.StatusCode())+" status from POST "+url),
				r,
			)
		}
		return nil
	case err := <-errCh:
		return err
	}
}

type ApparatusNetError struct {
	err      error
	response *resty.Response
}

func NewApparatusNetError(err error, response *resty.Response) ApparatusNetError {
	return ApparatusNetError{
		err:      err,
		response: response,
	}
}

func (e ApparatusNetError) Error() string {
	return e.err.Error()
}

func (e ApparatusNetError) Timeout() bool {
	netErr, ok := e.err.(net.Error)
	if ok {
		return netErr.Timeout()
	}
	return false
}

func (e ApparatusNetError) NotRunning() bool {
	if e.response != nil && (e.response.StatusCode() > 0 || e.Timeout()) {
		return false
	}
	return true
}
