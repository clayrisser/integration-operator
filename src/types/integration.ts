import { KubernetesObject, V1JobSpec } from '@kubernetes/client-node';
import { KustomizationSpec } from 'kustomize-operator';

export interface IntegrationPlugSpec {
  cleanup?: boolean; // bool `json:"cleanup,omitempty"`
  kustomization?: KustomizationSpec; // KustomizationSpec `json:"kustomization,omitempty" yaml:"kustomization,omitempty"`
  replications?: Replication[]; // []*Replication `json:"replications,omitempty"`
  socket?: IntegrationPlugSpecSocket; // IntegrationPlugSpecSocket `json:"socket,omitempty"`
  replicationAppendName?: string; // string `json:"replicationAppendName,omitempty"`
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
  hooks?: IntegrationSocketSpecHook[]; // []*IntegrationSocketSpecHook `json:"hooks,omitempty"`
  replications?: Replication[]; // []*Replication `json:"replications,omitempty"`
  wait?: IntegrationPlugSpecWait; // IntegrationPlugSpecWait `json:"wait,omitempty"`
  replicationAppendName?: string; // string `json:"replicationAppendName,omitempty"`
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
  group?: string; // string `json:"group,omitempty"`
  kind?: string; // string `json:"kind,omitempty"`
  name?: string; // string `json:"name,omitempty"`
  version?: string; // string `json:"version,omitempty"`
}
