/**
 * File: /coupler/handlers.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 26-06-2021 10:55:12
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Silicon Hills LLC (c) Copyright 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package coupler

import (
	"context"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
)

type Config map[string]string

type Handlers struct {
}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) HandlePlugCreated(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.PlugCreated(plug)
}

func (h *Handlers) HandlePlugCoupled(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.PlugCoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugUpdated(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.PlugUpdated(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugDecoupled(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandlePlugDeleted(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.PlugDeleted(plug)
}

func (h *Handlers) HandlePlugBroken(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.PlugBroken(plug)
}

func (h *Handlers) HandleSocketCreated(
	ctx *context.Context,
	socket *integrationv1alpha2.Socket,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.SocketCreated(socket)
}

func (h *Handlers) HandleSocketCoupled(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.SocketCoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketUpdated(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.SocketUpdated(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketDecoupled(
	ctx *context.Context,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig)
}

func (h *Handlers) HandleSocketDeleted(
	ctx *context.Context,
	socket *integrationv1alpha2.Socket,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.SocketDeleted(socket)
}

func (h *Handlers) HandleSocketBroken(
	ctx *context.Context,
	socket *integrationv1alpha2.Socket,
) error {
	eventUtil := util.NewEventUtil(ctx)
	return eventUtil.SocketBroken(socket)
}
