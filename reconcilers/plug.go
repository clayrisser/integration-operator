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

	joinedCondition := meta.FindStatusCondition(plug.Status.Conditions, "Joined")
	if plug.Generation <= 1 && joinedCondition == nil {
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
		if strings.Index(err.Error(), "not found") > -1 {
			plug.Status.Phase = integrationv1alpha2.PendingPhase
			meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
				Message:            "waiting for socket to be created",
				ObservedGeneration: plug.Generation,
				Reason:             "SocketNotCreated",
				Status:             "False",
				Type:               "Joined",
			})
		} else {
			plug.Status.Phase = integrationv1alpha2.FailedPhase
			meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
				Message:            err.Error(),
				ObservedGeneration: plug.Generation,
				Reason:             "Error",
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

	isJoined := meta.FindStatusCondition(plug.Status.Conditions, "Joined").Status != "True"

	plug.Status.Phase = integrationv1alpha2.PendingPhase
	meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
		Message:            "coupling to socket",
		ObservedGeneration: plug.Generation,
		Reason:             "CouplingInProcess",
		Status:             "False",
		Type:               "Joined",
	})
	err = client.Status().Update(ctx, plug)

	var plugConfig []byte
	if plug.Spec.ConfigEndpoint != "" {
		plugConfig, err = coupler.GlobalCoupler.GetConfig(plug.Spec.ConfigEndpoint)
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
	}

	var socketConfig []byte
	if socket.Spec.ConfigEndpoint != "" {
		socketConfig, err = coupler.GlobalCoupler.GetConfig(socket.Spec.ConfigEndpoint)
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
	}

	if isJoined {
		err = coupler.GlobalCoupler.JoinedPlug(plug, socket, socketConfig)
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
		err = coupler.GlobalCoupler.JoinedSocket(plug, socket, plugConfig)
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
	} else {
		err = coupler.GlobalCoupler.ChangedPlug(plug, socket, socketConfig)
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
		err = coupler.GlobalCoupler.ChangedSocket(plug, socket, socketConfig)
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
	}

	plug.Status.Phase = integrationv1alpha2.SucceededPhase
	meta.SetStatusCondition(&plug.Status.Conditions, metav1.Condition{
		Message:            "coupling succeeded",
		ObservedGeneration: plug.Generation,
		Reason:             "CouplingSucceeded",
		Status:             "True",
		Type:               "Joined",
	})
	err = client.Status().Update(ctx, plug)
	if err != nil {
		return err
	}
	if isJoined {
		socket.Status.PlugsCoupledCount = socket.Status.PlugsCoupledCount + 1
		err = client.Status().Update(ctx, socket)
		if err != nil {
			return err
		}
	}

	return nil
}
