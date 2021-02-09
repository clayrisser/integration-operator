import { KubernetesObject } from '@kubernetes/client-node';
import { HashMap } from '~/types';

export default class ResourceTracker<T = KubernetesObject> {
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
