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

import * as k8s from '@kubernetes/client-node';
import chalk from 'chalk';
import newRegExp from 'newregexp';
import stripAnsi from 'strip-ansi';
import { KustomizationResource } from 'kustomize-operator';
import { ResourceMeta } from '@dot-i/k8s-operator';
import { ReplicationService, KubectlService, Output } from '~/services';
import {
  HashMap,
  IntegrationPlugResource,
  IntegrationPlugSpecMergeConfigmaps,
  IntegrationPlugSpecMergeSecrets,
  IntegrationPlugSpecWaitResource,
  IntegrationPlugStatus,
  IntegrationPlugStatusPhase,
  IntegrationSocketResource,
  IntegrationSocketSpecHook,
  Replication
} from '~/types';
import {
  KustomizeResourceGroup,
  KustomizeResourceKind,
  KustomizeResourceVersion,
  ResourceKind,
  ResourceVersion
} from '~/integrationOperator';
import Controller from './controller';

export default class IntegrationPlug extends Controller {
  private coreV1Api: k8s.CoreV1Api;

  private customObjectsApi: k8s.CustomObjectsApi;

  private kubeConfig: k8s.KubeConfig;

  private batchV1Api: k8s.BatchV1Api;

  private kubectl = new KubectlService();

  constructor(groupnameprefix: string, kind: string) {
    super(groupnameprefix, kind);
    this.kubeConfig = new k8s.KubeConfig();
    this.kubeConfig.loadFromDefault();
    this.batchV1Api = this.kubeConfig.makeApiClient(k8s.BatchV1Api);
    this.coreV1Api = this.kubeConfig.makeApiClient(k8s.CoreV1Api);
    this.customObjectsApi = this.kubeConfig.makeApiClient(k8s.CustomObjectsApi);
  }

  async deleted(
    plugResource: IntegrationPlugResource,
    _meta: ResourceMeta,
    _oldPlugResource?: IntegrationPlugResource
  ) {
    const socketResource = await this.getSocketResource(plugResource);
    if (!socketResource) return null;
    await this.callHook(Hook.BeforeCleanup, plugResource, socketResource);
    await this.callHook(Hook.Cleanup, plugResource, socketResource);
    await this.callHook(Hook.AfterCleanup, plugResource, socketResource);
    return null;
  }

