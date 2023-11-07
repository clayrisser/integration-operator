/**
 * File: /util/resource.go
 * Project: integration-operator
 * File Created: 17-10-2023 13:49:54
 * Author: Clay Risser
 * -----
 * BitSpur (c) Copyright 2021 - 2023
 *
 * Licensed under the GNU Affero General Public License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.gnu.org/licenses/agpl-3.0.en.html
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * You can be released from the requirements of the license by purchasing
 * a commercial license. Buying such a license is mandatory as soon as you
 * develop commercial activities involving this software without disclosing
 * the source code of your own applications.
 */

package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	integrationv1beta1 "gitlab.com/bitspur/rock8s/integration-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

type ResourceUtil struct {
	client *kubernetes.Clientset
	ctx    context.Context
}

func NewResourceUtil(ctx context.Context) *ResourceUtil {
	return &ResourceUtil{
		client: kubernetes.NewForConfigOrDie(ctrl.GetConfigOrDie()),
		ctx:    ctx,
	}
}

func (u *ResourceUtil) PlugCreated(plug *integrationv1beta1.Plug) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1beta1.CreatedWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		nil,
		nil,
		nil,
		nil,
		nil,
		plug.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugCoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1beta1.CoupledWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		nil,
		nil,
		plug.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugUpdated(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1beta1.UpdatedWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		nil,
		nil,
		plug.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugDecoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1beta1.DecoupledWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		nil,
		nil,
		plug.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) PlugDeleted(
	plug *integrationv1beta1.Plug,
) error {
	resources := u.filterResources(plug.Spec.Resources, integrationv1beta1.DeletedWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, plug.Namespace, EnsureServiceAccount(plug.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		nil,
		nil,
		nil,
		nil,
		nil,
		plug.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketCreated(socket *integrationv1beta1.Socket) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1beta1.CreatedWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		nil,
		socket,
		nil,
		nil,
		nil,
		nil,
		socket.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketCoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1beta1.CoupledWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		nil,
		nil,
		socket.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketUpdated(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1beta1.UpdatedWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		nil,
		nil,
		socket.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketDecoupled(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1beta1.DecoupledWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		plug,
		socket,
		plugConfig,
		socketConfig,
		nil,
		nil,
		socket.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) SocketDeleted(
	socket *integrationv1beta1.Socket,
) error {
	resources := u.filterResources(socket.Spec.Resources, integrationv1beta1.DeletedWhen)
	kubectlUtil := NewKubectlUtil(u.ctx, socket.Namespace, EnsureServiceAccount(socket.Spec.ServiceAccountName))
	if err := u.ProcessResources(
		nil,
		socket,
		nil,
		nil,
		nil,
		nil,
		socket.Namespace,
		resources,
		kubectlUtil,
	); err != nil {
		return err
	}
	return nil
}

func (u *ResourceUtil) GetResource(
	namespace string,
	objRef kustomizeTypes.Target,
	kubectlUtil *KubectlUtil,
) (*unstructured.Unstructured, error) {
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
		objRef.APIVersion = objRef.Group + "/" + objRef.Version
	}
	var buff bytes.Buffer
	err = t.Execute(&buff, objRef)
	if err != nil {
		return nil, err
	}
	body := buff.Bytes()
	return kubectlUtil.Get(body)
}

func (u *ResourceUtil) ProcessResources(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	plugResult *Result,
	socketResult *Result,
	namespace string,
	resources []*integrationv1beta1.ResourceAction,
	kubectlUtil *KubectlUtil,
) error {
	for _, resource := range resources {
		templates := []string{}
		if resource.Template != nil {
			templates = append(templates, string(resource.Template.Raw))
		}
		if resource.Templates != nil {
			for _, template := range *resource.Templates {
				templates = append(templates, string(template.Raw))
			}
		}
		if resource.StringTemplate != "" {
			templates = append(templates, resource.StringTemplate)
		}
		if resource.StringTemplates != nil {
			templates = append(templates, *resource.StringTemplates...)
		}
		for _, template := range templates {
			templatedResource, err := u.templateResource(
				plug,
				socket,
				plugConfig,
				socketConfig,
				plugResult,
				socketResult,
				namespace,
				template,
			)
			if strings.TrimSpace(templatedResource) == "" {
				return nil
			}
			if err != nil {
				return err
			}
			do := resource.Do
			if do == "" {
				do = integrationv1beta1.ApplyDo
			}
			if resource.Do == integrationv1beta1.ApplyDo {
				if err := kubectlUtil.Apply([]byte(templatedResource)); err != nil {
					return err
				}
			} else if resource.Do == integrationv1beta1.DeleteDo {
				if err := kubectlUtil.Delete([]byte(templatedResource)); err != nil {
					return err
				}
			} else if resource.Do == integrationv1beta1.RecreateDo {
				kubectlUtil.Delete([]byte(templatedResource))
				if err := kubectlUtil.Apply([]byte(templatedResource)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (u *ResourceUtil) filterResources(
	resources []*integrationv1beta1.Resource,
	when integrationv1beta1.When,
) []*integrationv1beta1.ResourceAction {
	filteredResources := []*integrationv1beta1.ResourceAction{}
	if resources == nil {
		return filteredResources
	}
	if when == "" {
		when = integrationv1beta1.CoupledWhen
	}
	for _, resource := range resources {
		if WhenInWhenSlice(when, resource.When) {
			filteredResources = append(filteredResources, &integrationv1beta1.ResourceAction{
				Do:              resource.Do,
				StringTemplate:  resource.StringTemplate,
				StringTemplates: resource.StringTemplates,
				Template:        resource.Template,
				Templates:       resource.Templates,
			})
		}
	}
	return filteredResources
}

func (u *ResourceUtil) templateResource(
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	plugResult *Result,
	socketResult *Result,
	namespace string,
	body string,
) (string, error) {
	data, err := u.buildTemplateData(plug, socket, plugConfig, socketConfig, plugResult, socketResult)
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
	plug *integrationv1beta1.Plug,
	socket *integrationv1beta1.Socket,
	plugConfig *Config,
	socketConfig *Config,
	plugResult *Result,
	socketResult *Result,
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
	if plugResult != nil {
		dataMap["plugResult"] = plugResult
	}
	if socketResult != nil {
		dataMap["socketResult"] = socketResult
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
