/**
 * File: /data.go
 * Project: integration-operator
 * File Created: 23-06-2021 22:11:01
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 14-08-2022 14:34:43
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Risser Labs LLC (c) Copyright 2021
 */

package util

import (
	"context"

	integrationv1alpha2 "gitlab.com/risserlabs/internal/integration-operator/api/v1alpha2"
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
