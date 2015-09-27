---
layout: "intro"
page_title: "The Successor to Vagrant"
sidebar_current: "vagrant"
description: |-
  TODO
---

# The Successor to Vagrant

Otto is the successor to [Vagrant](https://vagrantup.com).

The creators of Otto are also the creators of Vagrant. After working
on Vagrant for over six years, we've learned a lot and we believe Otto
is a superior tool for development plus so much more. Vagrant will still
fill an important role for some users, but for the majority of developers,
Otto will replace Vagrant over time.

While we believe Otto is a successor to Vagrant, Vagrant development will
continue for years to come. Otto is built on Vagrant, so improvements to
 Vagrant will benefit Otto users as well.

### Improvements Over Vagrant

**Application level vs. machine level configuration**. The Vagrantfile describes
the _machine_ for development. The Appfile in Otto describes the _application_.
The first thing you configure in a Vagrantfile is the box, which is the
machine image Vagrant uses. The first thing you configure in an Appfile is
the application type ("php", "ruby", etc.). Using this knowledge, Otto
automatically builds a development environment for that application type.
This environment can be further customized through more advanced Appfile
directives, but out of the box Otto usually just works.

**Dependencies as first-class feature**. In a modern microservice oriented
world, a Vagrantfile is a difficult mechanism to describe a proper development
environment. Multi-VM is too heavy for microservices, and you would have to
manually install the full dependency tree (including dependencies of dependencies)
to make a complete development environment. Dependencies are a first-class
feature in Otto. You only need to specify direct dependencies, and Otto
automatically installs and configures all dependencies, including
dependencies of dependencies. Even with hundreds of microservices, the
Appfile to configure Otto is simple, readable, and intuitive.

**Deployment**. Otto can deploy your application. Users of Vagrant for years
have wanted a way to deploy their Vagrant environments to production.
Unfortunately, the Vagrantfile doesn't contain enough information to
build a proper production environment with industry best practices. An
Appfile is made to encode this knowledge, and deployment is a single
command away.

**Performance**. We took what we learned from Vagrant, and optimized the most common
operations. For example, `vagrant status` is an operation that takes
several seconds. `otto status` finishes in milliseconds.

### The Future of Vagrant

Vagrant is a mature, battle-hardened piece of technology. It has been
in use by millions of users for many years. We don't want to reinvent the
wheel, so we've taken the best parts of Vagrant and used them within Otto to
manage development environments automatically for the user.

Otto builds on top of Vagrant to make operations such as SSH take
milliseconds, automatically assign addresses,
and more. In addition to all the development environment features, Otto does a
lot more to enable deployment.

In addition to Otto using Vagrant, Vagrant will continue to be the best
tool for managing highly customized virtualized or containerized environments.
Vagrant is a great way to test configuration management, obscure operating
systems, etc. These use cases will not go away.

Due to the above, we're committed to continuing to improve Vagrant and releasing
new versions of Vagrant for years to come. But for the everyday developer,
Otto should replace Vagrant over time.
