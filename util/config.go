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
		if bSecretData, err := json.Marshal(secret.Data); err != nil {
			if secretData, err := u.jsonToHashMap(bSecretData); err != nil {
				for key, value := range secretData {
					plugConfig[key] = value
				}
			}
		}
	}
	if plug.Spec.Config != nil {
		if bConfig, err := json.Marshal(plug.Spec.Config); err != nil {
			if config, err := u.jsonToHashMap(bConfig); err != nil {
				for key, value := range config {
					plugConfig[key] = value
				}
			}
		}
	}
	if plug.Spec.ConfigConfigMapName != "" {
		configMap, err := u.getPlugConfigConfigMap(plug)
		if err != nil {
			return nil, err
		}
		if bConfigMapData, err := json.Marshal(configMap.Data); err != nil {
			if configMapData, err := u.jsonToHashMap(bConfigMapData); err != nil {
				for key, value := range configMapData {
					plugConfig[key] = value
				}
			}
		}
	}
	if plug.Spec.Apparatus != nil {
		body, err := u.apparartusUtil.GetPlugConfig(plug)
		if err != nil {
			return nil, err
		}
		plugConfig, err = u.jsonToHashMap(body)
		if err != nil {
			return nil, err
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
		if bSecretData, err := json.Marshal(secret.Data); err != nil {
			if secretData, err := u.jsonToHashMap(bSecretData); err != nil {
				for key, value := range secretData {
					socketConfig[key] = value
				}
			}
		}
	}
	if socket.Spec.Config != nil {
		if bConfig, err := json.Marshal(socket.Spec.Config); err != nil {
			if config, err := u.jsonToHashMap(bConfig); err != nil {
				for key, value := range config {
					socketConfig[key] = value
				}
			}
		}
	}
	if socket.Spec.ConfigConfigMapName != "" {
		configMap, err := u.getSocketConfigConfigMap(socket)
		if err != nil {
			return nil, err
		}
		if bConfigMapData, err := json.Marshal(configMap.Data); err != nil {
			if configMapData, err := u.jsonToHashMap(bConfigMapData); err != nil {
				for key, value := range configMapData {
					socketConfig[key] = value
				}
			}
		}
	}
	if socket.Spec.Apparatus != nil {
		body, err := u.apparartusUtil.GetSocketConfig(socket)
		if err != nil {
			return nil, err
		}
		socketConfig, err = u.jsonToHashMap(body)
		if err != nil {
			return nil, err
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
