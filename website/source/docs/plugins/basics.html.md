---
layout: "docs"
page_title: "Basics - Plugins"
sidebar_current: "docs-plugins-basics"
description: |-
  Basic guide to how plugins work in Otto.
---

# Plugin Basics

This page documents the basics of how the plugin system in Otto
works, and how to setup a basic development environment for plugin development
if you're writing an Otto plugin.

If you're only looking for how to use plugins, see the
[plugin installation page](/docs/plugins/install.html). If you're interested
in developing plugins but haven't read the guide on installation, read
that as well since it'll be important to understand how to install the
plugin you create.

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

Otto provides high-level libraries for making the creation of plugins
simple and to ensure a common behavior.

## Developing a Plugin

Developing a plugin is simple. The only knowledge necessary to write
a plugin is basic command-line skills and basic knowledge of the
[Go programming language](http://golang.org).

-> **Note:** A common pitfall is not properly setting up a
`$GOPATH`. This can lead to strange errors. You can read more about
this [here](https://golang.org/doc/code.html) to familiarize
yourself.

Create a new Go project somewhere in your `$GOPATH`. If you're a
GitHub user, we recommend creating the project in the directory
`$GOPATH/src/github.com/USERNAME/otto-plugin-NAME`, where `USERNAME`
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
	plugin.Serve(&plugin.ServeOpts{
		AppFunc: AppFactory,
	})
}
```

And that's basically it! You'll have to change the argument given to
`plugin.Serve` to be your actual plugin, but that is the only change
you'll have to make. See the GoDoc or the
[example plugin](https://github.com/hashicorp/otto-example-app-plugin)
for more details.

Instead of going into extreme technical detail here, we've uploaded a
really basic [example app plugin](https://github.com/hashicorp/otto-example-app-plugin).
Please use that as a guide for developing your own plugin combined with
the GoDoc of Otto itself.
