---
layout: "docs"
page_title: "Infra Types"
sidebar_current: "docs-infra"
description: |-
  The infrastructure type in the Appfile tells Otto what kind of
  infrastructure you'd like to use to deploy your application
---

# Infrastructure Types

The infrastructure type in the [Appfile](/docs/concepts/appfile.html)
tells Otto what kind of infrastructure you'd like to use to deploy your
application. Otto uses this information to create a base infrastructure, to
build your application for that infrastructure, and finally to deploy your
application into that infrastructure.

This section documents the infrastructure types that Otto supports by default
as well as how to work with those infrastructure types.

Use the navigation to the left to read about the available infrastructure types.

A few notes that are common to all infrastructure types can be found below.

## Credentials

For most infrastructure types, Otto needs API credentials to be able to manage
resources for your application.

Otto will ask for any credentials it needs the first time it needs to interact
with infrastructure.

Otto maintains a local encrypted cache file of these credentials, which it will
also create on first run. You will be prompted for a password to decrypt this
file anytime Otto interacts with your infrastructure. To skip the prompt, you
can specify this password via the `OTTO_CREDS_PASSWORD` environment variable.

