package util

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"

	ctrl "sigs.k8s.io/controller-runtime"
)

type SocketService struct {
	client *client.Client
	ctx    *context.Context
	socket *integrationv1alpha2.Socket
	req    *ctrl.Request
	update *Update
}

func NewSocketUtil(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	socket *integrationv1alpha2.Socket,
) *SocketService {
	return &SocketService{
		client: client,
		ctx:    ctx,
		socket: socket,
		req:    req,
		update: NewUpdate(99),
	}
}

func (s *SocketService) Get() (*integrationv1alpha2.Socket, error) {
	client := *s.client
	ctx := *s.ctx
	socket := &integrationv1alpha2.Socket{}
	if err := client.Get(ctx, types.NamespacedName{
		Name:      s.socket.Name,
		Namespace: s.socket.Namespace,
	}, socket); err != nil {
		return nil, err
	}
	return socket.DeepCopy(), nil
}

func (s *SocketService) Update(socket *integrationv1alpha2.Socket) {
	s.update.ScheduleSocketUpdate(s.client, s.ctx, nil, socket)
}

func (s *SocketService) UpdateStatus(socket *integrationv1alpha2.Socket) {
	s.update.ScheduleSocketUpdateStatus(s.client, s.ctx, nil, socket)
}

func (s *SocketService) CommonUpdateStatus(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	appendPlug *integrationv1alpha2.Plug,
	message string,
	reset bool,
) error {
	socket, err := s.Get()
	if err != nil {
		return err
	}
	if reset {
		socket.Status = integrationv1alpha2.SocketStatus{}
	}
	if appendPlug != nil {
		socket.Status.CoupledPlugs = append(socket.Status.CoupledPlugs, integrationv1alpha2.CoupledPlug{
			APIVersion: appendPlug.APIVersion,
			Kind:       appendPlug.Kind,
			Name:       appendPlug.Name,
			Namespace:  appendPlug.Namespace,
			UID:        appendPlug.UID,
		})
		joinedCondition, err := s.GetJoinedCondition()
		if err != nil {
			return err
		}
		if (*joinedCondition).Reason == string(SocketReadyStatusCondition) {
			joinedStatusCondition = SocketReadyStatusCondition
		}
	}
	if joinedStatusCondition != "" {
		joinedStatus := false
		coupledPlugsCount := len(socket.Status.CoupledPlugs)
		if message == "" {
			if joinedStatusCondition == SocketCreatedStatusCondition {
				message = "socket created"
			} else if joinedStatusCondition == ErrorStatusCondition {
				message = "unknown error"
			} else if joinedStatusCondition == SocketReadyStatusCondition {
				message = "socket ready with " + fmt.Sprint(coupledPlugsCount) + " plugs coupled"
			}
		}
		if joinedStatusCondition == SocketReadyStatusCondition && coupledPlugsCount > 0 {
			joinedStatus = true
		}
		c := metav1.Condition{
			Message:            message,
			ObservedGeneration: socket.Generation,
			Reason:             string(joinedStatusCondition),
			Status:             "False",
			Type:               "Joined",
		}
		if joinedStatus {
			c.Status = "True"
		}
		meta.SetStatusCondition(&socket.Status.Conditions, c)
	}
	if phase != "" {
		socket.Status.Phase = phase
	}
	s.UpdateStatus(socket)
	return nil
}

func (s *SocketService) UpdateStatusSimple(
	phase integrationv1alpha2.Phase,
	joinedStatusCondition StatusCondition,
	appendPlug *integrationv1alpha2.Plug,
) error {
	return s.CommonUpdateStatus(phase, joinedStatusCondition, appendPlug, "", false)
}

func (s *SocketService) UpdateStatusAppendPlug(
	appendPlug *integrationv1alpha2.Plug,
) error {
	return s.CommonUpdateStatus("", "", appendPlug, "", false)
}

func (s *SocketService) UpdateStatusPhase(
	phase integrationv1alpha2.Phase,
) error {
	return s.CommonUpdateStatus(phase, "", nil, "", false)
}

func (s *SocketService) UpdateStatusJoinedCondition(
	joinedStatusCondition StatusCondition,
	message string,
) error {
	return s.CommonUpdateStatus("", joinedStatusCondition, nil, message, false)
}

func (s *SocketService) UpdateStatusJoinedConditionError(
	err error,
) error {
	return s.CommonUpdateStatus("", ErrorStatusCondition, nil, err.Error(), false)
}

func (s *SocketService) UpdateStatusRemovePlug(
	plug *integrationv1alpha2.Plug,
) error {
	socket, err := s.Get()
	if err != nil {
		return err
	}
	coupledPlugs := []integrationv1alpha2.CoupledPlug{}
	for _, coupledPlug := range socket.Status.CoupledPlugs {
		if coupledPlug.UID != plug.UID {
			coupledPlugs = append(coupledPlugs, coupledPlug)
		}
	}
	socket.Status.CoupledPlugs = coupledPlugs
	joinedCondition, err := s.GetJoinedCondition()
	if (*joinedCondition).Reason == string(SocketReadyStatusCondition) {
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
	}
	s.UpdateStatus(socket)
	return nil
}

func (s *SocketService) GetJoinedCondition() (*metav1.Condition, error) {
	socket, err := s.Get()
	if err != nil {
		return nil, err
	}
	joinedCondition := meta.FindStatusCondition(socket.Status.Conditions, "Joined")
	return joinedCondition, nil
}
