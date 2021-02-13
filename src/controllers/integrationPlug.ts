import * as k8s from '@kubernetes/client-node';
import chalk from 'chalk';
import newRegExp from 'newregexp';
import stripAnsi from 'strip-ansi';
import { KustomizationResource } from 'kustomize-operator';
import { ResourceMeta } from '@dot-i/k8s-operator';
import { ReplicationService, KubectlService, Output } from '~/services';
import {
  HashMap,
  IntegrationPlugResource,
  IntegrationPlugSpecWaitResource,
  IntegrationPlugStatus,
  IntegrationPlugStatusPhase,
  IntegrationSocketResource,
  IntegrationSocketSpecHook,
  Replication
} from '~/types';
import {
  KustomizeResourceGroup,
  KustomizeResourceKind,
  KustomizeResourceVersion,
  ResourceKind,
  ResourceVersion
} from '~/integrationOperator';
import Controller from './controller';

export default class IntegrationPlug extends Controller {
  private coreV1Api: k8s.CoreV1Api;

  private customObjectsApi: k8s.CustomObjectsApi;

  private kubeConfig: k8s.KubeConfig;

  private batchV1Api: k8s.BatchV1Api;

  private kubectl = new KubectlService();

  private tracking: HashMap<string | true> = {};

  private runningHooks: Set<string> = new Set();

  constructor(groupnameprefix: string, kind: string) {
    super(groupnameprefix, kind);
    this.kubeConfig = new k8s.KubeConfig();
    this.kubeConfig.loadFromDefault();
    this.batchV1Api = this.kubeConfig.makeApiClient(k8s.BatchV1Api);
    this.coreV1Api = this.kubeConfig.makeApiClient(k8s.CoreV1Api);
    this.customObjectsApi = this.kubeConfig.makeApiClient(k8s.CustomObjectsApi);
  }

  async deleted(
    plugResource: IntegrationPlugResource,
    _meta: ResourceMeta,
    _oldPlugResource?: IntegrationPlugResource
  ) {
    const socketResource = await this.getSocketResource(plugResource);
    if (!socketResource) return null;
    await this.cleanupReplicatedPlugResources(plugResource, socketResource);
    if (plugResource.spec?.cleanup) {
      await this.callHook(Hook.BeforeCleanup, plugResource, socketResource);
      await this.callHook(Hook.Cleanup, plugResource, socketResource);
      await this.callHook(Hook.AfterCleanup, plugResource, socketResource);
    }
    this.unregisterTracking(plugResource);
    return null;
  }

