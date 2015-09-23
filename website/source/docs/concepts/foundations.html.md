---
layout: "docs"
page_title: "Foundations"
sidebar_current: "docs-concepts-foundations"
description: |-
  A foundation is what Otto considers to be a fundamental building block
  of a real infrastructure. Otto comes built-in with many foundations that
  are automatically setup for you.
---

# Foundations

A foundation is what Otto considers to be a fundamental building block
of a real infrastructure. Otto comes built-in with many foundations that
are automatically installed and configured for you.

Examples of foundations are: service discovery, security, scheduling, etc.
More foundations will be added, changed, and removed as time goes on and
the state of the art changes.

These are core features of an infrastructure that have grown to become
best practices within modern organizations. However, they're also usually
complicated to deploy and operate. But with Otto, Otto installs, configures,
and scales the foundations for you.

## Features

Otto 0.1 uses foundations to make the following automatically happen:

**Service discovery**. Dependencies are discoverable using DNS:
`<appname>.service.consul`. For example, if your application depends on
a PostgreSQL database, it will be available at `postgresql.service.consul`.
Otto makes this all happen automatically.

This is available both in development and production. In development,
to access a dependency, you use DNS the same way.

## Future Foundations
<a id="future-foundations"></a>

In the future, Otto will add more foundations to enable more features:

**Security**. Otto will install, bootstrap, and configure
[Vault](https://vaultproject.io) automatically. Otto will automatically
configure authentication for all applications so they can store and
access secrets. Any required credentials (such as for databases) will be
automatically requested from Vault and passed in to applications.

The learning curve for this is usually very high, but with Otto, we can
do it all automatically.

**Scheduling**. Otto will install, bootstrap, configure, and scale
[Nomad](https://nomadproject.io) to enable much faster deploys,
auto-scaling, application failover, and more. This will happen transparently:
you'll update your Otto version, update the infrastructure, run
`otto deploy`, and notice that everything happened MUCH faster.

Combining Otto with a scheduler foundation such as Nomad will also
dramatically lower the cost of infrastructure since Nomad can enable
much higher resource utilization, requiring less servers overall.

## Custom Foundations

Foundations are internally architected in such a way to allow them to
be configurable and replaceable. Otto 0.1 doesn't allow this functionality,
but it will be exposed in a future Otto version.

This will allow other projects to replace the built-in foundations and
provide similar features.
