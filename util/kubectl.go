/*
 * File: /util/kubectl.go
 * Project: integration-operator
 * File Created: 23-06-2021 22:52:53
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 27-06-2021 01:20:17
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Silicon Hills LLC (c) Copyright 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type KubectlUtil struct {
	ctx *context.Context
	cfg *rest.Config
}

func NewKubectlUtil(ctx *context.Context, cfg *rest.Config) *KubectlUtil {
	return &KubectlUtil{
		cfg: cfg,
		ctx: ctx,
	}
}

func (u *KubectlUtil) Create(body []byte) error {
	dr, obj, err := u.prepareDynamic(body)
	if err != nil {
		return err
	}
	if _, err := dr.Create(*u.ctx, obj, metav1.CreateOptions{
		FieldManager: "integration-operator",
	}); err != nil {
		return err
	}
	return nil
}

func (u *KubectlUtil) Update(body []byte) error {
	dr, obj, err := u.prepareDynamic(body)
	if err != nil {
		return err
	}
	if _, err := dr.Update(*u.ctx, obj, metav1.UpdateOptions{
		FieldManager: "integration-operator",
	}); err != nil {
		return err
	}
	return nil
}

func (u *KubectlUtil) Apply(body []byte) error {
	dr, obj, err := u.prepareDynamic(body)
	if err != nil {
		return err
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	if _, err = dr.Patch(*u.ctx, obj.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: "integration-operator",
	}); err != nil {
		return err
	}
	return nil
}

func (u *KubectlUtil) Delete(body []byte) error {
	dr, obj, err := u.prepareDynamic(body)
	if err != nil {
		return err
	}
	if err = dr.Delete(*u.ctx, obj.GetName(), metav1.DeleteOptions{}); err != nil {
		return err
	}
	return nil
}

func (u *KubectlUtil) Get(body []byte) (*unstructured.Unstructured, error) {
	dr, obj, err := u.prepareDynamic(body)
	if err != nil {
		return nil, err
	}
	return dr.Get(*u.ctx, obj.GetName(), metav1.GetOptions{})
}

// https://ymmt2005.hatenablog.com/entry/2020/04/14/An_example_of_using_dynamic_client_of_k8s.io/client-go
func (u *KubectlUtil) prepareDynamic(resource []byte) (dynamic.ResourceInterface, *unstructured.Unstructured, error) {
	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(u.cfg)
	if err != nil {
		return nil, nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(u.cfg)
	if err != nil {
		return nil, nil, err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(resource, nil, obj)
	if err != nil {
		return nil, nil, err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, nil, err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	return dr, obj, nil
}
