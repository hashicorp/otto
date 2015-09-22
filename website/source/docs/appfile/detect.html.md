---
layout: "docs"
page_title: "Detection - Appfile"
sidebar_current: "docs-appfile-detect"
description: |-
  As noted in the [Appfile concepts](/docs/concepts/appfile.html) page,
  everything within an Appfile is optional. Otto _detects_ smart defaults
  for your application, and fills in any missing information with these
  defaults. This page describes how Otto does this.
---

# Detection

As noted in the [Appfile concepts](/docs/concepts/appfile.html) page,
everything within an Appfile is optional. Otto _detects_ smart defaults
for your application, and fills in any missing information with these
defaults. This page describes how Otto does this.

The purpose of detection is to both make a fantastic getting started
experience with Otto, and to enable Otto to work with projects that don't
specifically configure it.

## Detected Properties

* **Application name** is detected based on the name of the directory
  the Appfile is in.

* **Application type** is detected using filename pattern matching. For
  example, the existence of a "Gemfile" signals to Otto that the application
  is likely a Ruby project. Detection fingerprints are documented within
  each individual [application type](/docs/apps) section.

* **Project name** is the detected application name.

* **Infrastructure** is currently hardcoded to be "aws" with the "simple"
  flavor. A future version of Otto may do something more clever here.

## Merging

If an explicitly written Appfile is missing any sections, then Otto merges
it with the detected defaults. For example, the following is a completely
valid Appfile:

```
application {
    name = "my-app"
}
```

Despite "type" missing, project configuration, and infrastructure
configuration, the above is a valid Appfile. This is because Otto will merge
this with the detected properties. Any explicitly written configurations
will take priority and overwrite any detected properties.
