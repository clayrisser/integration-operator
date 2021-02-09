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

export interface OperatorFrameworkPlugins {
  [key: string]: {};
}

export interface OperatorFrameworkResource {
  group: string;
  kind: string;
  version: string;
}

export interface OperatorFrameworkProject {
  domain: string;
  layout: string;
  plugins: OperatorFrameworkPlugins;
  projectName: string;
  repo: string;
  resources: OperatorFrameworkResource[];
  version: string;
}

export interface HashMap<T = any> {
  [key: string]: T;
}

export * from './integration';
