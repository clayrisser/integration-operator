/**
 * File: /config.go
 * Project: integration-operator
 * File Created: 23-06-2021 22:09:27
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 11:49:19
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"text/template"

	"github.com/Masterminds/sprig"
	integrationv1alpha2 "gitlab.com/bitspur/rock8s/integration-operator/api/v1alpha2"
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
		ctx:           ctx,
		dataUtil:      NewDataUtil(ctx),
		varUtil:       NewVarUtil(ctx),
	}
}

func (u *ConfigUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
	plugInterface *integrationv1alpha2.Interface,
	socket *integrationv1alpha2.Socket,
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
			result, err := u.plugLookup(plug, value, socket)
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
	plugConfig, err := u.ValidatePlugConfig(plug, plugInterface, plugConfig)
	if err != nil {
		return nil, err
	}
	return plugConfig, nil
}

func (u *ConfigUtil) GetSocketConfig(
	socket *integrationv1alpha2.Socket,
	socketInterface *integrationv1alpha2.Interface,
	plug *integrationv1alpha2.Plug,
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
			result, err := u.socketLookup(socket, value, plug)
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
	socketConfig, err := u.ValidateSocketConfig(socket, socketInterface, socketConfig)
	if err != nil {
		return nil, err
	}
	return socketConfig, nil
}

func (u *ConfigUtil) ValidatePlugConfig(
	plug *integrationv1alpha2.Plug,
	plugInterface *integrationv1alpha2.Interface,
	plugConfig map[string]string,
) (map[string]string, error) {
	if plugInterface == nil {
		return plugConfig, nil
	}
	var schema *integrationv1alpha2.InterfaceSpecSchema
	for _, s := range plugInterface.Spec.Schemas {
		if u.validVersion(plug.Spec.InterfaceVersions, s.Version) {
			schema = s
		}
	}
	if schema == nil {
		return plugConfig, errors.New("schema version " + schema.Version + " not supported for plug " + plug.Name)
	}
	for propertyName, property := range schema.PlugDefinition.Properties {
		if _, found := plugConfig[propertyName]; !found {
			if property.Required {
				return plugConfig, errors.New("config property " + propertyName + " is required for plug " + plug.Name)
			} else if property.Default != "" {
				plugConfig[propertyName] = property.Default
			}
		}
	}
	return plugConfig, nil
}

func (u *ConfigUtil) ValidateSocketConfig(
	socket *integrationv1alpha2.Socket,
	socketInterface *integrationv1alpha2.Interface,
	socketConfig map[string]string,
) (map[string]string, error) {
	if socketInterface == nil {
		return socketConfig, nil
	}
	var schema *integrationv1alpha2.InterfaceSpecSchema
	for _, s := range socketInterface.Spec.Schemas {
		if u.validVersion(socket.Spec.InterfaceVersions, s.Version) {
			schema = s
		}
	}
	if schema == nil {
		return socketConfig, errors.New("schema version " + schema.Version + " not supported for socket " + socket.Name)
	}
	for propertyName, property := range schema.SocketDefinition.Properties {
		if _, found := socketConfig[propertyName]; !found {
			if property.Required {
				return socketConfig, errors.New("config property " + propertyName + " is required for socket " + socket.Name)
			} else if property.Default != "" {
				socketConfig[propertyName] = property.Default
			}
		}
	}
	return socketConfig, nil
}

func (u *ConfigUtil) validVersion(versions string, version string) bool {
	if versions == "*" {
		return true
	}
	return versions == version
}

func (u *ConfigUtil) plugLookup(plug *integrationv1alpha2.Plug, mapper string, socket *integrationv1alpha2.Socket) (string, error) {
	data, err := u.buildPlugTemplateData(plug, socket)
	if err != nil {
		return "", err
	}
	return u.templateConfigMapper(&data, mapper)
}

func (u *ConfigUtil) socketLookup(socket *integrationv1alpha2.Socket, mapper string, plug *integrationv1alpha2.Plug) (string, error) {
	data, err := u.buildSocketTemplateData(socket, plug)
	if err != nil {
		return "", err
	}
	return u.templateConfigMapper(&data, mapper)
}

func (u *ConfigUtil) buildPlugTemplateData(plug *integrationv1alpha2.Plug, socket *integrationv1alpha2.Socket) (map[string]interface{}, error) {
	dataMap := map[string]interface{}{}
	if plug != nil {
		dataMap["plug"] = plug
	}
	if socket != nil {
		dataMap["socket"] = socket
	}
	plugData, err := u.dataUtil.GetPlugData(plug)
	if err != nil {
		return dataMap, err
	}
	if dataMap != nil {
		dataMap["plugData"] = plugData
	}
	socketData, err := u.dataUtil.GetSocketData(socket)
	if err != nil {
		return dataMap, err
	}
	if dataMap != nil {
		dataMap["socketData"] = socketData
	}
	if plug.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(plug.Namespace, plug.Spec.Vars)
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

func (u *ConfigUtil) buildSocketTemplateData(socket *integrationv1alpha2.Socket, plug *integrationv1alpha2.Plug) (map[string]interface{}, error) {
	dataMap := map[string]interface{}{}
	if socket != nil {
		dataMap["socket"] = socket
	}
	if plug != nil {
		dataMap["plug"] = plug
	}
	socketData, err := u.dataUtil.GetSocketData(socket)
	if err != nil {
		return dataMap, err
	}
	if dataMap != nil {
		dataMap["socketData"] = socketData
	}
	plugData, err := u.dataUtil.GetPlugData(plug)
	if err != nil {
		return dataMap, err
	}
	if dataMap != nil {
		dataMap["plugData"] = plugData
	}
	if socket.Spec.Vars != nil {
		varsMap, err := u.varUtil.GetVars(socket.Namespace, socket.Spec.Vars)
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

func (u *ConfigUtil) templateConfigMapper(
	data *map[string]interface{},
	mapper string,
) (string, error) {
	t, err := template.New("").Funcs(sprig.TxtFuncMap()).Delims("{%", "%}").Parse(mapper)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}
