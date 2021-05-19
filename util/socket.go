package util

import (
	"context"

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
