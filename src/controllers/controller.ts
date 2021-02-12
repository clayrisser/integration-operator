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
