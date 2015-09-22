---
layout: "intro"
page_title: "Customization"
sidebar_current: "gettingstarted-customization"
description: |-
  Customizations are a way to change the way Otto behaves by default with an application type.
---

# Customization

Otto has a lot of built-in knowledge. It was able to detect that our
example application is a Ruby application, it automatically installed
Ruby in our development environment, it automatically installed an
application server for deployment, etc. Otto is an opinionated tool, making
it incredibly powerful and easy to use out of the box.

However, real applications often diverge slightly in various ways.
_Customizations_ are a way to change the behavior of Otto.

## Customization Configuration

We'll make a very basic customization with our example
application and change the Ruby version. In the Appfile we created,
add a new section:

```
customization "ruby" {
    ruby_version = "2.1"
}
```

There can be multiple "customization" blocks in an Appfile. The
name of the block is the component to customize. These are mostly
depended on the application type and are documented as part of the
[application types](/docs/apps).

For Ruby, there is a "ruby_version" configuration in the "ruby"
customization type. This sets the Ruby version to install, and currently
defaults to "2.2".

Let's pretend our application doesn't work with Ruby 2.2, and request
that Otto install Ruby 2.1. This is what we've done in the customization
above.

## Applying the Change

Changes to the Appfile won't take effect until we recompile, so
run an `otto compile` to recompile now. After compiling, let's rebuild
our development environment with the latest Ruby version.

First, destroy the old environment with `otto dev destroy`. This
should only take a few seconds. Otto doesn't currently allow in-place
updates of changes to the development environment.

Once it is destroyed, run `otto dev` again. In a few minutes, you
should see output from Otto about the version of Ruby it is installing,
and it should be "2.1".

You can verify this by going into the development environment:

```
$ otto dev ssh
> ruby -v
TODO
```

It'll be left as an exercise to you to deploy the new Ruby version,
if you want. Note that you'll have to run an `otto build` again since the
Ruby version is built-in to the AMI that was built.

## Next

We showed a basic example of how a customization can alter the behavior
of Otto. Customizations are incredibly powerful, but they're still
high level enough for users of Otto to be productive without having
to know difficult details.

Next, we'll learn how to develop and deploy applications that
have [dependencies](/intro/getting-started/deps.html).
