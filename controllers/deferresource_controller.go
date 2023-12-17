/**
 * File: /controllers/deferresource_controller.go
 * Project: integration-operator
 * File Created: 17-12-2023 03:35:18
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

/*
Copyright 2023.

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
	"encoding/json"
	"os"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"gitlab.com/bitspur/rock8s/integration-operator/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// DeferResourceReconciler reconciles a DeferResource object
type DeferResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=integration.rock8s.com,resources=deferresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=deferresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=integration.rock8s.com,resources=deferresources/finalizers,verbs=update

func (r *DeferResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.V(1).Info("DeferResource Reconcile")
	deferResourceUtil := util.NewDeferResourceUtil(&r.Client, ctx, &req, &integrationv1beta1.NamespacedName{
		Name:      req.NamespacedName.Name,
		Namespace: req.NamespacedName.Namespace,
	})
	deferResource, err := deferResourceUtil.Get()
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	kubectlUtil := util.NewKubectlUtil(
		ctx, deferResource.Namespace,
		util.EnsureServiceAccount(deferResource.Spec.ServiceAccountName),
	)

	if deferResource.Spec.Timeout > 0 {
		if time.Since(deferResource.CreationTimestamp.Time) < time.Duration(deferResource.Spec.Timeout)*time.Second {
			return ctrl.Result{
				RequeueAfter: time.Duration(deferResource.Spec.Timeout)*time.Second - time.Since(deferResource.CreationTimestamp.Time),
			}, nil
		}
	}

	if deferResource.Spec.WaitFor != nil {
		for _, waitFor := range *deferResource.Spec.WaitFor {
			body, err := json.Marshal(waitFor)
			if err != nil {
				return deferResourceUtil.Error(err, deferResource)
			}
			if _, err := kubectlUtil.Get(body); err != nil {
				if k8serrors.IsNotFound(err) {
					return ctrl.Result{Requeue: true}, nil
				}
				return deferResourceUtil.Error(err, deferResource)
			}
		}
	}

	return deferResourceUtil.ApplyResource(deferResource, kubectlUtil)
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeferResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
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
