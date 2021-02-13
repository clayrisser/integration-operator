import Operator, { ResourceEventType } from '@dot-i/k8s-operator';
import ora from 'ora';
import { Config } from '~/config';
import { Controller, IntegrationPlug, IntegrationSocket } from '~/controllers';
import { OperatorService, ResourceTrackerService } from '~/services';

const logger = console;

export default class IntegrationOperator extends Operator {
  spinner = ora();

  private resourceTrackerService = new ResourceTrackerService();

  private operatorService = new OperatorService();

  constructor(protected config: Config) {
    super(logger);
  }

  protected async init() {
    this.watchController(
      ResourceKind.IntegrationPlug,
      new IntegrationPlug(
        ResourceGroup.Integration,
        ResourceKind.IntegrationPlug
      )
    );
  }

  protected watchController(
    resourceKind: ResourceKind,
    controller: Controller
  ) {
    this.watchResource(
      this.operatorService.getGroupName(ResourceGroup.Integration),
      ResourceVersion.V1alpha1,
      this.operatorService.kind2plural(resourceKind),
      async (e) => {
        // spawn as non blocking process
        (async () => {
          const {
            oldResource,
            newResource
          } = this.resourceTrackerService.rotateResource(e.object);
          try {
            switch (e.type) {
              case ResourceEventType.Added:
                await controller.added(newResource, e.meta, oldResource);
                await controller.addedOrModified(
                  newResource,
                  e.meta,
                  oldResource
                );
                return;
              case ResourceEventType.Deleted:
                this.resourceTrackerService.resetResource(e.object);
                await controller.deleted(newResource, e.meta, oldResource);
                return;
              case ResourceEventType.Modified:
                await controller.modified(newResource, e.meta, oldResource);
                await controller.addedOrModified(
                  newResource,
                  e.meta,
                  oldResource
                );
                return;
            }
          } catch (err) {
            this.spinner.fail(this.operatorService.getErrorMessage(err));
            if (this.config.debug) logger.error(err);
          }
        })().catch(logger.error);
      }
    ).catch(logger.error);
  }
}

export enum ResourceGroup {
  Integration = 'integration'
}

export enum ResourceKind {
  IntegrationPlug = 'IntegrationPlug',
  IntegrationSocket = 'IntegrationSocket'
}

export enum ResourceVersion {
  V1alpha1 = 'v1alpha1'
}

export enum KustomizeResourceGroup {
  Kustomize = 'kustomize'
}

export enum KustomizeResourceKind {
  Kustomization = 'Kustomization'
}

export enum KustomizeResourceVersion {
  V1alpha1 = 'v1alpha1'
}
