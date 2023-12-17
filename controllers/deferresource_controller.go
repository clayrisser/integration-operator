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

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
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
	log := log.FromContext(ctx)
	var deferResource integrationv1beta1.DeferResource
	if err := r.Get(ctx, req.NamespacedName, &deferResource); err != nil {
		log.Error(err, "unable to fetch DeferResource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if deferResource.Spec.WaitFor != nil {
		for _, target := range *deferResource.Spec.WaitFor {
			var targetResource unstructured.Unstructured
			targetResource.SetGroupVersionKind(schema.GroupVersionKind{
				Group:   target.Gvk.Group,
				Version: target.Gvk.Version,
				Kind:    target.Gvk.Kind,
			})
			if err := r.Get(ctx, types.NamespacedName{Name: target.Name}, &targetResource); err != nil {
				log.Info("waitFor target does not exist", "target", target.Name)
				return ctrl.Result{RequeueAfter: time.Duration(deferResource.Spec.Timeout) * time.Second}, nil
			}
		}
	}

	// Create the resource from spec.resource after the waitFor targets exist and the timeout has passed
	var resource unstructured.Unstructured
	if err := json.Unmarshal(deferResource.Spec.Resource.Raw, &resource); err != nil {
		log.Error(err, "unable to unmarshal resource from spec")
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, &resource); err != nil {
		log.Error(err, "unable to create resource", "resource", resource.GetName())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
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
