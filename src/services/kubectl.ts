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
import { ExecaChildProcess, Options } from 'execa';
import { KubernetesObject } from '@kubernetes/client-node';
import { Readable } from 'stream';
import CommandService, { RunCallback } from './command';

export default class KubectlService extends CommandService {
  command = 'kubectl';

  private getStdin(
    stdin: string | KubernetesObject | KubernetesObject[]
  ): string {
    if (Array.isArray(stdin)) {
      return this.resources2String(stdin);
    }
    if (typeof stdin === 'object') {
      return this.resources2String([stdin]);
    }
    return stdin as string;
  }

  async help(options?: Options, cb?: RunCallback) {
    return this.run('--help', options, cb);
  }

  async apply(
    applyOptions: Partial<ApplyOptions> | string = {},
    options?: Options
  ) {
    let { stdin } = applyOptions as Partial<GetOptions>;
    if (typeof applyOptions === 'string') stdin = applyOptions;
    const { file, stdout } = {
      ...((stdin ? { file: '-' } : {}) as Partial<GetOptions>),
      ...(typeof applyOptions === 'string'
        ? ({} as Partial<GetOptions>)
        : applyOptions)
    };
    return this.run(
      ['apply', ...(file ? ['-f', file] : [])],
      options,
      (p: ExecaChildProcess) => {
        if (stdin) {
          const stream = Readable.from([this.getStdin(stdin)]);
          if (p.stdin) stream.pipe(p.stdin);
        }
        if (stdout) {
          p.stderr?.pipe(process.stderr);
          p.stdout?.pipe(process.stdout);
        }
      }
    );
  }

  async get<T = any>(
    getOptions: Partial<GetOptions> | string = {},
    options?: Options
  ): Promise<T> {
    const { stdout } = getOptions as Partial<GetOptions>;
    let { stdin } = getOptions as Partial<GetOptions>;
    if (typeof getOptions === 'string') stdin = getOptions;
    const { file, output, ignoreNotFound } = {
      ignoreNotFound: true,
      ...((stdin ? { file: '-' } : {}) as Partial<GetOptions>),
      ...(typeof getOptions === 'string'
        ? ({} as Partial<GetOptions>)
        : getOptions)
    };
    return this.run<T>(
      [
        'get',
        ...(file ? ['-f', file] : []),
        ...(output ? ['-o', output] : []),
        ...(ignoreNotFound ? ['--ignore-not-found'] : [])
      ],
      options,
      (p: ExecaChildProcess) => {
        if (stdin) {
          const stream = Readable.from([this.getStdin(stdin)]);
          if (p.stdin) stream.pipe(p.stdin);
        }
        if (stdout) {
          p.stderr?.pipe(process.stderr);
          p.stdout?.pipe(process.stdout);
        }
      }
    );
  }

  async delete(
    deleteOptions: Partial<DeleteOptions> | string = {},
    options?: Options
  ) {
    let { stdin } = deleteOptions as Partial<GetOptions>;
    if (typeof deleteOptions === 'string') stdin = deleteOptions;
    const { file, stdout } = {
      ...((stdin ? { file: '-' } : {}) as Partial<GetOptions>),
      ...(typeof deleteOptions === 'string'
        ? ({} as Partial<GetOptions>)
        : deleteOptions)
    };
    return this.run(
      ['delete', ...(file ? ['-f', file] : [])],
      options,
      (p: ExecaChildProcess) => {
        if (stdin) {
          const stream = Readable.from([this.getStdin(stdin)]);
          if (p.stdin) stream.pipe(p.stdin);
        }
        if (stdout) {
          p.stderr?.pipe(process.stderr);
          p.stdout?.pipe(process.stdout);
        }
      }
    );
  }

  resources2String(resources: KubernetesObject[]): string {
    return resources
      .map((resource: KubernetesObject) => YAML.stringify(resource))
      .join('---\n');
  }

  string2Resources(resourcesStr: string): KubernetesObject[] {
    return `\n${resourcesStr}\n`
      .split(/\n---+\n/)
      .map(
        (resourceStr: string) =>
          YAML.parse(resourceStr.trim()) as KubernetesObject
      );
  }
}

export interface GetOptions {
  file?: string;
  ignoreNotFound?: boolean;
  output?: Output;
  stdin?: string | KubernetesObject | KubernetesObject[];
  stdout?: boolean;
}

export interface ApplyOptions {
  file?: string;
  stdin?: string | KubernetesObject | KubernetesObject[];
  stdout?: boolean;
}

export interface DeleteOptions {
  file?: string;
  stdin?: string | KubernetesObject | KubernetesObject[];
  stdout?: boolean;
}

export enum Output {
  Yaml = 'yaml',
  Json = 'json'
}
