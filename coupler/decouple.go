package coupler

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

	if socket != nil {
		coupledPlugs := []integrationv1alpha2.CoupledPlug{}
		for _, coupledPlug := range socket.Status.CoupledPlugs {
			if coupledPlug.UID != plug.UID {
				coupledPlugs = append(coupledPlugs, coupledPlug)
			}
		}
		socket.Status.CoupledPlugs = coupledPlugs
		condition := metav1.Condition{
			Message:            "socket ready with " + fmt.Sprint(len(coupledPlugs)) + " plugs coupled",
			ObservedGeneration: socket.Generation,
			Reason:             "SocketReady",
			Status:             "False",
			Type:               "Joined",
		}
		if len(coupledPlugs) > 0 {
			condition.Status = "True"
		}
		meta.SetStatusCondition(&socket.Status.Conditions, condition)
		if err = client.Status().Update(ctx, socket); err != nil {
			return err
		}
	}

	controllerutil.RemoveFinalizer(plug, integrationv1alpha2.PlugFinalizer)
	if err := client.Update(ctx, plug); err != nil {
		return err
	}

	return nil
}
