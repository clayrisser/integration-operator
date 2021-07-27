/**
 * File: /config/main.go
 * Project: integration-operator
 * File Created: 27-06-2021 02:53:17
 * Author: Clay Risser <email@clayrisser.com>
 * -----
 * Last Modified: 30-06-2021 11:58:19
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
