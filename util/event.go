package util

import (
	"context"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
)

type EventUtil struct {
	apparatusUtil *ApparatusUtil
	resourceUtil  *ResourceUtil
}

func NewEventUtil(ctx *context.Context) *EventUtil {
	return &EventUtil{
		apparatusUtil: NewApparatusUtil(ctx),
		resourceUtil:  NewResourceUtil(ctx),
	}
}

func (u *EventUtil) PlugCreated(plug *integrationv1alpha2.Plug) error {
	return u.apparatusUtil.PlugCreated(plug)
}

func (u *EventUtil) PlugCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.apparatusUtil.PlugCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.apparatusUtil.PlugUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.apparatusUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDeleted(
	plug *integrationv1alpha2.Plug,
) error {
	return u.apparatusUtil.PlugDeleted(plug)
}

func (u *EventUtil) PlugBroken(
	plug *integrationv1alpha2.Plug,
) error {
	return u.apparatusUtil.PlugBroken(plug)
}

func (u *EventUtil) SocketCreated(socket *integrationv1alpha2.Socket) error {
	return u.apparatusUtil.SocketCreated(socket)
}

func (u *EventUtil) SocketCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.apparatusUtil.SocketCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.apparatusUtil.SocketUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	return u.apparatusUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDeleted(
	socket *integrationv1alpha2.Socket,
) error {
	return u.apparatusUtil.SocketDeleted(socket)
}

func (u *EventUtil) SocketBroken(
	socket *integrationv1alpha2.Socket,
) error {
	return u.apparatusUtil.SocketBroken(socket)
}
