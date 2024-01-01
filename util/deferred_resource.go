/**
 * File: /util/deferred_resource.go
 * Project: integration-operator
 * File Created: 17-12-2023 03:49:19
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

package util

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeferredResourceUtil struct {
	apparatusUtil  *ApparatusUtil
	client         *client.Client
	ctx            context.Context
	namespacedName types.NamespacedName
	req            *ctrl.Request
	resultUtil     *ResultUtil
}

func NewDeferredResourceUtil(
	client *client.Client,
	ctx context.Context,
	req *ctrl.Request,
	namespacedName *integrationv1beta1.NamespacedName,
) *DeferredResourceUtil {
	operatorNamespace := GetOperatorNamespace()
	return &DeferredResourceUtil{
		apparatusUtil:  NewApparatusUtil(ctx),
		client:         client,
		ctx:            ctx,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
		resultUtil:     NewResultUtil(ctx),
	}
}

func (u *DeferredResourceUtil) Get() (*integrationv1beta1.DeferredResource, error) {
	client := *u.client
	ctx := u.ctx
	deferredResource := &integrationv1beta1.DeferredResource{}
	if err := client.Get(ctx, u.namespacedName, deferredResource); err != nil {
		return nil, err
	}
	return deferredResource.DeepCopy(), nil
}

func (u *DeferredResourceUtil) Update(deferredResource *integrationv1beta1.DeferredResource, requeue bool) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Update(ctx, deferredResource); err != nil {
		return u.Error(err, deferredResource)
	}
	if requeue {
		return ctrl.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	return ctrl.Result{}, nil
}

func (u *DeferredResourceUtil) UpdateStatus(
	deferredResource *integrationv1beta1.DeferredResource,
	requeue bool,
) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Status().Update(ctx, deferredResource); err != nil {
		if strings.Contains(err.Error(), registry.OptimisticLockErrorMsg) {
			return ctrl.Result{Requeue: true}, nil
		}
		return ctrl.Result{}, err
	}
	if requeue {
		return ctrl.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	return ctrl.Result{}, nil
}

func (u *DeferredResourceUtil) Delete(deferredResource *integrationv1beta1.DeferredResource) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Delete(ctx, deferredResource); err != nil {
		return u.Error(err, deferredResource)
	}
	return ctrl.Result{}, nil
}

func (u *DeferredResourceUtil) GetResolvedCondition(
	deferredResource *integrationv1beta1.DeferredResource,
) (*metav1.Condition, error) {
	if deferredResource == nil {
		var err error
		deferredResource, err = u.Get()
		if err != nil {
			return nil, err
		}
	}
	coupledCondition := meta.FindStatusCondition(deferredResource.Status.Conditions, string(DeferredResourceConditionTypeResolved))
	return coupledCondition, nil
}

func (u *DeferredResourceUtil) Error(
	err error,
	deferredResource *integrationv1beta1.DeferredResource,
) (ctrl.Result, error) {
	e := err
	if deferredResource == nil {
		var err error
		deferredResource, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	return u.UpdateErrorStatus(e, deferredResource)
}

func (u *DeferredResourceUtil) UpdateErrorStatus(
	err error,
	deferredResource *integrationv1beta1.DeferredResource,
) (ctrl.Result, error) {
	e := err
	if deferredResource == nil {
		var err error
		deferredResource, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if err = u.setErrorStatus(e, deferredResource); err != nil {
		return ctrl.Result{}, err
	}
	if _, err := u.UpdateStatus(deferredResource, true); err != nil {
		return ctrl.Result{}, err
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, e
}

func (u *DeferredResourceUtil) UpdateResolvedStatus(
	conditionResolvedReason DeferredResourceConditionResolvedReason,
	deferredResource *integrationv1beta1.DeferredResource,
	appliedResource *unstructured.Unstructured,
	message string,
	requeue bool,
) (ctrl.Result, error) {
	if deferredResource == nil {
		var err error
		deferredResource, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if deferredResource != nil {
		u.setResolvedStatus(deferredResource, appliedResource)
	}
	if conditionResolvedReason != "" {
		u.setResolvedStatusCondition(conditionResolvedReason, message, deferredResource)
	}
	return u.UpdateStatus(deferredResource, requeue)
}

func (u *DeferredResourceUtil) ApplyResource(
	deferredResource *integrationv1beta1.DeferredResource,
	kubectlUtil *KubectlUtil,
) (ctrl.Result, error) {
	resource, err := u.getResource(deferredResource, kubectlUtil)
	if err != nil {
		return ctrl.Result{}, err
	}
	if appliedResource, err := kubectlUtil.Get(resource); err == nil {
		return u.UpdateResolvedStatus(DeferredResourceSuccess, deferredResource, appliedResource, "", false)
	}
	err = kubectlUtil.Apply(resource)
	if err != nil {
		return ctrl.Result{}, err
	}
	var appliedResource *unstructured.Unstructured
	maxRetries := 5
	retryInterval := time.Second * 2
	for i := 0; i < maxRetries; i++ {
		appliedResource, err = kubectlUtil.Get(resource)
		if err == nil {
			break
		}
		time.Sleep(retryInterval)
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	return u.UpdateResolvedStatus(DeferredResourceSuccess, deferredResource, appliedResource, "", false)
}

func (u *DeferredResourceUtil) DeleteResource(
	deferredResource *integrationv1beta1.DeferredResource,
	kubectlUtil *KubectlUtil,
) error {
	resource, err := u.getResource(deferredResource, kubectlUtil)
	if err != nil {
		return err
	}
	err = kubectlUtil.Delete(resource)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (u *DeferredResourceUtil) getResource(
	deferredResource *integrationv1beta1.DeferredResource,
	kubectlUtil *KubectlUtil,
) ([]byte, error) {
	var resource unstructured.Unstructured
	if err := json.Unmarshal(deferredResource.Spec.Resource.Raw, &resource); err != nil {
		return nil, err
	}
	resource.SetNamespace(deferredResource.Namespace)
	raw, err := json.Marshal(resource.Object)
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (u *DeferredResourceUtil) setResolvedStatusCondition(
	conditionResolvedReason DeferredResourceConditionResolvedReason,
	message string,
	deferredResource *integrationv1beta1.DeferredResource,
) {
	resolvedStatus := false
	if message == "" {
		if conditionResolvedReason == DeferredResourcePending {
			message = "pending"
		} else if conditionResolvedReason == DeferredResourceSuccess {
			message = "success"
		} else if conditionResolvedReason == DeferredResourceError {
			message = "unknown error"
		}
	}
	if conditionResolvedReason != DeferredResourceError {
		deferredResource.Status.Conditions = []metav1.Condition{}
	}
	if conditionResolvedReason == DeferredResourceSuccess {
		resolvedStatus = true
	}
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: deferredResource.Generation,
		Reason:             string(conditionResolvedReason),
		Status:             "False",
		Type:               string(DeferredResourceConditionTypeResolved),
	}
	if resolvedStatus {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&deferredResource.Status.Conditions, condition)
}

func (u *DeferredResourceUtil) setErrorStatus(err error, deferredResource *integrationv1beta1.DeferredResource) error {
	e := err
	if e == nil {
		return nil
	}
	if deferredResource == nil {
		return nil
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return nil
	}
	message := e.Error()
	resolvedCondition, err := u.GetResolvedCondition(deferredResource)
	if err != nil {
		return err
	}
	if resolvedCondition != nil {
		u.setResolvedStatusCondition(DeferredResourceError, "failed", deferredResource)
	}
	failedCondition := metav1.Condition{
		Message:            message,
		ObservedGeneration: deferredResource.Generation,
		Reason:             "Error",
		Status:             "True",
		Type:               string(DeferredResourceConditionTypeFailed),
	}
	meta.SetStatusCondition(&deferredResource.Status.Conditions, failedCondition)
	return nil
}

func (u *DeferredResourceUtil) setResolvedStatus(
	deferredResource *integrationv1beta1.DeferredResource,
	appliedResource *unstructured.Unstructured,
) {
	if appliedResource != nil {
		deferredResource.Status.OwnerReference = metav1.OwnerReference{
			APIVersion: appliedResource.GetAPIVersion(),
			Kind:       appliedResource.GetKind(),
			Name:       appliedResource.GetName(),
			UID:        appliedResource.GetUID(),
		}
	}
}
