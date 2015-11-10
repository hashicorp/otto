---
layout: "app_ruby"
page_title: "Development - Ruby App Type"
sidebar_current: "docs-ruby-dev"
description: |-
  The development environment built for Ruby applications is built for
  general Ruby development with a lean towards web development.
---

# Development

The development environment built for Ruby applications is built for
general Ruby development with a lean towards web development.

Please see the [customizations](/docs/apps/ruby/customization.html)
page for details on how to customize some of the behavior on this page.

## Pre-Installed Software

  * **Ruby and RubyGems** - The version of Ruby is determined based on
      [customizations](/docs/apps/ruby/customization.html). See that page
      for defaults as well.
  * **Bundler**
  * **Node.JS** - A common requirement for web development.
  * **Git, Mercurial** - Useful for Bundler
  * **PostgreSQL dev headers** - A common requirement for web development.

## Common Issues and Solutions

**I can't access my web application!** When you start your web server,
such as via `rails s`, make sure it is bound to `0.0.0.0`. For Rails,
the command is `rails s -b 0.0.0.0`. If this isn't specified, many
servers default to binding to localhost. This means that only
that machine can access it. Since the Otto development environment is in
a virtual machine, you have to configure it to bind to `0.0.0.0` to
allow your host machine to access it.
