---
layout: "docs"
page_title: "Import - Appfile"
sidebar_current: "docs-appfile-import"
description: |-
  The `import` statement can be used within an Appfile to import fragments
  of an Appfile from other sources, including the local filesystem and
  remote URLs.
---

# Importing Appfile Fragments

The `import` statement can be used within an Appfile to import fragments
of an Appfile from other sources, including the local filesystem and
remote URLs.

Imports are an important concept for
[collaborating](/docs/concepts/collaboration.html)
on applications with Otto. They are the recommended way to share
configurations for Appfiles, especially the `project` and
`infrastructure` blocks.

This page assumes you're familiar with the
[Appfile syntax](/docs/appfile/syntax.html) already.

## Example

Imports look like the following:

```
import "github.com/hashicorp/otto-shared/database" {}
```

## Description

`import` blocks fetch and copy contents from another Appfile into this
Appfile. While it is similar to copying and pasting the contents of another
Appfile, this isn't exactly correct, since imported contents are merged
after the entire Appfile is loaded. This means the location of the
`import` statement doesn't matter.

Multiple `import` statements can be specified. In this case, their contents
are merged in the order they were specified within the original Appfile.

Due to the syntax of HCL, you must specify a trailing `{}` at the end of
the import statement. There is no inner configuration allowed for imports.

The URL allowed for imports is identical to the
[allowed sources for dependencies](/docs/appfile/dep-sources.html).
Within the source, Otto loads the "Appfile" and merges it. You cannot
specify an alternate filename. To store multiple Appfiles in a single
source, use folders.

## Syntax

The full syntax is:

```
import URL {}
```
