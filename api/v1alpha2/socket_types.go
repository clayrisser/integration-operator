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
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

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

	// meta
	Meta string `json:"meta,omitempty"`

	// config mapper
	ConfigMapper string `json:"configMapper,omitempty"`

	// config endpoint
	ConfigEndpoint string `json:"configEndpoint,omitempty"`
}

// SocketStatus defines the observed state of Socket
type SocketStatus struct {
	// Conditions represent the latest available observations of an object's state
	Conditions []metav1.Condition `json:"conditions"`

	// integration socket phase (Pending, Succeeded, Failed, Unknown)
	Phase Phase `json:"phase,omitempty"`

	// socket is ready for coupling
	Ready bool `json:"ready,omitempty"`

	// number of plugs coupled to this socket
	PlugsCoupledCount int `json:"plugCoupledCount,omitempty"`
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