  async added(
    plugResource: IntegrationPlugResource,
    _meta: ResourceMeta,
    oldPlugResource?: IntegrationPlugResource
  ): Promise<any> {
    if (
      plugResource.metadata?.generation ===
      oldPlugResource?.metadata?.generation
    ) {
      return null;
    }
    this.registerTracking(plugResource);
    const socketResource = await this.getSocketResource(plugResource);
    if (!socketResource) {
      const message = `${this.operatorService.getFullName({
        kind: ResourceKind.IntegrationSocket,
        name: plugResource.spec?.socket?.name || ''
      })} does not exist in namespace ${chalk.blueBright.bold(
        plugResource.spec?.socket?.namespace
      )}`;
      this.spinner.fail(message);
      await this.updateStatus(
        {
          message,
          phase: IntegrationPlugStatusPhase.Failed,
          ready: false
        },
        plugResource
      );
      return null;
    }
    try {
      await this.beginApply(plugResource, socketResource);
      await Promise.all([
        this.callHook(Hook.BeforeCreate, plugResource, socketResource),
        this.callHook(Hook.BeforeCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.apply(plugResource, socketResource);
      const [createResult, createOrUpdateResult] = await Promise.all([
        this.callHook(Hook.Create, plugResource, socketResource),
        this.callHook(Hook.CreateOrUpdate, plugResource, socketResource)
      ]);
      if (plugResource?.spec?.kustomization) {
        await this.applyKustomization(plugResource);
      }
      await Promise.all([
        this.callHook(Hook.AfterCreate, plugResource, socketResource),
        this.callHook(Hook.AfterCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.endApply(plugResource, socketResource, [
        ...createResult,
        ...createOrUpdateResult
      ]);
    } catch (err) {
      if (this.isTracking(plugResource)) {
        await this.updateStatus(
          {
            message: this.operatorService.getErrorMessage(err),
            phase: IntegrationPlugStatusPhase.Failed,
            ready: false
          },
          plugResource
        );
      }
      throw err;
    }
    return null;
  }

  async modified(
    plugResource: IntegrationPlugResource,
    _meta: ResourceMeta,
    oldPlugResource?: IntegrationPlugResource
  ): Promise<any> {
    if (
      plugResource.metadata?.generation ===
      oldPlugResource?.metadata?.generation
    ) {
      return null;
    }
    this.registerTracking(plugResource);
    const socketResource = await this.getSocketResource(plugResource);
    if (!socketResource) {
      const message = `${this.operatorService.getFullName({
        kind: ResourceKind.IntegrationSocket,
        name: plugResource.spec?.socket?.name || ''
      })} does not exist in namespace ${chalk.blueBright.bold(
        plugResource.spec?.socket?.namespace
      )}`;
      this.spinner.fail(message);
      await this.updateStatus(
        {
          message,
          phase: IntegrationPlugStatusPhase.Failed,
          ready: false
        },
        plugResource
      );
      return null;
    }
    try {
      await this.beginApply(plugResource, socketResource);
      await Promise.all([
        this.callHook(Hook.BeforeUpdate, plugResource, socketResource),
        this.callHook(Hook.BeforeCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.apply(plugResource, socketResource);
      const [updateResult, createOrUpdateResult] = await Promise.all([
        this.callHook(Hook.Update, plugResource, socketResource),
        this.callHook(Hook.CreateOrUpdate, plugResource, socketResource)
      ]);
      if (plugResource?.spec?.kustomization) {
        await this.applyKustomization(plugResource);
      }
      await Promise.all([
        this.callHook(Hook.AfterUpdate, plugResource, socketResource),
        this.callHook(Hook.AfterCreateOrUpdate, plugResource, socketResource)
      ]);
      await this.endApply(plugResource, socketResource, [
        ...updateResult,
        ...createOrUpdateResult
      ]);
    } catch (err) {
      if (this.isTracking(plugResource)) {
        await this.updateStatus(
          {
            message: this.operatorService.getErrorMessage(err),
            phase: IntegrationPlugStatusPhase.Failed,
            ready: false
          },
          plugResource
        );
      }
      throw err;
    }
    return null;
  }

  private async beginApply(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    const message = `integrating with ${this.operatorService.getFullName({
      resource: socketResource
    })}`;
    this.spinner.info(
      `${this.operatorService.getFullName({
        resource: plugResource
      })} is ${message}`
    );
    await this.updateStatus(
      {
        message,
        phase: IntegrationPlugStatusPhase.Pending,
        ready: false
      },
      plugResource
    );
  }

  private async apply(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    if (!this.isTracking(plugResource)) return;
    await this.waitForResources(plugResource, socketResource);
    await this.replicateSocketResources(plugResource, socketResource);
    await this.replicatePlugResources(plugResource, socketResource);
  }

  private async endApply(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource,
    hookResults: HookResult[]
  ) {
    if (!this.isTracking(plugResource)) return;
    const statusMessage = hookResults
      .map(
        ({ name, namespace, message, hookName }: HookResult) =>
          `${message} in ${this.operatorService.getFullName({
            kind: 'Job',
            apiVersion: 'batch/v1',
            name,
            ns: namespace
          })} for hook ${hookName}`
      )
      .join('\n');
    const message = `successfully integrated with ${this.operatorService.getFullName(
      {
        resource: socketResource
      }
    )}`;
    this.spinner.succeed(
      `${this.operatorService.getFullName({
        resource: plugResource
      })} has ${message}`
    );
    await this.updateStatus(
      {
        message: statusMessage || message,
        phase: IntegrationPlugStatusPhase.Succeeded,
        ready: true
      },
      plugResource
    );
  }

  private registerTracking(
    plugResource: IntegrationPlugResource,
    waitingOn?: string
  ) {
    this.tracking[
      `${plugResource.metadata?.name || ''}.${
        plugResource.metadata?.namespace || ''
      }`
    ] = waitingOn || true;
  }

  private unregisterTracking(plugResource: IntegrationPlugResource) {
    const waitingOn = this.tracking[
      `${plugResource.metadata?.name || ''}.${
        plugResource.metadata?.namespace || ''
      }`
    ];
    delete this.tracking[
      `${plugResource.metadata?.name || ''}.${
        plugResource.metadata?.namespace || ''
      }`
    ];
    if (waitingOn !== true) {
      this.spinner.info(
        `stopped waiting on ${waitingOn} for ${this.operatorService.getFullName(
          {
            resource: plugResource
          }
        )}`
      );
    }
  }

  private isTracking(plugResource: IntegrationPlugResource) {
    return !!this.tracking[
      `${plugResource.metadata?.name || ''}.${
        plugResource.metadata?.namespace || ''
      }`
    ];
  }

  private async getResources(
    resources: k8s.KubernetesObject[]
  ): Promise<
    (k8s.KubernetesObject & {
      status?: { [key: string]: any; phase?: string };
    })[]
  > {
    return (
      (
        await this.kubectl.get<k8s.KubernetesListObject<k8s.KubernetesObject>>({
          stdin: resources,
          output: Output.Json
        })
      )?.items || []
    );
  }

  private async waitForResources(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource,
    timeout = 60000,
    timeLeft?: number
  ) {
    if (!this.isTracking(plugResource)) return;
    this.registerTracking(plugResource, 'resources');
    const waitTime = Math.max(5000, timeout / 10);
    if (typeof timeLeft !== 'number') timeLeft = timeout;
    const resources = await this.getResources(
      (socketResource.spec?.wait?.resources || []).map<k8s.KubernetesObject>(
        (waitResource: IntegrationPlugSpecWaitResource) => {
          if (!waitResource.version) {
            throw new Error('resource version is not defined');
          }
          if (!waitResource.kind) {
            throw new Error('resource kind is not defined');
          }
          if (!waitResource.name) {
            throw new Error('resource name is not defined');
          }
          return {
            apiVersion: this.operatorService.getApiVersion(
              waitResource.version,
              waitResource.group
            ),
            kind: waitResource.kind,
            metadata: {
              name: waitResource.name,
              namespace: plugResource.metadata?.namespace
            }
          };
        }
      )
    );
    const foundAllResources = (
      socketResource.spec?.wait?.resources || []
    ).reduce(
      (ready: boolean, waitResource: IntegrationPlugSpecWaitResource) => {
        if (!ready) return ready;
        const resource = resources.find((resource: k8s.KubernetesObject) => {
          return (
            waitResource.name === resource.metadata?.name &&
            waitResource.kind === resource.kind &&
            this.operatorService.getApiVersion(
              waitResource.version!,
              waitResource.group
            ) === resource.apiVersion
          );
        });
        if (typeof resource === 'undefined' || !resource) return false;
        return !!(
          !waitResource.statusPhases?.length ||
          (waitResource.statusPhases || []).find(
            (statusPhase: string) => statusPhase === resource?.status?.phase
          )
        );
      },
      true
    );
    if (foundAllResources || !this.isTracking(plugResource)) return;
    if (timeLeft <= 0) {
      throw new Error(
        `failed to find some resources for ${this.operatorService.getFullName({
          resource: plugResource
        })}`
      );
    }
    this.spinner.info(
      `waiting ${timeLeft}ms on resources for ${this.operatorService.getFullName(
        {
          resource: plugResource
        }
      )}`
    );
    await new Promise((r) => setTimeout(r, waitTime));
    await this.waitForResources(
      plugResource,
      socketResource,
      timeout,
      timeLeft - waitTime
    );
    if (this.isTracking(plugResource)) this.registerTracking(plugResource);
  }

  private async callHook(
    hookName: Hook,
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    if (!this.isTracking(plugResource)) return [];
    this.registerTracking(plugResource, 'jobs');
    this.runningHooks.add(hookName);
    const ns = plugResource.metadata?.namespace!;
    const filteredHooks = (socketResource.spec?.hooks || []).filter(
      (hook: IntegrationSocketSpecHook) => {
        return hook.name === hookName;
      }
    );
    const append = socketResource.spec?.appendName || 'socket';
    const result = await Promise.all(
      filteredHooks.map(async (hook: IntegrationSocketSpecHook, i: number) => {
        const name = `${plugResource.metadata
          ?.name!}-${hookName}-${i.toString()}${append ? `-${append}` : ''}`;
        if (typeof hook.job === 'undefined') {
          throw new Error(`hook ${hookName} job is undefined`);
        }
        const job = await this.createOrUpdateJob(
          name,
          ns,
          hook.job,
          plugResource
        );
        await this.waitForJobToFinish(plugResource, job, hook.timeout);
        if (!this.isTracking(plugResource)) {
          return {
            hookName: hook.name!,
            message: '',
            name: job.metadata?.name!,
            namespace: job.metadata?.namespace!
          };
        }
        const logs = await this.getJobLogs(job);
        let message = '';
        if (hook.messageRegex) {
          const messageMatches = logs.match(newRegExp(hook.messageRegex));
          message = [...(messageMatches || [])].join('\n').trim();
        }
        return {
          hookName: hook.name!,
          message,
          name: job.metadata?.name!,
          namespace: job.metadata?.namespace!
        };
      })
    );
    this.runningHooks.delete(hookName);
    if (this.isTracking(plugResource) && !this.runningHooks.size) {
      this.registerTracking(plugResource);
    }
    return result;
  }

  private async waitForJobToFinish(
    plugResource: IntegrationPlugResource,
    job: k8s.V1Job,
    timeout = 60000,
    timeLeft?: number
  ) {
    if (!this.isTracking(plugResource)) return;
    const waitTime = Math.max(5000, timeout / 10);
    if (typeof timeLeft !== 'number') timeLeft = timeout;
    const jobStatus = (
      await this.batchV1Api.readNamespacedJobStatus(
        job.metadata?.name!,
        job.metadata?.namespace!
      )
    ).body.status;
    if (jobStatus?.succeeded || !this.isTracking(plugResource)) {
      return;
    }
    if (timeLeft <= 0) {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: job
        })} timed out before completing`
      );
    }
    this.spinner.info(
      `waiting ${timeLeft}ms for ${this.operatorService.getFullName({
        resource: job
      })} to complete`
    );
    await new Promise((r) => setTimeout(r, waitTime));
    await this.waitForJobToFinish(
      plugResource,
      job,
      timeout,
      timeLeft - waitTime
    );
  }

  private async getJobLogs(job: k8s.V1Job): Promise<string> {
    const pods = (
      await this.coreV1Api.listNamespacedPod(job.metadata?.namespace || '')
    ).body;
    const podName =
      pods.items.find(
        (pod: k8s.V1Pod) =>
          pod.metadata?.labels?.['job-name'] === job.metadata?.name
      )?.metadata?.name || '';
    return (
      (
        await this.coreV1Api.readNamespacedPodLog(
          podName,
          job.metadata?.namespace || '',
          undefined,
          false
        )
      ).body || ''
    ).toString();
  }

  private async cleanupReplicatedPlugResources(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    const fromNamespace = plugResource.metadata?.namespace;
    const toNamespace = socketResource.metadata?.namespace;
    if (typeof fromNamespace === 'undefined') {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: plugResource
        })} ns is undefined`
      );
    }
    if (typeof toNamespace === 'undefined') {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: socketResource
        })} ns is undefined`
      );
    }
    const replicationService = new ReplicationService(fromNamespace);
    await Promise.all(
      (
        plugResource.spec?.replications || []
      ).map(async (replication: Replication) =>
        replicationService.cleanupToResources(
          replication,
          toNamespace,
          plugResource.spec?.appendName || 'plug'
        )
      )
    );
  }

  private async replicatePlugResources(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    if (!this.isTracking(plugResource)) return;
    const fromNamespace = plugResource.metadata?.namespace;
    const toNamespace = socketResource.metadata?.namespace;
    if (typeof fromNamespace === 'undefined') {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: plugResource
        })} ns is undefined`
      );
    }
    if (typeof toNamespace === 'undefined') {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: socketResource
        })} ns is undefined`
      );
    }
    const replicationService = new ReplicationService(fromNamespace);
    await Promise.all(
      (
        plugResource.spec?.replications || []
      ).map(async (replication: Replication) =>
        replicationService.apply(
          replication,
          toNamespace,
          plugResource.spec?.appendName || 'plug'
        )
      )
    );
  }

  private async replicateSocketResources(
    plugResource: IntegrationPlugResource,
    socketResource: IntegrationSocketResource
  ) {
    if (!this.isTracking(plugResource)) return;
    const fromNamespace = socketResource.metadata?.namespace;
    const toNamespace = plugResource.metadata?.namespace;
    if (typeof fromNamespace === 'undefined') {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: socketResource
        })} ns is undefined`
      );
    }
    if (typeof toNamespace === 'undefined') {
      throw new Error(
        `${this.operatorService.getFullName({
          resource: plugResource
        })} ns is undefined`
      );
    }
    const replicationService = new ReplicationService(fromNamespace);
    await Promise.all(
      (
        socketResource.spec?.replications || []
      ).map(async (replication: Replication) =>
        replicationService.apply(
          replication,
          toNamespace,
          socketResource.spec?.appendName || 'socket',
          plugResource
        )
      )
    );
  }

  private async getSocketResource(
    plugResource: IntegrationPlugResource
  ): Promise<IntegrationSocketResource | null> {
    if (
      !plugResource.metadata?.name ||
      !plugResource.metadata?.namespace ||
      !plugResource.spec?.socket?.name ||
      !this.isTracking(plugResource)
    ) {
      return null;
    }
    const namespace =
      plugResource.spec?.socket?.namespace || plugResource.metadata.namespace;
    try {
      const socketResource = (
        await this.customObjectsApi.getNamespacedCustomObject(
          this.group,
          ResourceVersion.V1alpha1,
          namespace,
          this.operatorService.kind2plural(ResourceKind.IntegrationSocket),
          plugResource.spec.socket.name
        )
      ).body as IntegrationSocketResource;
      return socketResource;
    } catch (err) {
      if (err.statusCode !== 404) throw err;
      return null;
    }
  }

  private async updateStatus(
    plugStatus: IntegrationPlugStatus,
    plugResource: IntegrationPlugResource
  ): Promise<void> {
    if (
      !plugResource.metadata?.name ||
      !plugResource.metadata.namespace ||
      !this.isTracking(plugResource)
    ) {
      return;
    }
    plugStatus.message = stripAnsi(plugStatus.message || '');
    await this.customObjectsApi.patchNamespacedCustomObjectStatus(
      this.group,
      ResourceVersion.V1alpha1,
      plugResource.metadata.namespace,
      this.plural,
      plugResource.metadata.name,
      [
        {
          op: 'replace',
          path: '/status',
          value: plugStatus
        }
      ],
      undefined,
      undefined,
      undefined,
      {
        headers: { 'Content-Type': 'application/json-patch+json' }
      }
    );
  }

  private async applyKustomization(
    plugResource: IntegrationPlugResource
  ): Promise<void> {
    if (
      !plugResource.metadata?.name ||
      !plugResource.metadata.namespace ||
      !this.isTracking(plugResource)
    ) {
      return;
    }
    try {
      await this.customObjectsApi.getNamespacedCustomObject(
        this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        ),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        this.operatorService.kind2plural(KustomizeResourceKind.Kustomization),
        plugResource.metadata.name
      );
      await this.customObjectsApi.patchNamespacedCustomObject(
        this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        ),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        this.operatorService.kind2plural(KustomizeResourceKind.Kustomization),
        plugResource.metadata.name,
        [
          {
            op: 'replace',
            path: '/spec',
            value: plugResource.spec?.kustomization
          }
        ],
        undefined,
        undefined,
        undefined,
        {
          headers: { 'Content-Type': 'application/json-patch+json' }
        }
      );
    } catch (err) {
      if (err.statusCode !== 404) throw err;
      const ns = plugResource.metadata.namespace;
      const kustomizationResource: KustomizationResource = {
        apiVersion: `${this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        )}/${KustomizeResourceVersion.V1alpha1}`,
        kind: KustomizeResourceKind.Kustomization,
        metadata: {
          name: plugResource.metadata.name,
          namespace: ns,
          ownerReferences: [
            this.operatorService.getOwnerReference(plugResource, ns)
          ]
        },
        spec: plugResource.spec?.kustomization
      };
      await this.customObjectsApi.createNamespacedCustomObject(
        this.operatorService.getGroupName(
          KustomizeResourceGroup.Kustomize,
          'siliconhills.dev'
        ),
        KustomizeResourceVersion.V1alpha1,
        plugResource.metadata.namespace,
        this.operatorService.kind2plural(KustomizeResourceKind.Kustomization),
        kustomizationResource
      );
    }
  }

  private async createOrUpdateJob(
    name: string,
    ns: string,
    spec: k8s.V1JobSpec,
    owner?: k8s.KubernetesObject
  ) {
    try {
      let job = (await this.batchV1Api.readNamespacedJob(name, ns)).body;
      job = (
        await this.batchV1Api.patchNamespacedJob(
          name,
          ns,
          [
            {
              op: 'replace',
              path: '/spec',
              value: {
                ...(job.spec || {}),
                ...spec,
                template: {
                  ...(job.spec?.template || {}),
                  ...(spec.template || {}),
                  metadata: {
                    ...(job.spec?.template?.metadata || {}),
                    ...(spec.template?.metadata || {}),
                    labels: {
                      ...(job.spec?.template?.metadata?.labels || {}),
                      ...(spec.template?.metadata?.labels || {})
                    }
                  }
                }
              }
            }
          ],
          undefined,
          undefined,
          undefined,
          undefined,
          {
            headers: { 'Content-Type': 'application/json-patch+json' }
          }
        )
      ).body;
      this.spinner.info(
        `updated ${this.operatorService.getFullName({
          resource: job
        })}`
      );
      return job;
    } catch (err) {
      if (err.statusCode !== 404) throw err;
      const job = (
        await this.batchV1Api.createNamespacedJob(ns, {
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
          },
          spec
        })
      ).body;
      this.spinner.info(
        `created ${this.operatorService.getFullName({
          resource: job
        })}`
      );
      return job;
    }
  }
}

export enum Hook {
  AfterCleanup = 'after-cleanup',
  AfterCreate = 'after-create',
  AfterCreateOrUpdate = 'after-create-or-update',
  AfterUpdate = 'after-update',
  BeforeCleanup = 'before-cleanup',
  BeforeCreate = 'before-create',
  BeforeCreateOrUpdate = 'before-create-or-update',
  BeforeUpdate = 'before-update',
  Cleanup = 'cleanup',
  Create = 'create',
  CreateOrUpdate = 'create-or-update',
  Update = 'update'
}

export interface HookResult {
  hookName: string;
  message: string;
  name: string;
  namespace: string;
}
