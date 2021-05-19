package util

import (
	"os"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	"k8s.io/apimachinery/pkg/types"
)

func Default(value string, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

func EnsureNamespacedName(
	partialNamespacedName *integrationv1alpha2.NamespacedName,
	defaultNamespace string,
) types.NamespacedName {
	return types.NamespacedName{
		Name:      partialNamespacedName.Name,
		Namespace: Default(partialNamespacedName.Namespace, defaultNamespace),
	}
}

func GetOperatorNamespace() string {
	operatorNamespace := os.Getenv("POD_NAMESPACE")
	if operatorNamespace == "" {
		operatorNamespace = "kube-system"
	}
	return operatorNamespace
}
