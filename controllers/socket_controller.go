/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/coupler"
	"github.com/silicon-hills/integration-operator/util"
)

// SocketReconciler reconciles a Socket object
type SocketReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=integration.siliconhills.dev,resources=sockets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=integration.siliconhills.dev,resources=sockets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=integration.siliconhills.dev,resources=sockets/finalizers,verbs=update

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
	log.Info("RECONCILING SOCKET")
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
		return socketUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketCreatedStatusCondition, nil)
	}

	socketInterfaceUtil := util.NewInterfaceUtil(&r.Client, &ctx, &req, &log, &socket.Spec.Interface)
	_, err = socketInterfaceUtil.Get()
	if err != nil {
		return socketUtil.Error(err)
	}

	for _, connectedPlug := range socket.Status.CoupledPlugs {
		plugUtil := util.NewPlugUtil(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
			Name:      connectedPlug.Name,
			Namespace: connectedPlug.Namespace,
		}, util.GlobalPlugMutex)
		plug, err := plugUtil.Get()
		if err != nil {
			return socketUtil.Error(err)
		}
		if _, err := coupler.GlobalCoupler.Couple(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
			Name:      plug.Name,
			Namespace: plug.Namespace,
		}); err != nil {
			return socketUtil.Error(err)
		}
	}
	return socketUtil.UpdateStatusSimple(integrationv1alpha2.ReadyPhase, util.SocketCoupledStatusCondition, nil)
}

func filterSocketPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return !e.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SocketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&integrationv1alpha2.Socket{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		WithEventFilter(filterSocketPredicate()).
		Complete(r)
}
