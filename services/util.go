package services

import (
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	"k8s.io/apimachinery/pkg/types"
)

type UtilService struct {
	s *Services
}

func NewUtilService(services *Services) *UtilService {
	return &UtilService{s: services}
}

func (u *UtilService) EnsureNamespacedName(partialNamespacedName *integrationv1alpha2.NamespacedName) types.NamespacedName {
	namespacedName := types.NamespacedName{
		Name:      partialNamespacedName.Name,
		Namespace: partialNamespacedName.Namespace,
	}
	if partialNamespacedName.Namespace == "" {
		namespacedName.Namespace = u.s.req.Namespace
	}
	return namespacedName
}
