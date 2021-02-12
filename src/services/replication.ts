import chalk from 'chalk';
import ora from 'ora';
import {
  KubernetesListObject,
  KubernetesObject
} from '@kubernetes/client-node';
import { Replication } from '~/types';
import { ResourceGroup } from '~/integrationOperator';
import Kubectl, { Output } from './kubectl';
import OperatorService from './operator';

export default class ReplicationService {
  private spinner = ora();

  private kubectl = new Kubectl();

  private operatorService = new OperatorService();

  private replicatedFromAnnotationKey: string;

  constructor(private fromNamespace: string) {
    this.replicatedFromAnnotationKey = `${this.operatorService.getGroupName(
      ResourceGroup.Integration
    )}/replicated-from`;
  }

  async apply(
    replication: Replication,
    toNamespace: string,
    owner?: KubernetesObject
  ) {
    const fromResource = await this.getFromResource(replication);
    if (!fromResource) {
      throw new Error(
        `${this.operatorService.getFullName({
          group: replication.group,
          kind: replication.kind,
          name: replication.name
        })} not found in ns ${chalk.blueBright.bold(this.fromNamespace)}`
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

  async cleanupToResources(replicationFrom: Replication, toNamespace: string) {
    const resources =
      (
        await this.kubectl.get<KubernetesListObject<KubernetesObject>>({
          stdin: {
            apiVersion: this.operatorService.getApiVersion(
              replicationFrom.version || '',
              replicationFrom.group
            ),
            kind: replicationFrom?.kind,
            metadata: {
              name: replicationFrom?.name,
              namespace: toNamespace
            }
          },
          output: Output.Json
        })
      )?.items || [];
    await Promise.all(
      resources.map(async (resource: KubernetesObject) => {
        if (
          (resource.metadata?.annotations || {})[
            this.replicatedFromAnnotationKey
          ] === `${replicationFrom?.name}.${this.fromNamespace}`
        ) {
          await this.kubectl.delete({
            stdin: {
              apiVersion: resource.apiVersion,
              kind: resource.kind,
              metadata: {
                name: resource.metadata?.name,
                namespace: resource.metadata?.namespace
              }
            }
          });
          this.spinner.succeed(
            `deleted replicated resource ${this.operatorService.getFullName({
              resource
            })}`
          );
        }
      })
    );
  }

  private async replicateTo(
    fromResource: KubernetesObject,
    toNamespace: string,
    owner?: KubernetesObject,
    ownerReferences = false
  ) {
    if (typeof fromResource.metadata?.namespace === 'undefined') {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: fromResource
        })} ns is not defined`
      );
    }
    const annotations = fromResource.metadata?.annotations || {};
    annotations[
      this.replicatedFromAnnotationKey
    ] = `${fromResource.metadata.name}.${fromResource.metadata.namespace}`;
    await this.kubectl.apply({
      stdin: {
        ...fromResource,
        metadata: {
          annotations,
          labels: fromResource.metadata?.labels,
          name: fromResource.metadata?.name || '',
          namespace: toNamespace,
          ...(typeof owner !== 'undefined' &&
          owner.metadata?.namespace === toNamespace &&
          ownerReferences
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
            namespace: this.fromNamespace
          }
        },
        output: Output.Json
      })
    )?.items || [])?.[0];
  }
}
