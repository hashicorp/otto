---
layout: "docs"
page_title: "Appfile"
sidebar_current: "docs-concepts-appfile"
description: |-
  Before performing any operation with Vault, the connecting client must be authenticated.
---

# Appfile

The "Appfile" is the file that Otto uses as a source of configuration
for an application. It tells Otto just enough information for Otto to
manage the application from development to deployment.

An example Appfile is shown below. This Appfile is completely valid.

```
application {
    name = "example"
    type = "ruby"
}

customization "ruby" {
    ruby_version = "2.0"
}
```

The complete syntax and sections of an Appfile are
[documented in the Appfile section](/docs/appfile).

The goal of this page is to explain the purpose of an Appfile,
what Otto does with the Appfile, etc.

## Everything is Optional

Appfiles are completely optional. If no Appfile exists, Otto will
inspect the project and attempt to detect what kind of application it is.
It uses this detection to generate an Appfile internally. This generated
Appfile is never saved externally in a project.

In the Otto [getting started guide](/intro/getting-started), we develop
and deploy a complete application without an Appfile.

In addition to the Appfile itself being optional, all the individual
blocks ("application", "customization", etc.) within an Appfile are also
optional. If an Appfile exists and a block such as "application" is missing,
Otto will just use the detected Appfile data and merge it in.

You can play with this on your own by creating a directory with an
identifiable application file such as "index.php", and playing around with
omitting various parts of an Appfile and seeing what Otto generates.

## Source vs. Compiled

The "Appfile" itself is the source code form of Otto configuration.
Otto takes this file and [compiles](/docs/concepts/compile.html) it to
an internal representation that is used by all the Otto subcommands,
such as `otto dev`.

For people writing Appfiles, the important concept to understand here is
that modifying an Appfile has no affect on the behavior of Otto until
that Appfile is compiled.

This is a very useful feature, since modifications to an Appfile that
might be pushed to version control don't affect other developers using
Otto until they decide to recompile.

More details about compilation are covered on the
[compilation concepts](/docs/concepts/compile.html) page.
