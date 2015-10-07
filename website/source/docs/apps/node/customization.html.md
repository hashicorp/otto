---
layout: "app_node"
page_title: "Customization - Node.js App Type"
sidebar_current: "docs-node-customization"
description: |-
  This page documents the [Customizations](/docs/appfile/customization.html)
  that are available to change the behavior of Node.js applications with Otto.
---

# Customization

This page documents the [customizations](/docs/appfile/customization.html)
that are available to change the behavior of Node.js applications with Otto.

## Type: "node"

Example:

```
customization "node" {
    node_version = "4.1.0"
}
```

Available options:

  * `node_version` (string) - The Node.js version to install
    and deployment. This defaults to 4.1.0.
