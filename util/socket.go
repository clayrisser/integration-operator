/**
 * File: /util/socket.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 01-07-2021 16:37:38
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
	"fmt"
	"strings"
	"sync"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/config"
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
	ctx            *context.Context
	log            *logr.Logger
	mutex          *sync.Mutex
	namespacedName types.NamespacedName
	req            *ctrl.Request
}

func NewSocketUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	namespacedName *integrationv1alpha2.NamespacedName,
	mutex *sync.Mutex,
) *SocketUtil {
	operatorNamespace := GetOperatorNamespace()
	if mutex == nil {
		mutex = &sync.Mutex{}
	}
	return &SocketUtil{
		apparatusUtil:  NewApparatusUtil(ctx),
		client:         client,
		ctx:            ctx,
		log:            log,
		mutex:          mutex,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
	}
}

func (u *SocketUtil) Get() (*integrationv1alpha2.Socket, error) {
	client := *u.client
	ctx := *u.ctx
	socket := &integrationv1alpha2.Socket{}
	if err := client.Get(ctx, u.namespacedName, socket); err != nil {
		return nil, err
	}
	return socket.DeepCopy(), nil
}

func (u *SocketUtil) Update(socket *integrationv1alpha2.Socket) error {
	client := *u.client
	ctx := *u.ctx
	u.mutex.Lock()
	if err := client.Update(ctx, socket); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *SocketUtil) UpdateStatus(
	socket *integrationv1alpha2.Socket,
	requeue bool,
	exponentialBackoff bool,
) error {
	client := *u.client
	ctx := *u.ctx
	if !exponentialBackoff ||
		socket.Status.LastUpdate.IsZero() ||
		config.StartTime.Unix() > socket.Status.LastUpdate.Unix() {
		socket.Status.LastUpdate = metav1.Now()
	}
	socket.Status.Requeued = requeue
	u.mutex.Lock()
	if err := client.Status().Update(ctx, socket); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *SocketUtil) CoupledPlugExists(coupledPlugs []*integrationv1alpha2.CoupledPlug, plugUid types.UID) bool {
	coupledPlugExits := false
	for _, coupledPlug := range coupledPlugs {
		if coupledPlug.UID == plugUid {
			coupledPlugExits = true
		}
	}
	return coupledPlugExits
}

func (u *SocketUtil) GetCoupledCondition() (*metav1.Condition, error) {
	socket, err := u.Get()
	if err != nil {
		return nil, err
	}
	coupledCondition := meta.FindStatusCondition(socket.Status.Conditions, "Coupled")
	return coupledCondition, nil
}

func (u *SocketUtil) Error(err error) (ctrl.Result, error) {
	stashedErr := err
	log := *u.log
	socket, err := u.Get()
	if err != nil {
		log.Error(nil, stashedErr.Error())
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: config.MaxRequeueDuration,
		}, err
	}
	requeueAfter := CalculateExponentialRequireAfter(
		socket.Status.LastUpdate,
		2,
	)
	if u.apparatusUtil.NotRunning(stashedErr) {
		started, err := u.apparatusUtil.StartFromSocket(socket)
		if err != nil {
			return u.Error(err)
		}
		if started {
			if err := u.UpdateStatus(socket, true, true); err != nil {
				if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
					return ctrl.Result{
						Requeue:      true,
						RequeueAfter: requeueAfter,
					}, nil
				}
				return u.Error(err)
			}
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: requeueAfter,
			}, nil
		}
	}
	log.Error(nil, stashedErr.Error())
	if strings.Index(stashedErr.Error(), registry.OptimisticLockErrorMsg) <= -1 {
		if _, err := u.UpdateErrorStatus(stashedErr, true); err != nil {
			if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
				return ctrl.Result{
					Requeue:      true,
					RequeueAfter: requeueAfter,
				}, nil
			}
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: requeueAfter,
			}, err
		}
	}
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: requeueAfter,
	}, nil
}

func (u *SocketUtil) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	coupledStatusCondition StatusCondition,
	appendPlug *integrationv1alpha2.Plug,
	requeue bool,
) (ctrl.Result, error) {
	socket, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if phase != "" {
		u.setPhaseStatus(socket, phase)
	}
	if appendPlug != nil {
		if err := u.appendCoupledPlugStatus(socket, appendPlug); err != nil {
			return u.Error(err)
		}
	}
	if coupledStatusCondition != "" {
		u.setCoupledStatusCondition(socket, coupledStatusCondition, "")
	}
	if err := u.UpdateStatus(socket, requeue, false); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) UpdateErrorStatus(err error, requeue bool) (ctrl.Result, error) {
	stashedErr := err
	socket, err := u.Get()
	if err != nil {
		return ctrl.Result{}, err
	}
	u.setErrorStatus(socket, stashedErr)
	if err := u.UpdateStatus(socket, requeue, true); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) UpdateStatusRemovePlug(plugUid types.UID, requeue bool) (ctrl.Result, error) {
	socket, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if err := u.removeCoupledPlugStatus(socket, plugUid); err != nil {
		return u.Error(err)
	}
	if err := u.UpdateStatus(socket, requeue, false); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) UpdateStatusAppendPlug(
	plug *integrationv1alpha2.Plug,
	requeue bool,
) (ctrl.Result, error) {
	socket, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if err := u.appendCoupledPlugStatus(socket, plug); err != nil {
		return u.Error(err)
	}
	if err := u.UpdateStatus(socket, requeue, false); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *SocketUtil) appendCoupledPlugStatus(
	socket *integrationv1alpha2.Socket,
	plug *integrationv1alpha2.Plug,
) error {
	if !u.CoupledPlugExists(socket.Status.CoupledPlugs, plug.UID) {
		socket.Status.CoupledPlugs = append(socket.Status.CoupledPlugs, &integrationv1alpha2.CoupledPlug{
			APIVersion: plug.APIVersion,
			Kind:       plug.Kind,
			Name:       plug.Name,
			Namespace:  plug.Namespace,
			UID:        plug.UID,
		})
	}
	u.setCoupledStatusCondition(socket, SocketCoupledStatusCondition, "")
	return nil
}

func (u *SocketUtil) setCoupledStatusCondition(
	socket *integrationv1alpha2.Socket,
	coupledStatusCondition StatusCondition,
	message string,
) {
	u.setReadyStatus(socket, false)
	coupledStatus := false
	coupledPlugsCount := len(socket.Status.CoupledPlugs)
	if message == "" {
		if coupledStatusCondition == SocketCreatedStatusCondition {
			message = "socket created"
		} else if coupledStatusCondition == ErrorStatusCondition {
			message = "unknown error"
		} else if coupledStatusCondition == SocketCoupledStatusCondition {
			message = "socket ready with " + fmt.Sprint(coupledPlugsCount) + " plugs coupled"
		} else if coupledStatusCondition == SocketEmptyStatusCondition {
			message = "socket ready with 0 plugs coupled"
		}
	}
	if coupledStatusCondition == SocketCoupledStatusCondition {
		if coupledPlugsCount > 0 {
			coupledStatus = true
		} else {
			coupledStatusCondition = SocketEmptyStatusCondition
		}
	}
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: socket.Generation,
		Reason:             string(coupledStatusCondition),
		Status:             "False",
		Type:               "Coupled",
	}
	if coupledStatus {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&socket.Status.Conditions, condition)
	if coupledStatusCondition == SocketCoupledStatusCondition || coupledStatusCondition == SocketEmptyStatusCondition {
		u.setReadyStatus(socket, true)
	}
}

func (u *SocketUtil) setPhaseStatus(
	socket *integrationv1alpha2.Socket,
	phase integrationv1alpha2.Phase,
) {
	if phase != integrationv1alpha2.FailedPhase {
		socket.Status.Message = ""
	}
	socket.Status.Phase = phase
}

func (u *SocketUtil) setReadyStatus(socket *integrationv1alpha2.Socket, ready bool) {
	socket.Status.Ready = ready
}

func (u *SocketUtil) setErrorStatus(socket *integrationv1alpha2.Socket, err error) {
	message := err.Error()
	coupledCondition, err := u.GetCoupledCondition()
	if err == nil {
		coupledCondition = nil
	}
	if coupledCondition != nil {
		u.setCoupledStatusCondition(socket, ErrorStatusCondition, message)
	}
	socket.Status.Phase = integrationv1alpha2.FailedPhase
	socket.Status.Message = message
}

func (u *SocketUtil) removeCoupledPlugStatus(
	socket *integrationv1alpha2.Socket,
	plugUid types.UID,
) error {
	coupledPlugs := []*integrationv1alpha2.CoupledPlug{}
	for _, coupledPlug := range socket.Status.CoupledPlugs {
		if coupledPlug.UID != plugUid {
			coupledPlugs = append(coupledPlugs, coupledPlug)
		}
	}
	socket.Status.CoupledPlugs = coupledPlugs
	coupledCondition, err := u.GetCoupledCondition()
	if err != nil {
		return err
	}
	if (*coupledCondition).Reason == string(SocketCoupledStatusCondition) {
		u.setCoupledStatusCondition(socket, SocketCoupledStatusCondition, "")
	}
	return nil
}

var GlobalSocketMutex *sync.Mutex = &sync.Mutex{}
