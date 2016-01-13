---
layout: "app_docker_external"
page_title: "Build & Deploy - Docker (External) App Type"
sidebar_current: "docs-docker-deploy"
description: |-
  Deployment of "docker-external" currently starts a Docker container
  on a single server. This isn't an ideal way to run Docker containers,
---

# Build & Deploy

Deployment of "docker-external" currently starts a Docker container
on a single server. This isn't an ideal way to run Docker containers,
and a future version of Otto will vastly improve this by leaning on
[Nomad](https://www.nomadproject.io).

This highlights a major benefit of Otto: as deployment practices improve,
future version of Otto can adopt the new practices, and upgrade to an
improved process for you.
