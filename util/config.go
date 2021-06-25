package util

import (
	"context"
	"encoding/json"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigUtil struct {
	client        *kubernetes.Clientset
	ctx           *context.Context
	apparatusUtil *ApparatusUtil
}

func NewConfigUtil(
	client *client.Client,
	ctx *context.Context,
) *ConfigUtil {
	return &ConfigUtil{
		apparatusUtil: NewApparatusUtil(ctx),
		client:        kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
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
		bSecretData, err := json.Marshal(secret.Data)
		if err != nil {
			return nil, err
		}
		secretData, err := jsonToHashMap(bSecretData)
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
		config, err := jsonToHashMap(bConfig)
		if err != nil {
			return nil, err
		}
		for key, value := range config {
			plugConfig[key] = value
		}
	}
	if plug.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().Secrets(plug.Namespace).Get(
			*u.ctx,
			plug.Spec.ConfigConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		bConfigMapData, err := json.Marshal(configMap.Data)
		if err != nil {
			return nil, err
		}
		configMapData, err := jsonToHashMap(bConfigMapData)
		if err != nil {
			return nil, err
		}
		for key, value := range configMapData {
			plugConfig[key] = value
		}
	}
	if plug.Spec.Apparatus != nil {
		body, err := u.apparatusUtil.GetPlugConfig(plug)
		if err != nil {
			return nil, err
		}
		apparatusPlugConfig, err := jsonToHashMap(body)
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
		bSecretData, err := json.Marshal(secret.Data)
		if err != nil {
			return nil, err
		}
		secretData, err := jsonToHashMap(bSecretData)
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
		config, err := jsonToHashMap(bConfig)
		if err != nil {
			return nil, err
		}
		for key, value := range config {
			socketConfig[key] = value
		}
	}
	if socket.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().Secrets(socket.Namespace).Get(
			*u.ctx,
			socket.Spec.ConfigConfigMapName,
			metav1.GetOptions{},
		)
		if err != nil {
			return nil, err
		}
		bConfigMapData, err := json.Marshal(configMap.Data)
		if err != nil {
			return nil, err
		}
		configMapData, err := jsonToHashMap(bConfigMapData)
		if err != nil {
			return nil, err
		}
		for key, value := range configMapData {
			socketConfig[key] = value
		}
	}
	if socket.Spec.Apparatus != nil {
		body, err := u.apparatusUtil.GetSocketConfig(socket)
		if err != nil {
			return nil, err
		}
		apparatusSocketConfig, err := jsonToHashMap(body)
		if err != nil {
			return nil, err
		}
		for key, value := range apparatusSocketConfig {
			socketConfig[key] = value
		}
	}
	return socketConfig, nil
}
