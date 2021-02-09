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
import { KustomizationResource } from 'kustomize-operator';
import { ResourceMeta } from '@dot-i/k8s-operator';
import { IntegrationPlugResource, IntegrationSocketResource } from '~/types';
import { kind2plural, getGroupName } from '~/util';
import {
  KustomizeResourceGroup,
  KustomizeResourceKind,
  KustomizeResourceVersion
} from '~/integrationOperator';
import Controller from './controller';

export default class ExternalMongo extends Controller {
  coreV1Api: k8s.CoreV1Api;

  customObjectsApi: k8s.CustomObjectsApi;

  kubeConfig: k8s.KubeConfig;

  constructor(groupnameprefix: string, kind: string) {
    super(groupnameprefix, kind);
    this.kubeConfig = new k8s.KubeConfig();
    this.kubeConfig.loadFromDefault();
    this.coreV1Api = this.kubeConfig.makeApiClient(k8s.CoreV1Api);
    this.customObjectsApi = this.kubeConfig.makeApiClient(k8s.CustomObjectsApi);
  }

  async addedOrModified(
    resource: IntegrationPlugResource,
    _meta: ResourceMeta,
    oldResource?: IntegrationPlugResource
  ): Promise<any> {
    if (resource.metadata?.generation === oldResource?.metadata?.generation) {
      return null;
    }
    const socketResource = await this.getSocket();
    await this.replicateConfigmaps();
    await this.replicateSecrets();
    await this.replicateResources();
    return null;
  }

  async replicateConfigmaps() {}

  async replicateSecrets() {}

  async replicateResources() {}

  async getSocket(): Promise<IntegrationSocketResource> {
    return {} as IntegrationSocketResource;
  }

  async applyKustomization(resource: IntegrationPlugResource): Promise<void> {
    if (!resource.metadata?.name || !resource.metadata.namespace) return;
    try {
      await this.customObjectsApi.getNamespacedCustomObject(
        getGroupName(KustomizeResourceGroup.Kustomize, 'siliconhills.dev'),
        KustomizeResourceVersion.V1alpha1,
        resource.metadata.namespace,
        kind2plural(KustomizeResourceKind.Kustomization),
        resource.metadata.name
      );
      await this.customObjectsApi.patchNamespacedCustomObject(
        getGroupName(KustomizeResourceGroup.Kustomize, 'siliconhills.dev'),
        KustomizeResourceVersion.V1alpha1,
        resource.metadata.namespace,
        kind2plural(KustomizeResourceKind.Kustomization),
        resource.metadata.name,
        [
          {
            op: 'replace',
            path: '/spec',
            value: resource.spec?.kustomization
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
          name: resource.metadata.name,
          namespace: resource.metadata.namespace
        },
        spec: resource.spec?.kustomization
      };
      await this.customObjectsApi.createNamespacedCustomObject(
        getGroupName(KustomizeResourceGroup.Kustomize, 'siliconhills.dev'),
        KustomizeResourceVersion.V1alpha1,
        resource.metadata.namespace,
        kind2plural(KustomizeResourceKind.Kustomization),
        kustomizationResource
      );
    }
  }
}
