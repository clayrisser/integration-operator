package reconcilers

import (
	"github.com/silicon-hills/integration-operator/services"
)

type Reconcilers struct {
	Plug      *PlugReconciler
	Socket    *SocketReconciler
	Interface *InterfaceReconciler
}

func NewReconcilers() *Reconcilers {
	s := services.NewServices()
	return &Reconcilers{
		Interface: NewInterfaceReconciler(s),
		Plug:      NewPlugReconciler(s),
		Socket:    NewSocketReconciler(s),
	}
}
