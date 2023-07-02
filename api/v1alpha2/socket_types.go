/**
 * File: /socket_types.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 11:58:32
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

const SocketFinalizer = "integration.rock8s.com/finalizer"

// SocketSpec defines the desired state of Socket
type SocketSpec struct {
	// interface
	Interface NamespacedName `json:"interface,omitempty"`

	// interface versions
	InterfaceVersions string `json:"interfaceVersions,omitempty"`

	// limit
	Limit int32 `json:"limit,omitempty"`

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
	Epoch string `json:"epoch,omitempty"`

	// validation
	Validation *SocketSpecValidation `json:"validation,omitempty"`
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

	// integration socket phase (Pending, Succeeded, Failed, Unknown)
	Phase Phase `json:"phase,omitempty"`

	// socket is ready for coupling
	Ready bool `json:"ready,omitempty"`

	// last update time
	LastUpdate metav1.Time `json:"lastUpdate,omitempty"`

	// status message
	Message string `json:"message,omitempty"`

	// plugs coupled to socket
	CoupledPlugs []*CoupledPlug `json:"coupledPlugs,omitempty"`

	// requeued
	Requeued bool `json:"requeued,omitempty"`
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
