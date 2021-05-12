package reconcilers

import (
	"github.com/silicon-hills/integration-operator/services"
)

type Reconcilers struct {
	Plug *PlugReconciler
}

func NewReconcilers() *Reconcilers {
	s := services.NewServices()
	return &Reconcilers{
		Plug: NewPlugReconciler(s),
	}
}
