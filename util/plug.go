/*
 * File: /util/plug.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 30-06-2021 13:40:14
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

type PlugUtil struct {
	apparatusUtil  *ApparatusUtil
	client         *client.Client
	ctx            *context.Context
	log            *logr.Logger
	mutex          *sync.Mutex
	namespacedName types.NamespacedName
	req            *ctrl.Request
}

func NewPlugUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	namespacedName *integrationv1alpha2.NamespacedName,
	mutex *sync.Mutex,
) *PlugUtil {
	operatorNamespace := GetOperatorNamespace()
	if mutex == nil {
		mutex = &sync.Mutex{}
	}
	return &PlugUtil{
		apparatusUtil:  NewApparatusUtil(ctx),
		client:         client,
		ctx:            ctx,
		log:            log,
		mutex:          mutex,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
	}
}

func (u *PlugUtil) Get() (*integrationv1alpha2.Plug, error) {
	client := *u.client
	ctx := *u.ctx
	plug := &integrationv1alpha2.Plug{}
	if err := client.Get(ctx, u.namespacedName, plug); err != nil {
		return nil, err
	}
	return plug.DeepCopy(), nil
}

func (u *PlugUtil) Update(plug *integrationv1alpha2.Plug) error {
	client := *u.client
	ctx := *u.ctx
	u.mutex.Lock()
	if err := client.Update(ctx, plug); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *PlugUtil) UpdateStatus(
	plug *integrationv1alpha2.Plug,
	requeue bool,
	exponentialBackoff bool,
) error {
	client := *u.client
	ctx := *u.ctx
	if !exponentialBackoff ||
		plug.Status.LastUpdate.IsZero() ||
		config.StartTime.Unix() > plug.Status.LastUpdate.Unix() {
		plug.Status.LastUpdate = metav1.Now()
	}
	plug.Status.Requeued = requeue
	u.mutex.Lock()
	if err := client.Status().Update(ctx, plug); err != nil {
		u.mutex.Unlock()
		return err
	}
	u.mutex.Unlock()
	return nil
}

func (u *PlugUtil) GetCoupledCondition() (*metav1.Condition, error) {
	plug, err := u.Get()
	if err != nil {
		return nil, err
	}
	coupledCondition := meta.FindStatusCondition(plug.Status.Conditions, "coupled")
	return coupledCondition, nil
}

func (u *PlugUtil) IsCoupled(
	plug *integrationv1alpha2.Plug,
	coupledCondition *metav1.Condition,
) (bool, error) {
	if coupledCondition == nil {
		var err error
		coupledCondition, err = u.GetCoupledCondition()
		if err != nil {
			return false, err

		}
	}
	return coupledCondition != nil && plug.Status.Phase == integrationv1alpha2.SucceededPhase &&
		coupledCondition.Reason == string(CouplingSucceededStatusCondition), nil
}

func (u *PlugUtil) Error(err error) (ctrl.Result, error) {
	stashedErr := err
	log := *u.log
	plug, err := u.Get()
	if err != nil {
		log.Error(nil, stashedErr.Error())
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: config.MaxRequeueDuration,
		}, err
	}
	requeueAfter := CalculateExponentialRequireAfter(
		plug.Status.LastUpdate,
		2,
	)
	if u.apparatusUtil.NotRunning(stashedErr) {
		started, err := u.apparatusUtil.Start(
			plug.Spec.Apparatus,
			plug.Name,
			plug.Namespace,
			string(plug.UID),
		)
		if err != nil {
			return u.Error(err)
		}
		if started {
			if err := u.UpdateStatus(plug, true, true); err != nil {
				if strings.Index(err.Error(), registry.OptimisticLockErrorMsg) > -1 {
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
			if strings.Index(err.Error(), registry.OptimisticLockErrorMsg) > -1 {
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

func (u *PlugUtil) UpdateErrorStatus(err error, requeue bool) (ctrl.Result, error) {
	stashedErr := err
	plug, err := u.Get()
	if err != nil {
		return ctrl.Result{}, err
	}
	u.setErrorStatus(plug, stashedErr)
	if err := u.UpdateStatus(plug, requeue, true); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (u *PlugUtil) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	coupledStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
	requeue bool,
) (ctrl.Result, error) {
	plug, err := u.Get()
	if err != nil {
		return u.Error(err)
	}
	if coupledStatusCondition != "" {
		u.setCoupledStatusCondition(plug, coupledStatusCondition, "")
	}
	if socket != nil {
		u.setCoupledSocketStatus(plug, socket)
	}
	if phase != "" {
		u.setPhaseStatus(plug, phase)
	}
	if coupledStatusCondition == SocketNotCreatedStatusCondition ||
		coupledStatusCondition == SocketNotReadyStatusCondition {
		requeueAfter := CalculateExponentialRequireAfter(
			plug.Status.LastUpdate,
			2,
		)
		if err := u.UpdateStatus(plug, requeue, false); err != nil {
			return u.Error(err)
		}
		return ctrl.Result{Requeue: true, RequeueAfter: requeueAfter}, nil
	}
	if err := u.UpdateStatus(plug, requeue, false); err != nil {
		return u.Error(err)
	}
	return ctrl.Result{}, nil
}

func (u *PlugUtil) setPhaseStatus(
	plug *integrationv1alpha2.Plug,
	phase integrationv1alpha2.Phase,
) {
	if phase != integrationv1alpha2.FailedPhase {
		plug.Status.Message = ""
	}
	plug.Status.Phase = phase
}

func (u *PlugUtil) setCoupledStatusCondition(
	plug *integrationv1alpha2.Plug,
	coupledStatusCondition StatusCondition,
	message string,
) {
	coupledStatus := false
	if message == "" {
		if coupledStatusCondition == PlugCreatedStatusCondition {
			message = "plug created"
		} else if coupledStatusCondition == SocketNotCreatedStatusCondition {
			message = "waiting for socket to be created"
		} else if coupledStatusCondition == SocketNotReadyStatusCondition {
			message = "waiting for socket to be ready"
		} else if coupledStatusCondition == CouplingInProcessStatusCondition {
			message = "coupling to socket"
		} else if coupledStatusCondition == CouplingSucceededStatusCondition {
			message = "coupling succeeded"
		} else if coupledStatusCondition == ErrorStatusCondition {
			message = "unknown error"
		}
	}
	if coupledStatusCondition == CouplingSucceededStatusCondition {
		coupledStatus = true
	}
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: plug.Generation,
		Reason:             string(coupledStatusCondition),
		Status:             "False",
		Type:               "coupled",
	}
	if coupledStatus {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&plug.Status.Conditions, condition)
}

func (u *PlugUtil) setErrorStatus(plug *integrationv1alpha2.Plug, err error) {
	message := err.Error()
	coupledCondition, err := u.GetCoupledCondition()
	if err == nil {
		coupledCondition = nil
	}
	if coupledCondition != nil {
		u.setCoupledStatusCondition(plug, ErrorStatusCondition, message)
	}
	plug.Status.Phase = integrationv1alpha2.FailedPhase
	plug.Status.Message = message
}

func (u *PlugUtil) setCoupledSocketStatus(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
) {
	plug.Status.CoupledSocket = &integrationv1alpha2.CoupledSocket{
		APIVersion: socket.APIVersion,
		Kind:       socket.Kind,
		Name:       socket.Name,
		Namespace:  socket.Namespace,
		UID:        socket.UID,
	}
}

var GlobalPlugMutex *sync.Mutex = &sync.Mutex{}
