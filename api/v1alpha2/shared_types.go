/*
 * File: /api/v1alpha2/shared_types.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 27-06-2021 05:35:18
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Silicon Hills LLC (c) Copyright 2021
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha2

import (
	v1 "k8s.io/api/core/v1"
)

type Phase string

const (
	FailedPhase    Phase = "Failed"
	PendingPhase   Phase = "Pending"
	ReadyPhase     Phase = "Ready"
	SucceededPhase Phase = "Succeeded"
	UnknownPhase   Phase = "Unknown"
)

type When string

const (
	BrokenWhen    When = "broken"
	CoupledWhen   When = "coupled"
	CreatedWhen   When = "created"
	DecoupledWhen When = "decoupled"
	DeletedWhen   When = "deleted"
	UpdatedWhen   When = "updated"
)

type Do string

const (
	ApplyDo  Do = "apply"
	DeleteDo Do = "delete"
)

type Resource struct {
	Do       Do     `json:"do,omitempty"`
	Resource string `json:"resource,omitempty"`
	When     When   `json:"when,omitempty"`
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

	// containers
	Containers []*v1.Container `json:"containers,omitempty"`
}
