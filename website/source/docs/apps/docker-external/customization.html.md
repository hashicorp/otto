---
layout: "app_docker_external"
page_title: "Customization - Docker (External) App Type"
sidebar_current: "docs-docker-customization"
description: |-
  This page documents the [customizations](/docs/appfile/customization.html)
  that are available to change the behavior of Docker applications with Otto.
---

# Customization

This page documents the [customizations](/docs/appfile/customization.html)
that are available to change the behavior of Docker applications with Otto.

## Type: "docker"

Example:

```
customization "docker" {
    image = "mongo:3.0"
}
```

Available options:

  * `image` (string) - The Docker image to run (along with any tags).
    This will defaulit to the application name.

  * `run_args` (string) - Raw arguments to pass to `docker run`.
