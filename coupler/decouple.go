package coupler

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
)

func (c *Coupler) Decouple(client client.Client, ctx context.Context, req ctrl.Request, result *ctrl.Result, plug *integrationv1alpha2.Plug) error {
	socket := &integrationv1alpha2.Socket{}
	err := client.Get(ctx, c.s.Util.EnsureNamespacedName(&plug.Spec.Socket, req.Namespace), socket)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	if err := GlobalCoupler.Departed(plug, socket); err != nil {
		return err
	}

	coupledPlugs := []integrationv1alpha2.CoupledPlug{}
	for _, coupledPlug := range socket.Status.CoupledPlugs {
		if coupledPlug.UID != plug.UID {
			coupledPlugs = append(coupledPlugs, coupledPlug)
		}
	}
	socket.Status.CoupledPlugs = coupledPlugs
	if err = client.Status().Update(ctx, socket); err != nil {
		return err
	}

	return nil
}
