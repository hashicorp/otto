---
layout: "intro"
page_title: "Otto vs. Docker"
sidebar_current: "vs-other-docker"
description: |-
  Comparison between Otto and Docker.
---

# Otto vs. Docker

The Docker ecosystem is comprised of multiple separate tools. Otto doesn't
overlap with most of these, so we're going to name specific tools that Otto
may be confused with.

## Single Workflow

Docker is distributed as multiple tools, each of which needs to be learned
separately, configured separately, and run separately.

Otto exposes a single simple workflow: `otto dev`, `otto deploy`, etc.
that orchestrates all the configuration and execution of every other
piece of software for you.

The Otto workflow is not bound by the underlying technology: it supports
a containerized workload as easily as a completely virtualized workload.
More importantly, it easily supports hybrid environments. Docker tooling
is tied directly to Docker containers, Docker hosts, etc. Otto is built
to support any underlying technology paradigms: virtualization, containerization,
etc.

In addition to this: the Docker tools on their own aren't enough to build
a complete infrastructure. Issues such as scheduling, service discovery,
security, etc. still require external solutions that you have to learn
and deploy yourself. Otto automates the installation of all of this.

## Development

For **development**, Otto can be compared with running Docker directly
in addition to Docker Compose. Otto has a number of differences when compared to
these with just development in mind.

With Docker Compose, you must list the full list of dependencies an
application has, including dependencies of dependencies. With Otto, you
list only your [immediate dependencies](/docs/concepts/deps.html). Otto
parses the dependencies to find their dependencies automatically, and so on.
This makes it much easier to work in microservice environments since
dependencies may change often and changing the dependencies of your application
won't negatively impact downstream consumers.

For development, Docker Compose recommends you use the "build" option
to build a Dockerfile and run that Dockerfile. This makes for a slow
feedback cycle: edit files then rebuild the compose environment. Otto
recognizes this and creates a mutable environment designed for fast updates.
For example, if you're working on a PHP project, just save the PHP page,
refresh your browser, and the change is visible.

Otto uses a single virtual machine for the development environment.
Multiple dependencies are automatically installed onto this single
virtual machine. When working with Docker, you're also likely using
a virtual machine to run Docker, so this is similar. The only time this
isn't true is if you're developing on Linux directly.

## Deployment

For **deployment**, Otto can be compared with Docker Machine, perhaps with
Docker Swarm.

Docker Machine only spins up compute resources for Docker. Otto is able
to spin up a complete production-ready infrastructure: VPC, routing tables,
security groups, networking rules, etc. Otto can add compute resources
to this for Docker as well as non-Docker resources. Under the covers,
Otto uses [Terraform](https://terraform.io) to manage infrastructure,
a technology in use by many large companies and proven at large scale.

Docker Machine can also intall a Docker Swarm cluster. Otto will install
and bootstrap a [Consul](https://consul.io) cluster for service discovery
and configuration, and future versions of Otto will automatically install
[Vault](https://vaultproject.io) for security and [Nomad](https://nomadproject.io)
for deployment and cluster management. The benefit of the stack setup by Otto
is that it supports Docker containers as a first class
deployment mechanism, but also supports any other application from
custom server images to standalone JARs, executables, and more.

Docker Swarm clusters multiple Docker hosts together. Otto will use
[Nomad](https://nomadproject.io) to achieve the same thing. Nomad is able
to deploy Docker containers as well as other applications. The benefit to
using Otto is that the user doesn't need to make this decision: the deploy
will just work. The distributed system that handles this is automatically
bootstrapped, configured, secured, and scaled by Otto.

In addition to the above, we'll repeat that Otto is a single tool and workflow.
It uses other software under the covers, but Otto hides that complexity
from the end user behind a [simple configuration](/docs/appfile) and
a [simple workflow](/docs/concepts/workflow.html).
