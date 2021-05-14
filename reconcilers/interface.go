package reconcilers

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	"github.com/silicon-hills/integration-operator/services"
)

type InterfaceReconciler struct {
	s *services.Services
}

func NewInterfaceReconciler(s *services.Services) *InterfaceReconciler {
	return &InterfaceReconciler{s: s}
}

func (p *InterfaceReconciler) Reconcile(client client.Client, ctx context.Context, req ctrl.Request, result *ctrl.Result, intrface *integrationv1alpha2.Interface) error {
	return nil
}
