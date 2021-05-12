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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"

	"sigs.k8s.io/controller-runtime/pkg/predicate"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/coupler"
	"github.com/silicon-hills/integration-operator/services"
)

// PlugReconciler reconciles a Plug object
type PlugReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=integration.siliconhills.dev,resources=plugs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=integration.siliconhills.dev,resources=plugs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=integration.siliconhills.dev,resources=plugs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Plug object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *PlugReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("plug", req.NamespacedName)
	s := services.NewServices()
	result := ctrl.Result{}
	operatorNamespace := s.Util.GetOperatorNamespace()

	plug := &integrationv1alpha2.Plug{}
	err := r.Get(ctx, req.NamespacedName, plug)
	if err != nil {
		return result, nil
	}

	if plug.Generation <= 1 {
		err = coupler.GlobalCoupler.CreatedPlug(plug)
		if err != nil {
			return result, nil
		}
	}

	socket := &integrationv1alpha2.Socket{}
	err = r.Get(ctx, s.Util.EnsureNamespacedName(&plug.Spec.Socket, req.Namespace), socket)
	if err != nil {
		return result, nil
	}

	plugInterface := &integrationv1alpha2.Interface{}
	err = r.Get(ctx, s.Util.EnsureNamespacedName(&plug.Spec.Interface, operatorNamespace), plugInterface)
	if err != nil {
		return result, nil
	}

	socketInterface := &integrationv1alpha2.Interface{}
	err = r.Get(ctx, s.Util.EnsureNamespacedName(&socket.Spec.Interface, req.Namespace), socketInterface)
	if err != nil {
		return result, nil
	}

	if plug.Generation <= 1 {
		coupler.GlobalCoupler.Joined(plug, socket, []byte(""))
	} else {
		coupler.GlobalCoupler.ChangedPlug(plug, socket, []byte(""))
	}

	return result, nil
}

func filterPlugPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			coupler.GlobalCoupler.Departed()
			return !e.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlugReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&integrationv1alpha2.Plug{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		WithEventFilter(filterPlugPredicate()).
		Complete(r)
}
