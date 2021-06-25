/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const SocketFinalizer = "integration.siliconhills.dev/finalizer"

// SocketSpec defines the desired state of Socket
type SocketSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// interface
	Interface NamespacedName `json:"interface,omitempty"`

	// interface versions
	InterfaceVersions string `json:"interfaceVersions,omitempty"`

	// limit
	Limit int32 `json:"limit,omitempty"`

	// namspace scope
	NamespaceScope string `json:"namespaceScope,omitempty"`

	// A var is a name (e.g. FOO) associated
	// with a field in a specific resource instance.  The field must
	// contain a value of type string/bool/int/float, and defaults to the name field
	// of the instance.  Any appearance of "$(FOO)" in the object
	// spec will be replaced, after the final
	// value of the specified field has been determined.
	Vars []kustomizeTypes.Var `json:"vars,omitempty" yaml:"vars,omitempty"`

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
	Apparatus SpecApparatus `json:"apparatus,omitempty"`

	// Resources
	Resources []*Resource `json:"resources,omitempty"`
}

// SocketStatus defines the observed state of Socket
type SocketStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// integration socket phase (Pending, Succeeded, Failed, Unknown)
	Phase Phase `json:"phase,omitempty"`

	// socket is ready for coupling
	Ready bool `json:"ready,omitempty"`

	// last update time
	LastUpdate metav1.Time `json:"lastUpdate,omitempty"`

	// status message
	Message string `json:"message,omitempty"`

	// plugs coupled to socket
	CoupledPlugs []CoupledPlug `json:"coupledPlugs,omitempty"`
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
