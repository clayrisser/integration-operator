/**
 * File: /util/socket.go
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
	"fmt"
	"strings"
	"sync"
	"time"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SocketUtil struct {
	apparatusUtil  *ApparatusUtil
	client         *client.Client
	ctx            context.Context
	namespacedName types.NamespacedName
	req            *ctrl.Request
}

func NewSocketUtil(
	client *client.Client,
	ctx context.Context,
	req *ctrl.Request,
	namespacedName *integrationv1beta1.NamespacedName,
) *SocketUtil {
	operatorNamespace := GetOperatorNamespace()
	return &SocketUtil{
		apparatusUtil:  NewApparatusUtil(ctx),
		client:         client,
		ctx:            ctx,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
	}
}

func (u *SocketUtil) Get() (*integrationv1beta1.Socket, error) {
	client := *u.client
	ctx := u.ctx
	socket := &integrationv1beta1.Socket{}
	if err := client.Get(ctx, u.namespacedName, socket); err != nil {
		return nil, err
	}
	return socket.DeepCopy(), nil
}

func (u *SocketUtil) Update(socket *integrationv1beta1.Socket, requeue bool) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Update(ctx, socket); err != nil {
		return u.Error(err, socket)
	}
	return ctrl.Result{Requeue: requeue}, nil
}

func (u *SocketUtil) UpdateStatus(
	socket *integrationv1beta1.Socket,
	requeue bool,
) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Status().Update(ctx, socket); err != nil {
		if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	return ctrl.Result{Requeue: requeue}, nil
}

func (u *SocketUtil) Delete(socket *integrationv1beta1.Socket) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Delete(ctx, socket); err != nil {
		return u.Error(err, socket)
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) GetCoupledCondition(
	socket *integrationv1beta1.Socket,
) (*metav1.Condition, error) {
	if socket == nil {
		var err error
		socket, err = u.Get()
		if err != nil {
			return nil, err
		}
	}
	coupledCondition := meta.FindStatusCondition(socket.Status.Conditions, string(ConditionTypeCoupled))
	return coupledCondition, nil
}

func (u *SocketUtil) CoupledPlugExists(coupledPlugs []*integrationv1beta1.CoupledPlug, plugUid types.UID) bool {
	for _, coupledPlug := range coupledPlugs {
		if coupledPlug.UID == plugUid {
			return true
		}
	}
	return false
}

func (u *SocketUtil) Error(err error, socket *integrationv1beta1.Socket) (ctrl.Result, error) {
	e := err
	if socket == nil {
		var err error
		socket, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if u.apparatusUtil.NotRunning(err) {
		requeueAfter := time.Duration(time.Second.Nanoseconds() * 10)
		started, err := u.apparatusUtil.StartFromSocket(socket, &requeueAfter)
		if err != nil {
			return u.UpdateErrorStatus(err, socket)
		}
		if started {
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: requeueAfter,
			}, nil
		}
	}
	result, err := u.UpdateErrorStatus(e, socket)
	if strings.Contains(e.Error(), "result property") &&
		strings.Contains(e.Error(), "is required") {
		return ctrl.Result{Requeue: true}, nil
	}
	return result, err
}

func (u *SocketUtil) UpdateCoupledStatus(
	conditionCoupledReason ConditionCoupledReason,
	socket *integrationv1beta1.Socket,
	appendPlug *integrationv1beta1.Plug,
	requeue bool,
) (ctrl.Result, error) {
	if socket == nil {
		var err error
		socket, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if appendPlug != nil {
		if err := u.appendCoupledPlugStatus(socket, appendPlug); err != nil {
			return u.Error(err, socket)
		}
	}
	if conditionCoupledReason != "" {
		u.setCoupledStatusCondition(conditionCoupledReason, "", socket)
	}
	return u.UpdateStatus(socket, requeue)
}

func (u *SocketUtil) UpdateErrorStatus(
	err error,
	socket *integrationv1beta1.Socket,
) (ctrl.Result, error) {
	e := err
	if socket == nil {
		var err error
		socket, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if err = u.setErrorStatus(e, socket); err != nil {
		return ctrl.Result{}, err
	}
	if _, err := u.UpdateStatus(socket, true); err != nil {
		return ctrl.Result{}, err
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, e
}

func (u *SocketUtil) UpdateRemoveCoupledPlugStatus(
	plugUid types.UID,
	socket *integrationv1beta1.Socket,
	requeue bool,
) (ctrl.Result, error) {
	if socket == nil {
		var err error
		socket, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	removedCoupledPlug, err := u.RemoveCoupledPlugStatus(plugUid, socket)
	if err != nil {
		return u.Error(err, socket)
	}
	if removedCoupledPlug {
		return u.UpdateStatus(socket, requeue)
	}
	return ctrl.Result{Requeue: requeue}, nil
}

func (u *SocketUtil) UpdateAppendCoupledPlugStatus(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	requeue bool,
) (ctrl.Result, error) {
	if socket == nil {
		var err error
		socket, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if err := u.appendCoupledPlugStatus(socket, plug); err != nil {
		return u.Error(err, socket)
	}
	return u.UpdateStatus(socket, requeue)
}

func (u *SocketUtil) RemoveCoupledPlugStatus(
	plugUid types.UID,
	socket *integrationv1beta1.Socket,
) (bool, error) {
	if socket == nil {
		var err error
		socket, err = u.Get()
		if err != nil {
			return false, err
		}
	}
	coupledPlugs := []*integrationv1beta1.CoupledPlug{}
	removedCoupledPlug := false
	for _, coupledPlug := range socket.Status.CoupledPlugs {
		if coupledPlug.UID == plugUid {
			removedCoupledPlug = true
		} else {
			coupledPlugs = append(coupledPlugs, coupledPlug)
		}
	}
	socket.Status.CoupledPlugs = coupledPlugs
	coupledCondition, err := u.GetCoupledCondition(socket)
	if err != nil {
		return false, err
	}
	if (*coupledCondition).Reason == string(SocketCoupled) {
		u.setCoupledStatusCondition(SocketCoupled, "", socket)
	}
	return removedCoupledPlug, nil
}

func (u *SocketUtil) appendCoupledPlugStatus(
	socket *integrationv1beta1.Socket,
	plug *integrationv1beta1.Plug,
) error {
	if !u.CoupledPlugExists(socket.Status.CoupledPlugs, plug.UID) {
		socket.Status.CoupledPlugs = append(socket.Status.CoupledPlugs, &integrationv1beta1.CoupledPlug{
			APIVersion: plug.APIVersion,
			Kind:       plug.Kind,
			Name:       plug.Name,
			Namespace:  plug.Namespace,
			UID:        plug.UID,
		})
	}
	u.setCoupledStatusCondition(SocketCoupled, "", socket)
	return nil
}

func (u *SocketUtil) setCoupledStatusCondition(
	conditionCoupledReason ConditionCoupledReason,
	message string,
	socket *integrationv1beta1.Socket,
) {
	coupledStatus := false
	coupledPlugsCount := len(socket.Status.CoupledPlugs)
	if message == "" {
		if conditionCoupledReason == SocketCreated {
			message = "socket created"
		} else if conditionCoupledReason == Error {
			message = "unknown error"
		} else if conditionCoupledReason == SocketCoupled {
			message = fmt.Sprint(coupledPlugsCount)
			if coupledPlugsCount == 1 {
				message += " plug coupled"
			} else {
				message += " plugs coupled"
			}
		} else if conditionCoupledReason == SocketEmpty {
			message = "0 plugs coupled"
		}
	}
	if conditionCoupledReason != Error {
		socket.Status.Conditions = []metav1.Condition{}
	}
	if conditionCoupledReason == SocketCoupled {
		if coupledPlugsCount > 0 {
			coupledStatus = true
		} else {
			conditionCoupledReason = SocketEmpty
		}
	}
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: socket.Generation,
		Reason:             string(conditionCoupledReason),
		Status:             "False",
		Type:               string(ConditionTypeCoupled),
	}
	if coupledStatus {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&socket.Status.Conditions, condition)
}

func (u *SocketUtil) setErrorStatus(err error, socket *integrationv1beta1.Socket) error {
	e := err
	if e == nil {
		return nil
	}
	if socket == nil {
		return nil
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return nil
	}
	message := e.Error()
	coupledCondition, err := u.GetCoupledCondition(socket)
	if err != nil {
		return err
	}
	if coupledCondition != nil {
		u.setCoupledStatusCondition(Error, "coupling failed", socket)
	}
	failedCondition := metav1.Condition{
		Message:            message,
		ObservedGeneration: socket.Generation,
		Reason:             "Error",
		Status:             "True",
		Type:               string(ConditionTypeFailed),
	}
	meta.SetStatusCondition(&socket.Status.Conditions, failedCondition)
	return nil
}

var GlobalSocketMutex *sync.Mutex = &sync.Mutex{}
