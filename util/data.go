/**
 * File: /util/data.go
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type DataUtil struct {
	client *kubernetes.Clientset
	ctx    *context.Context
}

func NewDataUtil(ctx *context.Context) *DataUtil {
	return &DataUtil{
		client: kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
	}

}

func (u *DataUtil) GetPlugData(plug *integrationv1beta1.Plug) (map[string]string, error) {
	plugData := make(map[string]string)
	if plug.Spec.DataSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(plug.Namespace).Get(
			*u.ctx,
			plug.Spec.DataSecretName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range secret.Data {
			plugData[key] = string(value)
		}
	}
	if plug.Spec.Data != nil {
		for key, value := range plug.Spec.Data {
			plugData[key] = value
		}
	}
	if plug.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().ConfigMaps(plug.Namespace).Get(
			*u.ctx,
			plug.Spec.ConfigConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range configMap.Data {
			plugData[key] = value
		}
	}
	return plugData, nil
}

func (u *DataUtil) GetSocketData(socket *integrationv1beta1.Socket) (map[string]string, error) {
	socketData := make(map[string]string)
	if socket.Spec.DataSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(socket.Namespace).Get(
			*u.ctx,
			socket.Spec.DataSecretName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range secret.Data {
			socketData[key] = string(value)
		}
	}
	if socket.Spec.Data != nil {
		for key, value := range socket.Spec.Data {
			socketData[key] = value
		}
	}
	if socket.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().ConfigMaps(socket.Namespace).Get(
			*u.ctx,
			socket.Spec.DataConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range configMap.Data {
			socketData[key] = value
		}
	}
	return socketData, nil
}
