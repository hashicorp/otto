---
layout: "docs"
page_title: "Collaboration"
sidebar_current: "docs-concepts-collab"
description: |-
  Development and especially deployment is a collaborative process. Typically
  multiple developers are working on and deploying an application. It is
  important to understand how Otto can be used in a team environment.
---

# Collaboration

Development and especially deployment is a collaborative process. Typically
multiple developers are working on and deploying an application. It is
important to understand how Otto can be used in a team environment.

## Appfiles in Version Control

First, Appfiles themselves should go into version control.

This allows anyone on the team to clone a repository and immediately
start using Otto to develop and deploy without any prior knowledge. Typically
one developer sets up the initial Appfile, and everybody on the team benefits.

The real power behind this is even people who know very little about an
application can easily get it up and running. For example: a designer, a
new employee, a manager, or even just an experienced engineer working on
a very old project.

## Dependencies as URLs

[Dependencies](/docs/concepts/deps.html) are configured in Otto using
URLs. These URLs can be to files, GitHub, BitBucket, and more. This makes it
easy for one team to point to another team's application as a dependency,
and automatically get those updates as they occur.

In addition to this, an application only needs to specify its immediate
dependencies. For example if application "foo" depends on "bar", and "bar"
depends on "baz", then the developers of "foo" only need to know that
they depend on "bar". Otto automatically determines the fully tree of
dependencies.

Both of these features are very important for collaboration.

Dependencies as URLs keeps version control and the Appfile as the source
of truth. And only needing to know immediate dependencies limits the complexity
of working with other teams, and also gives the other teams flexibility to
change their dependencies without negatively affecting consumers of their
application.

## Imports for Shared Configuration

Appfiles support [import statements](#)
to import configuration from many sources. This allows organizations to
create shared repositories of Appfile fragments that can be used by teams.

For example, if an organization made a standardized Ruby Appfile, a team
could potentially make an Appfile just like this:

```
import "github.com/myorg/otto-shared/ruby" {}

application {
    name = "my-app"
    type = "ruby"
}
```

These import fragments are downloaded as part of the
[compilation step](/docs/concepts/compile.html), so just like many other
features, changes to upstream imports don't affect a developer working
with Otto until they recompile.

## Shared Infrastructure

Infrastructures with the same names (as specified in the Appfile) are
only created once. This allows multiple applications to share a single
infrastructure. Multiple teams can then deploy their applications together
into the same place.

Paired with Appfile imports above to share the fragment to configure
infrastructures, this makes it very easy to work on multiple applications
that are meant to co-exist in the same infrastructure.

~> **NOTE:** Otto 0.1 doesn't support [shared directories](/docs/concepts/directory.html)
   yet, making this point basically impossible at the time. This is a top
   priority feature to ship in an upcoming release and will be addressed soon.
