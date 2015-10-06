---
layout: "docs"
page_title: "Infra Type - DigitalOcean"
sidebar_current: "docs-infra-digitalocean"
description: |-
  The DigitalOcean infrastructure type allows Otto to deploy applications to
  DigitalOcean Cloud Hosting.
---

# Infrastructure Type: DigitalOcean

The DigitalOcean infrastructure type is currently only available for Ruby Apps.
It allows Otto to deploy applications to [DigitalOcean Cloud Hosting](https://www.digitalocean.com/).

## Credentials

Otto needs DigitalOcean API credentials in order to be able to manage resources on DigitalOcean
for you. Otto will ask you for these credentials during its first run it does
not not have any. Otto [stores these credentials in an encrypted cache
file](/docs/infra/index.html#credentials) for subsequent runs.

You can avoid being prompted for these credentials by providing them via
environment variables instead. Each field lists the environment variable that
can be used to set it.

Here is the list of credentials Otto needs for the DigitalOcean infrastructure type:

 * __DigitalOcean token__ - the token required by the DigitalOcean API.
   (Env var: `DIGITALOCEAN_TOKEN`)
 * __SSH Public Key Path__ - a path to an SSH public key that Otto will grant
   access to any instances it creates in this infrastructure (Env var:
   `TF_DO_SSH_PUBLIC_KEY_PATH`)

## Flavors

Otto currently supports one infrastructure "flavors" for DigitalOcean.

### Flavor: "simple"

This is the default infrastructure that Otto uses. It's
meant for demonstration purposes, as it skips over certain security and
redundancy considerations in favor of getting users up and running quickly and
cheaply.

It consists of the following resources:

 * A Simple droplet configured to assign a public IP address.
 * A key pair using the SSH public key you provide, which will be added to any
   instances launched into your infrastructure.
 * By default, it also includes a simple flavor of the ["Consul"
   foundation](/docs/foundations/consul.html).
