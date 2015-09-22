---
layout: "app_docker_external"
page_title: "Dev - Docker (External) App Type"
sidebar_current: "docs-docker-dev"
description: |-
  Because the "docker-external" type expects the Docker images to exist
  and be built externally, the development environment is only meant as a way
  to test that starting the container works properly for downstream applications
  that may depend on this application.
---

# Development

Because the "docker-external" type expects the Docker images to exist
and be built externally, the development environment is only meant as a way
to test that starting the container works properly for downstream applications
that may depend on this application.

The primary purpose of this application type is to make it incredibly
easy to depend and deploy anything in the Docker ecosystem with Otto,
it isn't meant as a Docker image development environment. We have future
plans for this, but it isn't currently available.

When you run `otto dev`, you'll get an environment with the container
running as configured in the Appfile. The use case for this is to just verify
that it works for downstream applications using this Docker image as
a [dependency](/docs/concepts/deps.html).
