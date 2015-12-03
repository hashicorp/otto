---
layout: "docs"
page_title: "Customization - Appfile"
sidebar_current: "docs-appfile-custom"
description: |-
  Customization blocks change the behavior of Otto for a specific application
  type or infrastructure.
---

# Customization Configuration

Customization blocks change the default behavior of Otto for application types,
dependencies, infrastructure, and more.

This page assumes you're familiar with the
[Appfile syntax](/docs/appfile/syntax.html) already.

## Example

Customization configuration looks like the following:

```
customization {
    ruby_version = "2.1"
}

customization "app" {
    ruby_version = "2.2.3"
}
```

## Description

`customization` blocks configure custom behavior for an Appfile that
deviates from the built-in defaults. Multiple customization blocks
can be specified. If customization blocks have to be merged, Otto does this
at a per key level.

Customization blocks can be named. The example above shows a customization
block that is both unnamed and named ("app"). The name of the customization
block becomes a filter for what the customization applies to.

The available names of customization blocks are well defined. If no name
is specified, "app" is assumed. The names can be in the following format,
where all capital letters are placeholders:

  * "" (blank) - Equivalent to "app". See "app" below.

  * "app" - Applies the customization to the app that this Appfile defines.

  * "infra" - Applies the customization to the infrastructure created by
      this application on deploy.

Within the customization blocks, the available options are dependent on
the application type or infrastructure type itself. See the respective
documentation for a reference. For example, see [app types](/docs/apps/index.html)
for a list of app types and their available customizations.

## Syntax

The full syntax is:

```
customization [FILTER] {
    CONFIG ...
}
```

where `CONFIG` is:

```
KEY = VALUE

KEY {
    CONFIG
}
```
