---
layout: "docs"
page_title: "Application - Appfile"
sidebar_current: "docs-appfile-app"
description: |-
  The application configuration tells Otto the name of your
  application, the type, and any dependencies it may have.
---

# Application Configuration

The application configuration tells Otto the name of your
application, the type, and any dependencies it may have.

This page assumes you're familiar with the
[Appfile syntax](/docs/appfile/syntax.html) already.

## Example

The application configuration looks like the following:

```
application {
    name = "my-app"
    type = "ruby"
}
```

## Description

The `application` block tells Otto about the application that the
Appfile describes. Only one application block can exist in an Appfile.

More fine-grained configuration of the application is done using
[customization](/docs/appfile/customization.html) blocks.

The `application` block allows the following keys to be set:

  * `name` (string) - The name of the application. This doesn't have
      to be unique across all your other apps. This will be used for
      service discovery.

  * `type` (string) - The type of the application. The list of types
      is available in the [app types](/docs/apps) section.

-------------

Within a resource, you can specify zero or more **dependencies**.

Within the dependency, the following keys are allowed:

  * `source` (string) - The URL of the dependency. The full list of
      acceptable URL types is documented on the
      [dependency sources](/docs/appfile/dep-sources.html) page.

## Syntax

The full syntax is:

```
application {
	name = NAME
	type = TYPE

	[DEPENDENCY ...]
}
```

where `DEPENDENCY` is:

```
dependency {
	source = SOURCE
}
```
