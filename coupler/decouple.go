/*
 * File: /coupler/decouple.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 26-06-2021 10:55:03
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

	"github.com/go-logr/logr"
	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/silicon-hills/integration-operator/util"
	"k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (c *Coupler) Decouple(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	log *logr.Logger,
	plugNamespacedName *integrationv1alpha2.NamespacedName,
) (ctrl.Result, error) {
	configUtil := util.NewConfigUtil(ctx)

	plugUtil := util.NewPlugUtil(client, ctx, req, log, plugNamespacedName, util.GlobalPlugMutex)
	plug, err := plugUtil.Get()
	if err != nil {
		return ctrl.Result{}, err
	}

	socketUtil := util.NewSocketUtil(client, ctx, req, log, &plug.Spec.Socket, util.GlobalSocketMutex)
	socket, err := socketUtil.Get()
	if err != nil {
		if !errors.IsNotFound(err) {
			return plugUtil.Error(err)
		}
	}

	var plugConfig map[string]string
	if plug.Spec.Apparatus.Endpoint != "" {
		plugConfig, err = configUtil.GetPlugConfig(plug)
		if err != nil {
			return plugUtil.Error(err)
		}
	}
	var socketConfig map[string]string
	if socket != nil && socket.Spec.Apparatus.Endpoint != "" {
		socketConfig, err = configUtil.GetSocketConfig(socket)
		if err != nil {
			return plugUtil.Error(err)
		}
	}

	if err := GlobalCoupler.DecoupledPlug(plug, socket, plugConfig, socketConfig); err != nil {
		return plugUtil.Error(err)
	}
	if err := GlobalCoupler.DecoupledSocket(plug, socket, plugConfig, socketConfig); err != nil {
		return plugUtil.Error(err)
	}

	if socket != nil {
		if _, err := socketUtil.UpdateStatusRemovePlug(plug); err != nil {
			return plugUtil.Error(err)
		}
	}

	return ctrl.Result{}, nil
}
