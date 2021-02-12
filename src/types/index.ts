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
