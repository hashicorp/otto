---
layout: "docs"
page_title: "Plugins - Appfile"
sidebar_current: "docs-plugins"
description: |-
  Otto was architected around a plugin-based architecture. All
  application types and infrastructure types
  are plugins, even the core built-in types.
---

# Plugins

Otto was architected around a plugin-based architecture. All
[application types](/docs/apps) and [infrastructure types](/docs/infra)
are plugins, even the core built-in types.

This section of the documentation gives a high-level overview of how to write
plugins for Otto. It does not hold your hand through the process, however, and
expects a relatively high level of understanding of Go, Otto semantics,
Unix, etc.

**Plugins are not yet available with Otto 0.1.** We'll expose the plugin
interface publicly with Otto 0.2.

~> **Advanced topic!** Plugin development is a highly advanced topic in
   Otto, and is not required knowledge for day-to-day usage. If you don't plan
   on writing any plugins, we recommend not reading this section of the
   documentation.