  async added(
    plugResource: IntegrationPlugResource,
    _meta: ResourceMeta,
    oldPlugResource?: IntegrationPlugResource
  ): Promise<any> {
    if (
      plugResource.metadata?.generation ===
      oldPlugResource?.metadata?.generation
    ) {
      return null;
    }
    const socketResource = await this.getSocketResource(plugResource);
    if (!socketResource) {
      const message = `${this.operatorService.getFullName({
        kind: ResourceKind.IntegrationSocket,
        name: plugResource.spec?.socket?.name || ''
      })} does not exist in namespace ${chalk.blueBright.bold(
        plugResource.spec?.socket?.namespace
      )}`;
      this.spinner.fail(message);
      await this.updateStatus(
        {
          message,
          phase: IntegrationPlugStatusPhase.Failed,
          ready: false
        },
        plugResource
      );
      return null;
    }
    try {
      await this.beginApply(plugResource, socketResource);
      await Promise.all([
        this.callHook(Hook.BeforeCreate, plugResource, socketResource),
        this.callHook(Hook.BeforeCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.apply(plugResource, socketResource);
      const [createResult, createOrUpdateResult] = await Promise.all([
        this.callHook(Hook.Create, plugResource, socketResource),
        this.callHook(Hook.CreateOrUpdate, plugResource, socketResource)
      ]);
      if (plugResource?.spec?.kustomization) {
        await this.applyKustomization(plugResource);
      }
      await Promise.all([
        this.callHook(Hook.AfterCreate, plugResource, socketResource),
        this.callHook(Hook.AfterCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.endApply(plugResource, socketResource, [
        ...createResult,
        ...createOrUpdateResult
      ]);
    } catch (err) {
      await this.updateStatus(
        {
          message: this.operatorService.getErrorMessage(err),
          phase: IntegrationPlugStatusPhase.Failed,
          ready: false
        },
        plugResource
      );
      throw err;
    }
    return null;
  }

  async modified(
    plugResource: IntegrationPlugResource,
    _meta: ResourceMeta,
    oldPlugResource?: IntegrationPlugResource
  ): Promise<any> {
    if (
      plugResource.metadata?.generation ===
      oldPlugResource?.metadata?.generation
    ) {
      return null;
    }
    const socketResource = await this.getSocketResource(plugResource);
    if (!socketResource) {
      const message = `${this.operatorService.getFullName({
        kind: ResourceKind.IntegrationSocket,
        name: plugResource.spec?.socket?.name || ''
      })} does not exist in namespace ${chalk.blueBright.bold(
        plugResource.spec?.socket?.namespace
      )}`;
      this.spinner.fail(message);
      await this.updateStatus(
        {
          message,
          phase: IntegrationPlugStatusPhase.Failed,
          ready: false
        },
        plugResource
      );
      return null;
    }
    try {
      await this.beginApply(plugResource, socketResource);
      await Promise.all([
        this.callHook(Hook.BeforeUpdate, plugResource, socketResource),
        this.callHook(Hook.BeforeCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.apply(plugResource, socketResource);
      const [updateResult, createOrUpdateResult] = await Promise.all([
        this.callHook(Hook.Update, plugResource, socketResource),
        this.callHook(Hook.CreateOrUpdate, plugResource, socketResource)
      ]);
      if (plugResource?.spec?.kustomization) {
        await this.applyKustomization(plugResource);
      }
      await Promise.all([
        this.callHook(Hook.AfterUpdate, plugResource, socketResource),
        this.callHook(Hook.AfterCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.endApply(plugResource, socketResource, [
        ...updateResult,
        ...createOrUpdateResult
      ]);
    } catch (err) {
      await this.updateStatus(
        {
          message: this.operatorService.getErrorMessage(err),
          phase: IntegrationPlugStatusPhase.Failed,
          ready: false
        },
        plugResource
      );
      throw err;
    }
    return null;
  }

  private async beginApply(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    const message = `integrating with ${this.operatorService.getFullName({
      resource: socketResource
    })}`;
    this.spinner.info(
      `${this.operatorService.getFullName({
        resource: plugResource
      })} is ${message}`
    );
    await this.updateStatus(
      {
        message,
        phase: IntegrationPlugStatusPhase.Pending,
        ready: false
      },
      plugResource
    );
  }

  private async apply(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    await this.waitForResources(plugResource, socketResource);
    await this.copyAndMergeConfigmaps(plugResource, socketResource);
    await this.copyAndMergeSecrets(plugResource, socketResource);
    await this.replicateSocketResources(socketResource);
    await this.replicatePlugResources(plugResource);
  }

  private async endApply(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource,
    hookResults: HookResult[]
  ) {
    const statusMessage = hookResults
      .map(
        ({ name, namespace, message, hookName }: HookResult) =>
          `${message} in ${this.operatorService.getFullName({
            kind: 'Job',
            apiVersion: 'batch/v1',
            name,
            ns: namespace
          })} for hook ${hookName}`
      )
      .join('\n');
    const message = `successfully integrated with ${this.operatorService.getFullName(
      {
        resource: socketResource
      }
    )}`;
    this.spinner.succeed(
      `${this.operatorService.getFullName({
        resource: plugResource
      })} has ${message}`
    );
    await this.updateStatus(
      {
        message: statusMessage || message,
        phase: IntegrationPlugStatusPhase.Succeeded,
        ready: true
      },
      plugResource
    );
  }

  private async getResources(
    resources: k8s.KubernetesObject[]
  ): Promise<
    (k8s.KubernetesObject & {
      status?: { [key: string]: any; phase?: string };
    })[]
  > {
    return (
      (
        await this.kubectl.get<k8s.KubernetesListObject<k8s.KubernetesObject>>({
          stdin: resources,
          output: Output.Json
        })
      )?.items || []
    );
  }

  private async waitForResources(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource,
    timeout = 60000,
    timeLeft?: number
  ) {
    const waitTime = Math.max(5000, timeout / 10);
    if (typeof timeLeft !== 'number') timeLeft = timeout;
    const resources = await this.getResources(
      (socketResource.spec?.wait?.resources || []).map<k8s.KubernetesObject>(
        (waitResource: IntegrationPlugSpecWaitResource) => {
          if (!waitResource.version) {
            throw new Error('resource version is not defined');
          }
          if (!waitResource.kind) {
            throw new Error('resource kind is not defined');
          }
          if (!waitResource.name) {
            throw new Error('resource name is not defined');
          }
          return {
            apiVersion: this.operatorService.getApiVersion(
              waitResource.version,
              waitResource.group
            ),
            kind: waitResource.kind,
            metadata: {
              name: waitResource.name,
              namespace: plugResource.metadata?.namespace
            }
          };
        }
      )
    );
    const foundAllResources = (
      socketResource.spec?.wait?.resources || []
    ).reduce(
      (ready: boolean, waitResource: IntegrationPlugSpecWaitResource) => {
        if (!ready) return ready;
        const resource = resources.find((resource: k8s.KubernetesObject) => {
          return (
            waitResource.name === resource.metadata?.name &&
            waitResource.kind === resource.kind &&
            this.operatorService.getApiVersion(
              waitResource.version!,
              waitResource.group
            ) === resource.apiVersion
          );
        });
        if (typeof resource === 'undefined' || !resource) return false;
        return !!(
          !waitResource.statusPhases?.length ||
          (waitResource.statusPhases || []).find(
            (statusPhase: string) => statusPhase === resource?.status?.phase
          )
        );
      },
      true
    );
    if (!foundAllResources) {
      if (timeLeft < 0) {
        throw new Error(
          `failed to find some resources for ${this.operatorService.getFullName(
            {
              resource: plugResource
            }
          )}`
        );
      }
      this.spinner.info(
        `waiting ${timeLeft}ms on resources for ${this.operatorService.getFullName(
          {
            resource: plugResource
          }
        )}`
      );
      await new Promise((r) => setTimeout(r, waitTime));
      await this.waitForResources(
        plugResource,
        socketResource,
        timeout,
        timeLeft - waitTime
      );
    }
  }

  private async callHook(
    hookName: Hook,
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    const ns = plugResource.metadata?.namespace!;
    const filteredHooks = (socketResource.spec?.hooks || []).filter(
      (hook: IntegrationSocketSpecHook) => {
        return hook.name === hookName;
      }
    );
    return Promise.all(
      filteredHooks.map(async (hook: IntegrationSocketSpecHook, i: number) => {
        const job = (
          await this.batchV1Api.createNamespacedJob(ns, {
            metadata: {
              name: `${plugResource.metadata
                ?.name!}-${hookName}-${i.toString()}`,
              namespace: ns,
              ownerReferences: [
                this.operatorService.getOwnerReference(plugResource, ns)
              ]
            },
            spec: hook.job
          })
        ).body;
        this.spinner.succeed(
          `created ${this.operatorService.getFullName({
            resource: job
          })}`
        );
        await this.waitForJobToFinish(job, hook.timeout);
        const logs = await this.getJobLogs(job);
        let message = 'completed';
        if (hook.messageRegex) {
          const messageMatches = logs.match(newRegExp(hook.messageRegex));
          message = [...(messageMatches || [])].join('\n');
        }
        return {
          hookName: hook.name!,
          message,
          name: job.metadata?.name!,
          namespace: job.metadata?.namespace!
        };
      })
    );
  }

  private async waitForJobToFinish(
    job: k8s.V1Job,
    timeout = 60000,
    timeLeft?: number
  ) {
    const waitTime = Math.max(5000, timeout / 10);
    if (typeof timeLeft !== 'number') timeLeft = timeout;
    const jobStatus = (
      await this.batchV1Api.readNamespacedJobStatus(
        job.metadata?.name!,
        job.metadata?.namespace!
      )
    ).body.status;
    if (jobStatus?.succeeded) return;
    await new Promise((r) => setTimeout(r, waitTime));
    await this.waitForJobToFinish(job, timeout, timeLeft - waitTime);
  }

  private async getJobLogs(job: k8s.V1Job): Promise<string> {
    const pods = (
      await this.coreV1Api.listNamespacedPod(job.metadata?.namespace || '')
    ).body;
    const podName =
      pods.items.find(
        (pod: k8s.V1Pod) =>
          pod.metadata?.labels?.['job-name'] === job.metadata?.name
      )?.metadata?.name || '';
    return (
      await this.coreV1Api.readNamespacedPodLog(
        podName,
        job.metadata?.namespace || '',
        undefined,
        false
      )
    ).body;
  }

  private async replicatePlugResources(plugResource: IntegrationPlugResource) {
    const replicationService = new ReplicationService(
      plugResource.metadata?.namespace!
    );
    await Promise.all(
      (
        plugResource.spec?.replications || []
      ).map(async (replication: Replication) =>
        replicationService.apply(replication)
      )
    );
  }

  private async replicateSocketResources(
    socketResource: IntegrationSocketResource
  ) {
    const replicationService = new ReplicationService(
      socketResource.metadata?.namespace!
    );
    await Promise.all(
      (
        socketResource.spec?.replications || []
      ).map(async (replication: Replication) =>
        replicationService.apply(replication)
      )
    );
  }

  private async copyAndMergeConfigmaps(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    const configmapPostfix = plugResource?.spec?.configmapPostfix;
    const postfix = configmapPostfix ? `-${configmapPostfix}` : '';
    const mergeConfigmapsData = await Promise.all(
      (plugResource.spec?.mergeConfigmaps || []).map<
        Promise<[string | null, HashMap<string>]>
      >(
        async ({
          to,
          from
        }: IntegrationPlugSpecMergeConfigmaps): Promise<
          [string | null, HashMap<string>]
        > => {
          let plugConfigmap: k8s.V1ConfigMap | null = null;
          if (from) {
            plugConfigmap = (
              await this.coreV1Api.readNamespacedConfigMap(
                from,
                plugResource.metadata?.namespace!
              )
            ).body;
          }
          return [to || null, plugConfigmap?.data || {}];
        }
      )
    );
    await Promise.all(
      (socketResource.spec?.configmaps || []).map(
        async (configmapName: string) => {
          const socketConfigmap = (
            await this.coreV1Api.readNamespacedConfigMap(
              configmapName,
              socketResource.metadata?.namespace!
            )
          ).body;
          let mergedData = socketConfigmap?.data || {};
          mergeConfigmapsData.forEach(
            ([mergeConfigmapToName, mergeConfigmapFromData]: [
              string | null,
              HashMap<string>
            ]) => {
              if (mergeConfigmapToName === socketConfigmap?.metadata?.name) {
                mergedData = {
                  ...mergedData,
                  ...mergeConfigmapFromData
                };
              }
            }
          );
          await this.createOrUpdateConfigMap(
            configmapName + postfix,
            plugResource.metadata?.namespace!,
            mergedData,
            plugResource
          );
        }
      )
    );
  }

  private async copyAndMergeSecrets(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    const secretPostfix = plugResource?.spec?.secretPostfix;
    const postfix = secretPostfix ? `-${secretPostfix}` : '';
    const mergeSecretsStringData = await Promise.all(
      (plugResource.spec?.mergeSecrets || []).map<
        Promise<[string | null, HashMap<string>]>
      >(
        async ({
          to,
          from
        }: IntegrationPlugSpecMergeSecrets): Promise<
          [string | null, HashMap<string>]
        > => {
          let plugSecret: k8s.V1Secret | null = null;
          if (from) {
            plugSecret = (
              await this.coreV1Api.readNamespacedSecret(
                from,
                plugResource.metadata?.namespace!
              )
            ).body;
          }
          return [
            to || null,
            this.operatorService.base64DecodeSecretData(plugSecret?.data)
          ];
        }
      )
    );
    await Promise.all(
      (socketResource.spec?.secrets || []).map(async (secretName: string) => {
        const socketSecret = (
          await this.coreV1Api.readNamespacedSecret(
            secretName,
            socketResource.metadata?.namespace!
          )
        ).body;
        let mergedStringData = this.operatorService.base64DecodeSecretData(
          socketSecret?.data
        );
        mergeSecretsStringData.forEach(
          ([mergeSecretToName, mergeSecretFromStringData]: [
            string | null,
            HashMap<string>
          ]) => {
            if (mergeSecretToName === socketSecret?.metadata?.name) {
              mergedStringData = {
                ...mergedStringData,
                ...mergeSecretFromStringData
              };
            }
          }
        );
        await this.createOrUpdateSecret(
          secretName + postfix,
          plugResource.metadata?.namespace!,
          mergedStringData,
          plugResource
        );
      })
    );
  }

  private async createOrUpdateSecret(
    name: string,
    namespace: string,
    data: HashMap<string>,
    owner?: k8s.KubernetesObject
  ) {
    try {
      await this.coreV1Api.readNamespacedSecret(name, namespace);
      await this.coreV1Api.patchNamespacedSecret(
        name,
        namespace,
        [
          {
            op: 'replace',
            path: '/stringData',
            value: data
          }
        ],
        undefined,
        undefined,
        undefined,
        undefined,
        {
          headers: { 'Content-Type': 'application/json-patch+json' }
        }
      );
    } catch (err) {
      if (err.statusCode !== 404) throw err;
      await this.coreV1Api.createNamespacedSecret(namespace, {
        metadata: {
          name,
          namespace,
          ...(typeof owner !== 'undefined'
            ? {
                ownerReferences: [
                  this.operatorService.getOwnerReference(owner, namespace)
                ]
              }
            : {})
        },
        stringData: data
      });
    }
  }

  private async createOrUpdateConfigMap(
    name: string,
    namespace: string,
    data: HashMap<string>,
    owner?: k8s.KubernetesObject
  ) {
    try {
      await this.coreV1Api.readNamespacedConfigMap(name, namespace);
      await this.coreV1Api.patchNamespacedConfigMap(
        name,
        namespace,
        [
          {
            op: 'replace',
            path: '/data',
            value: data
          }
        ],
        undefined,
        undefined,
        undefined,
        undefined,
        {
          headers: { 'Content-Type': 'application/json-patch+json' }
        }
      );
    } catch (err) {
      if (err.statusCode !== 404) throw err;
      await this.coreV1Api.createNamespacedConfigMap(namespace, {
        metadata: {
          name,
          namespace,
          ...(typeof owner !== 'undefined'
            ? {
                ownerReferences: [
                  this.operatorService.getOwnerReference(owner, namespace)
                ]
              }
            : {})
        },
        data
      });
    }
  }

  private async getSocketResource(
    plugResource: IntegrationPlugResource
  ): Promise<IntegrationSocketResource | null> {
    if (
      !plugResource.metadata?.name ||
      !plugResource.metadata?.namespace ||
      !plugResource.spec?.socket?.name
    ) {
      return null;
    }
    const namespace =
      plugResource.spec?.socket?.namespace || plugResource.metadata.namespace;
    try {
      const socketResource = (
        await this.customObjectsApi.getNamespacedCustomObject(
          this.group,
          ResourceVersion.V1alpha1,
          namespace,
          this.operatorService.kind2plural(ResourceKind.IntegrationSocket),
          plugResource.spec.socket.name
        )
      ).body as IntegrationSocketResource;
      return socketResource;
    } catch (err) {
      if (err.statusCode !== 404) throw err;
      return null;
    }
  }

  private async updateStatus(
    plugStatus: IntegrationPlugStatus,
    plugResource: IntegrationPlugResource
  ): Promise<void> {
    if (!plugResource.metadata?.name || !plugResource.metadata.namespace) {
      return;
    }
    plugStatus.message = stripAnsi(plugStatus.message || '');
    await this.customObjectsApi.patchNamespacedCustomObjectStatus(
      this.group,
      ResourceVersion.V1alpha1,
      plugResource.metadata.namespace,
      this.plural,
      plugResource.metadata.name,
      [
        {
          op: 'replace',
          path: '/status',
          value: plugStatus
        }
      ],
      undefined,
      undefined,
      undefined,
      {
        headers: { 'Content-Type': 'application/json-patch+json' }
      }
    );
  }

  private async applyKustomization(
    plugResource: IntegrationPlugResource
  ): Promise<void> {
    if (!plugResource.metadata?.name || !plugResource.metadata.namespace) {
      return;
    }
    try {
      await this.customObjectsApi.getNamespacedCustomObject(
        this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        ),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        this.operatorService.kind2plural(KustomizeResourceKind.Kustomization),
        plugResource.metadata.name
      );
      await this.customObjectsApi.patchNamespacedCustomObject(
        this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        ),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        this.operatorService.kind2plural(KustomizeResourceKind.Kustomization),
        plugResource.metadata.name,
        [
          {
            op: 'replace',
            path: '/spec',
            value: plugResource.spec?.kustomization
          }
        ],
        undefined,
        undefined,
        undefined,
        {
          headers: { 'Content-Type': 'application/json-patch+json' }
        }
      );
    } catch (err) {
      if (err.statusCode !== 404) throw err;
      const ns = plugResource.metadata.namespace;
      const kustomizationResource: KustomizationResource = {
        apiVersion: `${this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        )}/${KustomizeResourceVersion.V1alpha1}`,
        kind: KustomizeResourceKind.Kustomization,
        metadata: {
          name: plugResource.metadata.name,
          namespace: ns,
          ownerReferences: [
            this.operatorService.getOwnerReference(plugResource, ns)
          ]
        },
        spec: plugResource.spec?.kustomization
      };
      await this.customObjectsApi.createNamespacedCustomObject(
        this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        ),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        this.operatorService.kind2plural(KustomizeResourceKind.Kustomization),
        kustomizationResource
      );
    }
  }
}

export enum Hook {
  AfterCleanup = 'after-cleanup',
  AfterCreate = 'after-create',
  AfterCreateOrUpdate = 'after-create-or-update',
  AfterUpdate = 'after-update',
  BeforeCleanup = 'before-cleanup',
  BeforeCreate = 'before-create',
  BeforeCreateOrUpdate = 'before-create-or-update',
  BeforeUpdate = 'before-update',
  Cleanup = 'cleanup',
  Create = 'create',
  CreateOrUpdate = 'create-or-update',
  Update = 'update'
}

export interface HookResult {
  hookName: string;
  message: string;
  name: string;
  namespace: string;
}
