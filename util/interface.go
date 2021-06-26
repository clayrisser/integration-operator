/*
 * File: /util/interface.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 26-06-2021 10:53:34
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

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InterfaceUtil struct {
	client         *client.Client
	ctx            *context.Context
	namespacedName types.NamespacedName
	req            *ctrl.Request
	log            *logr.Logger
}

func NewInterfaceUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	namespacedName *integrationv1alpha2.NamespacedName,
) *InterfaceUtil {
	operatorNamespace := GetOperatorNamespace()
	return &InterfaceUtil{
		client:         client,
		ctx:            ctx,
		log:            log,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
	}
}

func (u *InterfaceUtil) Get() (*integrationv1alpha2.Interface, error) {
	client := *u.client
	ctx := *u.ctx
	intrface := &integrationv1alpha2.Interface{}
	if err := client.Get(ctx, u.namespacedName, intrface); err != nil {
		return nil, err
	}
	return intrface.DeepCopy(), nil
}
