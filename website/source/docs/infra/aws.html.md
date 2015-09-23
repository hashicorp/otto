---
layout: "docs"
page_title: "Infra Type - AWS"
sidebar_current: "docs-infra-aws"
description: |-
  The AWS infrastructure type allows Otto to deploy applications to
  Amazon Web Services.
---

# Infrastructure Type: AWS

The AWS infrastructure type allows Otto to deploy applications to [Amazon
Web Services](https://aws.amazon.com/).

## Credentials

Otto needs AWS API credentials in order to be able to manage resources on AWS
for you. Otto will ask you for these credentials during its first run it does
not not have any. Otto [stores these credentials in an encrypted cache
file](/docs/infra/index.html#credentials) for subsequent runs.

You can avoid being prompted for these credentials by providing them via
environment variables instead. Each field lists the environment variable that
can be used to set it.

Here is the list of credentials Otto needs for the AWS infrastructure type:

 * __AWS Access Key__ - the identifier portion of the standard AWS API key pair
   (Env var: `AWS_ACCESS_KEY_ID`)
 * __AWS Secret Key__ - the secret portion of the standard AWS API key pair
   (Env var: `AWS_SECRET_ACCESS_KEY`)
 * __SSH Public Key Path__ - a path to an SSH public key that Otto will grant
   access to any instances it creates in this infrastructure (Env var:
   `AWS_SSH_PUBLIC_KEY_PATH`)

## Flavors

Otto currently supports two infrastructure "flavors", both of which
involve AWS, but represent different styles of deployment.

### Flavor: "simple"

This is the default infrastructure that Otto uses. It's
meant for demonstration purposes, as it skips over certain security and
redundancy considerations in favor of getting users up and running quickly and
cheaply.

It consists of the following resources:

 * A VPC with one public subnet, configured to assign a public IP address and
   DNS hostname to any instance launched into it.
 * A key pair using the SSH public key you provide, which will be added to any
   instances launched into your infrastructure.
 * By default, it also includes a simple flavor of the ["Consul"
   foundation](/docs/foundations/consul.html).

### Flavor: "vpc-public-private"

This flavor sets up a more production-like infrastructure in AWS. It involves
more resources, which means it is more expensive to run and it takes longer to
create, but it yields more robust and performant application deployments.

It has the following resources:

 * A VPC with one public and one private subnet.
 * A NAT instance launched into the public subnet, configured to provide egress
   network connectivity from the private subnet.
 * A Bastion host launched into the public subnet, so private instances can be
   accessed via SSH.
 * By default, it also includes a prouduction cluster from the ["Consul"
   foundation](/docs/foundations/consul.html).
