---
layout: "app_custom"
page_title: "Custom - App Types"
sidebar_current: "docs-custom-index"
description: |-
  Otto has two mechanisms that can be used to extend Otto to support
  new application types: the "custom" application type and writing
  a custom application type plugin.
---

# Custom App Type

**Type**: `custom`

While Otto has built-ins for many of the most popular languages and
frameworks, we understand that the variety of application development
and deployment styles varies widely.

Otto has two mechanisms that can be used to extend Otto to support
new application types: the "custom" application type and writing
a custom application type plugin.

This page documents the "custom" application type. Application type plugins
are much more powerful but also more complicated to use. We recommend
prototyping with the "custom" application type if possible, and then
codifying that into a plugin. To learn how to write an application type
plugin, see the [plugins page here](/docs/plugins/app.html).

The "custom" application type does almost nothing, and relies on heavy use
of [customizations](/docs/apps/custom/customization.html) to tell Otto
how to behave. This gives the user the ultimate power, but also a high
level of complexity. It is meant for advanced users.

If you find yourself using a common custom application type often,
consider extracting it into a real [plugin](/docs/plugins/app.html).

## Example

An example Appfile using the "custom" application type is shown below:

```
application {
    name = "my-app"
    type = "custom"
}

customization {
    dev_vagrantfile = "./Vagrantfile"
    terraform = "./terraform"
}
```

As you can see with the "custom" application type, the Appfile must specify
all the lower level configurations that Otto will use for various steps.
For development, a [Vagrantfile](https://docs.vagrantup.com/v2/vagrantfile/index.html)
must be specified. For deployment, a
[Terraform module](https://www.terraform.io/docs/modules/index.html)
must be specified.

More [customizations](/docs/apps/custom/customization.html) are available
to fine-tune other steps as well.
