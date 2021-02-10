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
import { Replicate } from '~/services';
import { kind2plural, getGroupName } from '~/util';
import {
  HashMap,
  IntegrationPlugResource,
  IntegrationPlugSpecMergeConfigmaps,
  IntegrationPlugSpecMergeSecrets,
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
  coreV1Api: k8s.CoreV1Api;

  customObjectsApi: k8s.CustomObjectsApi;

  kubeConfig: k8s.KubeConfig;

  batchV1Api: k8s.BatchV1Api;

  spinner = ora();

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
    await this.callHook(Hook.Cleanup, plugResource, socketResource);
    return null;
  }

  async addedOrModified(
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
      await this.callHook(Hook.Before, plugResource, socketResource);
      await this.copyAndMergeConfigmaps(plugResource, socketResource);
      await this.copyAndMergeSecrets(plugResource, socketResource);
      await this.replicateSocketResources(socketResource);
      await this.replicatePlugResources(plugResource);
      const integrateHookResults = await this.callHook(
        Hook.Integrate,
        plugResource,
        socketResource
      );
      const statusMessage = integrateHookResults
        .map(
          ({ name, namespace, message }: HookResult) =>
            `${name} ${namespace} ${message}`
        )
        .join('\n');
      if (plugResource?.spec?.kustomization) {
        await this.applyKustomization(plugResource);
      }
      await this.callHook(Hook.After, plugResource, socketResource);
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
              namespace: ns
            },
            spec: hook.job
          })
        ).body;
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
            mergedData
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
          mergedStringData
        );
      })
    );
  }

  private async createOrUpdateSecret(
    name: string,
    namespace: string,
    data: HashMap<string>
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
          namespace
        },
        stringData: data
      });
    }
  }

  private async createOrUpdateConfigMap(
    name: string,
    namespace: string,
    data: HashMap<string>
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
          namespace
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
      const kustomizationResource: KustomizationResource = {
        apiVersion: `${getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        )}/${KustomizeResourceVersion.V1alpha1}`,
        kind: KustomizeResourceKind.Kustomization,
        metadata: {
          name: plugResource.metadata.name,
          namespace: plugResource.metadata.namespace
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
}

export enum Hook {
  After = 'after',
  Before = 'before',
  Cleanup = 'cleanup',
  Integrate = 'integrate'
}

export interface HookResult {
  message: string;
  name: string;
  namespace: string;
}
