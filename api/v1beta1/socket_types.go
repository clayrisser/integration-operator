/**
 * File: /api/v1beta1/socket_types.go
 * Project: integration-operator
 * File Created: 17-10-2023 10:50:35
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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

// SocketSpec defines the desired state of Socket
type SocketSpec struct {
	// interface
	Interface *Interface `json:"interface,omitempty"`

	// limit
	Limit int32 `json:"limit,omitempty"`

	// vars
	Vars []*kustomizeTypes.Var `json:"vars,omitempty" yaml:"vars,omitempty"`

	// result vars
	ResultVars []*kustomizeTypes.Var `json:"resultVars,omitempty" yaml:"vars,omitempty"`

	// data
	Data map[string]string `json:"data,omitempty"`

	// data configmap name
	DataConfigMapName string `json:"dataConfigMapName,omitempty"`

	// data secret name
	DataSecretName string `json:"dataSecretName,omitempty"`

	// config
	Config map[string]string `json:"config,omitempty"`

	// config configmap name
	ConfigConfigMapName string `json:"configConfigMapName,omitempty"`

	// config secret name
	ConfigSecretName string `json:"configSecretName,omitempty"`

	// config template
	ConfigTemplate map[string]string `json:"configTemplate,omitempty"`

	// result
	Result map[string]string `json:"result,omitempty"`

	// result configmap name
	ResultConfigMapName string `json:"resultConfigMapName,omitempty"`

	// result secret name
	ResultSecretName string `json:"resultSecretName,omitempty"`

	// result template
	ResultTemplate map[string]string `json:"resultTemplate,omitempty"`

	// ServiceAccountName is the name of the ServiceAccount to use to run integrations.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,8,opt,name=serviceAccountName"`

	// apparatus
	Apparatus *SpecApparatus `json:"apparatus,omitempty"`

	// resources
	Resources []*Resource `json:"resources,omitempty"`

	// result resources
	ResultResources []*ResourceAction `json:"resultResources,omitempty"`

	// change epoch to force an update
	Epoch string `json:"epoch,omitempty"`

	// validation
	Validation *SocketSpecValidation `json:"validation,omitempty"`
}

type Interface struct {
	// config interface
	Config *ConfigInterface `json:"config,omitempty"`

	// result interface
	Result *ResultInterface `json:"result,omitempty"`
}

type ConfigInterface struct {
	// plug config properties
	Plug map[string]*SchemaProperty `json:"plug,omitempty"`

	// socket config properties
	Socket map[string]*SchemaProperty `json:"socket,omitempty"`
}

type ResultInterface struct {
	// plug result properties
	Plug map[string]*SchemaProperty `json:"plug,omitempty"`

	// socket result properties
	Socket map[string]*SchemaProperty `json:"socket,omitempty"`
}

type SchemaProperty struct {
	Default     string `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

type SocketSpecValidation struct {
	// namespace whitelist
	NamespaceWhitelist []string `json:"namespaceWhitelist,omitempty"`

	// namespace blacklist
	NamespaceBlacklist []string `json:"namespaceBlacklist,omitempty"`
}

// SocketStatus defines the observed state of Socket
type SocketStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// plugs coupled to socket
	CoupledPlugs []*CoupledPlug `json:"coupledPlugs,omitempty"`
}

type CoupledPlug struct {
	// API version of the plug
	APIVersion string `json:"apiVersion"`

	// Kind of the plug
	Kind string `json:"kind"`

	// Name of the plug
	Name string `json:"name"`

	// Namespace of the plug
	Namespace string `json:"namespace"`

	// UID of the plug
	UID types.UID `json:"uid"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Socket is the Schema for the sockets API
type Socket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SocketSpec   `json:"spec,omitempty"`
	Status SocketStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SocketList contains a list of Socket
type SocketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Socket `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Socket{}, &SocketList{})
}
