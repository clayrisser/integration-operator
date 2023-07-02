/**
 * File: /groupversion_info.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 02-07-2023 11:49:19
 * Modified By: Clay Risser <email@clayrisser.com>
 * -----
 * BitSpur (c) Copyright 2021
 */

// Package v1alpha2 contains API Schema definitions for the integration v1alpha2 API group
//+kubebuilder:object:generate=true
//+groupName=integration.rock8s.com
package v1alpha2

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "integration.rock8s.com", Version: "v1alpha2"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)
