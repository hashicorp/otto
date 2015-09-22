---
layout: "docs"
page_title: "Customization - Appfile"
sidebar_current: "docs-appfile-custom"
description: |-
  Customization blocks change the behavior of Otto for a specific application
  type or infrastructure.
---

# Customization Configuration

Customization blocks change the behavior of Otto for a specific application
type or infrastructure.

This page assumes you're familiar with the
[Appfile syntax](/docs/appfile/syntax.html) already.

## Example

Customization configuration looks like the following:

```
customization "ruby" {
    ruby_version = "2.1"
}
```

## Description

`customization` blocks configure custom behavior for an Appfile that
deviates from the built-in defaults. Multiple customization blocks
can be specified, but only one per type (such as "ruby" in the
example above).

The available types of a customization block are defined by the
application as well as the infrastructure. See the respective documentation
for those for a list of customization types.
The contents of a customization block are defined by the type itself
within the same documentation.

## Syntax

The full syntax is:

```
customization TYPE {
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
