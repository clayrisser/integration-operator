import execa, { ExecaChildProcess, ExecaReturnValue, Options } from 'execa';
import { HashMap } from '~/types';

const logger = console;

export default abstract class CommandService {
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
