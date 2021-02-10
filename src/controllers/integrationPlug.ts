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
import newRegExp from 'newregexp';
import ora from 'ora';
import { KustomizationResource } from 'kustomize-operator';
import { ResourceMeta } from '@dot-i/k8s-operator';
import { Replicate, Kubectl } from '~/services';
import {
  kind2plural,
  getGroupName,
  resources2String,
  getApiVersion
} from '~/util';
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

  private spinner = ora();

  private kubectl = new Kubectl();

  constructor(groupnameprefix: string, kind: string) {
    super(groupnameprefix, kind);
    this.kubeConfig = new k8s.KubeConfig();
    this.kubeConfig.loadFromDefault();
    this.batchV1Api = this.kubeConfig.makeApiClient(k8s.BatchV1Api);
    this.coreV1Api = this.kubeConfig.makeApiClient(k8s.CoreV1Api);
    this.customObjectsApi = this.kubeConfig.makeApiClient(k8s.CustomObjectsApi);
  }

  private static base64DecodeSecretData(
    data: HashMap<string> = {}
  ): HashMap<string> {
    return Object.entries(data).reduce(
      (
        stringData: HashMap<string>,
        [key, base64EncodedValue]: [string, string]
      ) => {
        stringData[key] = Buffer.from(base64EncodedValue, 'base64').toString(
          'utf8'
        );
        return stringData;
      },
      {}
    );
  }

  async deleted(
    plugResource: IntegrationPlugResource,
    _meta: ResourceMeta,
    _oldPlugResource?: IntegrationPlugResource
  ) {
    const socketResource = await this.getSocketResource(plugResource);
    if (!socketResource) {
      this.spinner.warn(
        `integrationsocket/${plugResource.spec?.socket?.name} does not exist in namespace ${plugResource.spec?.socket?.namespace}`
      );
      return null;
    }
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
      this.spinner.warn(
        `integrationsocket/${plugResource.spec?.socket?.name} does not exist in namespace ${plugResource.spec?.socket?.namespace}`
      );
      return null;
    }
    try {
      await this.updateStatus(
        {
          message: 'pending',
          phase: IntegrationPlugStatusPhase.Pending,
          ready: false
        },
        plugResource
      );
      await Promise.all([
        this.callHook(Hook.BeforeCreate, plugResource, socketResource),
        this.callHook(Hook.BeforeCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.waitForResources(plugResource, socketResource);
      await this.copyAndMergeConfigmaps(plugResource, socketResource);
      await this.copyAndMergeSecrets(plugResource, socketResource);
      await this.replicateSocketResources(socketResource);
      await this.replicatePlugResources(plugResource);
      const [createResult, createOrUpdateResult] = await Promise.all([
        this.callHook(Hook.Create, plugResource, socketResource),
        this.callHook(Hook.CreateOrUpdate, plugResource, socketResource)
      ]);
      const statusMessage = [...createResult, ...createOrUpdateResult]
        .map(
          ({ name, namespace, message }: HookResult) =>
            `${name} ${namespace} ${message}`
        )
        .join('\n');
      if (plugResource?.spec?.kustomization) {
        await this.applyKustomization(plugResource);
      }
      await Promise.all([
        this.callHook(Hook.AfterCreate, plugResource, socketResource),
        this.callHook(Hook.AfterCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.updateStatus(
        {
          message: statusMessage,
          phase: IntegrationPlugStatusPhase.Succeeded,
          ready: true
        },
        plugResource
      );
    } catch (err) {
      await this.updateStatus(
        {
          message: err.message?.toString() || '',
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
      this.spinner.warn(
        `integrationsocket/${plugResource.spec?.socket?.name} does not exist in namespace ${plugResource.spec?.socket?.namespace}`
      );
      return null;
    }
    try {
      await this.updateStatus(
        {
          message: 'pending',
          phase: IntegrationPlugStatusPhase.Pending,
          ready: false
        },
        plugResource
      );
      await Promise.all([
        this.callHook(Hook.BeforeUpdate, plugResource, socketResource),
        this.callHook(Hook.BeforeCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.waitForResources(plugResource, socketResource);
      await this.copyAndMergeConfigmaps(plugResource, socketResource);
      await this.copyAndMergeSecrets(plugResource, socketResource);
      await this.replicateSocketResources(socketResource);
      await this.replicatePlugResources(plugResource);
      const [updateResult, createOrUpdateResult] = await Promise.all([
        this.callHook(Hook.Update, plugResource, socketResource),
        this.callHook(Hook.CreateOrUpdate, plugResource, socketResource)
      ]);
      const statusMessage = [...updateResult, ...createOrUpdateResult]
        .map(
          ({ name, namespace, message }: HookResult) =>
            `${name} ${namespace} ${message}`
        )
        .join('\n');
      if (plugResource?.spec?.kustomization) {
        await this.applyKustomization(plugResource);
      }
      await Promise.all([
        this.callHook(Hook.AfterUpdate, plugResource, socketResource),
        this.callHook(Hook.AfterCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.updateStatus(
        {
          message: statusMessage,
          phase: IntegrationPlugStatusPhase.Succeeded,
          ready: true
        },
        plugResource
      );
    } catch (err) {
      await this.updateStatus(
        {
          message: err.message?.toString() || '',
          phase: IntegrationPlugStatusPhase.Failed,
          ready: false
        },
        plugResource
      );
      throw err;
    }
    return null;
  }

  private async getResources(
    resources: k8s.KubernetesObject[]
  ): k8s.KubernetesObject[] {
    const resourcesStr = resources2String(resources);
    const resources =
      (
        await this.kubectl.get<KubernetesListObject<KubernetesObject>>({
          stdin: resourcesStr,
          output: Output.Json
        })
      )?.items || [];
    return resources;
  }

  private async waitForResources(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    const timeout = socketResource?.spec?.wait?.timeout || 60000;
    const interval = Math.max(
      socketResource?.spec?.wait?.interval || 5000,
      timeout / 10
    );
    try {
      await this.getResources(
        (socketResource.spec?.wait?.resources || []).map<k8s.KubernetesObject>(
          (resource: IntegrationPlugSpecWaitResource) => ({
            apiVersion: getApiVersion(resource.version, resource.group),
            kind: resource.kind,
            metadata: {
              name: resource.name,
              namespace: plugResource.metadata?.namespace
            }
          })
        )
      );
    } catch (err) {
      await new Promise((r) => setTimeout(r, interval));
      await this.waitForResources(
        plugResource,
        socketResource,
        timeout - interval,
        interval
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
              ownerReferences: [this.getOwnerReference(plugResource, ns)]
            },
            spec: hook.job
          })
        ).body;
        await this.waitForJobToFinish(job, hook.timeout, hook.interval);
        const logs = await this.getJobLogs(job);
        let message = '';
        if (hook.messageRegex) {
          const messageMatches = logs.match(newRegExp(hook.messageRegex));
          message = [...(messageMatches || [])].join('\n');
        }
        return {
          name: job.metadata?.name!,
          namespace: job.metadata?.namespace!,
          message
        };
      })
    );
  }

  private async waitForJobToFinish(
    job: k8s.V1Job,
    timeout = 60000,
    interval = 5000
  ) {
    interval = Math.max(timeout / 10, interval);
    const jobStatus = (
      await this.batchV1Api.readNamespacedJobStatus(
        job.metadata?.name!,
        job.metadata?.namespace!
      )
    ).body.status;
    if (jobStatus?.succeeded) return;
    await new Promise((r) => setTimeout(r, interval));
    await this.waitForJobToFinish(job, timeout - interval, interval);
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
    const replicate = new Replicate(plugResource.metadata?.namespace!);
    await Promise.all(
      (
        plugResource.spec?.replications || []
      ).map(async (replication: Replication) => replicate.apply(replication))
    );
  }

  private async replicateSocketResources(
    socketResource: IntegrationSocketResource
  ) {
    const replicate = new Replicate(socketResource.metadata?.namespace!);
    await Promise.all(
      (
        socketResource.spec?.replications || []
      ).map(async (replication: Replication) => replicate.apply(replication))
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
            IntegrationPlug.base64DecodeSecretData(plugSecret?.data)
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
        let mergedStringData = IntegrationPlug.base64DecodeSecretData(
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
                ownerReferences: [this.getOwnerReference(owner, namespace)]
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
                ownerReferences: [this.getOwnerReference(owner, namespace)]
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
          kind2plural(ResourceKind.IntegrationSocket),
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
    if (!plugResource.metadata?.name || !plugResource.metadata.namespace)
      return;
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
    if (!plugResource.metadata?.name || !plugResource.metadata.namespace)
      return;
    try {
      await this.customObjectsApi.getNamespacedCustomObject(
        getGroupName(KustomizeResourceGroup.Kustomize, 'siliconhills.dev'),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        kind2plural(KustomizeResourceKind.Kustomization),
        plugResource.metadata.name
      );
      await this.customObjectsApi.patchNamespacedCustomObject(
        getGroupName(KustomizeResourceGroup.Kustomize, 'siliconhills.dev'),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        kind2plural(KustomizeResourceKind.Kustomization),
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
        apiVersion: `${getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        )}/${KustomizeResourceVersion.V1alpha1}`,
        kind: KustomizeResourceKind.Kustomization,
        metadata: {
          name: plugResource.metadata.name,
          namespace: ns,
          ownerReferences: [this.getOwnerReference(plugResource, ns)]
        },
        spec: plugResource.spec?.kustomization
      };
      await this.customObjectsApi.createNamespacedCustomObject(
        getGroupName(KustomizeResourceGroup.Kustomize, 'siliconhills.dev'),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        kind2plural(KustomizeResourceKind.Kustomization),
        kustomizationResource
      );
    }
  }

  private getOwnerReference(
    owner: k8s.KubernetesObject,
    childNamespace: string
  ) {
    const ownerNamespace = owner.metadata?.namespace;
    if (!childNamespace) {
      throw new Error(
        `cluster-scoped resource must not have a namespace-scoped owner, owner's namespace ${ownerNamespace}`
      );
    }
    if (ownerNamespace !== childNamespace) {
      throw new Error(
        `cross-namespace owner references are disallowed, owner's namespace ${ownerNamespace}, obj's namespace ${childNamespace}`
      );
    }
    return {
      apiVersion: owner?.apiVersion!,
      blockOwnerDeletion: true,
      controller: true,
      kind: owner?.kind!,
      name: owner?.metadata?.name!,
      uid: owner?.metadata?.uid!
    };
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
  message: string;
  name: string;
  namespace: string;
}
