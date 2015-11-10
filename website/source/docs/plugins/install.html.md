---
layout: "docs"
page_title: "Install - Plugins"
sidebar_current: "docs-plugins-install"
description: |-
  How to install Otto plugins.
---

# Plugin Installation

This page documents how to install a plugin.

Otto plugins are distributed as standalone applications. Otto uses
a multi-process plugin model. Otto plugins should be named in the
format of `otto-plugin-*` where `*` is anything. If the plugin doesn't
follow this naming convention, Otto will not be able to discover it. The
name itself doesn't actually matter, so if a plugin isn't following the
naming convention, just rename it to the proper format.

You can verify the plugin is an Otto plugin by executing it directly.
It should give you an error message similar to that below:

```
$ otto-app-ruby
This binary is an Otto plugin. These are not meant to be
executed directly. Please execute `otto`, which will load
any plugins automatically.
```

Once named properly, the easiest way to install a plugin is to place
it in a proper directory. Otto automatically discovers plugins
in the following directories in the given order.

1. The directory where Otto is, or the executable directory.

1. `~/.otto.d/plugins` on Unix systems or `%APPDATA%/otto.d/plugins` on Windows.

1. The current working directory.

That's it! As long as the plugin is named properly and is in one of those
directories, Otto will automatically load it.

## Conflicting Plugins

It is possible in a few scenarios for plugins to conflict. The behavior
in this case is well defined and documented here.

The names of plugins don't matter. If two plugins are named the same but
exist in two different directories that are auto-discovered, this doesn't
cause any issues.

If two plugins claim to support the same application type, the plugin
loaded later takes precedence. "Later" is defined by being found last in
the auto-discovery process. The auto-discovery process follows the directories
above in that order, and in alphabetical order within the directory itself.

If two plugins claim to auto-detect the same way (e.g. plugin A and B both
detect based on "foo.txt"), then the plugin loaded later takes precedence.
"Later" is defined in the paragraph above.
