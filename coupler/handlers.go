package coupler

import (
	"github.com/silicon-hills/integration-operator/util"
	"github.com/tidwall/gjson"
)

type Config map[string]string

type Handlers struct {
	apparatusUtil *util.ApparatusUtil
}

func NewHandlers() *Handlers {
	return &Handlers{
		apparatusUtil: util.NewApparatusUtil(),
	}
}

func (h *Handlers) HandlePlugCreated(plug gjson.Result) error {
	return h.apparatusUtil.PlugCreated(plug)
}

func (h *Handlers) HandlePlugCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return h.apparatusUtil.PlugCoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return h.apparatusUtil.PlugUpdated(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return h.apparatusUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugDeleted(
	plug gjson.Result,
) error {
	return h.apparatusUtil.PlugDeleted(plug)
}

func (h *Handlers) HandlePlugBroken(
	plug gjson.Result,
) error {
	return h.apparatusUtil.PlugBroken(plug)
}

func (h *Handlers) HandleSocketCreated(socket gjson.Result) error {
	return h.apparatusUtil.SocketCreated(socket)
}

func (h *Handlers) HandleSocketCoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return h.apparatusUtil.SocketCoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketUpdated(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return h.apparatusUtil.SocketUpdated(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketDecoupled(
	plug gjson.Result,
	socket gjson.Result,
	plugConfig map[string]string,
	socketConfig map[string]string,
) error {
	return h.apparatusUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketDeleted(
	socket gjson.Result,
) error {
	return h.apparatusUtil.SocketDeleted(socket)
}

func (h *Handlers) HandleSocketBroken(
	socket gjson.Result,
) error {
	return h.apparatusUtil.SocketBroken(socket)
}
