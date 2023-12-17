/**
 * File: /api/v1beta1/deferredresource_types.go
 * Project: integration-operator
 * File Created: 17-12-2023 11:14:58
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
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/kustomize/api/resid"
)

// DeferredResourceSpec defines the desired state of DeferredResource
type DeferredResourceSpec struct {
	// Timeout is the maximum time to wait before creating the resource
	Timeout int64 `json:"timeout,omitempty"`

	// WaitFor is a list of resources to wait for before creating the resource
	WaitFor *[]*WaitForTarget `json:"waitFor,omitempty"`

	// Resource is the resource to create after the defer is resolved
	Resource *apiextv1.JSON `json:"resource,omitempty"`
	// ServiceAccountName is the name of the ServiceAccount to use to create deferred resources from.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,8,opt,name=serviceAccountName"`
}

// DeferredResourceStatus defines the observed state of DeferredResource
type DeferredResourceStatus struct {
	Conditions     []metav1.Condition    `json:"conditions,omitempty"`
	OwnerReference metav1.OwnerReference `json:"ownerReferences,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DeferredResource is the Schema for the deferredresources API
type DeferredResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeferredResourceSpec   `json:"spec,omitempty"`
	Status DeferredResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeferredResourceList contains a list of DeferredResource
type DeferredResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeferredResource `json:"items"`
}

// Target refers to a kubernetes object by Group, Version, Kind and Name
// gvk.Gvk contains Group, Version and Kind
// APIVersion is added to keep the backward compatibility of using ObjectReference
// for Var.ObjRef
type WaitForTarget struct {
	resid.Gvk  `json:",inline,omitempty" yaml:",inline,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Name       string `json:"name,omitempty" yaml:"name,omitempty"`
}

func init() {
	SchemeBuilder.Register(&DeferredResource{}, &DeferredResourceList{})
}
