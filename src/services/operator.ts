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

import YAML from 'yaml';
import chalk from 'chalk';
import fs from 'fs-extra';
import path from 'path';
import { KubernetesObject } from '@kubernetes/client-node';
import { HashMap, OperatorFrameworkProject } from '~/types';

export default class OperatorService {
  project: OperatorFrameworkProject = YAML.parse(
    fs.readFileSync(path.resolve(__dirname, '../../PROJECT')).toString()
  );

  getFullName({
    apiVersion,
    group,
    hideGroup,
    kind,
    name,
    ns,
    resource
  }: {
    apiVersion?: string;
    group?: string;
    hideGroup?: boolean;
    kind?: string;
    name?: string;
    ns?: string;
    resource?: KubernetesObject;
  }) {
    if (typeof hideGroup === 'undefined') hideGroup = true;
    if (resource) {
      ({ kind, apiVersion } = resource);
      name = resource.metadata?.name;
      ns = resource.metadata?.namespace;
    }
    if (apiVersion) {
      const splitApiVersion = apiVersion.split('/');
      group = splitApiVersion.length > 1 ? splitApiVersion[0] : undefined;
    }
    if (hideGroup) group = undefined;
    return `${chalk.yellow.bold(
      `${kind ? `${this.getFullType(kind, group)}/` : ''}${name}`
    )}${ns ? ` in ns ${chalk.blueBright.bold(ns)}` : ''}`;
  }

  kind2plural(kind: string) {
    let lowercasedKind = kind.toLowerCase();
    if (lowercasedKind[lowercasedKind.length - 1] === 's') {
      return lowercasedKind;
    }
    if (lowercasedKind[lowercasedKind.length - 1] === 'o') {
      lowercasedKind = `${lowercasedKind}e`;
    }
    return `${lowercasedKind}s`;
  }

  getGroupName(groupNamePrefix: string, domain?: string) {
    return `${groupNamePrefix}.${domain || this.project.domain}`;
  }

  getApiVersion(version: string, group?: string): string {
    return `${group ? `${group}/` : ''}${version}`;
  }

  getFullType(kind: string, group?: string): string {
    return `${this.kind2plural(kind)}${group ? `.${group}` : ''}`;
  }

  getOwnerReference(owner: KubernetesObject, childNamespace: string) {
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

  base64DecodeSecretData(data: HashMap<string> = {}): HashMap<string> {
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

  getErrorMessage(err: any) {
    return [
      err.message || '',
      err.body?.message || err.response?.body?.message || ''
    ].join(': ');
  }
}
