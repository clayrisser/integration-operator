/*
 * File: /util/resource.go
 * Project: integration-operator
 * File Created: 23-06-2021 22:09:31
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 26-06-2021 23:55:54
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
	"bytes"
	"context"
	"encoding/json"
	"text/template"

	integrationv1alpha2 "github.com/silicon-hills/integration-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type ResourceUtil struct {
	client      *kubernetes.Clientset
	ctx         *context.Context
	kubectlUtil *KubectlUtil
}

func NewResourceUtil(ctx *context.Context) *ResourceUtil {
	return &ResourceUtil{
		client:      kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		kubectlUtil: NewKubectlUtil(ctx, &rest.Config{}),
	}
}

func (u *ResourceUtil) PlugCreated(plug *integrationv1alpha2.Plug) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.CreatedWhen)
	if err := u.ProcessResources(plug, nil, nil, nil, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.CoupledWhen)
	if err := u.ProcessResources(plug, socket, plugConfig, socketConfig, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.UpdatedWhen)
	if err := u.ProcessResources(plug, socket, plugConfig, socketConfig, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.DecoupledWhen)
	if err := u.ProcessResources(plug, socket, plugConfig, socketConfig, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugDeleted(
	plug *integrationv1alpha2.Plug,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.DeletedWhen)
	if err := u.ProcessResources(plug, nil, nil, nil, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugBroken(
	plug *integrationv1alpha2.Plug,
) error {
	resources := u.GetResources(plug.Spec.Resources, integrationv1alpha2.BrokenWhen)
	if err := u.ProcessResources(plug, nil, nil, nil, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketCreated(socket *integrationv1alpha2.Socket) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.CreatedWhen)
	if err := u.ProcessResources(nil, socket, nil, nil, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketCoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.CoupledWhen)
	if err := u.ProcessResources(plug, socket, plugConfig, socketConfig, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketUpdated(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.UpdatedWhen)
	if err := u.ProcessResources(plug, socket, plugConfig, socketConfig, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketDecoupled(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.DecoupledWhen)
	if err := u.ProcessResources(plug, socket, plugConfig, socketConfig, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketDeleted(
	socket *integrationv1alpha2.Socket,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.DeletedWhen)
	if err := u.ProcessResources(nil, socket, nil, nil, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketBroken(
	socket *integrationv1alpha2.Socket,
) error {
	resources := u.GetResources(socket.Spec.Resources, integrationv1alpha2.BrokenWhen)
	if err := u.ProcessResources(nil, socket, nil, nil, resources); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) GetResource(objRef kustomizeTypes.Target) (*unstructured.Unstructured, error) {
	const tpl = `
apiVersion: {{ .APIVersion }}
kind: {{ .Kind }}
meta:
  name: {{ .Name }}
  namespace: {{ .Namespace }}`
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return nil, err
	}
	if objRef.Group != "" && objRef.Version != "" {
		objRef.APIVersion = objRef.Group + objRef.Version
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, objRef)
	if err != nil {
		return nil, err
	}
	body := []byte(buff.String())
	return u.kubectlUtil.Get(body)
}

func (u *ResourceUtil) GetResources(
	resources []*integrationv1alpha2.Resource,
	when integrationv1alpha2.When,
) []*integrationv1alpha2.Resource {
	filteredResources := []*integrationv1alpha2.Resource{}
	if resources == nil {
		return resources
	}
	for _, resource := range resources {
		if resource.When == when {
			filteredResources = append(filteredResources, resource)
		}
	}
	return filteredResources
}

func (u *ResourceUtil) ProcessResources(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
	resources []*integrationv1alpha2.Resource,
) error {
	for _, resource := range resources {
		templatedResource, err := u.templateResource(
			plug,
			socket,
			plugConfig,
			socketConfig,
			resource.Resource,
		)
		if err != nil {
			return err
		}
		if resource.Do == integrationv1alpha2.ApplyDo {
			if err := u.kubectlUtil.Apply([]byte(templatedResource)); err != nil {
				return err
			}
		} else if resource.Do == integrationv1alpha2.DeleteDo {
			if err := u.kubectlUtil.Delete([]byte(templatedResource)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *ResourceUtil) templateResource(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
	body string,
) (string, error) {
	data, err := u.buildTemplateData(plug, socket, plugConfig, socketConfig)
	if err != nil {
		return "", err
	}
	t, err := template.New("").Parse(body)
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

func (u *ResourceUtil) buildTemplateData(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
) (map[string]interface{}, error) {
	dataMap := map[string]interface{}{}
	if plug != nil {
		dataMap["plug"] = plug
	}
	if socket != nil {
		dataMap["socket"] = socket
	}
	if plugConfig != nil {
		dataMap["plugConfig"] = plugConfig
	}
	if socketConfig != nil {
		dataMap["socketConfig"] = socketConfig
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
