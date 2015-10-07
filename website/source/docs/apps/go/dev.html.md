---
layout: "app_go"
page_title: "Development - Go App Type"
sidebar_current: "docs-go-dev"
description: |-
  The development environment built for Go applications is built for
  general Go development.
---

# Development

The development environment built for Go applications is built for
general Go development.

Please see the [customizations](/docs/apps/go/customization.html)
page for details on how to customize some of the behavior on this page.

## Pre-Installed Software

  * **Go**
  * **Git, Mercurial, Bazaar** - For `go get`

In addition to the above, Otto automatically creates and configures
the `GOPATH` environment variable and syncs your project into the
proper import path.

This lets you immediately SSH into with `otto dev ssh` and get started
with `go get ./...` and `go build`.
