---
layout: "docs"
page_title: "Plugins - Appfile"
sidebar_current: "docs-plugins"
description: |-
  Otto was architected around a plugin-based architecture. All
  application types and infrastructure types
  are plugins, even the core built-in types.
---

# App Plugins

**NOTE: Plugins are not yet available with Otto 0.1.** We'll expose the plugin
interface publicly with Otto 0.2. We've started documented the process
below, however, if you want to look forward to it.

Apps in Otto are responsible for development, build, and deploy of
an application. An example of an app type is [PHP](/docs/apps/php).

By writing an app plugin for Otto, you can add support to Otto for new
kinds of languages and frameworks. You can also override the defaults that
are built-in to Otto by replacing the built-in plugins.

The primary reasons to care about app plugins are:

  * You want to add support for a new language or framework.

  * You want to change the behavior of an existing plugin.

  * You want to support highly custom or legacy application types. This is
    very common in large organizations adopting Otto.

~> **Advanced topic!** Plugin development is a highly advanced topic in
   Otto, and is not required knowledge for day-to-day usage. If you don't plan
   on writing any plugins, we recommend not reading this section of the
   documentation.

If you're interested in provider development, then read on. The remainder of
this page will assume you're familiar with
[plugin basics](/docs/plugins/basics.html)
and that you already have a basic development environment setup.
