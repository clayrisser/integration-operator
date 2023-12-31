/**
 * File: /coupler/update.go
 * Project: integration-operator
 * File Created: 17-10-2023 16:09:52
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
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Update(
	client *client.Client,
	ctx context.Context,
	req *ctrl.Request,
	plugUtil *util.PlugUtil,
	socketUtil *util.SocketUtil,
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	recorder record.EventRecorder,
) error {
	configUtil := util.NewConfigUtil(ctx)
	if plug == nil {
		var err error
		plug, err = plugUtil.Get()
		if err != nil {
			return err
		}
	}
	if socket == nil {
		var err error
		socket, err = socketUtil.Get()
		if err != nil {
			return err
		}
	}

	if !socketUtil.CoupledPlugExists(socket.Status.CoupledPlugs, plug.UID) || plug.Status.CoupledSocket == nil {
		return nil
	}

	plugConfig, err := configUtil.GetPlugConfig(plug, socket)
	if err != nil {
		return err
	}
	socketConfig, err := configUtil.GetSocketConfig(plug, socket)
	if err != nil {
		socketUtil.Error(err, socket)
		return err
	}

	if err = UpdatedPlug(plug, socket, plugConfig, socketConfig, recorder); err != nil {
		return err
	}
	if err = UpdatedSocket(plug, socket, plugConfig, socketConfig, recorder); err != nil {
		socketUtil.Error(err, socket)
		return err
	}

	return nil
}
