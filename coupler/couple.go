package coupler

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Coupler) Couple(client client.Client, ctx context.Context, req ctrl.Request, result *ctrl.Result, plug *integrationv1alpha2.Plug) error {
	operatorNamespace := c.s.Util.GetOperatorNamespace()

	joinedCondition := meta.FindStatusCondition(plug.Status.Conditions, "Joined")
	if joinedCondition == nil {
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
		err = GlobalCoupler.CreatedPlug(plug)
		if err != nil {
			return nil
		}
	}

	plugInterface := &integrationv1alpha2.Interface{}
	err := client.Get(ctx, c.s.Util.EnsureNamespacedName(&plug.Spec.Interface, operatorNamespace), plugInterface)
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
	err = client.Get(ctx, c.s.Util.EnsureNamespacedName(&plug.Spec.Socket, req.Namespace), socket)
	if err != nil {
		if errors.IsNotFound(err) {
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
	err = client.Get(ctx, c.s.Util.EnsureNamespacedName(&socket.Spec.Interface, req.Namespace), socketInterface)
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

	condition := meta.FindStatusCondition(plug.Status.Conditions, "Joined")
	isJoined := condition != nil && condition.Status != "True"

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
		plugConfig, err = GlobalCoupler.GetConfig(plug.Spec.ConfigEndpoint)
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
		socketConfig, err = GlobalCoupler.GetConfig(socket.Spec.ConfigEndpoint)
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
		err = GlobalCoupler.JoinedPlug(plug, socket, socketConfig)
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
		err = GlobalCoupler.JoinedSocket(plug, socket, plugConfig)
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
		plug.Status.CoupledSocket = integrationv1alpha2.CoupledSocket{
			APIVersion: socket.APIVersion,
			Kind:       socket.Kind,
			Name:       socket.Name,
			Namespace:  socket.Namespace,
			UID:        socket.UID,
		}
		err = client.Status().Update(ctx, plug)
		if err != nil {
			return err
		}
		isCoupled := false
		for _, coupledPlug := range socket.Status.CoupledPlugs {
			if coupledPlug.UID == plug.UID {
				isCoupled = true
				break
			}
		}
		if !isCoupled {
			socket.Status.CoupledPlugs = append(socket.Status.CoupledPlugs, integrationv1alpha2.CoupledPlug{
				APIVersion: plug.APIVersion,
				Kind:       plug.Kind,
				Name:       plug.Name,
				Namespace:  plug.Namespace,
				UID:        plug.UID,
			})
			meta.SetStatusCondition(&socket.Status.Conditions, metav1.Condition{
				Message:            "socket ready with " + fmt.Sprint(len(socket.Status.CoupledPlugs)) + " plugs coupled",
				ObservedGeneration: socket.Generation,
				Reason:             "SocketReady",
				Status:             "True",
				Type:               "Joined",
			})
			err = client.Status().Update(ctx, socket)
			if err != nil {
				return err
			}
		}
	} else {
		err = GlobalCoupler.ChangedPlug(plug, socket, socketConfig)
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
		err = GlobalCoupler.ChangedSocket(plug, socket, socketConfig)
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

	return nil
}
