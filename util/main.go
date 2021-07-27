/**
 * File: /util/main.go
 * Project: integration-operator
 * File Created: 23-06-2021 09:14:26
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 26-06-2021 10:53:47
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

package util

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
)

type Condition struct {
	Message string
	Reason  string
	Status  bool
	Type    string
}

type StatusCondition string

const (
	CouplingInProcessStatusCondition StatusCondition = "CouplingInProcess"
	CouplingSucceededStatusCondition StatusCondition = "CouplingSucceeded"
	ErrorStatusCondition             StatusCondition = "Error"
	PlugCreatedStatusCondition       StatusCondition = "PlugCreated"
	SocketCoupledStatusCondition     StatusCondition = "SocketCoupled"
	SocketCreatedStatusCondition     StatusCondition = "SocketCreated"
	SocketEmptyStatusCondition       StatusCondition = "SocketEmpty"
	SocketNotCreatedStatusCondition  StatusCondition = "SocketNotCreated"
	SocketNotReadyStatusCondition    StatusCondition = "SocketNotReady"
)
