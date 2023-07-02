/**
 * File: /interface_types.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 11:49:19
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// InterfaceSpec defines the desired state of Interface
type InterfaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// schemas
	Schemas []*InterfaceSpecSchema `json:"schemas,omitempty"`
}

type InterfaceSpecSchema struct {
	// version
	Version string `json:"version,omitempty"`

	// plug definition
	PlugDefinition *SchemaDefinition `json:"plugDefinition,omitempty"`

	// socket definition
	SocketDefinition *SchemaDefinition `json:"socketDefinition,omitempty"`
}

type SchemaDefinition struct {
	Description string                     `json:"description,omitempty"`
	Properties  map[string]*SchemaProperty `json:"properties,omitempty"`
}

type SchemaProperty struct {
	Default     string `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// InterfaceStatus defines the observed state of Interface
type InterfaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Interface is the Schema for the interfaces API
type Interface struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InterfaceSpec   `json:"spec,omitempty"`
	Status InterfaceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// InterfaceList contains a list of Interface
type InterfaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Interface `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Interface{}, &InterfaceList{})
}
