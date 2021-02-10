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

import { KubernetesObject, V1JobSpec } from '@kubernetes/client-node';
import { KustomizationSpec } from 'kustomize-operator';

export interface IntegrationPlugSpec {
  configmapPostfix?: string; // string `json:"configmapPostfix,omitempty"`
  kustomization?: KustomizationSpec; // KustomizationSpec `json:"kustomization,omitempty" yaml:"kustomization,omitempty"`
  mergeConfigmaps?: IntegrationPlugSpecMergeConfigmaps[]; // []*IntegrationPlugSpecMergeConfigmaps `json:"mergeConfigmaps,omitempty"`
  mergeSecrets?: IntegrationPlugSpecMergeSecrets[]; // []*IntegrationPlugSpecMergeSecrets `json:"mergeSecrets,omitempty"`
  replications?: Replication[]; // []*Replication `json:"replications,omitempty"`
  secretPostfix?: string; // string `json:"secretPostfix,omitempty"`
  socket?: IntegrationPlugSpecSocket; // IntegrationPlugSpecSocket `json:"socket,omitempty"`
}

export interface IntegrationPlugStatus {
  message?: string; // string `json:"message,omitempty"`
  phase?: IntegrationPlugStatusPhase; // string `json:"phase,omitempty"`
  ready?: boolean; // bool `json:"ready,omitempty"`
}

export interface IntegrationPlugResource extends KubernetesObject {
  spec?: IntegrationPlugSpec;
  status?: IntegrationPlugStatus;
}

export interface IntegrationSocketSpec {
  configmaps?: string[]; // []string `json:"configmaps,omitempty"`
  hooks?: IntegrationSocketSpecHook[]; // []*IntegrationSocketSpecHook `json:"hooks,omitempty"`
  replications?: Replication[]; // []*Replication `json:"replications,omitempty"`
  secrets?: string[]; // []string `json:"secrets,omitempty"`
  wait?: IntegrationPlugSpecWait; // IntegrationPlugSpecWait `json:"wait,omitempty"`
}

export interface IntegrationSocketStatus {}

export interface IntegrationSocketResource extends KubernetesObject {
  spec?: IntegrationSocketSpec;
  status?: IntegrationSocketStatus;
}

export enum IntegrationPlugStatusPhase {
  Failed = 'Failed',
  Pending = 'Pending',
  Succeeded = 'Succeeded',
  Unknown = 'Unknown'
}

export interface IntegrationPlugSpecWait {
  resources?: IntegrationPlugSpecWaitResource[]; // []*IntegrationPlugSpecWaitResource `json:"resources,omitempty"`
  timeout?: number; // int `json:"timeout,omitempty"`
}

export interface IntegrationPlugSpecWaitResource {
  group?: string; // string `json:"group,omitempty"`
  kind?: string; // string `json:"kind,omitempty"`
  name?: string; // string `json:"name,omitempty"`
  statusPhases?: string[]; // []string `json:"statusPhases,omitempty"`
  version?: string; // string `json:"version,omitempty"`
}

export interface IntegrationPlugSpecMergeConfigmaps {
  from?: string; // string `json:"from,omitempty"`
  to?: string; // string `json:"to,omitempty"`
}

export interface IntegrationPlugSpecMergeSecrets {
  from?: string; // string `json:"from,omitempty"`
  to?: string; // string `json:"to,omitempty"`
}

export interface IntegrationPlugSpecSocket {
  name?: string; // string `json:"name,omitempty"`
  namespace?: string; // string `json:"namespace,omitempty"`
}

export interface IntegrationSocketSpecHook {
  job?: V1JobSpec; // batchv1.JobSpec `json:"job,omitempty"`
  messageRegex?: string; // string `json:"messageRegex,omitempty"`
  name?: string; // string `json:"name,omitempty"`
  timeout?: number; // int `json:"timeout,omitempty"`
}

export interface Replication {
  from?: ReplicationFrom; // ReplicationFrom `json:"from,omitempty"`
  to?: ReplicationTo; // ReplicationTo `json:"to,omitempty"`
}

export interface ReplicationTo {
  name?: string; // string `json:"name,omitempty"`
  namespace?: string; // string `json:"namespace,omitempty"`
}

export interface ReplicationFrom {
  group?: string; // string `json:"group,omitempty"`
  kind?: string; // string `json:"kind,omitempty"`
  name?: string; // string `json:"name,omitempty"`
  version?: string; // string `json:"version,omitempty"`
}
