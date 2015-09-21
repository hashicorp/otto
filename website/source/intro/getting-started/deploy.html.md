---
layout: "intro"
page_title: "Deploy"
sidebar_current: "gettingstarted-deploy"
description: |-
  Deploy an application with Otto in the Otto getting started guide.
---

# Deploy

In the previous step, we got a development up and running with no
configuration and a couple simple commands. Now, let's deploy that application.

We'll deploy the application to [AWS](http://aws.amazon.com) for
the getting started guide since it is popular and generally well understood, but
Otto can deploy to many different infrastructure providers.

If you don't have an AWS account, [create one now](http://aws.amazon.com/free/).
For the getting started guide, we'll only be using resources which qualify
under the AWS [free-tier](http://aws.amazon.com/free/), meaning it will be free.
If you already have an AWS account, you may be charged some amount of money,
but it shouldn't be more than a few dollars at most.

~> **Warning!** If you're not using an account that qualifies under the AWS
free-tier, you may be charged to run these examples. The most you should be
charged should only be a few dollars, but we're not responsible for any
charges that may incur.

## Infrastructure

Before we can deploy, we need infrastructure to deploy to. In addition
to deploying the application itself, Otto also manages the infrastructure.
We'll see how Otto can share infrastructure between multiple applications
in a future step.

"Infrastructure" within Otto refers to the underlying resources necessary
to run applications. For AWS, this means a VPC, proper routing tables,
a subnet, etc. If you're not familiar with this, that's okay! The point
of Otto is to know this, so you don't have to.

To build the infrastructure, run `otto infra`.

This step is likely going to ask you for permission to install
[Terraform](https://terraform.io). Otto uses Terraform under the covers
to build and manage the infrastructure. If you say yes, Otto will install
this for you. If you already have Terraform installed, you won't be
asked.

After installing Terraform, Otto will ask you for your AWS access
credentials. These are available from
[this page](https://console.aws.amazon.com/iam/home?#security_credential).

Next, Otto will go forward and build you an infrastructure. This will
take a few minutes. During this time, you can attempt to read the output,
which will be fairly verbose. You'll see Otto creating a VPC, a subnet,
some routing tables, and even launching a micro EC2 instance.

TODO

## Build

## Deploy

## Next

In this step, you learned how easy it is to use Otto. You experienced
Otto compilation for the first time and also saw how Otto works with
_zero configuration_. You hopefully are beginning to sense the power of
Otto, even if we've only covered development so far.

Next, we'll [deploy this application](/intro/getting-started/deploy.html)
to a real cloud environment.
