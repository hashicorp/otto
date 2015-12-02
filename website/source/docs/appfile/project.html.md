---
layout: "docs"
page_title: "Project - Appfile"
sidebar_current: "docs-appfile-project"
description: |-
  The project configuration binds an application to a project and
  binds a project to an infrastructure.
---

# Project Configuration

The project configuration binds an application to a project and
binds a project to an infrastructure.

This page assumes you're familiar with the
[Appfile syntax](/docs/appfile/syntax.html) already.

## Example

The project configuration looks like the following:

```
project {
    name = "my-app"
    infrastructure = "production"
}
```

## Description

The `project` block ties an application to a project. A project is
a group of applications that share the same infrastructure.

The `project` block allows the following keys to be set:

  * `name` (string) - The name of the project. This must be unique
      to an infrastructure.

  * `infrastructure` (string) - The name of the infrastructure that this
    project should be deployed onto by default. This should match the name of a
    configured [infrastructure](/docs/appfile/infra.html). In the example
    above, the infrastructure is named "production".

For people with multiple applications, the `project` block is usually
shared via [imports](/docs/appfile/import.html) in the Appfile.

## Syntax

The full syntax is:

```
project {
	name = NAME
	infrastructure = TYPE
}
```
