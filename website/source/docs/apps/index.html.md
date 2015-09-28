---
layout: "docs"
page_title: "App Types - Appfile"
sidebar_current: "docs-apps"
description: |-
  The application type in the Appfile
  is the most important piece of information that tells Otto how to function.
  Otto uses this information to know how to develop, build, and deploy
  the application.
---

# App Types

The application type in the [Appfile](/docs/concepts/appfile.html)
is the most important piece of information that tells Otto how to function.
Otto uses this information to know how to develop, build, and deploy
the application.

This section documents the application types that Otto supports by default
as well as how to work with those application types. If you're a Ruby
developer, for example, you should read the [Ruby](/docs/apps/ruby/index.html)
section.

In addition to the built-in types, Otto makes it possible to create your
own [custom types](/docs/apps/custom).

Use the navigation to the left to read about the available application types.

## Documentation Format

Each application type will follow a similar documentation format. The sidebar
will contain the following sections:

  * **Detection** documents how the application type is automatically
    detected in the absense of an explicit Appfile.

  * **Development** documents the development environment and how to work
    in development for this application type.

  * **Build and Deploy** documents the build and deploy process for
    the application. This may have subsections for different infrastructures.

  * **Customizations** documents the customizations that are available
    for this application type.
