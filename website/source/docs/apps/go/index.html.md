---
layout: "app_go"
page_title: "Go - App Types"
sidebar_current: "docs-go-index"
description: |-
  The Go application type is used to develop general Go-based applications.
---

# Go App Type

**Type:** `go`

The Go application type is used to develop general Go-based applications.

## GOPATH Detection

Go is very particular about the directory structure of code so that
loading dependencies and building binaries happens in a consistent way.
See the Go [documentation on the GOPATH](https://golang.org/doc/code.html#GOPATH)
for more information.

Without any configuration, Otto will attempt to discover the proper import path
automatically. This works for getting started with Go and Otto, but we
recommend [specifying the import path explicitly](/docs/apps/go/customization.html).

To discover the import path, Otto inspects your environment for a GOPATH
environment variable. If it is found, it determines the subdirectory (if any)
that your application is in and uses that for the import path. If any
of this fails, then Otto will put your application at a default location
outside of the GOPATH. It is highly likely that this result in broken builds.

## Godep

If you're using [Godep](https://github.com/tools/godep), then this will
continue to work well with Otto. No special changes are necessary. The
workflow for development is the same as outside: use `godep go` instead of
`go`.

Otto will detect if you're using Godep by seeing the `Godep` folder and will
automatically install it for you, so it should be available on the PATH.
