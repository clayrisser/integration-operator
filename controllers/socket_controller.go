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
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/coupler"
	"github.com/silicon-hills/integration-operator/services"
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
	_ = r.Log.WithValues("socket", req.NamespacedName)
	s := services.NewServices()
	result := ctrl.Result{}
	socket := &integrationv1alpha2.Socket{}
	err := r.Get(ctx, req.NamespacedName, socket)
	if err != nil {
		if errors.IsNotFound(err) {
			return result, nil
		}
		return result, err
	}

	operatorNamespace := s.Util.GetOperatorNamespace()

	joinedCondition := meta.FindStatusCondition(socket.Status.Conditions, "Joined")
	if joinedCondition == nil {
		socket.Status.Phase = integrationv1alpha2.PendingPhase
		socket.Status.Ready = false
		meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
			Message:            "socket created",
			ObservedGeneration: socket.Generation,
			Reason:             "SocketCreated",
			Status:             "False",
			Type:               "Joined",
		})
		err := r.Status().Update(ctx, socket)
		if err != nil {
			return result, err
		}
		err = coupler.GlobalCoupler.CreatedSocket(socket)
		if err != nil {
			return result, nil
		}
	}

	socketInterface := &integrationv1alpha2.Interface{}
	err = r.Get(ctx, s.Util.EnsureNamespacedName(&socket.Spec.Interface, operatorNamespace), socketInterface)
	if err != nil {
		socket.Status.Phase = integrationv1alpha2.FailedPhase
		socket.Status.Ready = false
		meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
			Message:            err.Error(),
			ObservedGeneration: socket.Generation,
			Reason:             "Error",
			Status:             "False",
			Type:               "Joined",
		})
		err = r.Status().Update(ctx, socket)
		if err != nil {
			return result, err
		}
		return result, nil
	}

	socket.Status.Phase = integrationv1alpha2.ReadyPhase
	socket.Status.Ready = true
	meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
		Message:            "socket ready",
		ObservedGeneration: socket.Generation,
		Reason:             "SocketReady",
		Status:             "False",
		Type:               "Joined",
	})
	err = r.Status().Update(ctx, socket)
	if err != nil {
		return result, err
	}

	time.Sleep(time.Second * 5)

	// TODO: maybe ignore if plug not found
	for _, connectedPlug := range socket.Status.CoupledPlugs {
		plug := &integrationv1alpha2.Plug{}
		err := r.Get(ctx, types.NamespacedName{
			Name:      connectedPlug.Name,
			Namespace: connectedPlug.Namespace,
		}, plug)
		if err != nil {
			socket.Status.Phase = integrationv1alpha2.FailedPhase
			socket.Status.Ready = false
			meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
				Message:            err.Error(),
				ObservedGeneration: socket.Generation,
				Reason:             "Error",
				Status:             "False",
				Type:               "Joined",
			})
			err = r.Status().Update(ctx, socket)
			if err != nil {
				return result, err
			}
			return result, nil
		}
		err = coupler.GlobalCoupler.Couple(r.Client, ctx, req, &result, plug)
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func filterSocketPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
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
