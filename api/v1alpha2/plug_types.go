/**
 * File: /api/v1alpha2/plug_types.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 23-08-2021 20:55:29
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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const PlugFinalizer = "integration.siliconhills.dev/finalizer"

// PlugSpec defines the desired state of Plug
type PlugSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// socket
	Socket NamespacedName `json:"socket,omitempty"`

	// interface
	Interface NamespacedName `json:"interface,omitempty"`

	// interface versions
	InterfaceVersions string `json:"interfaceVersions,omitempty"`

	// A var is a name (e.g. FOO) associated
	// with a field in a specific resource instance.  The field must
	// contain a value of type string/bool/int/float, and defaults to the name field
	// of the instance.  Any appearance of "$(FOO)" in the object
	// spec will be replaced, after the final
	// value of the specified field has been determined.
	Vars []*kustomizeTypes.Var `json:"vars,omitempty" yaml:"vars,omitempty"`

	// data
	Data map[string]string `json:"data,omitempty"`

	// data config map name
	DataConfigMapName string `json:"dataConfigMapName,omitempty"`

	// data secret name
	DataSecretName string `json:"dataSecretName,omitempty"`

	// config
	Config map[string]string `json:"config,omitempty"`

	// config config map name
	ConfigConfigMapName string `json:"configConfigMapName,omitempty"`

	// config secret name
	ConfigSecretName string `json:"configSecretName,omitempty"`

	// config mapper
	ConfigMapper map[string]string `json:"configMapper,omitempty"`

	// apparatus
	Apparatus *SpecApparatus `json:"apparatus,omitempty"`

	// resources
	Resources []*Resource `json:"resources,omitempty"`

	// change epoch to force an update
	Epoch int `json:"epoch,omitempty"`
}

// PlugStatus defines the observed state of Plug
type PlugStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// integration plug phase (Pending, Succeeded, Failed, Unknown)
	Phase Phase `json:"phase,omitempty"`

	// last update time
	LastUpdate metav1.Time `json:"lastUpdate,omitempty"`

	// status message
	Message string `json:"message,omitempty"`

	// socket coupled to plug
	CoupledSocket *CoupledSocket `json:"coupledSocket,omitempty"`

	// requeued
	Requeued bool `json:"requeued,omitempty"`
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
