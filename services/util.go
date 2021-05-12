package services

import (
	"os"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	"k8s.io/apimachinery/pkg/types"
)

type UtilService struct {
	s *Services
}

func NewUtilService(services *Services) *UtilService {
	return &UtilService{s: services}
}

func (u *UtilService) Default(value string, defaultValue string) string {
	if value == "" {
		value = defaultValue
	}
	return value
}

func (u *UtilService) EnsureNamespacedName(
	partialNamespacedName *integrationv1alpha2.NamespacedName,
	defaultNamespace string,
) types.NamespacedName {
	return types.NamespacedName{
		Name:      partialNamespacedName.Name,
		Namespace: u.Default(partialNamespacedName.Namespace, defaultNamespace),
	}
}

func (u *UtilService) GetOperatorNamespace() string {
	operatorNamespace := os.Getenv("POD_NAMESPACE")
	if operatorNamespace == "" {
		operatorNamespace = "kube-system"
	}
	return operatorNamespace
}
