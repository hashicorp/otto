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
    vagrantfile = "./Vagrantfile.tpl"
}
```

For the Appfile above, running `otto dev` will run Vagrant in the directory
of the Appfile against the given Vagrantfile.

The Vagrantfile is rendered as [a template](/docs/apps/custom/template.html)
if the path ends with `.tpl`.

It is important at the very least to specify the Vagrant shared folder
should be the working directory. An example Vagrantfile config is shown
below to do this:

```
config.vm.synced_folder '{{ path.working }}', "/vagrant"
```

If you don't do this, then the directory that `/vagrant` sees will be
the compiled directory, which is very likely not what you want.
