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

### Type: "ruby"

Example:

```
customization "ruby" {
    ruby_version = "2.1"
}
```

Available options:

  * `ruby_version` (string) - The Ruby version to install for development
    and deployment. This defaults to "detect", but can be any specific Ruby
    version. If "detect" is specified, the Ruby version will be detected.
    See the "ruby version detection" section below.

### Ruby Version Detection

By default, Otto will attempt to automatically detect the proper Ruby
version to use. If no Ruby version is detected, it will default to some
recent version (we try to keep this up to date but it depends on the
release process of Otto itself).

To detect the Ruby version, Otto will inspect your `Gemfile` to look
for a Ruby version specification. If one is found, that version will be
installed.

The detected version, if any, will be output during `otto compile`.

The detected version is installed using [ruby-install](https://github.com/postmodern/ruby-install).
We believe this is the best practice of installing Ruby and uses the
official Ruby source to compile. Due to the dev layers system, the compilation
only needs to happen once per system.
