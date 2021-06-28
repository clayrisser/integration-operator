/*
 * File: /coupler/update.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 28-06-2021 17:50:40
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
	"errors"

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *Coupler) Update(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	plugNamespacedName *integrationv1alpha2.NamespacedName,
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	socketInterface *integrationv1alpha2.Interface,
) (ctrl.Result, error) {
	configUtil := util.NewConfigUtil(ctx)

	plugUtil := util.NewPlugUtil(client, ctx, req, log, plugNamespacedName, util.GlobalPlugMutex)
	if plug == nil {
		var err error
		plug, err = plugUtil.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	socketUtil := util.NewSocketUtil(client, ctx, req, log, &plug.Spec.Socket, util.GlobalSocketMutex)
	if socket == nil {
		var err error
		socket, err = socketUtil.Get()
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotCreatedStatusCondition, nil, false)
			}
			return plugUtil.Error(err)
		}
	}
	if !socket.Status.Ready {
		return plugUtil.UpdateStatusSimple(integrationv1alpha2.PendingPhase, util.SocketNotReadyStatusCondition, nil, false)
	}

	plugIsCoupled, err := plugUtil.IsCoupled(plug, nil)
	if err != nil {
		return plugUtil.Error(err)
	}
	if !socketUtil.CoupledPlugExists(socket.Status.CoupledPlugs, plug.UID) || !plugIsCoupled {
		return ctrl.Result{}, nil
	}

	plugInterfaceUtil := util.NewInterfaceUtil(client, ctx, req, log, &plug.Spec.Interface)
	plugInterface, err := plugInterfaceUtil.Get()
	if err != nil {
		return plugUtil.Error(err)
	}
	socketInterfaceUtil := util.NewInterfaceUtil(client, ctx, req, log, &socket.Spec.Interface)
	if socketInterface == nil {
		socketInterface, err = socketInterfaceUtil.Get()
		if err != nil {
			return plugUtil.Error(err)
		}
	}
	if plugInterface.UID != socketInterface.UID {
		return plugUtil.Error(errors.New("plug and socket interface do not match"))
	}
	// TODO: validate interface

	plugConfig, err := configUtil.GetPlugConfig(plug)
	if err != nil {
		return plugUtil.Error(err)
	}
	socketConfig, err := configUtil.GetSocketConfig(socket)
	if err != nil {
		return plugUtil.Error(err)
	}

	err = GlobalCoupler.UpdatedPlug(plug, socket, plugConfig, socketConfig)
	if err != nil {
		return plugUtil.Error(err)
	}
	err = GlobalCoupler.UpdatedSocket(plug, socket, plugConfig, socketConfig)
	if err != nil {
		result, err := plugUtil.Error(err)
		if _, err := socketUtil.UpdateErrorStatus(err, true); err != nil {
			return plugUtil.Error(err)
		}
		return result, err
	}

	return ctrl.Result{}, nil
}
