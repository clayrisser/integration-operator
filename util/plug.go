/**
 * File: /util/plug.go
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
	"strings"
	"sync"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PlugUtil struct {
	apparatusUtil  *ApparatusUtil
	client         *client.Client
	ctx            context.Context
	namespacedName types.NamespacedName
	req            *ctrl.Request
	resultUtil     *ResultUtil
}

func NewPlugUtil(
	client *client.Client,
	ctx context.Context,
	req *ctrl.Request,
	namespacedName *integrationv1beta1.NamespacedName,
) *PlugUtil {
	operatorNamespace := GetOperatorNamespace()
	return &PlugUtil{
		apparatusUtil:  NewApparatusUtil(ctx),
		client:         client,
		ctx:            ctx,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
		resultUtil:     NewResultUtil(ctx),
	}
}

func (u *PlugUtil) Get() (*integrationv1beta1.Plug, error) {
	client := *u.client
	ctx := u.ctx
	plug := &integrationv1beta1.Plug{}
	if err := client.Get(ctx, u.namespacedName, plug); err != nil {
		return nil, err
	}
	return plug.DeepCopy(), nil
}

func (u *PlugUtil) Update(plug *integrationv1beta1.Plug, requeue bool) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Update(ctx, plug); err != nil {
		return u.Error(err, plug)
	}
	if requeue {
		return ctrl.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	return ctrl.Result{}, nil
}

func (u *PlugUtil) UpdateStatus(
	plug *integrationv1beta1.Plug,
	requeue bool,
) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Status().Update(ctx, plug); err != nil {
		if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	if requeue {
		return ctrl.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	return ctrl.Result{}, nil
}

func (u *PlugUtil) Delete(plug *integrationv1beta1.Plug) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Delete(ctx, plug); err != nil {
		return u.Error(err, plug)
	}
	return ctrl.Result{}, nil
}

func (u *PlugUtil) GetCoupledCondition() (*metav1.Condition, error) {
	plug, err := u.Get()
	if err != nil {
		return nil, err
	}
	coupledCondition := meta.FindStatusCondition(plug.Status.Conditions, string(ConditionTypeCoupled))
	return coupledCondition, nil
}

func (u *PlugUtil) Error(
	err error,
	plug *integrationv1beta1.Plug,
) (ctrl.Result, error) {
	e := err
	if plug == nil {
		var err error
		plug, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if u.apparatusUtil.NotRunning(err) {
		started, err := u.apparatusUtil.StartFromPlug(plug)
		if err != nil {
			return ctrl.Result{}, err
		}
		if started {
			return ctrl.Result{Requeue: true}, nil
		}
	}
	result, err := u.UpdateErrorStatus(e, plug)
	if strings.Contains(e.Error(), "result property") &&
		strings.Contains(e.Error(), "is required") {
		return ctrl.Result{Requeue: true}, nil
	}
	return result, err
}

func (u *PlugUtil) UpdateErrorStatus(
	err error,
	plug *integrationv1beta1.Plug,
) (ctrl.Result, error) {
	e := err
	if plug == nil {
		var err error
		plug, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if err = u.setErrorStatus(e, plug); err != nil {
		return ctrl.Result{}, err
	}
	if _, err := u.UpdateStatus(plug, true); err != nil {
		return ctrl.Result{}, err
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, e
}

func (u *PlugUtil) UpdateCoupledStatus(
	conditionCoupledReason ConditionCoupledReason,
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	requeue bool,
) (ctrl.Result, error) {
	if plug == nil {
		var err error
		plug, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if socket != nil {
		u.setCoupledSocketStatus(plug, socket)
	}
	if conditionCoupledReason != "" {
		u.setCoupledStatusCondition(conditionCoupledReason, "", plug)
	}
	return u.UpdateStatus(plug, requeue)
}

func (u *PlugUtil) UpdateResultStatus(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig Config,
	socketConfig Config,
) (ctrl.Result, error) {
	coupledResult, err := u.resultUtil.GetResult(plug, socket, plugConfig, socketConfig)
	if err != nil {
		return u.Error(err, plug)
	}
	if err := u.resultUtil.SocketTemplateResultResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		coupledResult.Plug,
		coupledResult.Socket,
	); err != nil {
		return u.Error(err, plug)
	}
	if err := u.resultUtil.PlugTemplateResultResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		coupledResult.Plug,
		coupledResult.Socket,
	); err != nil {
		return u.Error(err, plug)
	}
	coupledResultStatus := integrationv1beta1.CoupledResultStatus{
		ObservedGeneration: plug.Generation,
	}
	coupledResultStatus.Plug = coupledResult.Plug
	coupledResultStatus.Socket = coupledResult.Socket
	plug.Status.CoupledResult = &coupledResultStatus
	return u.UpdateCoupledStatus(CouplingSucceeded, plug, socket, false)
}

func (u *PlugUtil) setCoupledStatusCondition(
	conditionCoupledReason ConditionCoupledReason,
	message string,
	plug *integrationv1beta1.Plug,
) {
	coupledStatus := false
	if message == "" {
		if conditionCoupledReason == PlugCreated {
			message = "plug created"
		} else if conditionCoupledReason == SocketNotCreated {
			message = "waiting for socket to be created"
		} else if conditionCoupledReason == CouplingInProcess {
			message = "coupling to socket"
		} else if conditionCoupledReason == CouplingSucceeded {
			message = "coupling succeeded"
		} else if conditionCoupledReason == UpdatingInProcess {
			message = "updating coupling"
		} else if conditionCoupledReason == Error {
			message = "unknown error"
		}
	}
	if conditionCoupledReason != Error {
		plug.Status.Conditions = []metav1.Condition{}
	}
	if conditionCoupledReason == CouplingSucceeded {
		coupledStatus = true
	}
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: plug.Generation,
		Reason:             string(conditionCoupledReason),
		Status:             "False",
		Type:               string(ConditionTypeCoupled),
	}
	if coupledStatus {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&plug.Status.Conditions, condition)
}

func (u *PlugUtil) setErrorStatus(err error, plug *integrationv1beta1.Plug) error {
	e := err
	if e == nil {
		return nil
	}
	if plug == nil {
		return nil
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return nil
	}
	message := e.Error()
	coupledCondition, err := u.GetCoupledCondition()
	if err != nil {
		return err
	}
	if coupledCondition != nil {
		u.setCoupledStatusCondition(Error, "coupling failed", plug)
	}
	failedCondition := metav1.Condition{
		Message:            message,
		ObservedGeneration: plug.Generation,
		Reason:             "Error",
		Status:             "True",
		Type:               string(ConditionTypeFailed),
	}
	meta.SetStatusCondition(&plug.Status.Conditions, failedCondition)
	return nil
}

func (u *PlugUtil) setCoupledSocketStatus(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
) {
	plug.Status.CoupledSocket = &integrationv1beta1.CoupledSocket{
		APIVersion: socket.APIVersion,
		Kind:       socket.Kind,
		Name:       socket.Name,
		Namespace:  socket.Namespace,
		UID:        socket.UID,
	}
}

var GlobalPlugMutex *sync.Mutex = &sync.Mutex{}
