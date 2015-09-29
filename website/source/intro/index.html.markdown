---
layout: "intro"
page_title: "Introduction"
sidebar_current: "what"
description: |-
  TODO
---

# What is Otto?

Otto knows how to develop and deploy any application on any cloud
platform, all controlled with a single consistent workflow to maximize
the productivity of you and your team.

## A Developer's Dream

Otto automatically builds a development environment tailored specifically
for your application, with zero or minimal configuration.

Otto isolates all your applications into their own local virtualized
development environments. It detects the type of application you're developing and
configures that development environment for you. For example, if you're
developing a PHP application, Otto will automatically install PHP and
other related tools for you.

Otto supports application dependencies as a first class feature. This makes
developing microservices a breeze. Otto automatically downloads, installs,
configures, and starts dependencies in your development environments for you.

Developers that are new to a team or are context switching to a new
project can get up and running in a single command with minimal friction:
`otto dev`.

## Best-in-Class Infrastructure, Automatically

Otto automatically builds an infrastructure and deploys your application
using industry standard tooling and best practices, so you don't have to.

There are hundreds of how-to guides to deploy applications. Unfortunately,
all of these how-to guides are usually copying and pasting outdated information
onto a server to barely get your application running.

Deploying an application properly with industry best practices in security,
scalability, monitoring, and more requires an immense amount of domain
knowledge. The solution to each of these problems usually requires
mastering a new tool.

Due to this complexity, many developers completely ignore best practices
and stick to the simple how-to guides.

Otto solves all these problems and automatically manages all of the
various software solutions for you to have a best-in-class
infrastructure. You only need to learn Otto, and Otto does the rest.

## Key Features

The key features of Otto are:

* **Automatic development environments**: Otto detects your application
  type and builds a development environment tailored specifically for that
  application, with zero or minimal configuration. If your application depends
  on other services (such as a database), it'll automatically configure and
  start those services in your development environment for you.

* **Built for Microservices**: Otto understands dependencies and versioning
  and can automatically deploy and configure an application and all
  of its dependencies for any environment. An application only needs to
  tell Otto its immediate dependencies; dependencies of dependencies are
  automatically detected and configured.

* **Deployment**: Otto knows how to deploy applications as well develop
  them. Whether your application is a modern microservice, a legacy
  monolith, or something in between, Otto can deploy your application to any
  environment.

* **Docker**: Otto can use Docker to download and start dependencies
  for development to simplify microservices. Applications can be containerized
  automatically to make deployments easier without changing the developer
  workflow.

* **Production-hardened tooling**: Otto uses production-hardened tooling to
  build development environments ([Vagrant](https://vagrantup.com)),
  launch servers ([Terraform](https://terraform.io)), configure
  services ([Consul](https://consul.io)), and more. Otto builds on
  tools that power the world's largest websites.
  Otto automatically installs and manages all of this tooling, so you don't
  have to.

## Next Steps

Continue onwards with the [getting started guide](/intro/getting-started/install.html)
to see how easy Otto makes it to develop and deploy a real application.
