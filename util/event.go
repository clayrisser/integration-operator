/**
 * File: /util/event.go
 * Project: integration-operator
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
	"fmt"

	"github.com/go-logr/logr"
	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type EventUtil struct {
	apparatusUtil *ApparatusUtil
	resourceUtil  *ResourceUtil
	logger        logr.Logger
}

func NewEventUtil(
	ctx context.Context,
) *EventUtil {
	return &EventUtil{
		apparatusUtil: NewApparatusUtil(ctx),
		resourceUtil:  NewResourceUtil(ctx),
		logger:        log.FromContext(ctx),
	}
}

func (u *EventUtil) PlugCreated(
	plug *integrationv1beta1.Plug,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("plug %s/%s created", plug.Name, plug.Namespace))
	if err := u.apparatusUtil.PlugCreated(plug); err != nil {
		return err
	}
	if err := u.resourceUtil.PlugCreated(plug); err != nil {
		return err
	}
	recorder.Event(plug, "Normal", "PlugCreated", fmt.Sprintf("plug %s/%s created", plug.Name, plug.Namespace))
	return nil
}

func (u *EventUtil) PlugCoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("plug %s/%s coupled", plug.Name, plug.Namespace))
	if err := u.apparatusUtil.PlugCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	if err := u.resourceUtil.PlugCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	recorder.Event(plug, "Normal", "PlugCoupled", fmt.Sprintf("plug %s/%s coupled", plug.Name, plug.Namespace))
	return nil
}

func (u *EventUtil) PlugUpdated(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("plug %s/%s updated", plug.Name, plug.Namespace))
	if err := u.apparatusUtil.PlugUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	if err := u.resourceUtil.PlugUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	recorder.Event(plug, "Normal", "PlugUpdated", fmt.Sprintf("plug %s/%s updated", plug.Name, plug.Namespace))
	return nil
}

func (u *EventUtil) PlugDecoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("plug %s/%s decoupled", plug.Name, plug.Namespace))
	if err := u.apparatusUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	if err := u.resourceUtil.PlugDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	recorder.Event(plug, "Normal", "PlugDecoupled", fmt.Sprintf("plug %s/%s decoupled", plug.Name, plug.Namespace))
	return nil
}

func (u *EventUtil) PlugDeleted(
	plug *integrationv1beta1.Plug,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("plug %s/%s deleted", plug.Name, plug.Namespace))
	if err := u.apparatusUtil.PlugDeleted(plug); err != nil {
		return err
	}
	if err := u.resourceUtil.PlugDeleted(plug); err != nil {
		return err
	}
	recorder.Event(plug, "Normal", "PlugDeleted", fmt.Sprintf("plug %s/%s deleted", plug.Name, plug.Namespace))
	return nil
}

func (u *EventUtil) SocketCreated(
	socket *integrationv1beta1.Socket,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("socket %s/%s created", socket.Name, socket.Namespace))
	if err := u.apparatusUtil.SocketCreated(socket); err != nil {
		return err
	}
	if err := u.resourceUtil.SocketCreated(socket); err != nil {
		return err
	}
	recorder.Event(socket, "Normal", "SocketCreated", fmt.Sprintf("socket %s/%s created", socket.Name, socket.Namespace))
	return nil
}

func (u *EventUtil) SocketCoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("socket %s/%s coupled", socket.Name, socket.Namespace))
	if err := u.apparatusUtil.SocketCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	if err := u.resourceUtil.SocketCoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	recorder.Event(socket, "Normal", "SocketCoupled", fmt.Sprintf("socket %s/%s coupled", socket.Name, socket.Namespace))
	return nil
}

func (u *EventUtil) SocketUpdated(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("socket %s/%s updated", socket.Name, socket.Namespace))
	if err := u.apparatusUtil.SocketUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	if err := u.resourceUtil.SocketUpdated(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	recorder.Event(socket, "Normal", "SocketUpdated", fmt.Sprintf("socket %s/%s updated", socket.Name, socket.Namespace))
	return nil
}

func (u *EventUtil) SocketDecoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("socket %s/%s decoupled", socket.Name, socket.Namespace))
	if err := u.apparatusUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	if err := u.resourceUtil.SocketDecoupled(plug, socket, plugConfig, socketConfig); err != nil {
		return err
	}
	recorder.Event(socket, "Normal", "SocketDecoupled", fmt.Sprintf("socket %s/%s decoupled", socket.Name, socket.Namespace))
	return nil
}

func (u *EventUtil) SocketDeleted(
	socket *integrationv1beta1.Socket,
	recorder record.EventRecorder,
) error {
	u.logger.Info(fmt.Sprintf("socket %s/%s deleted", socket.Name, socket.Namespace))
	if err := u.apparatusUtil.SocketDeleted(socket); err != nil {
		return err
	}
	if err := u.resourceUtil.SocketDeleted(socket); err != nil {
		return err
	}
	recorder.Event(socket, "Normal", "SocketDeleted", fmt.Sprintf("socket %s/%s deleted", socket.Name, socket.Namespace))
	return nil
}
