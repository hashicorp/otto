---
layout: "app_php"
page_title: "Build & Deploy - PHP App Type"
sidebar_current: "docs-php-deploy"
description: |-
  Otto defaults to assuming your PHP application is a public-facing web
  application, and deploys it with this assumption.
---

# Build & Deploy

Otto defaults to assuming your PHP application is a public-facing web
application, and deploys it with this assumption.

This page documents all the common deployment choices made for all
infrastructures. The sidebar on the left can be used to view
infrastructure-specific choices that are made for certain infrastructure
targets.

## Common Points

Below is an unordered list of common points about the build and deploy
process. Please see the [customizations](/docs/apps/php/customization.html)
page for a list of behavior that can be changed.

  * The application is deployed behind [Apache](https://httpd.apache.org/).

  * The same list of PHP modules made available for
    [development](/docs/apps/php/dev.html) are also installed in the deployed
    environment.

  * If a `composer.json` file is detected, then Otto installs Composer
    and runs `composer install --no-dev` during the build process to get all
    the dependencies for your application.
