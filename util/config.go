package util

import (
	"context"
	"encoding/json"
	"fmt"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigUtil struct {
	apparartusUtil *ApparatusUtil
	client         *client.Client
	ctx            *context.Context
}

func NewConfigUtil(
	client *client.Client,
	ctx *context.Context,
) *ConfigUtil {
	return &ConfigUtil{
		apparartusUtil: NewApparatusUtil(),
		client:         client,
		ctx:            ctx,
	}
}

func (u *ConfigUtil) GetPlugConfig(
	plug *integrationv1alpha2.Plug,
) (map[string]string, error) {
	plugConfig := make(map[string]string)
	if plug.Spec.ConfigSecretName != "" {
		secret, err := u.getPlugConfigSecret(plug)
		if err != nil {
			return nil, err
		}
		bSecretData, err := json.Marshal(secret.Data)
		if err != nil {
			return nil, err
		}
		secretData, err := u.jsonToHashMap(bSecretData)
		if err != nil {
			return nil, err
		}
		for key, value := range secretData {
			plugConfig[key] = value
		}
	}
	if plug.Spec.Config != nil {
		bConfig, err := json.Marshal(plug.Spec.Config)
		if err != nil {
			return nil, err
		}
		config, err := u.jsonToHashMap(bConfig)
		if err != nil {
			return nil, err
		}
		for key, value := range config {
			plugConfig[key] = value
		}
	}
	if plug.Spec.ConfigConfigMapName != "" {
		configMap, err := u.getPlugConfigConfigMap(plug)
		if err != nil {
			return nil, err
		}
		bConfigMapData, err := json.Marshal(configMap.Data)
		if err != nil {
			return nil, err
		}
		configMapData, err := u.jsonToHashMap(bConfigMapData)
		if err != nil {
			return nil, err
		}
		for key, value := range configMapData {
			plugConfig[key] = value
		}
	}
	if plug.Spec.Apparatus != nil {
		body, err := u.apparartusUtil.GetPlugConfig(plug)
		if err != nil {
			return nil, err
		}
		apparatusPlugConfig, err := u.jsonToHashMap(body)
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
		secret, err := u.getSocketConfigSecret(socket)
		if err != nil {
			return nil, err
		}
		bSecretData, err := json.Marshal(secret.Data)
		if err != nil {
			return nil, err
		}
		secretData, err := u.jsonToHashMap(bSecretData)
		if err != nil {
			return nil, err
		}
		for key, value := range secretData {
			socketConfig[key] = value
		}
	}
	if socket.Spec.Config != nil {
		bConfig, err := json.Marshal(socket.Spec.Config)
		if err != nil {
			return nil, err
		}
		config, err := u.jsonToHashMap(bConfig)
		if err != nil {
			return nil, err
		}
		for key, value := range config {
			socketConfig[key] = value
		}
	}
	if socket.Spec.ConfigConfigMapName != "" {
		configMap, err := u.getSocketConfigConfigMap(socket)
		if err != nil {
			return nil, err
		}
		bConfigMapData, err := json.Marshal(configMap.Data)
		if err != nil {
			return nil, err
		}
		configMapData, err := u.jsonToHashMap(bConfigMapData)
		if err != nil {
			return nil, err
		}
		for key, value := range configMapData {
			socketConfig[key] = value
		}
	}
	if socket.Spec.Apparatus != nil {
		body, err := u.apparartusUtil.GetSocketConfig(socket)
		if err != nil {
			return nil, err
		}
		apparatusSocketConfig, err := u.jsonToHashMap(body)
		if err != nil {
			return nil, err
		}
		for key, value := range apparatusSocketConfig {
			socketConfig[key] = value
		}
	}
	return socketConfig, nil
}

func (u *ConfigUtil) getSocketConfigSecret(
	socket *integrationv1alpha2.Socket,
) (*v1.Secret, error) {
	client := *u.client
	ctx := *u.ctx
	secret := &v1.Secret{}
	if err := client.Get(ctx, types.NamespacedName{
		Name:      socket.Spec.ConfigSecretName,
		Namespace: socket.Namespace,
	}, socket); err != nil {
		return nil, err
	}
	return secret, nil
}

func (u *ConfigUtil) getPlugConfigSecret(
	plug *integrationv1alpha2.Plug,
) (*v1.Secret, error) {
	client := *u.client
	ctx := *u.ctx
	secret := &v1.Secret{}
	if err := client.Get(ctx, types.NamespacedName{
		Name:      plug.Spec.ConfigSecretName,
		Namespace: plug.Namespace,
	}, plug); err != nil {
		return nil, err
	}
	return secret, nil
}

func (u *ConfigUtil) getSocketConfigConfigMap(
	socket *integrationv1alpha2.Socket,
) (*v1.ConfigMap, error) {
	client := *u.client
	ctx := *u.ctx
	configMap := &v1.ConfigMap{}
	if err := client.Get(ctx, types.NamespacedName{
		Name:      socket.Spec.ConfigSecretName,
		Namespace: socket.Namespace,
	}, socket); err != nil {
		return nil, err
	}
	return configMap, nil
}

func (u *ConfigUtil) getPlugConfigConfigMap(
	plug *integrationv1alpha2.Plug,
) (*v1.ConfigMap, error) {
	client := *u.client
	ctx := *u.ctx
	configMap := &v1.ConfigMap{}
	if err := client.Get(ctx, types.NamespacedName{
		Name:      plug.Spec.ConfigSecretName,
		Namespace: plug.Namespace,
	}, plug); err != nil {
		return nil, err
	}
	return configMap, nil
}

func (u *ConfigUtil) jsonToHashMap(body []byte) (map[string]string, error) {
	hashMap := make(map[string]string)
	var obj map[string]interface{}
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}
	for key, value := range obj {
		hashMap[key] = fmt.Sprintf("%v", value)
	}
	return hashMap, nil
}
