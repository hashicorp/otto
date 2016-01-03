---
layout: "docs"
page_title: "ScriptPacks- Appfile"
sidebar_current: "docs-plugins-scriptpack"
description: |-
  ScriptPacks are libraries of shell functions to run with application
  types and other Otto plugins.
---

# ScriptPacks

ScriptPacks are fully self-contained shell scripting libraries that
other plugins such as [application types](/docs/plugins/app.html)
use to perform logic on guest machines.

The idea is that all Otto functionality that happens on a user's machine
should be written in Go, and all functionality that happens on a remote
machine (dev, deploy, etc.) should be implemented by calling ScriptPack
functions.

ScriptPacks are tested in isolation, enabling you to unit test specific
functionality without a long feedback loop. Furthermore, ScriptPacks
encourage reusability of functionality across Otto plugins.

~> **Advanced topic!** Plugin development is a highly advanced topic in
   Otto, and is not required knowledge for day-to-day usage. If you don't plan
   on writing any plugins, we recommend not reading this section of the
   documentation.

If you're interested in ScriptPack development and usage, then read on. The
remainder of this page will assume you're familiar with
[plugin basics](/docs/plugins/basics.html)
and that you already have a basic development environment setup.

## Skeleton

A basic skeleton for writing a new ScriptPack can be
[found in the Otto repository](https://github.com/hashicorp/otto/tree/master/builtin/scriptpack/skeleton). This is a fully functioning example that any new ScriptPacks
can be based off of. For a real example, see any ScriptPacks within the
Otto repository.

The basic ScriptPack folder structure contains:

  * **Go files:** We use Go as an API to work with ScriptPacks from
    Go. The shell files within the data directory do not require Go, however.
    These Go files let you use the ScriptPack as a Go library from your
    other Otto plugins.

  * **`data` directory:** This contains a set of shell scripts that contain
    all the functionality of the ScriptPack. This directory will be available
    at the env var `SCRIPTPACK_NAME_ROOT` where `NAME` is the name of the
    ScriptPack.

  * **`test` directory:** This contains
    [BATS](https://github.com/sstephenson/bats) tests for the various
    functions within the `data` directory. You can easily run these tests
    by using Otto itself (covered below).

## Developing ScriptPacks

ScriptPacks are developed using Otto itself. Otto has a
[scriptpack app type](/docs/apps/scriptpack). This app type introduces
CLI commands for easily interacting with ScriptPacks.

Bring up a ScriptPack dev environment with `otto dev`. This will start
a VM that is used for testing ScriptPacks. Next, use `otto dev scriptpack-test`
to test your ScriptPack:

    $ otto dev scriptpack-test test/foo.bats
    ...

This will use a container to isolate the test file and test your ScriptPack. By
using containers, we can run the tests on "clean" systems to test a variety
of cases quickly.

The Otto workflow enables a fast feedback cycle on potentially complex
shell functionality. This enables plugin developers to more quickly develop
complex [app type plugins](/docs/plugins/app.html).

## Dependencies

ScriptPacks can depend on other ScriptPacks. Within `main.go`, you can
see the skeleton ScriptPack depends on "stdlib":

    var ScriptPack = scriptpack.ScriptPack{
        ...
        Dependencies: []*scriptpack.ScriptPack{
            &stdlib.ScriptPack,
        },
    }


As you can see, dependencies are done using Go directly. By using Go,
the dependency is "locked" to exactly the contents of the ScriptPack at
compile-time. For each dependent scriptpack, an environment variable
becomes available to access its data. For example, for the dependency above,
`SCRIPTPACK_STDLIB_ROOT` can be used to get the path to the root of that
dependency.

## Using ScriptPacks

To use a ScriptPack, see the [app type plugins](/docs/plugins/app.html)
page. App type plugins have a higher level API to make using ScriptPacks
easier within Otto.

ScriptPacks can also be used directly with the
[scriptpack package API](https://github.com/hashicorp/otto/tree/master/scriptpack).
This is a fairly low level API that you shouldn't need to use directly
for Otto plugins, but if you want to use ScriptPacks outside of Otto you
may need it.
