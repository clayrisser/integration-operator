/**
 * File: /config/main.go
 * Project: integration-operator
 * File Created: 17-10-2023 14:02:00
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

package config

import (
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var MaxRequeueDuration time.Duration = time.Duration(float64(time.Hour.Nanoseconds() * 6))

var StartTime metav1.Time = metav1.Now()

var DebugPlugEndpoint = os.Getenv("DEBUG_PLUG_ENDPOINT")

var DebugSocketEndpoint = os.Getenv("DEBUG_SOCKET_ENDPOINT")
