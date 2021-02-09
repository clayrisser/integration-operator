/**
 * Copyright 2021 Silicon Hills LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IntegrationPlugSpec defines the desired state of IntegrationPlug
type IntegrationPlugSpec struct {
        // socket to integrate with
        Socket IntegrationPlugSpecSocket `json:"socket,omitempty"`

	// kustomization to apply after success
	Kustomization KustomizationSpec `json:"kustomization,omitempty" yaml:"kustomization,omitempty"`

        // postfix to apply to copied resource names
        ResourcePostfix string `json:"resourcePostfix,omitempty"`

        // configmaps to merge with copied configmaps
        MergeConfigmaps []*IntegrationPlugSpecMergeConfigmaps `json:"mergeConfigmaps,omitempty"`

        // secrets to merge with copied secrets
        MergeSecrets []*IntegrationPlugSpecMergeSecrets `json:"mergeSecrets,omitempty"`
}

// IntegrationPlugStatus defines the observed state of IntegrationPlug
type IntegrationPlugStatus struct {
        // integration connection message
        Message string `json:"message,omitempty"`

        // integration connection phase (Pending, Succeeded, Failed, Unknown)
        Phase string `json:"phase,omitempty"`

        // integration connection ready
        Ready bool `json:"ready,omitempty"`
}

type IntegrationPlugSpecMergeConfigmaps struct {
        // name of the configmap to merge from
        from string `json:"from,omitempty"`

        // name of the copied configmap to merge to
        to string `json:"to,omitempty"`
}

type IntegrationPlugSpecMergeSecrets struct {
        // name of the secret to merge from
        from string `json:"from,omitempty"`

        // name of the copied secret to merge to
        to string `json:"to,omitempty"`
}

type IntegrationPlugSpecSocket struct {
	// name of the socket
	Name string `json:"name,omitempty"`

	// namespace of the socket
	Namespace string `json:"namespace,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// IntegrationPlug is the Schema for the integrationplugs API
type IntegrationPlug struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IntegrationPlugSpec   `json:"spec,omitempty"`
	Status IntegrationPlugStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IntegrationPlugList contains a list of IntegrationPlug
type IntegrationPlugList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IntegrationPlug `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IntegrationPlug{}, &IntegrationPlugList{})
}
