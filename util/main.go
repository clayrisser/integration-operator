/**
 * File: /util/main.go
 * Project: integration-operator
 * File Created: 17-10-2023 13:49:54
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

package util

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
)

type ConditionCoupledReason string

const (
	CouplingInProcess ConditionCoupledReason = "CouplingInProcess"
	CouplingSucceeded ConditionCoupledReason = "CouplingSucceeded"
	Error             ConditionCoupledReason = "Error"
	PlugCreated       ConditionCoupledReason = "PlugCreated"
	SocketCoupled     ConditionCoupledReason = "SocketCoupled"
	SocketCreated     ConditionCoupledReason = "SocketCreated"
	SocketEmpty       ConditionCoupledReason = "SocketEmpty"
	SocketNotCreated  ConditionCoupledReason = "SocketNotCreated"
	UpdatingInProcess ConditionCoupledReason = "UpdatingInProcess"
)

type ConditionType string

const (
	ConditionTypeCoupled ConditionType = "Coupled"
	ConditionTypeFailed  ConditionType = "Failed"
)

type Config map[string]string

type Result map[string]string
