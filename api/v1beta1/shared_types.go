/**
 * File: /api/v1beta1/shared_types.go
 * Project: integration-operator
 * File Created: 17-10-2023 12:06:48
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
	v1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/kustomize/api/resid"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

const Finalizer = "integration.rock8s.com/finalizer"

type When string

const (
	CoupledWhen   When = "coupled"
	CreatedWhen   When = "created"
	DecoupledWhen When = "decoupled"
	DeletedWhen   When = "deleted"
	UpdatedWhen   When = "updated"
)

type Do string

const (
	ApplyDo    Do = "apply"
	DeleteDo   Do = "delete"
	RecreateDo Do = "recreate"
)

type ResourceAction struct {
	Do              Do                `json:"do,omitempty"`
	Template        *apiextv1.JSON    `json:"template,omitempty"`
	Templates       *[]*apiextv1.JSON `json:"templates,omitempty"`
	StringTemplate  string            `json:"stringTemplate,omitempty"`
	StringTemplates *[]string         `json:"stringTemplates,omitempty"`
}

type Resource struct {
	ResourceAction      `json:",inline"`
	RetainWhenDecoupled bool    `json:"retainWhenDecoupled,omitempty"`
	When                *[]When `json:"when,omitempty"`
}

type NamespacedName struct {
	// name
	Name string `json:"name"`

	// namespace
	Namespace string `json:"namespace,omitempty"`
}

type SpecApparatus struct {
	// endpoint
	Endpoint string `json:"endpoint,omitempty"`

	// terminate apparatus after idle for timeout in milliseconds
	IdleTimeout uint `json:"idleTimeout,omitempty"`

	// List of containers belonging to the apparatus.
	// Containers cannot currently be added or removed.
	// There must be at least one container in an apparatus.
	// Cannot be updated.
	// +patchMergeKey=name
	// +patchStrategy=merge
	Containers *[]v1.Container `json:"containers"`
}

// Var represents a variable whose value will be sourced
// from a field in a Kubernetes object.
type Var struct {
	// Value of identifier name e.g. FOO used in container args, annotations
	// Appears in pod template as $(FOO)
	Name string `json:"name" yaml:"name"`

	// ObjRef must refer to a Kubernetes resource under the
	// purview of this kustomization. ObjRef should use the
	// raw name of the object (the name specified in its YAML,
	// before addition of a namePrefix and a nameSuffix).
	ObjRef Target `json:"objref" yaml:"objref"`

	// FieldRef refers to the field of the object referred to by
	// ObjRef whose value will be extracted for use in
	// replacing $(FOO).
	// If unspecified, this defaults to fieldPath: $defaultFieldPath
	FieldRef kustomizeTypes.FieldSelector `json:"fieldref,omitempty" yaml:"fieldref,omitempty"`
}

// Target refers to a kubernetes object by Group, Version, Kind and Name
// gvk.Gvk contains Group, Version and Kind
// APIVersion is added to keep the backward compatibility of using ObjectReference
// for Var.ObjRef
type Target struct {
	APIVersion        string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	resid.Gvk         `json:",inline,omitempty" yaml:",inline,omitempty"`
	Name              string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace         string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	TemplateName      string `json:"templateName,omitempty" yaml:"templateName,omitempty"`
	TemplateNamespace string `json:"templateNamespace,omitempty" yaml:"templateNamespace,omitempty"`
}
