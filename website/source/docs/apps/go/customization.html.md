---
layout: "app_go"
page_title: "Customization - Go App Type"
sidebar_current: "docs-go-customization"
description: |-
  This page documents the customizations
  that are available to change the behavior of Go applications with Otto.
---

# Customization

This page documents the [customizations](/docs/appfile/customization.html)
that are available to change the behavior of Go applications with Otto.

## Type: "go"

Example:

```
customization "go" {
    go_version = "1.4.2"
}
```

Available options:

  * `go_version` (string) - The Go version to install for development
    and for building the application for deployment. This defaulits to 1.5.1.

  * `import_path` (string) - The import path of this application so Otto
    knows where to place it in the GOPATH. Example: "github.com/hashicorp/foo"
