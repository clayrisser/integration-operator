# integration-operator

> kubernetes operator to integrate deployments

[![integration-operator](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/rock8s)](https://artifacthub.io/packages/helm/rock8s/integration-operator)

Please ★ this repo if you found it useful ★ ★ ★

This operator takes inspiration from [Juju](https://juju.is) [Charm](https://juju.is/docs/sdk)
[Relations](https://juju.is/docs/sdk/relations) by [Canonical](https://canonical.com).

## Install

```sh
helm repo add rock8s https://charts.rock8s.com
helm install integration-operator rock8s/integration-operator --version 1.2.0 --namespace kube-system
```

## Develop

1. Install the custom resource definitions

```sh
./mkpm install
```

2. Start the operator

```sh
./mkpm dev
```

3. Create plugs and sockets

   You can start by taking a look at [config/samples](config/samples).

   ```sh
   kubectl apply -f config/samples
   ```

## Terminology

| Term            | Juju Equivalent | Definition                                                                           |
| --------------- | --------------- | ------------------------------------------------------------------------------------ |
| Integration     | Relation        | unite and connect applications through mutual communication and shared configuration |
| Plug            | Requires        | request from an application to integrate with another application                    |
| Socket          | Provides        | fulfils requests from applications trying to integrate                               |
| Interface       | Interface       | plug and socket schema required to connect                                           |
| Created Event   | Created Event   | event triggered when plug or socket created                                          |
| Updated Event   | Changed Event   | event triggered when plug or socket updated                                          |
| Coupled Event   | Joined Event    | event triggered when applications connected                                          |
| Decoupled Event | Detached Event  | event triggered when applications disconnected                                       |

## Architecture

### A simple analogy

The best way to explain the architecture is to think about how plugs and sockets work in the real world.

Let's say I have a laptop purchased in the United States. In order to power my laptop, I need to **integrate** it with the power grid.
Since the laptop was purchased in the United States, the **interface** of the **plug** is Type A as illustrated below.

![Type A](https://gitlab.com/bitspur/rock8s/integration-operator/-/raw/main/images/typea.png)

This means the **socket** I connect to must be also be Type A.

Now, let's say I travel to India and the only **socket** available to me is Type D as illustrated below.

![Type D](https://gitlab.com/bitspur/rock8s/integration-operator/-/raw/main/images/typed.png)

Since the **socket** interface does not match the **plug** interface, I cannot integrate my laptop with the power grid in India. Of course
this can be overcome with converters, but that is beyond the scope of this analogy.

### A real example

Let's say I have an express application that needs to **integrate** with a mongo database. The express deployment will have a **plug** with
a mongo **interface** and the mongo deployment will have a **socket** with a mongo **interface**. If the **interface** of the **socket** is
a postgres **interface** then the integration will fail. In other words, you cannot connect a mongo **plug** to a postgres **socket**. That
would be like trying to plug a US Type A **plug** into an Indian Type D **socket**. You can only connect a mongo **plug** to a mongo **socket**.

## Concepts

### Socket

A socket is a custom kubernetes resource that fulfills integration requests from other applications.
It carries out the following tasks:

- defines the interface for the configuration and result of the plug and socket
- provides the configuration for the socket
- provides the result for the socket
- templates any resources within the socket's namespace
- executes any apparatuses within the socket's namespace
- templates result resources within the socket's namespace

**Example:**

```yaml
apiVersion: v1
kind: Socket
metadata:
  name: foo
  namespace: foo-namespace
spec:
  config:
    hello: world
```

### Plug

A plug is a custom kubernetes resource that initiates an integration request with another application.
It does not define its own interface as it utilizes the interface defined by the socket.
The plug carries out the following tasks:

- couples to a socket
- provides the configuration for the plug
- provides the result for the plug
- templates any resources within the plug's namespace
- executes any apparatuses within the plug's namespace
- templates result resources within the plug's namespace

**Example:**

```yaml
apiVersion: v1
kind: Plug
metadata:
  name: bar
  namespace: bar-namespace
spec:
  socket:
    name: foo
    namespace: foo-namespace
  configSecretName: my-secret
```

### Data

The _data_ in the plug or socket is a flexible and unstructured form of information exchange. It is
used during the preliminary stages of the integration process, before the final _config_ is established.
Unlike _config_ and _result_, _data_ is not bound by an interface. It is used for exchanging or simplifying
preliminary details or any other information that might be necessary for generating the final _config_.

The _data_ can be supplied directly through the `data` field, and indirectly through the `dataConfigMapName` field
and `dataSecretName` field. The `data` field is a key-value pair that can be defined directly within the plug or
socket. If the `dataConfigMapName` or `dataSecretName` field is used, the _data_ will be retrieved from a ConfigMap
or Secret respectively.

It is important to know that _data_ is utilized exclusively by the `configTemplate` field, `resultTemplate` field, and
the `/config` endpoint of an apparatus. It enables the exchange of information between plugs and sockets before the
final _config_ is established. This process prevents potential recursive issues that could arise if the _config_ of
the plug and socket were interdependent. As such, _data_ serves as an initial medium for information exchange,
facilitating the creation of the final _config_ for the integration process.

**Example:**

_this is a simplified incomplete example, only including necessary fields_

```yaml
spec:
  dataConfigMapName: my-configmap
  dataSecretName: my-secret
  data:
    username: admin
    password: secret
```

### Vars

The _vars_ allows the capture and insertion of values from one resource's field to another, functioning
similarly to vars in Kustomize. It is defined by the `vars` field. Like _data_, _vars_ can only be used
by the `configTemplate` field and the `/config` endpoint of an apparatus. Since _vars_ is used by _config_, the
lookup occurs before the _config_ is finalized.

In addition to the `vars` field, there is a separate field, known as `resultVars`, which is used by
the `resultTemplate` field. Since _resultVars_ is used by _result_, the lookup occurs after the integration has
been established or updated. This allows for the creation of _resultVars_ based on the results of the integration.

For more detailed information, please refer to the
[Kustomize Vars Documentation](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/vars/).

**Example:**

_this is a simplified incomplete example, only including necessary fields_

```yaml
spec:
  vars:
    - name: serviceAccountName
      objref:
        apiVersion: apps/v1
        kind: Deployment
        name: my-deployment
        namespace: default
      fieldref:
        fieldPath: spec.template.spec.serviceAccountName
  resultVars:
    - name: jobSuccessful
      objref:
        apiVersion: batch/v1
        kind: Job
        name: my-job
        namespace: default
      fieldref:
        fieldPath: status.successful
```

### Config

The _config_ is the most fundamental concept of the integrations, serving as a key-value data pair that enables secure
information exchange between the plug and socket. It contains essential details and information necessary for the
integration.

The _config_ can be supplied directly through the `config` field, or indirectly through the `configConfigMapName` field,
`configSecretName` field, `configTemplate` field and the `/config` endpoint of an apparatus. The `config` field is a
key-value pair that can be defined directly within the plug or socket. If the `configConfigMapName` or `configSecretName`
field is used, the _config_ will be retrieved from a ConfigMap or Secret respectively. If the `configTemplate` field is
used, the _config_ will be templated, allowing the composition of values from `vars`, `plugData`, `socketData`, `plug`
and `socket`. If the `/config` endpoint of an apparatus is used, the _config_ will come from the response payload. The
request body will contain `vars`, `plugData` and `socketData`. Please note that `plugConfig` and `socketConfig` will not
be available to the `configTemplate` field or the `/config` endpoint of an apparatus. All of these strategies for creating
the _config_ can be used in combination.

The _config_ is validated against the _config interface_ before the integration process begins. This ensures that the
_config_ contains all the necessary information, adheres to the correct format and enforces a contract between the
plug and socket integration

**Example:**

_this is a simplified incomplete example, only including necessary fields_

```yaml
spec:
  config:
    protocol: http
    port: "8080"
  configTemplate:
    hostname: "{% .vars.ingressHost %}"
  configConfigMapName: my-configmap
  configSecretName: my-secret
```

### Results

The _result_ serves as a key-value data pair that contains essential details and information after an integration
has been coupled or updated. It can be used in the `resultResources` field.

The _result_ can be supplied directly through the `result` field, or indirectly through the `resultConfigMapName`
or `resultSecretName` field. If the `resultConfigMapName` or `resultSecretName` field is used, the _result_ will be
retrieved from a ConfigMap or Secret respectively. If the `resultTemplate` field is used, the _result_ will be templated,
allowing the composition of values from `resultVars`, `plugData`, `socketData`, `plugConfig`, `socketConfig`, `plug`,
and `socket`. All of these strategies for creating the _result_ can be used in combination.

The _result_ is validated against the _result interface_ after the integration is coupled or updated. This ensures that
the _result_ contains all the necessary information, adheres to the correct format and enforces a contract between the
plug and socket integration.

**Example:**

_this is a simplified incomplete example, only including necessary fields_

```yaml
spec:
  result:
    hello: world
  resultTemplate:
    foo: "{% .plugConfig.foo %}"
  resultConfigMapName: my-configmap
  resultSecretName: my-secret
```

### Interface

The _interface_ validates the _config_ and _result_ against a defined schema, ensuring they contain all necessary
properties. The integration fails if the _interface_ requires a _config_ or _result_ that is missing. Any _config_
or _result_ provided that isn't defined in the _interface_ will be ignored. This guarantees that only properties
defined in the _interface_ are used during integration, preserving integrity and consistency. If no _interface_ is
provided, the _config_ and _result_ are not validated and can be any value. However, this is discouraged as it may
lead to inconsistencies and unexpected behavior during the integration process.

**Example:**

_this is a simplified incomplete example, only including necessary fields_

```yaml
kind: Socket
spec:
  interface:
    config:
      plug:
        hello:
          default: world
      socket:
        howdy:
          required: true
    result:
      socket:
        foo:
          required: true
      plug:
        bar: {}
```

### Resources

Resources are utilized during the integration process to template kubernetes resources. They are defined within the plug or
socket and can encompass any valid Kubernetes resource such as Jobs, Pods, Services, and more. These resources play a
pivotal role in executing the integration process.

Resource templates are defined using the `template` and `templates` fields. The `template` field is used for a single
resource template, while the `templates` field is used for multiple resource templates. These templates are defined in YAML
format.

The `stringTemplate` and `stringTemplates` fields are analogous to `template` and `templates`, but they accept resource
templates in string format. This is particularly useful when dealing with complex resource templates that require
conditional templating, such as wrapping a resource in an if statement.

The `do` field specifies the action to be performed on the resource. It can be `delete`, `apply`, or `recreate`.

The `when` field specifies the stage of the integration process when the resource action should be performed. It can
be `updated`, `coupled`, `decoupled`, `created`, or `deleted`.

The `preserveWhenDecoupled` field is a boolean that determines whether the resource should be preserved when the
integration is decoupled. If `true`, the resource will not be deleted during decoupling. If `false` or omitted, the
resource will be deleted unless the `when` field contains `decoupled`.

A unique field, `resultResources`, is used to create resources after the integration has been coupled or updated. The
templating of `resultResources` takes place after the integration process has been coupled or updated. This allows for
the creation of resources based on the results of the integration process.

The `resultResources` field is used to create resources after the integration has been coupled or updated. The templating
of `resultResources` takes place after the integration process has been coupled or updated. This allows for the creation
of resources based on the results of the integration process.

**Example:**

_this is a simplified incomplete example, only including necessary fields_

```yaml
spec:
  resources:
    - when: [coupled, updated]
      do: apply
      template:
        apiVersion: batch/v1
        kind: Job
        metadata:
          name: my-job
        spec:
          template:
            spec:
              containers:
                - name: my-job
                  image: my-job-image
                  env:
                    - name: HELLO
                      value: "{% .plugConfig.hello %}"
  resultResources:
    - do: apply
      stringTemplate: |
        {%- if (eq .result.resultJob "1") %}
        apiVersion: batch/v1
        kind: Job
        metadata:
          name: my-result-job
        spec:
          template:
            spec:
              containers:
                - name: my-result-job
                  image: my-result-job-image
        {%- endif %}
```

### Apparatus

The apparatus is a unique component that offers a unique approach to executing the integration process. Unlike resources,
which are primarily used for templating Kubernetes resources, the apparatus is a pod that operates a REST API. These APIs
are invoked at different stages of the integration process, passing data such as the `plug`, `socket`, `plugConfig`, and
`socketConfig` in the request body.

It's important to note that an apparatus and resources can be used together during the integration process. This
combination provides a flexible and robust integration process capable of handling a wide range of scenarios.

The apparatus pod is automatically cleaned up when it's not in use and will be created automatically when integrations
require it. The apparatus schema is the same as the schema used to define a pod.

An good example of an apparatus use case is the
[Keycloak Integration Apparatus](https://gitlab.com/bitspur/rock8s/keycloak-integration-apparatus). This apparatus is
necessary because the Keycloak integration involves interacting with the Keycloak API
via a TypeScript client, which would be challenging to accomplish using only resources. By constructing it
as an apparatus, we can leverage a NodeJS REST API to effectively communicate with Keycloak.

The apparatus controller, which can be programmed in any language due to its REST architecture, should implement the following endpoints:

| Method | Endpoint     | Description            | Request Body                                      |
| ------ | ------------ | ---------------------- | ------------------------------------------------- |
| `GET`  | `/ping`      | checks the health      |                                                   |
| `POST` | `/config`    | retrieves the config   | `vars`, `plug`, `socket`,`plugData`, `socketData` |
| `POST` | `/created`   | invoked when created   | `plug`, `socket`, `plugConfig`, `socketConfig`    |
| `POST` | `/coupled`   | invoked when coupled   | `plug`, `socket`, `plugConfig`, `socketConfig`    |
| `POST` | `/updated`   | invoked when updated   | `plug`, `socket`, `plugConfig`, `socketConfig`    |
| `POST` | `/decoupled` | invoked when decoupled | `plug`, `socket`, `plugConfig`, `socketConfig`    |
| `POST` | `/deleted`   | invoked when deleted   | `plug`, `socket`, `plugConfig`, `socketConfig`    |

**Example:**

_this is a simplified incomplete example, only including necessary fields_

```yaml
spec:
  apparatus:
    endpoint: /socket
    containers:
      - name: my-apparatus
        image: my-apparatus-image
        ports:
          - containerPort: 3000
            name: container
            protocol: TCP
```
