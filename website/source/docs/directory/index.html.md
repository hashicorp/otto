---
layout: "docs"
page_title: "Directory Backends - Appfile"
sidebar_current: "docs-directory"
description: |-
  Application dependencies and Appfile imports both take a URL as a parameter
  to tell Otto where to load the external Appfile.
---

# Directory Backends

The ["directory"](/docs/concepts/directory.html) is the name for the layer in
Otto that stores the state of all infrastructures, applications, deploys,
builds, etc.

Please read the [directory concepts](/docs/concepts/directory.html) page
for more information prior to reading about the various backends.

By default, Otto uses a local directory backend that stores data at
`~/.otto.d/directory/otto.db`. This is not well suited for sharing data
in a team and is not currently configurable.

For Otto 0.1, we focused on the single developer experience. In upcoming
releases, we'll be focusing more on teamwork. This will involve adding more
directory backends and making them more configurable.
