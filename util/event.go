/**
 * File: /util/event.go
 * Project: new
 * File Created: 17-10-2023 13:49:54
 * Author: Clay Risser
 * -----
 * BitSpur (c) Copyright 2021 - 2023
 *
 * Licensed under the GNU Affero General Public License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.gnu.org/licenses/agpl-3.0.en.html
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving this software without disclosing
 * the source code of your own applications.
 */

package util

import (
	"context"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
)

type EventUtil struct {
	apparatusUtil *ApparatusUtil
	resourceUtil  *ResourceUtil
}

func NewEventUtil(
	ctx *context.Context,
) *EventUtil {
	return &EventUtil{
		apparatusUtil: NewApparatusUtil(ctx),
		resourceUtil:  NewResourceUtil(ctx),
	}
}

func (u *EventUtil) PlugCreated(plug *integrationv1beta1.Plug) error {
	if err := u.apparatusUtil.PlugCreated(plug); err != nil {
		return err
	}
	return u.resourceUtil.PlugCreated(plug)
}

func (u *EventUtil) PlugCoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	if err := u.apparatusUtil.PlugCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.PlugCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugUpdated(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	if err := u.apparatusUtil.PlugUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.PlugUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDecoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	if err := u.apparatusUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) PlugDeleted(
	plug *integrationv1beta1.Plug,
) error {
	if err := u.apparatusUtil.PlugDeleted(plug); err != nil {
		return err
	}
	return u.resourceUtil.PlugDeleted(plug)
}

func (u *EventUtil) SocketCreated(socket *integrationv1beta1.Socket) error {
	if err := u.apparatusUtil.SocketCreated(socket); err != nil {
		return err
	}
	return u.resourceUtil.SocketCreated(socket)
}

func (u *EventUtil) SocketCoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	if err := u.apparatusUtil.SocketCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.SocketCoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketUpdated(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	if err := u.apparatusUtil.SocketUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.SocketUpdated(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDecoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	if err := u.apparatusUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	return u.resourceUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig)
}

func (u *EventUtil) SocketDeleted(
	socket *integrationv1beta1.Socket,
) error {
	if err := u.apparatusUtil.SocketDeleted(socket); err != nil {
		return err
	}
	return u.resourceUtil.SocketDeleted(socket)
}
