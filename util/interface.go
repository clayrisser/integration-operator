package util

import (
	"context"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	ctrl "sigs.k8s.io/controller-runtime"
)

type InterfaceUtil struct {
	client         *client.Client
	ctx            *context.Context
	namespacedName types.NamespacedName
	req            *ctrl.Request
	update         *Update
}

func NewInterfaceUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	namespacedName *integrationv1alpha2.NamespacedName,
) *InterfaceUtil {
	operatorNamespace := GetOperatorNamespace()
	return &InterfaceUtil{
		client:         client,
		ctx:            ctx,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
		update:         NewUpdate(99),
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
