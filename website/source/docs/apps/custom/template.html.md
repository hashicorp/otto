---
layout: "app_custom"
page_title: "Templates - Custom App Type"
sidebar_current: "docs-custom-template"
description: |-
  Some of the configurations for the custom app type are rendered as templates.
  This page documents some of the variables available for templates as well
  as the basic syntax.
---

# Templates

Some of the configurations for the custom app type are rendered as templates.
This page documents some of the variables available for templates as well
as the basic syntax.

## Syntax

The syntax of templates in Otto is the
[Django templating syntax](https://docs.djangoproject.com/en/1.8/topics/templates/#the-django-template-language).

## Available Variables

  * `name` (string) - Name of the application.
  * `dev_ip_address` (string) - IPv4 address allocated for this development
      environment.
  * `path.working` (string) - Directory where the Appfile is.
  * `path.compiled` (string) - Directory where the compiled data for the
      Appfile is.
