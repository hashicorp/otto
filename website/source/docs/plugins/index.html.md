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

Plugins were introduced in Otto 0.2. If you're using an earlier version of
Otto, plugins will not work. From Otto 0.2 onwards, any available plugins
can be used. If a plugin isn't compatible with your version of Otto, Otto
will give you a nice error message.

If you're only interested in using plugins, see the
[plugin installation](/docs/plugins/install.html) page. If you're interested
in developing plugins, see the [development basics](/docs/plugins/basics.html)
page.

~> **Advanced topic!** Plugin development is a highly advanced topic in
   Otto, and is not required knowledge for day-to-day usage. If you don't plan
   on writing any plugins, we recommend not reading this section of the
   documentation.
