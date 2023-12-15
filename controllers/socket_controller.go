/**
 * File: /controllers/socket_controller.go
 * Project: integration-operator
 * File Created: 17-10-2023 10:50:35
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
	"time"

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

// SocketReconciler reconciles a Socket object
type SocketReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=integration.rock8s.com,resources=sockets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=sockets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=sockets/finalizers,verbs=update

func (r *SocketReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.V(1).Info("Socket Reconcile")
	socketUtil := util.NewSocketUtil(&r.Client, ctx, &req, &integrationv1beta1.NamespacedName{
		Name:      req.NamespacedName.Name,
		Namespace: req.NamespacedName.Namespace,
	})
	socket, err := socketUtil.Get()
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if socket.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(socket, integrationv1beta1.Finalizer) {
			coupledPlugs := socket.Status.CoupledPlugs
			if len(coupledPlugs) > 0 {
				for _, coupledPlug := range coupledPlugs {
					if _, err := socketUtil.RemoveCoupledPlugStatus(coupledPlug.UID, socket); err != nil {
						return socketUtil.Error(err, socket)
					}
				}
				if _, err := socketUtil.UpdateCoupledStatus(util.SocketEmpty, socket, nil, true); err != nil {
					return socketUtil.Error(err, socket)
				}
				for _, coupledPlug := range coupledPlugs {
					plugUtil := util.NewPlugUtil(&r.Client, ctx, &req, &integrationv1beta1.NamespacedName{
						Name:      coupledPlug.Name,
						Namespace: coupledPlug.Namespace,
					}, socket)
					plug, err := plugUtil.Get()
					if err != nil {
						if errors.IsNotFound(err) {
							continue
						}
						return socketUtil.Error(err, socket)
					}
					if _, err := plugUtil.Delete(plug); err != nil {
						return socketUtil.Error(err, socket)
					}
				}
				return ctrl.Result{Requeue: true, RequeueAfter: 0}, nil
			}
			if socket == nil {
				return ctrl.Result{}, nil
			}
			if err := coupler.DeletedSocket(socket, r.Recorder); err != nil {
				return socketUtil.Error(err, socket)
			}
			controllerutil.RemoveFinalizer(socket, integrationv1beta1.Finalizer)
			return socketUtil.Update(socket, true)
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(socket, integrationv1beta1.Finalizer) {
		controllerutil.AddFinalizer(socket, integrationv1beta1.Finalizer)
		return socketUtil.Update(socket, true)
	}

	coupledCondition, err := socketUtil.GetCoupledCondition()
	if err != nil {
		return socketUtil.Error(err, socket)
	}
	if coupledCondition == nil {
		if err := coupler.CreatedSocket(socket, r.Recorder); err != nil {
			return socketUtil.Error(err, socket)
		}
		return socketUtil.UpdateCoupledStatus(util.SocketCreated, socket, nil, true)
	}

	setSocketStatus := false
	for _, coupledPlug := range socket.Status.CoupledPlugs {
		plugUtil := util.NewPlugUtil(&r.Client, ctx, &req, &integrationv1beta1.NamespacedName{
			Name:      coupledPlug.Name,
			Namespace: coupledPlug.Namespace,
		}, socket)
		if _, err := plugUtil.Get(); err != nil {
			if errors.IsNotFound(err) {
				if _, err := socketUtil.RemoveCoupledPlugStatus(coupledPlug.UID, socket); err != nil {
					return socketUtil.Error(err, socket)
				}
				setSocketStatus = true
				continue
			}
			return socketUtil.Error(err, socket)
		}
	}
	if setSocketStatus {
		return socketUtil.UpdateStatus(socket, true)
	}

	for _, coupledPlug := range socket.Status.CoupledPlugs {
		plugUtil := util.NewPlugUtil(&r.Client, ctx, &req, &integrationv1beta1.NamespacedName{
			Name:      coupledPlug.Name,
			Namespace: coupledPlug.Namespace,
		}, socket)
		plug, err := plugUtil.Get()
		if err != nil {
			return socketUtil.Error(err, socket)
		}
		plug.Spec.Epoch = strconv.FormatInt(time.Now().Unix(), 10)
		if _, err := plugUtil.Update(plug, false); err != nil {
			return socketUtil.Error(err, socket)
		}
	}

	return socketUtil.UpdateCoupledStatus(util.SocketCoupled, socket, nil, false)
}

func filterSocketPredicate() predicate.Predicate {
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
func (r *SocketReconciler) SetupWithManager(mgr ctrl.Manager) error {
	maxConcurrentReconciles := 3
	if value := os.Getenv("MAX_CONCURRENT_RECONCILES"); value != "" {
		if val, err := strconv.Atoi(value); err == nil {
			maxConcurrentReconciles = val
		}
	}
	return ctrl.NewControllerManagedBy(mgr).
		WithOptions(controller.Options{MaxConcurrentReconciles: maxConcurrentReconciles}).
		WithEventFilter(filterSocketPredicate()).
		For(&integrationv1beta1.Socket{}).
		Complete(r)
}
