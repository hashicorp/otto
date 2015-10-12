---
layout: "docs"
page_title: "Dependencies"
sidebar_current: "docs-concepts-deps"
description: |-
  Application dependencies are a first class feature in Otto.
---

# Dependencies

Application dependencies are a first class feature in Otto.

Modern applications often depend on many components. With the ever-growing
["microservices"](http://martinfowler.com/articles/microservices.html) trend,
it isn't abnormal for an application to depend on dozens of other applications.
A common difficulty with microservices is development environments and
orchestrating deploys. Otto solves both of these problems.

For development environments, Otto automatically installs, configures, and
starts any dependencies of an application. For deploys, Otto ensures that
all the dependencies are deployed, can deploy those dependencies for you,
and doesn't duplicate dependency deploys across multiple applications.

## Specifying Dependencies

Within the [Appfile](/docs/concepts/appfile.html), an application
can specify all of its dependencies. These dependencies can be other
applications, infrastructure components such as a queue, or external
services such as Datadog.

An application only needs to specify its immediate dependencies. If an
App "foo" depends on "bar", and "bar" depends on "baz", then the Appfile
for "foo" only needs to specify "bar." Otto is clever enough to discover that
"bar" depends on "baz" by inspecting the Appfiles of all the dependencies.

Dependencies are specified in the "application" block like so:

```
application {
    name = "example"
    type = "ruby"

    dependency {
        source = "github.com/hashicorp/otto/examples/mongodb"
    }
}
```

The "source" string can be local path, HTTP URL, Git URL, and many more.
The full reference to dependency sources can be seen in the
[dependency sources](/docs/appfile/dep-sources.html)
page.

## Compiling Dependencies

Dependencies are fetched during [compilation](/docs/concepts/compile.html).
They are not updated any other time. If an upstream dependency changes,
any downstream applications that want to bring in that change must
`otto compile`.

This functionality is very nice since it allows developers working on
a feature to not be impacted by potentially breaking changes upstream
until they're ready to.

Note that while downstream applications may not update their upstreams,
if the upstream dependency is already deployed, then any downstream
dependencies will see the new version of the dependency when deployed.

-> **A note on versioning:** Dependency versioning will be coming as a first class feature very
shortly into Otto, but isn't currently supported. For now, Otto will
always fetch whatever app is returned by the source URL.

## Communicating with Dependencies

Whether in development or production, communicating with dependencies
is the same. Otto automatically deploys and manages a
[Consul](https://consul.io) cluster for service discovery. To find
a dependency, just use DNS with `appname.service.consul`.

For example, our example above with MongoDB can be reached at
`mongodb.service.consul`. For ports, dependencies must use well known
ports. A future version of Otto will do automatic port assignment
and expose that information via environment variables to the
application.

For more information on how Otto automatically sets this up, see
the [concepts page on "foundations"](/docs/concepts/foundations.html).
