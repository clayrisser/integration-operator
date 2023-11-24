/**
 * File: /util/config.go
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
	"encoding/json"
	"errors"

	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ConfigUtil struct {
	apparatusUtil *ApparatusUtil
	client        *kubernetes.Clientset
	ctx           context.Context
	dataUtil      *DataUtil
	varUtil       *VarUtil
}

func NewConfigUtil(
	ctx context.Context,
) *ConfigUtil {
	return &ConfigUtil{
		apparatusUtil: NewApparatusUtil(ctx),
		client:        kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		ctx:           ctx,
		dataUtil:      NewDataUtil(ctx),
		varUtil:       NewVarUtil(ctx),
	}
}

func (u *ConfigUtil) GetPlugConfig(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
) (map[string]string, error) {
	plugConfig := make(map[string]string)
	if plug.Spec.ConfigSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(plug.Namespace).Get(
			u.ctx,
			plug.Spec.ConfigSecretName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range secret.Data {
			plugConfig[key] = string(value)
		}
	}
	if plug.Spec.Config != nil {
		for key, value := range plug.Spec.Config {
			plugConfig[key] = value
		}
	}
	if plug.Spec.ConfigTemplate != nil {
		for key, value := range plug.Spec.ConfigTemplate {
			result, err := u.plugConfigTemplateLookup(plug, value, socket)
			if err != nil {
				return nil, err
			}
			plugConfig[key] = result
		}
	}
	if plug.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().ConfigMaps(plug.Namespace).Get(
			u.ctx,
			plug.Spec.ConfigConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range configMap.Data {
			plugConfig[key] = string(value)
		}
	}
	if plug.Spec.Apparatus != nil {
		body, err := u.apparatusUtil.GetPlugConfig(plug, socket)
		if err != nil {
			return nil, err
		}
		apparatusPlugConfig, err := JsonToHashMap(body)
		if err != nil {
			return nil, err
		}
		for key, value := range apparatusPlugConfig {
			plugConfig[key] = value
		}
	}
	if socket.Spec.Interface == nil {
		return plugConfig, nil
	}
	plugConfig, err := u.ValidatePlugConfig(plug, socket.Spec.Interface.Config, plugConfig)
	if err != nil {
		return nil, err
	}
	return plugConfig, nil
}

func (u *ConfigUtil) GetSocketConfig(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
) (map[string]string, error) {
	socketConfig := make(map[string]string)
	if socket.Spec.ConfigSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(socket.Namespace).Get(
			u.ctx,
			socket.Spec.ConfigSecretName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range secret.Data {
			socketConfig[key] = string(value)
		}
	}
	if socket.Spec.Config != nil {
		for key, value := range socket.Spec.Config {
			socketConfig[key] = value
		}
	}
	if socket.Spec.ConfigTemplate != nil {
		for key, value := range socket.Spec.ConfigTemplate {
			result, err := u.socketConfigTemplateLookup(socket, value, plug)
			if err != nil {
				return nil, err
			}
			socketConfig[key] = result
		}
	}
	if socket.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().ConfigMaps(socket.Namespace).Get(
			u.ctx,
			socket.Spec.ConfigConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range configMap.Data {
			socketConfig[key] = value
		}
	}
	if socket.Spec.Apparatus != nil {
		body, err := u.apparatusUtil.GetSocketConfig(socket, plug)
		if err != nil {
			return nil, err
		}
		apparatusSocketConfig, err := JsonToHashMap(body)
		if err != nil {
			return nil, err
		}
		for key, value := range apparatusSocketConfig {
			socketConfig[key] = value
		}
	}
	socketConfig, err := u.ValidateSocketConfig(socket, socketConfig)
	if err != nil {
		return nil, err
	}
	return socketConfig, nil
}

func (u *ConfigUtil) ValidatePlugConfig(
	plug *integrationv1beta1.Plug,
	configInterface *integrationv1beta1.ConfigInterface,
	plugConfig map[string]string,
) (map[string]string, error) {
	if configInterface == nil {
		return plugConfig, nil
	}
	validatedPlugConfig := make(map[string]string)
	for propertyName, property := range configInterface.Plug {
		if value, found := plugConfig[propertyName]; found && value != "" {
			validatedPlugConfig[propertyName] = plugConfig[propertyName]
		} else {
			if property.Required {
				return plugConfig, errors.New("plug config property '" + propertyName + "' is required")
			} else if property.Default != "" {
				validatedPlugConfig[propertyName] = property.Default
			}
		}
	}
	return validatedPlugConfig, nil
}

func (u *ConfigUtil) ValidateSocketConfig(
	socket *integrationv1beta1.Socket,
	socketConfig map[string]string,
) (map[string]string, error) {
	if socket.Spec.Interface == nil {
		return socketConfig, nil
	}
	configInterface := socket.Spec.Interface.Config
	if configInterface == nil {
		return socketConfig, nil
	}
	validatedSocketConfig := make(map[string]string)
	for propertyName, property := range configInterface.Socket {
		if _, found := socketConfig[propertyName]; found {
			validatedSocketConfig[propertyName] = socketConfig[propertyName]
		} else {
			if property.Required {
				return socketConfig, errors.New("socket config property '" + propertyName + "' is required")
			} else if property.Default != "" {
				validatedSocketConfig[propertyName] = property.Default
			}
		}
	}
	return validatedSocketConfig, nil
}

func (u *ConfigUtil) plugConfigTemplateLookup(plug *integrationv1beta1.Plug, configTemplate string, socket *integrationv1beta1.Socket) (string, error) {
	data, err := u.buildPlugConfigTemplateData(*plug, socket)
	if err != nil {
		return "", err
	}
	return Template(&data, configTemplate)
}

func (u *ConfigUtil) socketConfigTemplateLookup(
	socket *integrationv1beta1.Socket,
	configTemplate string,
	plug *integrationv1beta1.Plug,
) (string, error) {
	data, err := u.buildSocketConfigTemplateData(*socket, plug)
	if err != nil {
		return "", err
	}
	return Template(&data, configTemplate)
}

func (u *ConfigUtil) buildPlugConfigTemplateData(
	plug integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
) (map[string]interface{}, error) {
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	dataMap := map[string]interface{}{}
	dataMap["plug"] = plug
	if socket != nil {
		dataMap["socket"] = socket
	}
	plugData, err := u.dataUtil.GetPlugData(&plug)
	if err != nil {
		return dataMap, err
	}
	dataMap["plugData"] = plugData
	socketData, err := u.dataUtil.GetSocketData(socket)
	if err != nil {
		return dataMap, err
	}
	dataMap["socketData"] = socketData
	if plug.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(plug.Namespace, plug.Spec.Vars, kubectlUtil, &plug, socket)
		if err != nil {
			return dataMap, err
		}
		dataMap["vars"] = varsMap
	}
	bData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(bData, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (u *ConfigUtil) buildSocketConfigTemplateData(
	socket integrationv1beta1.Socket,
	plug *integrationv1beta1.Plug,
) (map[string]interface{}, error) {
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	dataMap := map[string]interface{}{}
	dataMap["socket"] = socket
	if plug != nil {
		dataMap["plug"] = plug
	}
	socketData, err := u.dataUtil.GetSocketData(&socket)
	if err != nil {
		return dataMap, err
	}
	dataMap["socketData"] = socketData
	plugData, err := u.dataUtil.GetPlugData(plug)
	if err != nil {
		return dataMap, err
	}
	dataMap["plugData"] = plugData
	if socket.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(socket.Namespace, socket.Spec.Vars, kubectlUtil, plug, &socket)
		if err != nil {
			return dataMap, err
		}
		dataMap["vars"] = varsMap
	}
	bData, err := json.Marshal(dataMap)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(bData, &data); err != nil {
		return nil, err
	}
	return data, nil
}
