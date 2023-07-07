/**
 * File: /socket_controller.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 07-07-2023 08:25:17
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

package controllers

import (
	"context"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	integrationv1alpha2 "gitlab.com/bitspur/rock8s/integration-operator/api/v1alpha2"
	"gitlab.com/bitspur/rock8s/integration-operator/coupler"
	"gitlab.com/bitspur/rock8s/integration-operator/util"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// SocketReconciler reconciles a Socket object
type SocketReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=integration.rock8s.com,resources=sockets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=sockets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=sockets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Socket object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *SocketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("socket", req.NamespacedName)
	log.Info("R socket " + req.NamespacedName.String())
	socketUtil := util.NewSocketUtil(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
		Name:      req.NamespacedName.Name,
		Namespace: req.NamespacedName.Namespace,
	}, util.GlobalSocketMutex)
	socket, err := socketUtil.Get()
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if socket.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(socket, integrationv1alpha2.PlugFinalizer) {
			for _, connectedPlug := range socket.Status.CoupledPlugs {
				plug := &integrationv1alpha2.Plug{}
				err := r.Get(ctx, types.NamespacedName{
					Name:      connectedPlug.Name,
					Namespace: connectedPlug.Namespace,
				}, plug)
				if err != nil {
					if errors.IsNotFound(err) {
						continue
					}
					return socketUtil.Error(err)
				}
				if err := r.Delete(ctx, plug); err != nil {
					return socketUtil.Error(err)
				}
			}
			coupler.GlobalCoupler.DeletedSocket(socket)
			controllerutil.RemoveFinalizer(socket, integrationv1alpha2.PlugFinalizer)
			if err := r.Update(ctx, socket); err != nil {
				return socketUtil.Error(err)
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(socket, integrationv1alpha2.SocketFinalizer) {
		controllerutil.AddFinalizer(socket, integrationv1alpha2.SocketFinalizer)
		if err := r.Update(ctx, socket); err != nil {
			return socketUtil.Error(err)
		}
		return ctrl.Result{}, nil
	}

	coupledCondition, err := socketUtil.GetCoupledCondition()
	if err != nil {
		return socketUtil.Error(err)
	}
	if coupledCondition == nil {
		if err = coupler.GlobalCoupler.CreatedSocket(socket); err != nil {
			return socketUtil.Error(err)
		}
		return socketUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketCreatedStatusCondition, nil, true)
	}

	socketInterfaceUtil := util.NewInterfaceUtil(&r.Client, &ctx, &socket.Spec.Interface)
	socketInterface, err := socketInterfaceUtil.Get()
	if err != nil {
		return socketUtil.Error(err)
	}

	var requeueAfter time.Duration = 0
	for _, connectedPlug := range socket.Status.CoupledPlugs {
		plugUtil := util.NewPlugUtil(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
			Name:      connectedPlug.Name,
			Namespace: connectedPlug.Namespace,
		}, util.GlobalPlugMutex)
		plug, err := plugUtil.Get()
		if err != nil {
			if errors.IsNotFound(err) {
				if _, err := socketUtil.UpdateStatusRemovePlug(connectedPlug.UID, true); err != nil {
					return socketUtil.Error(err)
				}
				return ctrl.Result{Requeue: true}, nil
			}
			return socketUtil.Error(err)
		}

		if err != nil {
			return socketUtil.Error(err)
		}
		result, err := coupler.GlobalCoupler.Update(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
			Name:      plug.Name,
			Namespace: plug.Namespace,
		}, plug, socket, socketInterface)
		if err != nil {
			return socketUtil.Error(err)
		}
		if result.Requeue {
			requeueAfter = time.Duration(math.Min(
				float64(requeueAfter.Nanoseconds()),
				float64(result.RequeueAfter.Nanoseconds()),
			))
		}
	}

	result, err := socketUtil.UpdateStatusSimple(integrationv1alpha2.ReadyPhase, util.SocketCoupledStatusCondition, nil, false)
	if requeueAfter > 0 {
		result.Requeue = true
		result.RequeueAfter = requeueAfter
	}
	return result, err
}

func filterSocketPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			newSocket, ok := e.ObjectNew.(*integrationv1alpha2.Socket)
			if !ok {
				return false
			}
			if newSocket.Status.LastUpdate.IsZero() {
				return true
			}
			return e.ObjectNew.GetGeneration() > e.ObjectOld.GetGeneration() || newSocket.Status.Requeued
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return !e.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SocketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	maxConcurrentReconciles := 3
	if value := os.Getenv("MAX_CONCURRENT_RECONCILES"); value != "" {
		if val, err := strconv.Atoi(value); err == nil {
			maxConcurrentReconciles = val
		}
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&integrationv1alpha2.Socket{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: maxConcurrentReconciles}).
		WithEventFilter(filterSocketPredicate()).
		Complete(r)
}
