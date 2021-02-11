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
        kustomizeTypes "sigs.k8s.io/kustomize/api/types"
	batchv1 "k8s.io/api/batch/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IntegrationSocketSpec defines the desired state of IntegrationSocket
type IntegrationSocketSpec struct {
        // configuration for wait
        Wait IntegrationPlugSpecWait `json:"wait,omitempty"`

        // resources to replicate
        Replications []*Replication `json:"replications,omitempty"`

        // hooks that trigger during the lifecycle of an integration
        Hooks []*IntegrationSocketSpecHook `json:"hooks,omitempty"`
}

// IntegrationPlugSpecWait defines what to wait on before integrating
type IntegrationPlugSpecWait struct {
	// wait timeout in milliseconds
        Timeout int `json:"timeout,omitempty"`

        // resources to wait on
        Resources []*IntegrationPlugSpecWaitResource `json:"resources,omitempty"`
}

type IntegrationSocketSpecHook struct {
        // hook name
        Name string `json:"name,omitempty"`

        // job to run when hook is triggered
        Job batchv1.JobSpec `json:"job,omitempty"`

        // regex to find status message
        MessageRegex string `json:"messageRegex,omitempty"`

	// wait timeout in milliseconds
        Timeout int `json:"timeout,omitempty"`
}

type IntegrationPlugSpecWaitResource struct {
        // resource group to select
        Group string `json:"group,omitempty"`

        // resource version to select
        Version string `json:"version,omitempty"`

        // resource kind to select
        Kind string `json:"kind,omitempty"`

        // resource name to select
        Name string `json:"name,omitempty"`

        // resource status phases
        StatusPhases []string `json:"statusPhases,omitempty"`
}

// IntegrationSocketStatus defines the observed state of IntegrationSocket
type IntegrationSocketStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// IntegrationSocket is the Schema for the integrationsockets API
type IntegrationSocket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IntegrationSocketSpec   `json:"spec,omitempty"`
	Status IntegrationSocketStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IntegrationSocketList contains a list of IntegrationSocket
type IntegrationSocketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IntegrationSocket `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IntegrationSocket{}, &IntegrationSocketList{})
}
