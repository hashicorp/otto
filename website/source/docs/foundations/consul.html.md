---
layout: "docs"
page_title: "Consul Foundation"
sidebar_current: "docs-foundations-consul"
description: |-
  Otto's Consul foundation provides service discovery in your infrastructure.
---

# Consul Foundation

Otto's Consul foundation provides service discovery in to your applications and
infrastructure.

## Development

The Consul foundation installs a local Consul server in your development
environment and configures it to respond to `*.consul` DNS queries.

## Infrastructure

When the Consul foundation is active, Otto spins up a Consul cluster in your
infrastructure, ready for agents to connect. The size of this cluster varies by
[infrastructure type and flavor](/docs/infra/index.html).

 * __AWS (simple)__: This infrastructure type is meant for small demos, so Otto
   spins up a minimal one-node Consul server using a `t2.micro` instance type.
 * __AWS (vpc-public-private)__: Otto will spin up a 3-node highly available
   Consul cluster. Currently this also uses a `t2.micro` type, but in future
   versions the instance size will be configurable.

## Build & Deploy

Otto uses the Consul foundation during the build phase to ensure that a Consul
agent is installed and configured on each application instance that is built
and deployed. This agent will automatically join the Consul cluster in your
infrastructure and publish a service for your application.
