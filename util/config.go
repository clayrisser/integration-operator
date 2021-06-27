/*
 * File: /util/config.go
 * Project: integration-operator
 * File Created: 23-06-2021 22:09:27
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 27-06-2021 00:01:46
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

package util

import (
	"context"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ConfigUtil struct {
	apparatusUtil *ApparatusUtil
	client        *kubernetes.Clientset
	ctx           *context.Context
	dataUtil      *DataUtil
	varUtil       *VarUtil
}

func NewConfigUtil(
	ctx *context.Context,
) *ConfigUtil {
	return &ConfigUtil{
		apparatusUtil: NewApparatusUtil(ctx),
		client:        kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		dataUtil:      NewDataUtil(ctx),
		varUtil:       NewVarUtil(ctx),
	}
}

func (u *ConfigUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
) (map[string]string, error) {
	plugConfig := make(map[string]string)
	if plug.Spec.ConfigSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(plug.Namespace).Get(
			*u.ctx,
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
	if plug.Spec.ConfigMapper != nil {
		for key, value := range plug.Spec.ConfigMapper {
			result, err := u.plugLookup(plug, value)
			if err != nil {
				return nil, err
			}
			plugConfig[key] = result
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
			plugConfig[key] = string(value)
		}
	}
	if plug.Spec.Apparatus != nil {
		body, err := u.apparatusUtil.GetPlugConfig(plug)
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
	return plugConfig, nil
}

func (u *ConfigUtil) GetSocketConfig(
	socket *integrationv1alpha2.Socket,
) (map[string]string, error) {
	socketConfig := make(map[string]string)
	if socket.Spec.ConfigSecretName != "" {
		secret, err := u.client.CoreV1().Secrets(socket.Namespace).Get(
			*u.ctx,
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
	if socket.Spec.ConfigMapper != nil {
		for key, value := range socket.Spec.ConfigMapper {
			result, err := u.socketLookup(socket, value)
			if err != nil {
				return nil, err
			}
			socketConfig[key] = result
		}
	}
	if socket.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().ConfigMaps(socket.Namespace).Get(
			*u.ctx,
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
		body, err := u.apparatusUtil.GetSocketConfig(socket)
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
	return socketConfig, nil
}

func (u *ConfigUtil) plugLookup(plug *integrationv1alpha2.Plug, path string) (string, error) {
	plugLookup, err := u.buildPlugLookup(plug)
	if err != nil {
		return "", err
	}
	return plugLookup.Get(path).String(), nil
}

func (u *ConfigUtil) socketLookup(socket *integrationv1alpha2.Socket, path string) (string, error) {
	socketLookup, err := u.buildSocketLookup(socket)
	if err != nil {
		return "", err
	}
	return socketLookup.Get(path).String(), nil
}

func (u *ConfigUtil) buildPlugLookup(plug *integrationv1alpha2.Plug) (gjson.Result, error) {
	result := gjson.Parse("{}")

	resultStr, err := sjson.Set(result.String(), "resource", plug)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	dataMap, err := u.dataUtil.GetPlugData(plug)
	if err != nil {
		return result, err
	}
	resultStr, err = sjson.Set(result.String(), "data", dataMap)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	if plug.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(plug.Spec.Vars)
		if err != nil {
			return result, err
		}
		resultStr, err = sjson.Set(result.String(), "vars", varsMap)
		if err != nil {
			return result, err
		}
		result = gjson.Parse(resultStr)
	}

	return result, nil
}

func (u *ConfigUtil) buildSocketLookup(socket *integrationv1alpha2.Socket) (gjson.Result, error) {
	result := gjson.Parse("{}")

	resultStr, err := sjson.Set(result.String(), "resource", socket)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	dataMap, err := u.dataUtil.GetSocketData(socket)
	if err != nil {
		return result, err
	}
	resultStr, err = sjson.Set(result.String(), "data", dataMap)
	if err != nil {
		return result, err
	}
	result = gjson.Parse(resultStr)

	if socket.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(socket.Spec.Vars)
		if err != nil {
			return result, err
		}
		resultStr, err = sjson.Set(result.String(), "vars", varsMap)
		if err != nil {
			return result, err
		}
		result = gjson.Parse(resultStr)
	}

	return result, nil
}
