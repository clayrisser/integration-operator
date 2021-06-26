package util

import (
	"context"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

type ConfigUtil struct {
	apparatusUtil *ApparatusUtil
	client        *kubernetes.Clientset
	ctx           *context.Context
	varUtil       *VarUtil
}

func NewConfigUtil(
	ctx *context.Context,
) *ConfigUtil {
	return &ConfigUtil{
		apparatusUtil: NewApparatusUtil(ctx),
		client:        kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
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
			// TODO
			plugConfig[key] = value
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
		// TODO: get config from config mapper
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
