---
layout: "app_ruby"
page_title: "Customization - Ruby App Type"
sidebar_current: "docs-ruby-customization"
description: |-
  This page documents the [Customizations](/docs/appfile/customization.html)
  that are available to change the behavior of Ruby applications with Otto.
---

# Customization

This page documents the [customizations](/docs/appfile/customization.html)
that are available to change the behavior of Ruby applications with Otto.

## Type: "ruby"

Example:

```
customization "ruby" {
    ruby_version = "2.1"
}
```

Available options:

  * `ruby_version` (string) - The Ruby version to install for development
    and deployment. This defaults to 2.2. Note that for Ruby 1.9.3, you
    need to specify "1.9.1". This is a strange quirk due to upstream dependencies.
