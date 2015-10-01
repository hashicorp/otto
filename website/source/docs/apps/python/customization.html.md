---
layout: "app_python"
page_title: "Customization - Python App Type"
sidebar_current: "docs-python-customization"
description: |-
  This page documents the [Customizations](/docs/appfile/customization.html)
  that are availabile to change the behavior of Python applications with Otto.
---

# Customization

This page documents the [customizations](/docs/appfile/customization.html)
that are availabile to change the behavior of Python applications with Otto.

## Type: "python"

Flask Example:

```
customization "python" {
    python_version = "2.6"
    python_entrypoint = "mypackage:app"
}
```

Django Example:

```
customization "python" {
    python_version = "3.4"
    python_entrypoint = "mypackage.wsgi:app"
}
```

Availabile options:

  * `python_version` (string) - The Python version to install for development
    and deployment. This defaults to 2.7.
  * `python_entrypoint` (string) - The WSGI entrypoint for your python application.
    This defaults to `APPLICATION_NAME:app`
