package reconcilers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/coupler"

	"github.com/silicon-hills/integration-operator/services"
)

type SocketReconciler struct {
	s *services.Services
}

func NewSocketReconciler(s *services.Services) *SocketReconciler {
	return &SocketReconciler{s: s}
}

func (p *SocketReconciler) Reconcile(client client.Client, ctx context.Context, req ctrl.Request, result *ctrl.Result, socket *integrationv1alpha2.Socket) error {
	operatorNamespace := p.s.Util.GetOperatorNamespace()

	joinedCondition := meta.FindStatusCondition(socket.Status.Conditions, "Joined")
	if socket.Generation <= 1 && joinedCondition == nil {
		socket.Status.Phase = integrationv1alpha2.PendingPhase
		socket.Status.Ready = false
		meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
			Message:            "socket created",
			ObservedGeneration: socket.Generation,
			Reason:             "SocketCreated",
			Status:             "False",
			Type:               "Joined",
		})
		err := client.Status().Update(ctx, socket)
		if err != nil {
			return err
		}
		err = coupler.GlobalCoupler.CreatedSocket(socket)
		if err != nil {
			return nil
		}
	}

	socketInterface := &integrationv1alpha2.Interface{}
	err := client.Get(ctx, p.s.Util.EnsureNamespacedName(&socket.Spec.Interface, operatorNamespace), socketInterface)
	if err != nil {
		socket.Status.Phase = integrationv1alpha2.FailedPhase
		meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
			Message:            err.Error(),
			ObservedGeneration: socket.Generation,
			Reason:             "Error",
			Status:             "False",
			Type:               "Joined",
		})
		err = client.Status().Update(ctx, socket)
		if err != nil {
			return err
		}
		return nil
	}

	socket.Status.Phase = integrationv1alpha2.ReadyPhase
	socket.Status.Ready = true
	meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
		Message:            "socket ready",
		ObservedGeneration: socket.Generation,
		Reason:             "SocketReady",
		Status:             "False",
		Type:               "Joined",
	})
	err = client.Status().Update(ctx, socket)
	if err != nil {
		return err
	}

	return nil
}
