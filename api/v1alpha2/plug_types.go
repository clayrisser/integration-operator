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

// PlugSpec defines the desired state of Plug
type PlugSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// socket
	Socket PlugSpecSocket `json:"socket,omitempty"`

	// interface
	Interface PlugSpecSocket `json:"interface,omitempty"`

	// interface versions
	InterfaceVersions string `json:"interfaceVersions,omitempty"`

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

	// mapper
	Mapper string `json:"mapper,omitempty"`

	// joined mapper
	JoinedMapper string `json:"joinedMapper,omitempty"`

	// changed mapper
	ChangedMapper string `json:"changedMapper,omitempty"`

	// departed mapper
	DepartedMapper string `json:"departedMapper,omitempty"`

	// broken mapper
	BrokenMapper string `json:"brokenMapper,omitempty"`
}

type PlugSpecSocket struct {
	// socket name
	Name string `json:"name,omitempty"`

	// socket namespace
	Namespace string `json:"namespace,omitempty"`
}

type PlugSpecInterface struct {
	// interface name
	Name string `json:"name,omitempty"`

	// interface namespace
	Namespace string `json:"namespace,omitempty"`
}

// PlugStatus defines the observed state of Plug
type PlugStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
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
