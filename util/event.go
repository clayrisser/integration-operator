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
	if err := u.apparatusUtil.PlugCreated(plug); err != nil {
		return err
	}
	return u.resourceUtil.PlugCreated(plug)
}

func (u *EventUtil) PlugCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	if err := u.apparatusUtil.PlugCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.PlugCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	if err := u.apparatusUtil.PlugUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.PlugUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	if err := u.apparatusUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDeleted(
	plug *integrationv1alpha2.Plug,
) error {
	if err := u.apparatusUtil.PlugDeleted(plug); err != nil {
		return err
	}
	return u.resourceUtil.PlugDeleted(plug)
}

func (u *EventUtil) PlugBroken(
	plug *integrationv1alpha2.Plug,
) error {
	if err := u.apparatusUtil.PlugBroken(plug); err != nil {
		return err
	}
	return u.resourceUtil.PlugBroken(plug)
}

func (u *EventUtil) SocketCreated(socket *integrationv1alpha2.Socket) error {
	if err := u.apparatusUtil.SocketCreated(socket); err != nil {
		return err
	}
	return u.resourceUtil.SocketCreated(socket)
}

func (u *EventUtil) SocketCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	if err := u.apparatusUtil.SocketCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.SocketCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	if err := u.apparatusUtil.SocketUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.SocketUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	if err := u.apparatusUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDeleted(
	socket *integrationv1alpha2.Socket,
) error {
	if err := u.apparatusUtil.SocketDeleted(socket); err != nil {
		return err
	}
	return u.resourceUtil.SocketDeleted(socket)
}

func (u *EventUtil) SocketBroken(
	socket *integrationv1alpha2.Socket,
) error {
	if err := u.apparatusUtil.SocketBroken(socket); err != nil {
		return err
	}
	return u.resourceUtil.SocketBroken(socket)
}
