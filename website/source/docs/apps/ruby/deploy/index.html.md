---
layout: "app_ruby"
page_title: "Build & Deploy - Ruby App Type"
sidebar_current: "docs-ruby-deploy"
description: |-
  Otto defaults to assuming your Ruby application is a public-facing web
  application, and deploys it with this assumption.
---

# Build & Deploy

Otto defaults to assuming your Ruby application is a public-facing web
application, and deploys it with this assumption.

This page documents
all the common deployment choices made for all infrastructures. The sidebar
on the left can be used to view infrastructure-specific choices that are
made for certain infrastructure targets.

## Common Points

Below is an unordered list of common points about the build and deploy
process. Please see the [customizations](/docs/apps/ruby/customization.html)
page for a list of behavior that can be changed.

  * The application is deployed behind [Phusion Passenger](https://www.phusionpassenger.com/)
    and [Nginx](http://nginx.org/).

  * `bundle install --deployment` is used during the build process to
    get all the dependencies for your applicaiton.
