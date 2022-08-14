/**
 * File: /shared_types.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 14-08-2022 14:34:43
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * Risser Labs LLC (c) Copyright 2021
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
	ApplyDo    Do = "apply"
	DeleteDo   Do = "delete"
	RecreateDo Do = "recreate"
)

type Resource struct {
	Do       Do      `json:"do,omitempty"`
	Resource string  `json:"resource,omitempty"`
	When     *[]When `json:"when,omitempty"`
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
	Containers *[]v1.Container `json:"containers,omitempty"`
}
