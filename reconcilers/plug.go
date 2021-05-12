package reconcilers

import (
	"context"
	"strings"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/coupler"

	"github.com/silicon-hills/integration-operator/services"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PlugReconciler struct {
	s *services.Services
}

func NewPlugReconciler(s *services.Services) *PlugReconciler {
	return &PlugReconciler{s: s}
}

func (p *PlugReconciler) Reconcile(client client.Client, ctx context.Context, req ctrl.Request, result *ctrl.Result, plug *integrationv1alpha2.Plug) error {
	operatorNamespace := p.s.Util.GetOperatorNamespace()

	if plug.Generation <= 1 {
		plug.Status.Phase = integrationv1alpha2.PendingPhase
		meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
			Message:            "plug created",
			ObservedGeneration: plug.Generation,
			Reason:             "PlugCreated",
			Status:             "False",
			Type:               "Joined",
		})
		err := client.Status().Update(ctx, plug)
		if err != nil {
			return err
		}
		err = coupler.GlobalCoupler.CreatedPlug(plug)
		if err != nil {
			return nil
		}
	}

	plugInterface := &integrationv1alpha2.Interface{}
	err := client.Get(ctx, p.s.Util.EnsureNamespacedName(&plug.Spec.Interface, operatorNamespace), plugInterface)
	if err != nil {
		plug.Status.Phase = integrationv1alpha2.FailedPhase
		meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
			Message:            err.Error(),
			ObservedGeneration: plug.Generation,
			Reason:             "Error",
			Status:             "False",
			Type:               "Joined",
		})
		err = client.Status().Update(ctx, plug)
		if err != nil {
			return err
		}
		return nil
	}

	socket := &integrationv1alpha2.Socket{}
	err = client.Get(ctx, p.s.Util.EnsureNamespacedName(&plug.Spec.Socket, req.Namespace), socket)
	if err != nil {
		if strings.Index(err.Error(), "not found") <= -1 {
			plug.Status.Phase = integrationv1alpha2.FailedPhase
			meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
				Message:            err.Error(),
				ObservedGeneration: plug.Generation,
				Reason:             "Error",
				Status:             "False",
				Type:               "Joined",
			})
		} else {
			plug.Status.Phase = integrationv1alpha2.PendingPhase
			meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
				Message:            "waiting for socket to be created",
				ObservedGeneration: plug.Generation,
				Reason:             "SocketNotCreated",
				Status:             "False",
				Type:               "Joined",
			})
		}
		err = client.Status().Update(ctx, plug)
		if err != nil {
			return err
		}
		return nil
	}
	if !socket.Status.Ready {
		plug.Status.Phase = integrationv1alpha2.PendingPhase
		meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
			Message:            "waiting for socket to be ready",
			ObservedGeneration: plug.Generation,
			Reason:             "SocketNotReady",
			Status:             "False",
			Type:               "Joined",
		})
		err = client.Status().Update(ctx, plug)
		if err != nil {
			return err
		}
		return nil
	}

	socketInterface := &integrationv1alpha2.Interface{}
	err = client.Get(ctx, p.s.Util.EnsureNamespacedName(&socket.Spec.Interface, req.Namespace), socketInterface)
	if err != nil {
		plug.Status.Phase = integrationv1alpha2.FailedPhase
		meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
			Message:            err.Error(),
			ObservedGeneration: plug.Generation,
			Reason:             "Error",
			Status:             "False",
			Type:               "Joined",
		})
		err = client.Status().Update(ctx, plug)
		if err != nil {
			return err
		}
		return nil
	}

	if plug.Generation <= 1 {
		coupler.GlobalCoupler.Joined(plug, socket, []byte(""))
	} else {
		coupler.GlobalCoupler.ChangedPlug(plug, socket, []byte(""))
	}

	return nil
}
