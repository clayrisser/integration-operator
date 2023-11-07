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
	ResourceAction `json:",inline"`
	When           *[]When `json:"when,omitempty"`
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
