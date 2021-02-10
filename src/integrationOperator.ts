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
    this.watchController(
      ResourceKind.IntegrationSocket,
      new IntegrationSocket(
        ResourceGroup.Integration,
        ResourceKind.IntegrationSocket
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
            this.spinner.fail(
              [
                err.message || '',
                err.body?.message || err.response?.body?.message || ''
              ].join(': ')
            );
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
