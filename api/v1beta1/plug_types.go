/**
 * File: /api/v1beta1/plug_types.go
 * Project: integration-operator
 * File Created: 17-10-2023 10:50:57
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
)

// PlugSpec defines the desired state of Plug
type PlugSpec struct {
	// socket
	Socket NamespacedName `json:"socket,omitempty"`

	// vars
	Vars []*Var `json:"vars,omitempty" yaml:"vars,omitempty"`

	// result vars
	ResultVars []*Var `json:"resultVars,omitempty" yaml:"vars,omitempty"`

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
}

// PlugStatus defines the observed state of Plug
type PlugStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// socket coupled to plug
	CoupledSocket *CoupledSocket `json:"coupledSocket,omitempty"`

	// coupled result
	CoupledResult *CoupledResultStatus `json:"coupledResult,omitempty"`
}

type CoupledResult struct {
	// plug result
	Plug map[string]string `json:"plug,omitempty"`

	// socket result
	Socket map[string]string `json:"socket,omitempty"`
}

type CoupledResultStatus struct {
	CoupledResult `json:",inline"`

	// observed generation
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

type CoupledSocket struct {
	// API version of the socket
	APIVersion string `json:"apiVersion,omitempty"`

	// Kind of the socket
	Kind string `json:"kind,omitempty"`

	// Name of the socket
	Name string `json:"name,omitempty"`

	// Namespace of the socket
	Namespace string `json:"namespace,omitempty"`

	// UID of the socket
	UID types.UID `json:"uid,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Plug is the Schema for the plugs API
type Plug struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlugSpec   `json:"spec,omitempty"`
	Status PlugStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PlugList contains a list of Plug
type PlugList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Plug `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Plug{}, &PlugList{})
}
