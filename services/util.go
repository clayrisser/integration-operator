package services

import (
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	ctrl "sigs.k8s.io/controller-runtime"

	"k8s.io/apimachinery/pkg/types"
)

func EnsureNamespacedName(req ctrl.Request, partialNamespacedName *integrationv1alpha2.NamespacedName) types.NamespacedName {
	namespacedName := types.NamespacedName{
		Name:      partialNamespacedName.Name,
		Namespace: partialNamespacedName.Namespace,
	}
	if partialNamespacedName.Namespace == "" {
		namespacedName.Namespace = partialNamespacedName.Namespace
	}
	return namespacedName
}
