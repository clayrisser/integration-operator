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
import { KubernetesObject } from '@kubernetes/client-node';
import { ResourceMeta } from '@dot-i/k8s-operator';
import { OperatorService } from '~/services';

export default abstract class Controller {
  constructor(protected groupNamePrefix: string, protected kind: string) {
    this.group = this.operatorService.getGroupName(this.groupNamePrefix);
    this.plural = this.operatorService.kind2plural(this.kind);
  }

  protected operatorService = new OperatorService();

  protected group: string;

  protected plural: string;

  spinner = ora();

  async added(
    _resource: KubernetesObject,
    _meta: ResourceMeta,
    _oldResource?: KubernetesObject
  ): Promise<any> {}

  async addedOrModified(
    _resource: KubernetesObject,
    _meta: ResourceMeta,
    _oldResource?: KubernetesObject
  ): Promise<any> {}

  async deleted(
    _resource: KubernetesObject,
    _meta: ResourceMeta,
    _oldResource?: KubernetesObject
  ): Promise<any> {}

  async modified(
    _resource: KubernetesObject,
    _meta: ResourceMeta,
    _oldResource?: KubernetesObject
  ): Promise<any> {}
}
