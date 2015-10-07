---
layout: "app_custom"
page_title: "Customization - Custom App Type"
sidebar_current: "docs-custom-customization"
description: |-
  This page documents the customizations
  that are available to change the behavior of custom applications with Otto.
---

# Customization

This page documents the [customizations](/docs/appfile/customization.html)
that are available to change the behavior of the "custom" application
type with Otto.

## Type: "dev"

Example:

```
customization "dev" {
    vagrantfile = "./Vagrantfile"
}
```

Available options:

  * `vagrantfile` (string) - Path to a Vagrantfile to use for development.
    If this isn't specified, `otto dev` will not work for this application.

## Type: "dev-dep"

Example:

```
customization "dev-dep" {
    vagrantfile = "./Vagrantfile"
}
```

Available options:

  * `vagrantfile` (string) - Path to a Vagrantfile to use as a fragment
    that is embedded in other application's Vagrantfiles when this application
    is being used as a dependency.

## Type: "build"

Example:

```
customization "build" {
    packer = "./template.json"
}
```

Available options:

  * `packer` (string) - Path to a Packer template to execute.

## Type: "deploy"

Example:

```
customization "deploy" {
    terraform = "./tf-module"
}
```

Available options:

  * `terraform` (string) - Path to a Terraform module (directory) to use
    for deployment.
