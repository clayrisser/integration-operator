/**
 * Copyright 2021 Silicon Hills LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

type Replication struct {
        // resource to replicate from
        From ReplicationFrom `json:"from,omitempty"`

        // resource name to replicate to
        To ReplicationTo `json:"to,omitempty"`
}

type ReplicationTo struct {
        // resource name to replicate to
        Name string `json:"name,omitempty"`

        // namespace to replicate to
        Namespace string `json:"namespace,omitempty"`
}

type ReplicationFrom struct {
        // resource group to replicate from
        Group string `json:"group,omitempty"`

        // resource version to replicate from
        Version string `json:"version,omitempty"`

        // resource kind to replicate from
        Kind string `json:"kind,omitempty"`

        // resource name to replicate from
        Name string `json:"name,omitempty"`
}
