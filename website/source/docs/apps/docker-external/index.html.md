---
layout: "app_docker_external"
page_title: "Docker (External) - App Types"
sidebar_current: "docs-docker-index"
description: |-
  The "docker-external" application type is used for Docker images
  that are already built and in a registry somewhere.
---

# Docker (External) App Type

**Type:** `docker-external`

The "docker-external" application type is used for Docker images
that are already built and in a registry somewhere. The "external" is
there to signify that the image is built external to the Otto lifecycle,
and already exists.

The primary purpose of this application type is to make it incredibly
easy to depend and deploy anything in the Docker ecosystem with Otto.

## Example

Below is an example Appfile that uses the "docker-external" type:

```
application {
    name = "mongodb"
    type = "docker-external"
}

customization "docker" {
    image = "mongo:3.0"
    run_args = "-p 27017:27017"
}
```

This Appfile can be used as a dependency of another application, and
can be deployed as well.
