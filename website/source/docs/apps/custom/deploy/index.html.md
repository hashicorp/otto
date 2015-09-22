---
layout: "app_custom"
page_title: "Build & Deploy - Custom App Type"
sidebar_current: "docs-custom-deploy"
description: |-
  "Custom" application types use Packer and Terraform to build and deploy.
---

# Build & Deploy

"Custom" application type build & deploy processes are defined with
[customizations](/docs/apps/custom/customization.html) and use
Packer and Terraform.

-> **NOTE:** This page documents the "custom" application type. Theres are
   different from [application type plugins](/docs/plugins/app.html) which
   are a way to introduce new application types to Otto.

Both "build" and "deploy" customizations are optional. If neither are
specified, then Otto will not be able to build or deploy your application.
If only one is specified, it will only be able to do one of those operations.

## Example

```
application {
    name = "my-app"
    type = "custom"
}

customization "build" {
    packer = "./template.json"
}

customization "deploy" {
    terraform = "./tf-module"
}
```

## Build

When `otto build` is called, Otto will execute Packer against the
given Packer template. Artifacts will be stored in the
[directory](/docs/concepts/directory.html) and will be passed in via
variables to the deploy step.

## Deploy

When `otto deploy` is called, Otto will execute the configured Terraform
module. Any stored credentials from the infrastructure as well as outputs
will be passed in as Terraform variables.

If you specified a build step, then the artifact information will also
be passed in via Terraform variables. The exact variables that are
passed in are not yet documented, but will be soon.
