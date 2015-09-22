---
layout: "docs"
page_title: "Basics - Plugins"
sidebar_current: "docs-plugins-basics"
description: |-
  Basic guide to how plugins work in Otto.
---

# Plugin Basics

**NOTE: Plugins are not yet available with Otto 0.1.** We'll expose the plugin
interface publicly with Otto 0.2. We've started documented the process
below, however, if you want to look forward to it. The contents on this
page will likely change.

This page documents the basics of how the plugin system in Otto
works, and how to setup a basic development environment for plugin development
if you're writing an Otto plugin.

~> **Advanced topic!** Plugin development is a highly advanced
topic in Otto, and is not required knowledge for day-to-day usage.
If you don't plan on writing any plugins, we recommend not reading
this section of the documentation.

## How it Works

The plugin system for Otto is based on multi-process RPC.

Otto executes these binaries in a certain way and uses Unix domain
sockets or network sockets to perform RPC with the plugins.

If you try to execute a plugin directly, an error will be shown:

```
$ otto-app-ruby
This binary is an Otto plugin. These are not meant to be
executed directly. Please execute `otto`, which will load
any plugins automatically.
```

The code within the binaries must adhere to certain interfaces.
The network communication and RPC is handled automatically by higher-level
Otto libraries. The exact interface to implement is documented
in its respective documentation section.

## Installing a Plugin

The easiest way to install a plugin is to name it correctly, then place it
in the proper directory. To name a plugin correctly, make sure the binary is
named `otto-TYPE-NAME`. For example, `otto-app-ruby` is an "app" type plugin
named "ruby". Valid types for plugins are down this page more.

Once the plugin is named properly, Otto automatically discovers plugins in the
following directories in the given order. If a conflicting plugin is found
later, it will take precedence over one found earlier.

1. The directory where Otto is, or the executable directory.

1. ~/.otto.d/plugins on Unix systems or %APPDATA%/otto.d/plugins on Windows.

1. The current working directory.

## Developing a Plugin

Developing a plugin is simple. The only knowledge necessary to write
a plugin is basic command-line skills and basic knowledge of the
[Go programming language](http://golang.org).

-> **Note:** A common pitfall is not properly setting up a
<code>$GOPATH</code>. This can lead to strange errors. You can read more about
this [here](https://golang.org/doc/code.html) to familiarize
yourself.

Create a new Go project somewhere in your `$GOPATH`. If you're a
GitHub user, we recommend creating the project in the directory
`$GOPATH/src/github.com/USERNAME/otto-NAME`, where `USERNAME`
is your GitHub username and `NAME` is the name of the plugin you're
developing. This structure is what Go expects and simplifies things down
the road.

With the directory made, create a `main.go` file. This project will
be a binary so the package is "main":

```
package main

import (
	"github.com/hashicorp/otto/plugin"
)

func main() {
	plugin.Serve(new(MyPlugin))
}
```

And that's basically it! You'll have to change the argument given to
`plugin.Serve` to be your actual plugin, but that is the only change
you'll have to make. The argument should be a structure implementing
one of the plugin interfaces (depending on what sort of plugin
you're creating).

Otto plugins must follow a very specific naming convention of
`otto-TYPE-NAME`. For example, `otto-app-ruby`, which
tells Otto that the plugin is an app that serves Ruby applications.
