---
layout: "docs"
page_title: "Directory"
sidebar_current: "docs-concepts-directory"
description: |-
  The "directory" is the name for the layer in Otto that stores
  the state of all infrastructures, applications, deploys, builds, etc.
---

# Directory

The "directory" is the name for the layer in Otto that stores
the state of all infrastructures, applications, deploys, builds, etc.

When you run `otto status`, the data that it outputs is coming from
the directory. When an `otto deploy` happens, it verifies that all dependencies
are deployed by querying the directory. When you create a new Appfile
and compile it, that new application is inserted into the directory.

The primary purpose of the directory is to enable shared state between
multiple applications managed by Otto. By extension, the directory is
the primary mechanism for [collaboration](/docs/concepts/collaboration.html)
in Otto.

## Local Directory

By default when you use Otto, it uses a local directory. You can see
the database for the directory at `~/.otto.d/directory/otto.db`. Note
that this file will only exist once you've used Otto before.

Because the default directory is local, it is not suited for teamwork
and collaboration. If one developer on your team deploys an application,
and then you try to deploy an application, then Otto will deploy
two separate isolated instances. This is because the directory isn't
shared.

Do not attempt to use network storage to share the directory file. This
is unsafe and will not work.

There is more information on collaboration with Otto on the
[collaboration](/docs/concepts/collaboration.html) page.

## Shared Directories

For Otto 0.1, only the local directory is available. A future version of
Otto (in the very short term) will add support for remote directories, with
the first remote directory being part of HashiCorp's
[Atlas](https://atlas.hashicorp.com) offering. This directory service
on its own will be free.
