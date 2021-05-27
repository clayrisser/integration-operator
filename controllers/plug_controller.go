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
	log := r.Log.WithValues("plug", req.NamespacedName)
	log.Info("RECONCILING PLUG")
	plugUtil := util.NewPlugUtil(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
		Name:      req.NamespacedName.Name,
		Namespace: req.NamespacedName.Namespace,
	}, util.GlobalPlugMutex)
	plug, err := plugUtil.Get()
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if plug.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(plug, integrationv1alpha2.PlugFinalizer) {
			result, err := coupler.GlobalCoupler.Decouple(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
				Name:      plug.Name,
				Namespace: plug.Namespace,
			})
			if err != nil {
				return result, err
			}
			coupler.GlobalCoupler.DeletedPlug(plug)
			controllerutil.RemoveFinalizer(plug, integrationv1alpha2.PlugFinalizer)
			if err := plugUtil.Update(plug); err != nil {
				return ctrl.Result{}, err
			}
			return result, nil
		}
		return ctrl.Result{}, nil
	}
	if !controllerutil.ContainsFinalizer(plug, integrationv1alpha2.PlugFinalizer) {
		controllerutil.AddFinalizer(plug, integrationv1alpha2.PlugFinalizer)
		if err := plugUtil.Update(plug); err != nil {
			return plugUtil.Error(err)
		}
		return ctrl.Result{}, nil
	}
	result, err := coupler.GlobalCoupler.Couple(&r.Client, &ctx, &req, &log, &integrationv1alpha2.NamespacedName{
		Name:      plug.Name,
		Namespace: plug.Namespace,
	})
	return result, err
}

func filterPlugPredicate() predicate.Predicate {
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
func (r *PlugReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&integrationv1alpha2.Plug{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		WithEventFilter(filterPlugPredicate()).
		Complete(r)
}
