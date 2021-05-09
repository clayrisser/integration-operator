# integration-operator

> kubernetes operator to integrate deployments

Please ★ this repo if you found it useful ★ ★ ★

This operator takes inspiration from [Juju](https://juju.is) [Charm](https://juju.is/docs/sdk)
[Relations](https://juju.is/docs/sdk/relations) by [Canonical](https://canonical.com).

## Terminology

| Term        | Juju Equivalent | Definition                                                                                    |
| ----------- | --------------- | --------------------------------------------------------------------------------------------- |
| Integration | Relation        | means to unite and connect applications through mutual communication and shared configuration |
| Plug        | Requires        | request from an application to integrate with another application                             |
| Socket      | Provides        | fulfils requests from applications trying to integrate                                        |
| Interface   | Interface       | plug and socket schema required to connect                                                    |

## Architecture

### A simple analogy

The best way to explain the architecture is to think about how plugs and sockets work in the real world.

Let's say I have a laptop purchased in the United States. In order to power my laptop, I need to **integrate** it with the power grid.
Since the laptop was purchased in the United States, the **interface** of the **plug** is Type A as illustrated below.

![Type A](images/typea.png)

This means the **socket** I connect to must be also be Type A.

Now, let's say I travel to India and the only **socket** available to me is Type D as illustrated below.

![Type D](images/typed.png)

Since the **socket** interface does not match the **plug** interface, I cannot integrate my laptop with the power grid in India. Of course
this can be overcome with converters, but that is beyond the scope of this analogy.

### A real example

Let's say I have an express application that needs to **integrate** with a mongo database. This means the **plug** requires a **socket**
with a mongo **interface**. If the **interface** of the **socket** is a postgres **interface** then the integration will fail.
