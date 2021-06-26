package util

import (
	"context"

	"github.com/tidwall/gjson"
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

func (u *EventUtil) PlugCreated(plug gjson.Result) error {
	return u.apparatusUtil.PlugCreated(plug)
}

func (u *EventUtil) PlugCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.apparatusUtil.PlugCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.apparatusUtil.PlugUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.apparatusUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDeleted(
	plug gjson.Result,
) error {
	return u.apparatusUtil.PlugDeleted(plug)
}

func (u *EventUtil) PlugBroken(
	plug gjson.Result,
) error {
	return u.apparatusUtil.PlugBroken(plug)
}

func (u *EventUtil) SocketCreated(socket gjson.Result) error {
	return u.apparatusUtil.SocketCreated(socket)
}

func (u *EventUtil) SocketCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.apparatusUtil.SocketCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.apparatusUtil.SocketUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return u.apparatusUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDeleted(
	socket gjson.Result,
) error {
	return u.apparatusUtil.SocketDeleted(socket)
}

func (u *EventUtil) SocketBroken(
	socket gjson.Result,
) error {
	return u.apparatusUtil.SocketBroken(socket)
}
