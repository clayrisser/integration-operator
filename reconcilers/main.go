package reconcilers

import (
	"github.com/silicon-hills/integration-operator/services"
)

type Reconcilers struct {
	Plug   *PlugReconciler
	Socket *SocketReconciler
}

func NewReconcilers() *Reconcilers {
	s := services.NewServices()
	return &Reconcilers{
		Plug:   NewPlugReconciler(s),
		Socket: NewSocketReconciler(s),
	}
}
