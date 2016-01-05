---
layout: "app_java"
page_title: "Build & Deploy - Java App Type"
sidebar_current: "docs-java-deploy"
description: |-
  Otto defaults to assuming your Java application is a private service
  and deploys it with this assumption.
---

# Build & Deploy

Otto defaults to assuming your Java application is a private service
and deploys it with this assumption.

This page documents
all the common deployment choices made for all infrastructures. The sidebar
on the left can be used to view infrastructure-specific choices that are
made for certain infrastructure targets.

-> **NOTE:** Otto 0.1 doesn't yet support public-facing Java applications.
This will be coming very shortly.

## Build

The build process for Java involves installing Java, and calling `gradle build`
or `mvn compile` in  the project root with a custom output path so we can store
the binary.

The binary is configured to launch with no command line args for deploy.

A future version of Otto will allow customizing the launch parameters
so that things such as configuration files can be used.

## Deploy

To deploy, the Java process is simply launched.
