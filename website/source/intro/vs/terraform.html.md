---
layout: "intro"
page_title: "Otto vs. Terraform"
sidebar_current: "vs-other-terraform"
description: |-
  Comparison between Otto and Terraform.
---

# Otto vs. Terraform

[Terraform](https://terraform.io) is a tool for launching infrastructure.
It is also written by the same people who created Otto. Terraform is often
used as a deployment tool.

Terraform is a good solution for deployment. It can manage complex
infrastructures and safely orchestrates complicated operations. Terraform
isn't a development tool and makes no attempt to provide development-oriented
features.

Otto uses Terraform internally to manage infrastructure and power the deploys.
However, Otto automatically generates the Terraform configuration based
on the [Appfile](/docs/concepts/appfile.html) input, which is much simpler
than learning Terraform itself.

Terraform is a lower level tool that requires the user to have a deep
understanding of the cloud platform they want to deploy to, the resources
they'll need, and the configuration of those resources.

Otto requires little to no configuration and automatically can deploy
to multiple cloud platforms with complex automatically-generated
Terraform configurations using industry best practices.

Otto has a complete development experience and guides the developer
through the build process as well. Terraform expects deployable artifacts
(whether it is an AMI or a container) to be pre-built by some other
mechanism.
