/**
 * File: /util/deferResource.go
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

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apiserver/pkg/registry/generic/registry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeferResourceUtil struct {
	apparatusUtil  *ApparatusUtil
	client         *client.Client
	ctx            context.Context
	namespacedName types.NamespacedName
	req            *ctrl.Request
	resultUtil     *ResultUtil
}

func NewDeferResourceUtil(
	client *client.Client,
	ctx context.Context,
	req *ctrl.Request,
	namespacedName *integrationv1beta1.NamespacedName,
) *DeferResourceUtil {
	operatorNamespace := GetOperatorNamespace()
	return &DeferResourceUtil{
		apparatusUtil:  NewApparatusUtil(ctx),
		client:         client,
		ctx:            ctx,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
		resultUtil:     NewResultUtil(ctx),
	}
}

func (u *DeferResourceUtil) Get() (*integrationv1beta1.DeferResource, error) {
	client := *u.client
	ctx := u.ctx
	plug := &integrationv1beta1.DeferResource{}
	if err := client.Get(ctx, u.namespacedName, plug); err != nil {
		return nil, err
	}
	return plug.DeepCopy(), nil
}

func (u *DeferResourceUtil) Update(plug *integrationv1beta1.DeferResource, requeue bool) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Update(ctx, plug); err != nil {
		return u.Error(err, plug)
	}
	if requeue {
		return ctrl.Result{Requeue: true, RequeueAfter: 0}, nil
	}
	return ctrl.Result{}, nil
}

func (u *DeferResourceUtil) UpdateStatus(
	plug *integrationv1beta1.DeferResource,
	requeue bool,
) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Status().Update(ctx, plug); err != nil {
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

func (u *DeferResourceUtil) Delete(plug *integrationv1beta1.DeferResource) (ctrl.Result, error) {
	client := *u.client
	ctx := u.ctx
	if err := client.Delete(ctx, plug); err != nil {
		return u.Error(err, plug)
	}
	return ctrl.Result{}, nil
}

func (u *DeferResourceUtil) GetResolvedCondition(
	deferResource *integrationv1beta1.DeferResource,
) (*metav1.Condition, error) {
	if deferResource == nil {
		var err error
		deferResource, err = u.Get()
		if err != nil {
			return nil, err
		}
	}
	coupledCondition := meta.FindStatusCondition(deferResource.Status.Conditions, string(DeferResourceConditionTypeResolved))
	return coupledCondition, nil
}

func (u *DeferResourceUtil) Error(
	err error,
	deferResource *integrationv1beta1.DeferResource,
) (ctrl.Result, error) {
	e := err
	if deferResource == nil {
		var err error
		deferResource, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	result, err := u.UpdateErrorStatus(e, deferResource)
	if strings.Contains(e.Error(), "result property") &&
		strings.Contains(e.Error(), "is required") {
		return ctrl.Result{Requeue: true}, nil
	}
	return result, err
}

func (u *DeferResourceUtil) UpdateErrorStatus(
	err error,
	plug *integrationv1beta1.DeferResource,
) (ctrl.Result, error) {
	e := err
	if plug == nil {
		var err error
		plug, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if err = u.setErrorStatus(e, plug); err != nil {
		return ctrl.Result{}, err
	}
	if _, err := u.UpdateStatus(plug, true); err != nil {
		return ctrl.Result{}, err
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return ctrl.Result{Requeue: true}, nil
	}
	return ctrl.Result{}, e
}

func (u *DeferResourceUtil) UpdateResolvedStatus(
	conditionResolvedReason DeferResourceConditionResolvedReason,
	deferResource *integrationv1beta1.DeferResource,
	requeue bool,
) (ctrl.Result, error) {
	if deferResource == nil {
		var err error
		deferResource, err = u.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if deferResource != nil {
		targetResource, err := u.unmarshalTargetResource(deferResource)
		if err != nil {
			return ctrl.Result{}, err
		}
		u.setResolvedStatus(deferResource, targetResource)
	}
	if conditionResolvedReason != "" {
		u.setResolvedStatusCondition(conditionResolvedReason, "", deferResource)
	}
	return u.UpdateStatus(deferResource, requeue)
}

func (u *DeferResourceUtil) setResolvedStatusCondition(
	conditionResolvedReason DeferResourceConditionResolvedReason,
	message string,
	deferResource *integrationv1beta1.DeferResource,
) {
	resolvedStatus := false
	if message == "" {
		if conditionResolvedReason == DeferResourcePending {
			message = "pending"
		} else if conditionResolvedReason == DeferResourceSuccess {
			message = "success"
		} else if conditionResolvedReason == DeferResourceError {
			message = "unknown error"
		}
	}
	if conditionResolvedReason != DeferResourceError {
		deferResource.Status.Conditions = []metav1.Condition{}
	}
	if conditionResolvedReason == DeferResourceSuccess {
		resolvedStatus = true
	}
	condition := metav1.Condition{
		Message:            message,
		ObservedGeneration: deferResource.Generation,
		Reason:             string(conditionResolvedReason),
		Status:             "False",
		Type:               string(DeferResourceConditionTypeResolved),
	}
	if resolvedStatus {
		condition.Status = "True"
	}
	meta.SetStatusCondition(&deferResource.Status.Conditions, condition)
}

func (u *DeferResourceUtil) setErrorStatus(err error, deferResource *integrationv1beta1.DeferResource) error {
	e := err
	if e == nil {
		return nil
	}
	if deferResource == nil {
		return nil
	}
	if strings.Contains(e.Error(), registry.OptimisticLockErrorMsg) {
		return nil
	}
	message := e.Error()
	resolvedCondition, err := u.GetResolvedCondition(deferResource)
	if err != nil {
		return err
	}
	if resolvedCondition != nil {
		u.setResolvedStatusCondition(DeferResourceError, "failed", deferResource)
	}
	failedCondition := metav1.Condition{
		Message:            message,
		ObservedGeneration: deferResource.Generation,
		Reason:             "Error",
		Status:             "True",
		Type:               string(DeferResourceConditionTypeFailed),
	}
	meta.SetStatusCondition(&deferResource.Status.Conditions, failedCondition)
	return nil
}

func (u *DeferResourceUtil) setResolvedStatus(
	deferResource *integrationv1beta1.DeferResource,
	targetResource *unstructured.Unstructured,
) {
	if targetResource != nil {
		deferResource.Status.OwnerReference = metav1.OwnerReference{
			APIVersion: targetResource.GetAPIVersion(),
			Kind:       targetResource.GetKind(),
			Name:       targetResource.GetName(),
			UID:        targetResource.GetUID(),
		}
	}
}

func (u *DeferResourceUtil) unmarshalTargetResource(
	deferResource *integrationv1beta1.DeferResource,
) (*unstructured.Unstructured, error) {
	if deferResource.Spec.Resource != nil {
		targetResource := &unstructured.Unstructured{}
		err := json.Unmarshal([]byte(deferResource.Spec.Resource.Raw), targetResource)
		if err != nil {
			return nil, err
		}
		return targetResource, nil
	}
	return nil, nil
}
