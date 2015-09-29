---
layout: "intro"
page_title: "Build"
sidebar_current: "gettingstarted-build"
description: |-
  Build the application to prepare it for deployment in this getting started guide.
---

# Build

Now that we have an infrastructure, the next step is to prepare the
application for deployment by building it. The build step turns your
application into a deployable unit. This might be an Amazon Machine Image,
a Docker container, etc.

For our simple example, Otto will create an AMI to launch.

To build the application, run `otto build`.

This step is likely going to ask you for permission to install
[Packer](https://packer.io). Otto uses Packer under the covers
to build deployable artifacts. If you say yes, Otto will install
this for you. If you already have the latest version of Packer
installed, you won't be asked.

After installing Packer, Otto will once again ask you for your credentials.
This time, however, just enter the password you used to encrypt them
when building the infrastructure, and Otto will be able to read your
saved credentials.

Next, Otto will go forward and build an AMI that can be launched.
This will take a few minutes. Within the AMI, Otto installs Ruby, configures Passenger (a Ruby
application server), and more. Again, depending on your application type,
Otto will install different things here.

```
$ otto build
...

==> otto: Creating the AMI: otto-getting-started 1442990619
    otto: AMI: ami-4b19662e
==> otto: Waiting for AMI to become ready...
==> otto: Terminating the source AWS instance...
==> otto: Cleaning up any extra volumes...
==> otto: No volumes to clean up, skipping
==> otto: Deleting temporary security group...
==> otto: Deleting temporary keypair...
Build 'otto' finished.

==> Builds finished. The artifacts of successful builds are:
--> otto: AMIs were created:

us-east-1: ami-4b19662e
==> Storing build data in directory...
==> Build success!
    The build was completed successfully and stored within
    the directory service, meaning other members of your team
    don't need to rebuild this same version and can deploy it
    immediately.
```

**Congratulations!** You've just built an AMI that can be launched!

Imagine you're a developer that doesn't know how to properly configure
a server to serve an application you've built. In one command, Otto has
done this for you, for your specific type of application, and has
pre-configured the server with safe best practices.

## Status

You can see the status of your build at any point by running
`otto status`. Since the prior step, you can see that the build status
is now "BUILD READY":

```
$ otto status
==> App Info
    Application:    otto-getting-started (ruby)
    Project:        otto-getting-started
    Infrastructure: aws (simple)
==> Component Status
    Dev environment: CREATED
    Infra:           READY
    Build:           BUILD READY
    Deploy:          NOT DEPLOYED
```

## Build Speed

The granularity of a build is currently at the machine level and
creates an AMI (for AWS). In the very near future, the build step will
build a container that is deployed using [Nomad](https://nomadproject.io).

When this change comes, the build and deploy steps should become much
faster, creating a much tighter feedback loop to deploying applications.

This again shows the strength of Otto: as tooling and technology improve,
Otto can adopt these new best practices and manage them for you automatically.
As a developer, you can focus on your application, and Otto does the rest.

## Next

In this step, we built an AMI for our application. Along with the previous
step, we now have an infrastructure and an AMI. We're now ready to
[deploy the application](/intro/getting-started/deploy.html)!
