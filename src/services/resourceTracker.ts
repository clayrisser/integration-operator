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

import { KubernetesObject } from '@kubernetes/client-node';
import { HashMap } from '~/types';

export default class ResourceTrackerService<T = KubernetesObject> {
  oldResources: HashMap<T | undefined> = {};

  getResourceId(resource: T) {
    return `${(resource as KubernetesObject)?.apiVersion || ''}:${
      (resource as KubernetesObject)?.kind || ''
    }:${(resource as KubernetesObject)?.metadata?.namespace || ''}:${
      (resource as KubernetesObject)?.metadata?.name || ''
    }`;
  }

  rotateResource(resource: T): ResourcePair<T> {
    const id = this.getResourceId(resource);
    const oldResource = this.oldResources[id];
    const newResource = resource;
    this.oldResources[id] = resource;
    return { oldResource, newResource };
  }

  resetResource(resource: T) {
    const id = this.getResourceId(resource);
    delete this.oldResources[id];
  }
}

export interface ResourcePair<T = KubernetesObject> {
  oldResource?: T;
  newResource: T;
}
