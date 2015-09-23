---
layout: "intro"
page_title: "Introduction"
sidebar_current: "what"
description: |-
  TODO
---

# What is Otto?

Otto knows how to develop and deploy any application.

Otto automatically builds development environments tailored to your specific
application. When you're ready to deploy, Otto can create the infrastructure
and launch the application along with all of its dependencies. It is a single
tool to develop and deploy any application.

Deploying a relatively simple application today requires an
immense amount of domain knowledge: launching cloud servers, configuring
those servers, securing the servers, deploying new versions of your
application, etc. The solution to each of these problems usually requires
another tool you must learn and master if you intend to solve it properly.
And new trends such as microservices make all of this even more complicated.

Due to this complexity, many developers either completely ignore best practices
or turn to expensive completely managed solutions that often have limited
scalability and flexiblity.

Otto gives developers simple, single commands to develop and deploy
their application while using industry-standard tooling and best practices
under the covers.

Otto has built-in knowledge of many application types and infrastructure
architectures, and automatically generates the configuration necessary to
manage the entire lifecycle of that application. If you're a PHP developer,
for example, Otto knows how to install and configure Apache and PHP. If you
depend on a database, Otto can setup that database and configure firewalls
so that it is only accessible from your application. As time goes on, Otto
will get smarter, and will generate more and more advanced deployments,
all behind a single, simple deploy command from the developer.

TODO: This is not good.

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
  them. Whether your application is a Docker container or a legacy
  monolithic PHP application, Otto can deploy your application to any
  environment.

* **Battle-hardened tooling**: Otto uses battle-hardened tooling to
  build development environments ([Vagrant](https://vagrantup.com)),
  launch servers ([Terraform](https://terraform.io)), configure
  services ([Consul](https://consul.io)), and more. Each of the tools
  Otto builds on top of is in use by some of the world's largest websites.
  Otto automatically installs and manages all of this tooling, so you don't
  have to.

## Next Steps

Continue onwards with the [getting started guide](/intro/getting-started/install.html)
to see how easy Otto makes it to develop and deploy a real application.

