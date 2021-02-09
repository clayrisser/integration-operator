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

import { KubernetesObject } from '@kubernetes/client-node';
import { KustomizationSpec } from 'kustomize-operator';

export interface IntegrationPlugSpec {
  kustomization?: KustomizationSpec;
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

export interface IntegrationSocketSpec {}

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
