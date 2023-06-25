/**
 * File: /plug.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 25-06-2023 10:15:43
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Risser Labs LLC (c) Copyright 2021
 */

package util

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
	integrationv1alpha2 "gitlab.com/bitspur/internal/integration-operator/api/v1alpha2"
	"gitlab.com/bitspur/internal/integration-operator/config"
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

func (u *PlugUtil) SocketError(err error) error {
	if strings.Index(err.Error(), registry.OptimisticLockErrorMsg) <= -1 {
		if _, _err := u.UpdateErrorStatus(err, true); _err != nil {
			if strings.Contains(_err.Error(), registry.OptimisticLockErrorMsg) {
				return nil
			}
			return _err
		}
	}
	return nil
}

func (u *PlugUtil) Error(err error) (ctrl.Result, error) {
	log := *u.log
	plug, _err := u.Get()
	if _err != nil {
		log.Error(nil, err.Error())
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: config.MaxRequeueDuration,
		}, _err
	}
	requeueAfter := CalculateExponentialRequireAfter(
		plug.Status.LastUpdate,
		2,
	)
	if u.apparatusUtil.NotRunning(err) {
		successRequeueAfter := time.Duration(time.Second.Nanoseconds() * 10)
		started, _err := u.apparatusUtil.StartFromPlug(plug, &successRequeueAfter)
		if _err != nil {
			return u.Error(_err)
		}
		if started {
			if _err := u.UpdateStatus(plug, true, true); _err != nil {
				if strings.Contains(_err.Error(), registry.OptimisticLockErrorMsg) {
					return ctrl.Result{
						Requeue:      true,
						RequeueAfter: requeueAfter,
					}, nil
				}
				return u.Error(_err)
			}
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: successRequeueAfter,
			}, nil
		}
	}
	log.Error(nil, err.Error())
	if strings.Index(err.Error(), registry.OptimisticLockErrorMsg) <= -1 {
		if _, _err := u.UpdateErrorStatus(err, true); _err != nil {
			if strings.Contains(_err.Error(), registry.OptimisticLockErrorMsg) {
				return ctrl.Result{
					Requeue:      true,
					RequeueAfter: requeueAfter,
				}, nil
			}
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: requeueAfter,
			}, _err
		}
	}
	return ctrl.Result{
		Requeue:      true,
		RequeueAfter: requeueAfter,
	}, nil
}

func (u *PlugUtil) UpdateErrorStatus(err error, requeue bool) (ctrl.Result, error) {
	plug, _err := u.Get()
	if _err != nil {
		return ctrl.Result{}, _err
	}
	u.setErrorStatus(plug, err)
	if _err := u.UpdateStatus(plug, requeue, true); _err != nil {
		return ctrl.Result{}, _err
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
	coupledCondition, _err := u.GetCoupledCondition()
	if _err == nil {
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
