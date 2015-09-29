---
layout: "intro"
page_title: "Otto vs. Docker"
sidebar_current: "vs-other-docker"
description: |-
  Comparison between Otto and Docker.
---

# Otto vs. Docker

Docker is an ecosystem of tools including the Docker runtime, Compose,
Machine and Swarm. Each of these tools solves a different problem, but
they are all specific to Docker.

To develop using Docker, users must either develop on Linux or start a
virtual machine running Linux. Then a `docker-compose.yml` is constructed
with all the containers needed to run the full set of dependencies for your application.
This is used by Docker Compose to manage the lifecycle of development
containers.

Once ready for deployment, you must build the application, create
a Docker container, and upload it to Docker Hub or another registry.
To build your application you create a `Dockerfile` which contains
the compilation depedencies and emits a binary or compile artifact.
Another `Dockerfile` is used to create a minimal deployment container
with just this artifact, unless you want to deploy the container
that includes all the compilation dependencies.

Servers are provisioned with Docker Machine, which creates a server
that has the Docker runtime installed. Machine does not have
configuration file to specify all the machines needed, so it should
be invoked by a provisioning script.

If more than a single server is required, Docker Swarm is used to
cluster the servers together. This allows containers to be scheduled
to the cluster as if it were a single large machine. Docker Swarm
can be configured by Machine during server provisioning.

Once setup, containers must be provisioned on the cluster. This
can be done with the `docker` CLI pointing at the Swarm cluster.
For applications with dependencies, this should be wrapped in
a deployment script to ensure those are launched along with the
application.

Otto is a single tool and is much simpler to use. Otto uses an `Appfile`
to describe the application and any upstream dependencies. Otto
uses this same `Appfile` to setup development environments,
building, launch infrastructure, and deployment.

To create a development environment, `otto dev` is run. This
uses the `Appfile` to setup a virtual machine and downloads
any compilation or upstream dependencies.

Once ready, you can create the infrastructure for your application
by running `otto infra`. This provisions the servers needed to
run your application. This only needs to be done once to create
the infrastructure.

To build the application you run `otto build`. This packages
the application to be deployed. This can use Docker, but this
is a detail that developers do not need to worry about. The
`Appfile` already encodes what is needed and no `Dockerfile`
is necessary.

Lastly, `otto deploy` is used to deploy the application.
This will use the infrastructure setup by `otto infra` and
the build artifacts from `otto build`.

Otto is designed to be a single tool that manages the workflow
from development to production and requires only the `Appfile`.
It is meant to simplify the complex state of development today,
while using production hardened tools and industry best practices
automatically.
