---
layout: "app_custom"
page_title: "Development - Custom App Type"
sidebar_current: "docs-custom-dev"
description: |-
  "Custom" application types use Vagrant for development.
---

# Development

"Custom" application type development environments are defined by
a "dev" [customization](/docs/apps/custom/customization.html) and use
Vagrant.

-> **NOTE:** This page documents the "custom" application type. Theres are
   different from [application type plugins](/docs/plugins/app.html) which
   are a way to introduce new application types to Otto.

## Example

```
application {
    name = "my-app"
    type = "custom"
}

customization "dev" {
    vagrantfile = "./Vagrantfile"
}
```

For the Appfile above, running `otto dev` will run Vagrant in the directory
of the Appfile against the given Vagrantfile.
