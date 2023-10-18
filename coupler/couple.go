/**
 * File: /coupler/couple.go
 * Project: new
 * File Created: 17-10-2023 19:02:43
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

package coupler

import (
	"context"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"gitlab.com/bitspur/rock8s/integration-operator/util"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Couple(
	client *client.Client,
	ctx *context.Context,
	req *ctrl.Request,
	plugUtil *util.PlugUtil,
	socketUtil *util.SocketUtil,
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
) (ctrl.Result, error) {
	configUtil := util.NewConfigUtil(ctx)
	if plug == nil {
		var err error
		plug, err = plugUtil.Get()
		if err != nil {
			return ctrl.Result{}, err
		}
	}
	if socket == nil {
		var err error
		socket, err = socketUtil.Get()
		if err != nil {
			if k8serrors.IsNotFound(err) {
				return plugUtil.UpdateCoupledStatus(util.SocketNotCreated, plug, nil, true)
			}
			return plugUtil.Error(err, plug)
		}
	}

	if socketUtil.CoupledPlugExists(socket.Status.CoupledPlugs, plug.UID) && plug.Status.CoupledSocket != nil {
		if err := Update(client, ctx, req, plugUtil, socketUtil, plug, socket); err != nil {
			return plugUtil.Error(err, plug)
		}
		return ctrl.Result{}, nil
	}

	if err := util.Validate(plug, socket); err != nil {
		return plugUtil.Error(err, plug)
	}

	plugConfig, err := configUtil.GetPlugConfig(plug, socket)
	if err != nil {
		return plugUtil.Error(err, plug)
	}
	socketConfig, err := configUtil.GetSocketConfig(plug, socket)
	if err != nil {
		socketUtil.Error(err, socket)
		return plugUtil.Error(err, plug)
	}

	err = CoupledPlug(plug, socket, plugConfig, socketConfig)
	if err != nil {
		return plugUtil.Error(err, plug)
	}
	err = CoupledSocket(plug, socket, plugConfig, socketConfig)
	if err != nil {
		socketUtil.Error(err, socket)
		return plugUtil.Error(err, plug)
	}

	if _, err := socketUtil.UpdateAppendCoupledPlugStatus(plug, socket, false); err != nil {
		return plugUtil.Error(err, plug)
	}
	return plugUtil.UpdateCoupledStatus(util.CouplingSucceeded, plug, socket, false)
}
