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
import { KustomizationSpec, Selector } from 'kustomize-operator';

export interface IntegrationPlugSpec {
  kustomization?: KustomizationSpec; // KustomizationSpec `json:"kustomization,omitempty" yaml:"kustomization,omitempty"`
  mergeConfigmaps?: IntegrationPlugSpecMergeConfigmaps[]; // []*IntegrationPlugSpecMergeConfigmaps `json:"mergeConfigmaps,omitempty"`
  mergeSecrets?: IntegrationPlugSpecMergeSecrets[]; // []*IntegrationPlugSpecMergeSecrets `json:"mergeSecrets,omitempty"`
  resourcePostfix?: string; // string `json:"resourcePostfix,omitempty"`
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
  replications?: IntegrationSocketSpecReplication[]; // []*IntegrationSocketSpecReplication `json:"replications,omitempty"`
  secrets?: string[]; // []string `json:"secrets,omitempty"`
  wait?: IntegrationPlugSpecWait; // IntegrationPlugSpecWait `json:"wait,omitempty"`
  job?: V1JobSpec; // batchv1.JobSpec `json:"cleanupJob,omitempty"`
  cleanupJob?: V1JobSpec; // batchv1.JobSpec `json:"cleanupJob,omitempty"`
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
  interval?: number; // int `json:"interval,omitempty"`
  resources?: IntegrationPlugSpecWaitResource[]; // []*IntegrationPlugSpecWaitResource `json:"resources,omitempty"`
  timeout?: number; // int `json:"timeout,omitempty"`
}

export interface IntegrationSocketSpecReplication {
  from?: Selector; // kustomizeTypes.Selector `json:"from,omitempty"`
  to?: string; // string `json:"to,omitempty"`
}

export interface IntegrationPlugSpecWaitResource {
  selector?: Selector; // kustomizeTypes.Selector `json:"selector,omitempty"`
  statusPhases?: string[]; // []string `json:"statusPhases,omitempty"`
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
