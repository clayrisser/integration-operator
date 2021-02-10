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

import execa, { ExecaChildProcess, ExecaReturnValue, Options } from 'execa';
import { HashMap } from '../types';

const logger = console;

export default abstract class Command {
  protected config: CommandConfig;

  protected abstract command: string;

  protected execa = execa;

  constructor(config: Partial<CommandConfig> = {}) {
    this.config = {
      debug: false,
      ...config
    };
  }

  async run<T = Result>(
    args: string | string[] = [],
    options: Options = {},
    cb: RunCallback = (p: ExecaChildProcess) => {
      p.stderr?.pipe(process.stderr);
      p.stdout?.pipe(process.stdout);
    }
  ): Promise<T> {
    if (this.config.debug) {
      logger.debug('$', [this.command, ...args].join(' '));
    }
    if (!Array.isArray(args)) args = [args];
    const p = execa(this.command, args, options);
    cb(p);
    return this.smartParse(await p) as T;
  }

  smartParse(result: ExecaReturnValue<string>): Result {
    try {
      return JSON.parse(result.stdout);
    } catch (err) {
      return result.stdout as string;
    }
  }
}

export type Result = string | HashMap;

export interface CommandConfig {
  debug: boolean;
}

export type RunCallback = (p: ExecaChildProcess) => any;
