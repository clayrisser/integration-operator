package util

import (
	"context"
	"encoding/json"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
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

func (u *DataUtil) GetPlugData(plug *integrationv1alpha2.Plug) (map[string]string, error) {
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
		bSecretData, err := json.Marshal(secret.Data)
		if err != nil {
			return nil, err
		}
		secretData, err := jsonToHashMap(bSecretData)
		if err != nil {
			return nil, err
		}
		for key, value := range secretData {
			plugData[key] = value
		}
	}
	if plug.Spec.Data != nil {
		bConfig, err := json.Marshal(plug.Spec.Data)
		if err != nil {
			return nil, err
		}
		config, err := jsonToHashMap(bConfig)
		if err != nil {
			return nil, err
		}
		for key, value := range config {
			plugData[key] = value
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
			plugData[key] = value
		}
	}
	return plugData, nil
}

func (u *DataUtil) GetSocketData(socket *integrationv1alpha2.Socket) (map[string]string, error) {
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
		bSecretData, err := json.Marshal(secret.Data)
		if err != nil {
			return nil, err
		}
		secretData, err := jsonToHashMap(bSecretData)
		if err != nil {
			return nil, err
		}
		for key, value := range secretData {
			socketData[key] = value
		}
	}
	if socket.Spec.Data != nil {
		bConfig, err := json.Marshal(socket.Spec.Data)
		if err != nil {
			return nil, err
		}
		config, err := jsonToHashMap(bConfig)
		if err != nil {
			return nil, err
		}
		for key, value := range config {
			socketData[key] = value
		}
	}
	if socket.Spec.ConfigConfigMapName != "" {
		configMap, err := u.client.CoreV1().Secrets(socket.Namespace).Get(
			*u.ctx,
			socket.Spec.DataConfigMapName,
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
			socketData[key] = value
		}
	}
	return socketData, nil
}
