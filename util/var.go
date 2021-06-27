/*
 * File: /util/var.go
 * Project: integration-operator
 * File Created: 24-06-2021 22:10:01
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 27-06-2021 05:13:27
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

	"github.com/tidwall/gjson"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type VarUtil struct {
	client       *kubernetes.Clientset
	ctx          *context.Context
	resourceUtil *ResourceUtil
	kubectlUtil  *KubectlUtil
}

func NewVarUtil(ctx *context.Context) *VarUtil {
	return &VarUtil{
		client:       kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		kubectlUtil:  NewKubectlUtil(ctx),
		resourceUtil: NewResourceUtil(ctx),
	}
}

func (u *VarUtil) GetVars(namespace string, vars []kustomizeTypes.Var) (map[string]string, error) {
	resultMap := make(map[string]string)
	for _, v := range vars {
		varResult, err := u.GetVar(namespace, v)
		if err != nil {
			return nil, err
		}
		resultMap[v.Name] = varResult
	}
	return resultMap, nil
}

func (u *VarUtil) GetVar(namespace string, v kustomizeTypes.Var) (string, error) {
	resource, err := u.resourceUtil.GetResource(namespace, v.ObjRef)
	if err != nil {
		return "", err
	}
	bResource, err := json.Marshal(resource)
	if err != nil {
		return "", err
	}
	return gjson.Parse(string(bResource)).Get(v.FieldRef.FieldPath).String(), nil
}
