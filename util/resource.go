/**
 * File: /resource.go
 * Project: integration-operator
 * File Created: 23-07-2021 17:13:09
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
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	integrationv1alpha2 "gitlab.com/bitspur/rock8s/integration-operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type ResourceUtil struct {
	client      *kubernetes.Clientset
	kubectlUtil *KubectlUtil
}

func NewResourceUtil(ctx *context.Context) *ResourceUtil {
	return &ResourceUtil{
		client:      kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		kubectlUtil: NewKubectlUtil(ctx),
	}
}

func (u *ResourceUtil) PlugCreated(plug *integrationv1alpha2.Plug) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1alpha2.CreatedWhen)
	if err := u.ProcessResources(
		plug,
		nil,
		nil,
		nil,
		plug.Namespace,
		resources,
	); err != nil {
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
	resources := u.filterResources(plug.Spec.Resources, integrationv1alpha2.CoupledWhen)
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		plug.Namespace,
		resources,
	); err != nil {
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
	resources := u.filterResources(plug.Spec.Resources, integrationv1alpha2.UpdatedWhen)
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		plug.Namespace,
		resources,
	); err != nil {
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
	resources := u.filterResources(plug.Spec.Resources, integrationv1alpha2.DecoupledWhen)
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		plug.Namespace,
		resources,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugDeleted(
	plug *integrationv1alpha2.Plug,
) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1alpha2.DeletedWhen)
	if err := u.ProcessResources(
		plug,
		nil,
		nil,
		nil,
		plug.Namespace,
		resources,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugBroken(
	plug *integrationv1alpha2.Plug,
) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1alpha2.BrokenWhen)
	if err := u.ProcessResources(
		plug,
		nil,
		nil,
		nil,
		plug.Namespace,
		resources,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketCreated(socket *integrationv1alpha2.Socket) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1alpha2.CreatedWhen)
	if err := u.ProcessResources(
		nil,
		socket,
		nil,
		nil,
		socket.Namespace,
		resources,
	); err != nil {
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
	resources := u.filterResources(socket.Spec.Resources, integrationv1alpha2.CoupledWhen)
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		socket.Namespace,
		resources,
	); err != nil {
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
	resources := u.filterResources(socket.Spec.Resources, integrationv1alpha2.UpdatedWhen)
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		socket.Namespace,
		resources,
	); err != nil {
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
	resources := u.filterResources(socket.Spec.Resources, integrationv1alpha2.DecoupledWhen)
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		socket.Namespace,
		resources,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketDeleted(
	socket *integrationv1alpha2.Socket,
) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1alpha2.DeletedWhen)
	if err := u.ProcessResources(
		nil,
		socket,
		nil,
		nil,
		socket.Namespace,
		resources,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketBroken(
	socket *integrationv1alpha2.Socket,
) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1alpha2.BrokenWhen)
	if err := u.ProcessResources(
		nil,
		socket,
		nil,
		nil,
		socket.Namespace,
		resources,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) GetResource(namespace string, objRef kustomizeTypes.Target) (*unstructured.Unstructured, error) {
	const tpl = `
apiVersion: {{ .APIVersion }}
kind: {{ .Kind }}
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}`
	t, err := template.New("").Parse(tpl)
	if err != nil {
		return nil, err
	}
	if objRef.Namespace == "" {
		objRef.Namespace = namespace
	} else if objRef.Namespace != namespace {
		return nil, errors.New("var objRef namespace " + objRef.Namespace + " must be " + namespace)
	}
	if objRef.Group != "" && objRef.Version != "" {
		objRef.APIVersion = objRef.Group + objRef.Version
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, objRef)
	if err != nil {
		return nil, err
	}
	body := buff.Bytes()
	return u.kubectlUtil.Get(body)
}

func (u *ResourceUtil) ProcessResources(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
	namespace string,
	resources []*integrationv1alpha2.Resource,
) error {
	for _, resource := range resources {
		templatedResource, err := u.templateResource(
			plug,
			socket,
			plugConfig,
			socketConfig,
			namespace,
			resource.Resource,
		)
		if err != nil {
			return err
		}
		do := resource.Do
		if do == "" {
			do = integrationv1alpha2.ApplyDo
		}
		if resource.Do == integrationv1alpha2.ApplyDo {
			if err := u.kubectlUtil.Apply([]byte(templatedResource)); err != nil {
				return err
			}
		} else if resource.Do == integrationv1alpha2.DeleteDo {
			if err := u.kubectlUtil.Delete([]byte(templatedResource)); err != nil {
				return err
			}
		} else if resource.Do == integrationv1alpha2.RecreateDo {
			u.kubectlUtil.Delete([]byte(templatedResource))
			if err := u.kubectlUtil.Apply([]byte(templatedResource)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (u *ResourceUtil) filterResources(
	resources []*integrationv1alpha2.Resource,
	when integrationv1alpha2.When,
) []*integrationv1alpha2.Resource {
	filteredResources := []*integrationv1alpha2.Resource{}
	if resources == nil {
		return resources
	}
	if when == "" {
		when = integrationv1alpha2.CoupledWhen
	}
	for _, resource := range resources {
		if WhenInWhenSlice(when, resource.When) {
			filteredResources = append(filteredResources, resource)
		}
	}
	return filteredResources
}

func (u *ResourceUtil) templateResource(
	plug *integrationv1alpha2.Plug,
	socket *integrationv1alpha2.Socket,
	plugConfig *map[string]string,
	socketConfig *map[string]string,
	namespace string,
	body string,
) (string, error) {
	data, err := u.buildTemplateData(plug, socket, plugConfig, socketConfig)
	if err != nil {
		return "", err
	}
	t, err := template.New("").Funcs(sprig.TxtFuncMap()).Delims("{%", "%}").Parse(body)
	if err != nil {
		return "", err
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, data)
	if err != nil {
		return "", err
	}
	if buff.String() == "" {
		return "", errors.New("failed to parse template in namespace '" + namespace + "'")
	}
	obj := unstructured.Unstructured{}
	if _, _, err := decUnstructured.Decode(buff.Bytes(), nil, &obj); err != nil {
		return "", err
	}
	bJson, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	resultGjson := gjson.Parse(string(bJson)).Get("Object")
	parsedNamespace := resultGjson.Get("metadata.namespace").String()
	result := resultGjson.String()
	if parsedNamespace == "" {
		result, err = sjson.Set(result, "metadata.namespace", namespace)
		if err != nil {
			return "", err
		}
	} else if parsedNamespace != namespace {
		return "", errors.New("resource namespace " + parsedNamespace + " must be " + namespace)
	}
	return result, nil
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
