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
import fs from 'fs-extra';
import path from 'path';
import { OperatorFrameworkProject } from './types';

export const project: OperatorFrameworkProject = YAML.parse(
  fs.readFileSync(path.resolve(__dirname, '../PROJECT')).toString()
);

export function kind2plural(kind: string) {
  let lowercasedKind = kind.toLowerCase();
  if (lowercasedKind[lowercasedKind.length - 1] === 's') {
    return lowercasedKind;
  }
  if (lowercasedKind[lowercasedKind.length - 1] === 'o') {
    lowercasedKind = `${lowercasedKind}e`;
  }
  return `${lowercasedKind}s`;
}

export function getGroupName(groupNamePrefix: string, domain?: string) {
  return `${groupNamePrefix}.${domain || project.domain}`;
}
