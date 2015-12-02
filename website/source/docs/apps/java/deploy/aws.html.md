---
layout: "app_java"
page_title: "AWS - Build & Deploy - Java App Type"
sidebar_current: "docs-java-deploy-aws"
description: |-
  This page documents how the Java application builds and deploys on
  AWS infrastructure.
---

# Build & Deploy: AWS

This page documents how the Java application builds and deploys on
[AWS infrastructure](/docs/infra/aws).

The sections below are split into a section of commonalities between
the different infrastructure flavors, and then specific sections for
each infrastructure flavor.

Please see the [customizations](/docs/apps/java/customization.html)
page for a list of behavior that can be changed.

## Common

For all AWS flavors:

  * The build output is an AMI.

  * The deploy process launches at least one EC2 instance. The size
    of this instance varies by infrastructure flavor.

  * A custom security group just for that application is created. The
    exact rules of the security group vary by infrastructure flavor.

## Flavor: "simple"

For the "simple" AWS flavor:

  * A single `t2.micro` EC2 instance is launched to serve the application.

  * The security group allows only SSH access from the outside world.

## Flavor: "vpc-public-private"

For the "vpc-public-private" flavor:

  * One `t2.small` EC2 instance is launched to serve the application.

  * The EC2 instances are launched into the private subnet and can only
    be accessed for SSH via the bastion host.
