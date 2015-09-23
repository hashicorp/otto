---
layout: "intro"
page_title: "Appfile"
sidebar_current: "gettingstarted-appfile"
description: |-
  Learn about the Appfile in the Otto getting started guide.
---

# Appfile

Up to this point, we've developed and deployed a Ruby application
with _zero configuration_. For a getting started experience, it is hard
to beat that.

For real applications, Otto will likely need some minimal configuration.
The file Otto uses for configuration is the `Appfile`. The Appfile is
meant to describe everything required to develop and deploy an application.

On this page, we'll write a simple Appfile that mimics the behavior that Otto
automatically gave us with zero configuration. Then, in the following
getting started pages, we'll augment the Appfile to add new behavior
to our environment.

## A Complete Appfile

We'll start by showing the complete Appfile, and then we'll go over
each part in the following sections. Save the following file to the
root of the example application in a file named "Appfile":

```
application {
  name = "otto-getting-started"
  type = "ruby"
}

project {
  name = "otto-getting-started"
  infrastructure = "otto-getting-started"
}

infrastructure "otto-getting-started" {
  type = "aws"
  flavor = "simple"
}
```

This is functionally equivalent to the behavior of Otto with zero
configuration for our example application.

The syntax of the Appfile is [documented here](/docs/appfile).

~> **WARNING:** Make sure the name of the infrastructure is not
changed above. If you change this, Otto will "lose" your infrastructure
and you'll have to destroy it manually.

## Appfile Compilation

Changes to an Appfile don't take effect for any Otto commands such
as `otto dev` until it is recompiled. To recompile an Appfile, use
`otto compile`.

This feature is very nice since even if other members of a team
might have edited the Appfile, your environment isn't affected
until you decide to recompile.

With the Appfile above, recompiling won't have any adverse affect,
since we've created an Appfile that is functionally equivalent to what
Otto did with zero configuration.

## Appfile Sections

Now that we've written an Appfile, let's dive into the each part of it.

The **application** section defines properties related to the application.
In our Appfile, we specify the name of the application and the type.
The type of the application is what tells Otto how to build our
development environment, builds, deploys, etc. The full set of
[application types](/docs/apps/index.html) is documented.

In addition to the name and type of the application, this section
can include application dependencies, which we cover in an upcoming
step in the getting started guide.

The **project** section defines the project that this application is part
of. Projects currently aren't used for anything other than organization
within Otto. The only important configuration here is the "infrastructure"
setting which tells Otto what infrastructure to target.

The **infrastructure** section defines the infrastructure targets
that this application can be deployed to. This section can be repeated
multiple times. The name of an infrastructure must be unique _globally_
throughout your usage of Otto. If the name of an infrastructure matches
another Appfile, then Otto will assume you want to deploy to the same
infrastructure.

There is another section that an Appfile can contain which we haven't
used yet: **customization**. The customation sections change the behavior
of Otto. We'll see customizations in use in an upcoming step in the
getting started guide.

## Compile and Verify

Once you've saved the Appfile, run `otto compile` again followed by
`otto status`. The status should still look like the following:

```
==> App Info
    Application:    otto-getting-started (ruby)
    Project:        otto-getting-started
    Infrastructure: aws (simple)
==> Component Status
    Dev environment: CREATED
    Infra:           READY
    Build:           BUILT
    Deploy:          DEPLOYED
```

If any of the components are not ready, then the Appfile may have a typo
in it. Verify you copy and pasted the above properly and try again. Don't
forget to recompile!

## Zero Configuration vs. Appfile

So far, we've only recreated an Appfile that does exactly what Otto does
with zero configuration. What is the point?

For what we've done so far, there is no reason to have an Appfile.
However, Otto with zero configuration only does the bare minimum possible
for every application type, and is heavily opinionated. Writing an
Appfile allows you to be explicit in how you want Otto to behave.

In addition to that, features such as dependencies can't be used at all
without an Appfile. This is what we're going to learn next!

## Next

We've now added an Appfile and re-compiled.

Now that we have an Appfile that mimics our _zero configuration_ setup,
we can change how Otto behaves by [customizing our Appfile](/intro/getting-started/customization.html).
