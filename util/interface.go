/**
 * File: /interface.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 11:49:19
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

package util

import (
	"context"

	integrationv1alpha2 "gitlab.com/bitspur/rock8s/integration-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type InterfaceUtil struct {
	client         *client.Client
	ctx            *context.Context
	namespacedName types.NamespacedName
}

func NewInterfaceUtil(
	client *client.Client,
	ctx *context.Context,
	namespacedName *integrationv1alpha2.NamespacedName,
) *InterfaceUtil {
	operatorNamespace := GetOperatorNamespace()
	return &InterfaceUtil{
		client:         client,
		ctx:            ctx,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
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
