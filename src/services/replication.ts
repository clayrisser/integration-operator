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

import chalk from 'chalk';
import ora from 'ora';
import {
  KubernetesListObject,
  KubernetesObject
} from '@kubernetes/client-node';
import { Replication } from '~/types';
import Kubectl, { Output } from './kubectl';
import OperatorService from './operator';

export default class ReplicationService {
  private spinner = ora();

  private kubectl = new Kubectl();

  private operatorService = new OperatorService();

  constructor(private namespace: string) {}

  async apply(
    replication: Replication,
    toNamespace?: string,
    owner?: KubernetesObject
  ) {
    const fromResource = await this.getFromResource(replication);
    if (!fromResource) {
      throw new Error(
        `from resource ${this.operatorService.getFullName({
          resource: fromResource
        })} not found in ns ${chalk.blueBright.bold(this.namespace)}`
      );
    }
    await this.replicateTo(fromResource, toNamespace, owner);
    this.spinner.succeed(
      `replicated ${this.operatorService.getFullName({
        resource: fromResource
      })} to ${this.operatorService.getFullName({
        apiVersion: fromResource.apiVersion,
        name: fromResource.metadata?.name,
        kind: fromResource.kind,
        ns: toNamespace
      })}`
    );
  }

  private async replicateTo(
    fromResource: KubernetesObject,
    toNamespace?: string,
    owner?: KubernetesObject
  ) {
    if (typeof toNamespace === 'undefined' || !toNamespace) {
      const fullType = this.operatorService.getFullType(
        fromResource.kind || '',
        fromResource.apiVersion || ''
      );
      throw new Error(
        `cannot replicate ${fullType}/${
          fromResource.metadata?.name || ''
        } from namespace ${this.namespace} to ${fullType}/${
          fromResource.metadata?.name || ''
        } in namespace ${toNamespace}`
      );
    }
    await this.kubectl.apply({
      stdin: {
        ...fromResource,
        metadata: {
          name: fromResource.metadata?.name || '',
          namespace: toNamespace,
          ...(typeof owner !== 'undefined' &&
          owner.metadata?.namespace === toNamespace
            ? {
                ownerReferences: [
                  this.operatorService.getOwnerReference(owner, toNamespace)
                ]
              }
            : {})
        }
      },
      stdout: true
    });
  }

  private async getFromResource(
    replicationFrom: Replication
  ): Promise<KubernetesObject | undefined> {
    return ((
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
    )?.items || [])?.[0];
  }
}
