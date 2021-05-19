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

type PlugService struct {
	client *client.Client
	ctx    *context.Context
	plug   *integrationv1alpha2.Plug
	req    *ctrl.Request
	update *Update
}

func NewPlugUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	plug *integrationv1alpha2.Plug,
) *PlugService {
	return &PlugService{
		client: client,
		ctx:    ctx,
		plug:   plug,
		req:    req,
		update: NewUpdate(99),
	}
}

func (s *PlugService) Get() (*integrationv1alpha2.Plug, error) {
	client := *s.client
	ctx := *s.ctx
	plug := &integrationv1alpha2.Plug{}
	if err := client.Get(ctx, types.NamespacedName{
		Name:      s.plug.Name,
		Namespace: s.plug.Namespace,
	}, plug); err != nil {
		return nil, err
	}
	return plug.DeepCopy(), nil
}

func (s *PlugService) Update(plug *integrationv1alpha2.Plug) {
	s.update.SchedulePlugUpdate(s.client, s.ctx, nil, plug)
}

func (s *PlugService) UpdateStatus(plug *integrationv1alpha2.Plug) {
	s.update.SchedulePlugUpdateStatus(s.client, s.ctx, nil, plug)
}

func (s *PlugService) CommonUpdateStatus(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
	message string,
	reset bool,
) error {
	plug, err := s.Get()
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
	s.UpdateStatus(plug)
	return nil
}

func (s *PlugService) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	socket *integrationv1alpha2.Socket,
) error {
	return s.CommonUpdateStatus(phase, joinedStatusCondition, socket, "", false)
}

func (s *PlugService) UpdateStatusSocket(
	socket *integrationv1alpha2.Socket,
) error {
	return s.CommonUpdateStatus("", "", socket, "", false)
}

func (s *PlugService) UpdateStatusPhase(
	phase integrationv1alpha2.Phase,
) error {
	return s.CommonUpdateStatus(phase, "", nil, "", false)
}

func (s *PlugService) UpdateStatusJoinedCondition(
	joinedStatusCondition StatusCondition,
	message string,
) error {
	return s.CommonUpdateStatus("", joinedStatusCondition, nil, message, false)
}

func (s *PlugService) UpdateStatusJoinedConditionError(
	err error,
) error {
	return s.CommonUpdateStatus("", ErrorStatusCondition, nil, err.Error(), false)
}
