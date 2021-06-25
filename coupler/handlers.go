package coupler

import (
	"context"

	"github.com/silicon-hills/integration-operator/util"
	"github.com/tidwall/gjson"
)

type Config map[string]string

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HandlePlugCreated(
	ctx *context.Context,
	plug gjson.Result,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.PlugCreated(plug)
}

func (h *Handlers) HandlePlugCoupled(
	ctx *context.Context,
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.PlugCoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugUpdated(
	ctx *context.Context,
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.PlugUpdated(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugDecoupled(
	ctx *context.Context,
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugDeleted(
	ctx *context.Context,
	plug gjson.Result,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.PlugDeleted(plug)
}

func (h *Handlers) HandlePlugBroken(
	ctx *context.Context,
	plug gjson.Result,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.PlugBroken(plug)
}

func (h *Handlers) HandleSocketCreated(
	ctx *context.Context,
	socket gjson.Result,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.SocketCreated(socket)
}

func (h *Handlers) HandleSocketCoupled(
	ctx *context.Context,
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.SocketCoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketUpdated(
	ctx *context.Context,
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.SocketUpdated(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketDecoupled(
	ctx *context.Context,
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketDeleted(
	ctx *context.Context,
	socket gjson.Result,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.SocketDeleted(socket)
}

func (h *Handlers) HandleSocketBroken(
	ctx *context.Context,
	socket gjson.Result,
) error {
	apparatusUtil := util.NewApparatusUtil(ctx)
	return apparatusUtil.SocketBroken(socket)
}
