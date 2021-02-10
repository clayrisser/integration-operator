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

import ora from 'ora';
import {
  KubernetesListObject,
  KubernetesObject
} from '@kubernetes/client-node';
import { Replication, ReplicationFrom, ReplicationTo } from '~/types';
import Kubectl, { Output } from './kubectl';
import OperatorService from './operator';

export default class ReplicationService {
  private spinner = ora();

  private kubectl = new Kubectl();

  private operatorService = new OperatorService();

  constructor(private namespace: string) {}

  async apply(replication: Replication, owner?: KubernetesObject) {
    if (!replication.from) {
      throw new Error('replication from not defined');
    }
    if (!replication.to) {
      throw new Error('replication to not defined');
    }
    const fromResource = await this.getFromResource(replication.from);
    if (!fromResource) {
      throw new Error(
        `from resource ${this.operatorService.getFullType(
          replication.from.kind || '',
          replication.from.version || '',
          replication.from.group
        )}/${replication.from.name || ''} not found in namespace ${
          this.namespace
        }`
      );
    }
    const status = await this.replicateTo(fromResource, replication.to, owner);
    const fullType = this.operatorService.getFullType(
      fromResource.kind || '',
      fromResource.apiVersion || ''
    );
    this.spinner.succeed(
      `replicated resource ${fullType}/${
        fromResource.metadata?.name || ''
      } from namespace ${this.namespace} to ${fullType}/${
        status.name
      } in namespace ${status.namespace}`
    );
  }

  private async replicateTo(
    fromResource: KubernetesObject,
    replicationTo: ReplicationTo,
    owner?: KubernetesObject
  ) {
    const name = replicationTo.name || fromResource.metadata?.name;
    const ns = replicationTo.namespace;
    if (
      typeof name === 'undefined' ||
      typeof ns === 'undefined' ||
      !name ||
      !ns
    ) {
      const fullType = this.operatorService.getFullType(
        fromResource.kind || '',
        fromResource.apiVersion || ''
      );
      throw new Error(
        `cannot replicate ${fullType}/${
          fromResource.metadata?.name || ''
        } from namespace ${
          this.namespace
        } to ${fullType}/${name} in namespace ${ns}`
      );
    }
    await this.kubectl.apply({
      stdin: {
        ...fromResource,
        metadata: {
          name,
          namespace: ns,
          ...(typeof owner !== 'undefined' && owner.metadata?.namespace === ns
            ? {
                ownerReferences: [
                  this.operatorService.getOwnerReference(owner, ns)
                ]
              }
            : {})
        }
      },
      stdout: true
    });
    return { name, namespace: ns };
  }

  private async getFromResource(
    replicationFrom: ReplicationFrom
  ): Promise<KubernetesObject | null> {
    return (
      ((
        await this.kubectl.get<KubernetesListObject<KubernetesObject>>({
          stdin: {
            apiVersion: this.operatorService.getApiVersion(
              replicationFrom.version || '',
              replicationFrom.group
            ),
            kind: replicationFrom?.kind,
            metadata: {
              name: replicationFrom?.name,
              namespace: this.namespace
            }
          },
          output: Output.Json
        })
      )?.items || [])?.[0] || null
    );
  }
}
