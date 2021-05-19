package util

import (
	"context"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type PlugUtil struct {
	client         *client.Client
	ctx            *context.Context
	namespacedName types.NamespacedName
	req            *ctrl.Request
	update         *Update
}

func NewPlugUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	namespacedName *integrationv1alpha2.NamespacedName,
) *PlugUtil {
	operatorNamespace := GetOperatorNamespace()
	return &PlugUtil{
		client:         client,
		ctx:            ctx,
		namespacedName: EnsureNamespacedName(namespacedName, operatorNamespace),
		req:            req,
		update:         NewUpdate(99),
	}
}

func (u *PlugUtil) Get() (*integrationv1alpha2.Plug, error) {
	client := *u.client
	ctx := *u.ctx
	plug := &integrationv1alpha2.Plug{}
	if err := client.Get(ctx, u.namespacedName, plug); err != nil {
		return nil, err
	}
	return plug.DeepCopy(), nil
}

func (u *PlugUtil) Update(plug *integrationv1alpha2.Plug) {
	u.update.SchedulePlugUpdate(u.client, u.ctx, nil, plug)
}

func (u *PlugUtil) UpdateStatus(plug *integrationv1alpha2.Plug) {
	u.update.SchedulePlugUpdateStatus(u.client, u.ctx, nil, plug)
}

func (u *PlugUtil) CommonUpdateStatus(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
	message string,
	reset bool,
) error {
	plug, err := u.Get()
	if err != nil {
		return err
	}
	if reset {
		plug.Status = integrationv1alpha2.PlugStatus{}
	}
	if joinedStatusCondition != "" {
		joinedStatus := false
		if message == "" {
			if joinedStatusCondition == PlugCreatedStatusCondition {
				message = "plug created"
			} else if joinedStatusCondition == SocketNotCreatedStatusCondition {
				message = "waiting for socket to be created"
			} else if joinedStatusCondition == SocketNotReadyStatusCondition {
				message = "waiting for socket to be ready"
			} else if joinedStatusCondition == CouplingInProcessStatusCondition {
				message = "coupling to socket"
			} else if joinedStatusCondition == CouplingSucceededStatusCondition {
				message = "coupling succeeded"
			} else if joinedStatusCondition == ErrorStatusCondition {
				message = "unknown error"
			}
		}
		if joinedStatusCondition == CouplingSucceededStatusCondition {
			joinedStatus = true
		}
		c := metav1.Condition{
			Message:            message,
			ObservedGeneration: plug.Generation,
			Reason:             string(joinedStatusCondition),
			Status:             "False",
			Type:               "Joined",
		}
		if joinedStatus {
			c.Status = "True"
		}
		meta.SetStatusCondition(&plug.Status.Conditions, c)
	}
	if socket != nil {
		coupledSocket := integrationv1alpha2.CoupledSocket{
			APIVersion: socket.APIVersion,
			Kind:       socket.Kind,
			Name:       socket.Name,
			Namespace:  socket.Namespace,
			UID:        socket.UID,
		}
		plug.Status.CoupledSocket = coupledSocket
	}
	if phase != "" {
		plug.Status.Phase = phase
	}
	u.UpdateStatus(plug)
	return nil
}

func (u *PlugUtil) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
) error {
	return u.CommonUpdateStatus(phase, joinedStatusCondition, socket, "", false)
}

func (u *PlugUtil) UpdateStatusSocket(
	socket *integrationv1alpha2.Socket,
) error {
	return u.CommonUpdateStatus("", "", socket, "", false)
}

func (u *PlugUtil) UpdateStatusPhase(
	phase integrationv1alpha2.Phase,
) error {
	return u.CommonUpdateStatus(phase, "", nil, "", false)
}

func (u *PlugUtil) UpdateStatusJoinedCondition(
	joinedStatusCondition StatusCondition,
	message string,
) error {
	return u.CommonUpdateStatus("", joinedStatusCondition, nil, message, false)
}

func (u *PlugUtil) UpdateStatusJoinedConditionError(
	err error,
) error {
	return u.CommonUpdateStatus(integrationv1alpha2.FailedPhase, ErrorStatusCondition, nil, err.Error(), false)
}

func (u *PlugUtil) GetJoinedCondition(plug *integrationv1alpha2.Plug) (*metav1.Condition, error) {
	if plug == nil {
		var err error
		plug, err = u.Get()
		if err != nil {
			return nil, err
		}
	}
	joinedCondition := meta.FindStatusCondition(plug.Status.Conditions, "Joined")
	return joinedCondition, nil
}
