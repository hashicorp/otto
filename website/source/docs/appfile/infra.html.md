---
layout: "docs"
page_title: "Infrastructure - Appfile"
sidebar_current: "docs-appfile-infra"
description: |-
  The infrastructure configuration describes what infrastructures
  an application can be deployed to.
---

# Infrastructure Configuration

The infrastructure configuration describes what infrastructures
an application can be deployed to.

This page assumes you're familiar with the
[Appfile syntax](/docs/appfile/syntax.html) already.

## Example

The infrastructure configuration looks like the following:

```
infrastructure "production" {
    type = "aws"
    flavor = "vpc-public-private"
}
```

## Description

The `infrastructure` block configures a potential infrastructure
target for an application. Multiple infrastructure blocks can be
specified in an Appfile.

The [project](/docs/appfile/project.html) block ties a project (and
therefore an application) to a default infrastructure.

The `infrastructure` block allows the following keys to be set:

  * `type` (string) - The type of the infrastructure. The full list
      of available infrastructure types is [available here](/docs/infra).

  * `flavor` (string) - The flavor of the infrastructure. This will be
      documented on the page of the type of the infrastructure chosen.

## Syntax

The full syntax is:

```
infrastructure NAME {
	type = TYPE
	flavor = FLAVOR
}
```
