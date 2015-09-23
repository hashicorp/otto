---
layout: "intro"
page_title: "Infrastructure"
sidebar_current: "gettingstarted-infra"
description: |-
  Build the infrastructure to support your applications with Otto.
---

# Infrastructure

In the previous step, we got a development environment up and running with no
configuration and a couple simple commands. Now, let's deploy that application.

To deploy an application, Otto has three steps: start an infrastructure,
build the application, and launch the application. We'll go over each of
these steps on as a separate page in the getting started guide.

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

The first step before an application is deployed is to build an
infrastructure to deploy to. Within Otto, "infrastructure" refers to
a target cloud platform and the minimum resources necessary to run
applications.

Using AWS as an example, Otto creates a VPC, subnets, proper routing tables,
an internet gateway, and more. If you're not familiar with this, that's okay!
That is the point of Otto: to codify infrastructure best practices so
you don't have to know them on day one.

For each cloud platform, Otto knows multiple "flavors" of that
infrastructure. These flavors target different goals and users. In our AWS
example, a couple available flavors are "simple" and "vpc-public-private."
The "simple" flavor uses a minimal number of resources, sacrificing
scalability and fault tolerance for simplicity and cost. But
"vpc-public-private" automatically sets up private networks, bastion
hosts, NAT instances, etc. and is more scalable for a real long term
infrastructure.

By default, Otto defaults to the "simple" flavor of AWS, and that is the
flavor we'll use in this getting started guide. You can learn more about
infrastructures and flavors in the [documentation](/docs).

## Launch

To launch the infrastructure, run `otto infra`.

This step is likely going to ask you for permission to install
[Terraform](https://terraform.io). Otto uses Terraform under the covers
to build and manage the infrastructure. If you say yes, Otto will install
this for you. If you already have the latest version of Terraform installed,
you won't be asked.

After installing Terraform, Otto will ask you for your AWS access
credentials. These are available from
[this page](https://console.aws.amazon.com/iam/home?#security_credential).

Next, Otto will go forward and build you an infrastructure. This will
take a few minutes. During this time, you can attempt to read the output,
which will be fairly verbose. You'll see Otto creating a lot of cloud
resources.

TODO

The types and number of resources created are determined by the infrastructure
flavor, as mentioned above. For flavors such as "vpc-public-private," the
initial `infra` can take several minutes.

**Congratulations!** You've just launched an infrastructure with the
minimum necessary number of components to deploy one or many applications.

If you've never used AWS before, or you don't consider yourself an operator,
Otto just used industry best practices for simple applications to launch
and configure a robust infrastructure for you.

## Status

You can see the status of your infrastructure at any point by running
`otto status`:

```
$ otto status
==> App Info
    Application:    otto-getting-started (ruby)
    Project:        otto-getting-started
    Infrastructure: aws (simple)
==> Component Status
    Dev environment: CREATED
    Infra:           READY
    Build:           NOT BUILT
    Deploy:          NOT DEPLOYED
```

You can see that the "Infra" is "READY." This means that step is complete.

Note that even if an infrastructure is ready, you can always run `otto infra`
multiple times. Otto will only create infrastructure resources that don't exist.

## Next

In this step, you learned about how Otto manages infrastructures
and built a simple infrastructure that we can deploy applications to.

The goal of Otto is build real, production quality infrastructures
for hobbyists or professionals. Otto codifies industry best practices
and brings them to you in a single command.

Next, we'll [build this application](/intro/getting-started/build.html)
to prepare it for being deployed.
