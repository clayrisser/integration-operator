/**
 * File: /coupler/plug.go
 * Project: new
 * File Created: 17-10-2023 15:20:41
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
)

func CreatedPlug(
	plug *integrationv1beta1.Plug,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventUtil := util.NewEventUtil(&ctx)
	return eventUtil.PlugCreated(plug)
}

func DeletedPlug(
	plug *integrationv1beta1.Plug,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventUtil := util.NewEventUtil(&ctx)
	return eventUtil.PlugDeleted(plug)
}

func CoupledPlug(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig util.Config,
	socketConfig util.Config,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventUtil := util.NewEventUtil(&ctx)
	return eventUtil.PlugCoupled(plug, socket, &plugConfig, &socketConfig)
}

func UpdatedPlug(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig util.Config,
	socketConfig util.Config,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventUtil := util.NewEventUtil(&ctx)
	return eventUtil.PlugUpdated(plug, socket, &plugConfig, &socketConfig)
}

func DecoupledPlug(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig util.Config,
	socketConfig util.Config,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventUtil := util.NewEventUtil(&ctx)
	return eventUtil.PlugDecoupled(plug, socket, &plugConfig, &socketConfig)
}
