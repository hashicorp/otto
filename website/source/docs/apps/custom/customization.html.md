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

Example:

```
customization {
    dev_vagrantfile = "./Vagrantfile.tpl"
    terraform       = "./terraform"
}
```

Available options:

  * `dev_vagrantfile` (string) - Path to a Vagrantfile to use for development.
    If this isn't specified, `otto dev` will not work for this application.
    This Vagrantfile will be rendered as a [template](/docs/apps/custom/template.html)
    if the path ends with `.tpl`.

  * `dep_vagrantfile` (string) - Path to a Vagrantfile to use as a fragment
    that is embedded in other application's Vagrantfiles when this application
    is being used as a dependency.

  * `packer` (string) - Path to a Packer template to execute.

  * `terraform` (string) - Path to a Terraform module (directory) to use
    for deployment.
