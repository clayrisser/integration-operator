/**
 * File: /controllers/plug_controller.go
 * Project: integration-operator
 * File Created: 17-10-2023 10:50:57
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

package controllers

import (
	"context"
	"os"
	"strconv"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"gitlab.com/bitspur/rock8s/integration-operator/coupler"
	"gitlab.com/bitspur/rock8s/integration-operator/util"
)

// PlugReconciler reconciles a Plug object
type PlugReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=integration.rock8s.com,resources=plugs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=plugs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=plugs/finalizers,verbs=update

func (r *PlugReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.V(1).Info("Plug Reconcile")
	plugUtil := util.NewPlugUtil(&r.Client, ctx, &req, &integrationv1beta1.NamespacedName{
		Name:      req.NamespacedName.Name,
		Namespace: req.NamespacedName.Namespace,
	})
	plug, err := plugUtil.Get()
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	socketUtil := util.NewSocketUtil(&r.Client, ctx, &req, &integrationv1beta1.NamespacedName{
		Name:      plug.Spec.Socket.Name,
		Namespace: util.Default(plug.Spec.Socket.Namespace, req.NamespacedName.Namespace),
	})
	socket, err := socketUtil.Get()
	if err != nil && !errors.IsNotFound(err) {
		return plugUtil.Error(err, plug)
	}

	if plug.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(plug, integrationv1beta1.Finalizer) {
			if plug.Status.CoupledSocket != nil && socket != nil {
				if err := coupler.Decouple(&r.Client, ctx, &req, plugUtil, socketUtil, plug, socket, r.Recorder); err != nil {
					return plugUtil.Error(err, plug)
				}
			}
			if err := coupler.DeletedPlug(plug, r.Recorder); err != nil {
				return plugUtil.Error(err, plug)
			}
			controllerutil.RemoveFinalizer(plug, integrationv1beta1.Finalizer)
			return plugUtil.Update(plug, true)
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(plug, integrationv1beta1.Finalizer) {
		controllerutil.AddFinalizer(plug, integrationv1beta1.Finalizer)
		return plugUtil.Update(plug, true)
	}

	coupledCondition, err := plugUtil.GetCoupledCondition()
	if err != nil {
		return plugUtil.Error(err, plug)
	}
	if coupledCondition == nil {
		if err := coupler.CreatedPlug(plug, r.Recorder); err != nil {
			return plugUtil.Error(err, plug)
		}
		return plugUtil.UpdateCoupledStatus(util.PlugCreated, plug, nil, true)
	}

	return coupler.Couple(&r.Client, ctx, &req, plugUtil, socketUtil, plug, socket, r.Recorder)
}

func filterPlugPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectNew.GetDeletionTimestamp() != nil || e.ObjectNew.GetGeneration() > e.ObjectOld.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return !e.DeleteStateUnknown
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlugReconciler) SetupWithManager(mgr ctrl.Manager) error {
	maxConcurrentReconciles := 3
	if value := os.Getenv("MAX_CONCURRENT_RECONCILES"); value != "" {
		if val, err := strconv.Atoi(value); err == nil {
			maxConcurrentReconciles = val
		}
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: maxConcurrentReconciles}).
		WithEventFilter(filterPlugPredicate()).
		For(&integrationv1beta1.Plug{}).
		Complete(r)
}
