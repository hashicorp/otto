---
layout: "app_go"
page_title: "Build & Deploy - Go App Type"
sidebar_current: "docs-go-deploy"
description: |-
  Otto defaults to assuming your Go application is a private service
  and deploys it with this assumption.
---

# Build & Deploy

Otto defaults to assuming your Go application is a private service
and deploys it with this assumption.

This page documents
all the common deployment choices made for all infrastructures. The sidebar
on the left can be used to view infrastructure-specific choices that are
made for certain infrastructure targets.

-> **NOTE:** Otto 0.1 doesn't yet support public-facing Go applications.
This will be coming very shortly.

## Build

The build process for Go involves installing Go, setting up the GOPATH,
and calling `go build` in the project root with a custom output path so
we can store the binary.

The binary is configured to launch with no command line args for deploy.

A future version of Otto will allow customizing the launch parameters
so that things such as configuration files can be used.

## Deploy

To deploy, the Go process is simply launched.
