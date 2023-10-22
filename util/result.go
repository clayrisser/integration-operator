/**
 * File: /util/result.go
 * Project: integration-operator
 * File Created: 19-10-2023 09:39:05
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

type ResultUtil struct {
	client   *kubernetes.Clientset
	ctx      context.Context
	config   *ConfigUtil
	resource *ResourceUtil
}

func NewResultUtil(ctx context.Context) *ResultUtil {
	return &ResultUtil{
		client:   kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		ctx:      ctx,
		config:   NewConfigUtil(ctx),
		resource: NewResourceUtil(ctx),
	}
}

func (u *ResultUtil) GetResult(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig Config,
	socketConfig Config,
) (integrationv1beta1.CoupledResult, error) {
	coupledResult := integrationv1beta1.CoupledResult{}
	plugResult, err := u.getPlugResult(plug, socket, &plugConfig, &socketConfig)
	if err != nil {
		return coupledResult, err
	}
	coupledResult.Plug = plugResult
	socketResult, err := u.getSocketResult(plug, socket, &plugConfig, &socketConfig)
	if err != nil {
		return coupledResult, err
	}
	coupledResult.Socket = socketResult
	return coupledResult, nil
}

func (u *ResultUtil) PlugTemplateResultResources(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig Config,
	socketConfig Config,
	plugResult Result,
	socketResult Result,
) error {
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	if err := u.resource.ProcessResources(
		plug,
		socket,
		&plugConfig,
		&socketConfig,
		&plugResult,
		&socketResult,
		plug.Namespace,
		plug.Spec.ResultResources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResultUtil) SocketTemplateResultResources(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig Config,
	socketConfig Config,
	plugResult Result,
	socketResult Result,
) error {
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	if err := u.resource.ProcessResources(
		plug,
		socket,
		&plugConfig,
		&socketConfig,
		&plugResult,
		&socketResult,
		socket.Namespace,
		socket.Spec.ResultResources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResultUtil) getPlugResult(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) (map[string]string, error) {
	plugResult := make(map[string]string)
	if plug.Spec.ResultSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(plug.Namespace).Get(
			u.ctx,
			plug.Spec.ResultSecretName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range secret.Data {
			plugResult[key] = string(value)
		}
	}
	if plug.Spec.Result != nil {
		for key, value := range plug.Spec.Result {
			plugResult[key] = value
		}
	}
	if plug.Spec.ResultTemplate != nil {
		for key, value := range plug.Spec.ResultTemplate {
			result, err := u.plugResultTemplateLookup(plug, value, socket)
			if err != nil {
				return nil, err
			}
			plugResult[key] = result
		}
	}
	if plug.Spec.ResultConfigMapName != "" {
		configMap, err := u.client.CoreV1().ConfigMaps(plug.Namespace).Get(
			u.ctx,
			plug.Spec.ResultConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range configMap.Data {
			plugResult[key] = string(value)
		}
	}
	plugResult, err := u.validatePlugResult(plug, socket.Spec.Interface.Result, plugResult)
	if err != nil {
		return nil, err
	}
	return plugResult, nil
}

func (u *ResultUtil) getSocketResult(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) (map[string]string, error) {
	socketResult := make(map[string]string)
	if socket.Spec.ResultSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(socket.Namespace).Get(
			u.ctx,
			socket.Spec.ResultSecretName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range secret.Data {
			socketResult[key] = string(value)
		}
	}
	if socket.Spec.Result != nil {
		for key, value := range socket.Spec.Result {
			socketResult[key] = value
		}
	}
	if socket.Spec.ResultTemplate != nil {
		for key, value := range socket.Spec.ResultTemplate {
			result, err := u.socketResultTemplateLookup(socket, value, plug)
			if err != nil {
				return nil, err
			}
			socketResult[key] = result
		}
	}
	if socket.Spec.ResultConfigMapName != "" {
		configMap, err := u.client.CoreV1().ConfigMaps(socket.Namespace).Get(
			u.ctx,
			socket.Spec.ResultConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		for key, value := range configMap.Data {
			socketResult[key] = value
		}
	}
	socketResult, err := u.validateSocketResult(socket, socketResult)
	if err != nil {
		return nil, err
	}
	return socketResult, nil
}

func (u *ResultUtil) validatePlugResult(
	plug *integrationv1beta1.Plug,
	resultInterface *integrationv1beta1.ResultInterface,
	plugResult map[string]string,
) (map[string]string, error) {
	if resultInterface == nil || resultInterface.Plug == nil {
		return plugResult, nil
	}
	validatedPlugResult := make(map[string]string)
	for propertyName, property := range resultInterface.Plug {
		if _, found := plugResult[propertyName]; found {
			validatedPlugResult[propertyName] = plugResult[propertyName]
		} else {
			if property.Required {
				return plugResult, errors.New("plug result property '" + propertyName + "' is required")
			} else if property.Default != "" {
				validatedPlugResult[propertyName] = property.Default
			}
		}
	}
	return validatedPlugResult, nil
}

func (u *ResultUtil) validateSocketResult(
	socket *integrationv1beta1.Socket,
	socketResult map[string]string,
) (map[string]string, error) {
	resultInterface := socket.Spec.Interface.Result
	if resultInterface == nil || resultInterface.Socket == nil {
		return socketResult, nil
	}
	validatedSocketResult := make(map[string]string)
	for propertyName, property := range resultInterface.Socket {
		if _, found := socketResult[propertyName]; found {
			validatedSocketResult[propertyName] = socketResult[propertyName]
		} else {
			if property.Required {
				return socketResult, errors.New("socket result property '" + propertyName + "' is required")
			} else if property.Default != "" {
				validatedSocketResult[propertyName] = property.Default
			}
		}
	}
	return validatedSocketResult, nil
}

func (u *ResultUtil) plugResultTemplateLookup(plug *integrationv1beta1.Plug, mapper string, socket *integrationv1beta1.Socket) (string, error) {
	data, err := u.buildPlugResultTemplateData(*plug, socket)
	if err != nil {
		return "", err
	}
	return Template(&data, mapper)
}

func (u *ResultUtil) socketResultTemplateLookup(
	socket *integrationv1beta1.Socket,
	mapper string,
	plug *integrationv1beta1.Plug,
) (string, error) {
	data, err := u.buildSocketResultTemplateData(*socket, plug)
	if err != nil {
		return "", err
	}
	return Template(&data, mapper)
}

func (u *ResultUtil) buildPlugResultTemplateData(
	plug integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
) (map[string]interface{}, error) {
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	dataMap := map[string]interface{}{}
	dataMap["plug"] = plug
	if socket != nil {
		dataMap["socket"] = socket
	}
	plugData, err := u.config.dataUtil.GetPlugData(&plug)
	if err != nil {
		return dataMap, err
	}
	dataMap["plugData"] = plugData
	socketData, err := u.config.dataUtil.GetSocketData(socket)
	if err != nil {
		return dataMap, err
	}
	dataMap["socketData"] = socketData
	if plug.Spec.Vars != nil {
		varsMap, err := u.config.varUtil.GetVars(plug.Namespace, plug.Spec.Vars, kubectlUtil)
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

func (u *ResultUtil) buildSocketResultTemplateData(
	socket integrationv1beta1.Socket,
	plug *integrationv1beta1.Plug,
) (map[string]interface{}, error) {
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	dataMap := map[string]interface{}{}
	dataMap["socket"] = socket
	if plug != nil {
		dataMap["plug"] = plug
	}
	socketData, err := u.config.dataUtil.GetSocketData(&socket)
	if err != nil {
		return dataMap, err
	}
	dataMap["socketData"] = socketData
	plugData, err := u.config.dataUtil.GetPlugData(plug)
	if err != nil {
		return dataMap, err
	}
	dataMap["plugData"] = plugData
	if socket.Spec.Vars != nil {
		varsMap, err := u.config.varUtil.GetVars(socket.Namespace, socket.Spec.Vars, kubectlUtil)
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
