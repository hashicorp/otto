---
layout: "app_php"
page_title: "Customization - PHP App Type"
sidebar_current: "docs-php-customization"
description: |-
  This page documents the [Customizations](/docs/appfile/customization.html)
  that are available to change the behavior of PHP applications with Otto.
---

# Customization

This page documents the [customizations](/docs/appfile/customization.html)
that are available to change the behavior of PHP applications with Otto.

## Type: "php"

Example:

```
customization "php" {
  php_version="5.5"
}```

Available options

  * `php_version` (string) - The PHP version to install for development
    and deployment.  The currently supported versions for Otto are 5.4, 5.5 and
    5.6.  Support for PHP is provided by the PPA repositories created by
    [Ondřej Surý](https://launchpad.net/~ondrej). 