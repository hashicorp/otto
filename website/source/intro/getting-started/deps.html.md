---
layout: "intro"
page_title: "App Dependencies"
sidebar_current: "gettingstarted-deps"
description: |-
  In this step, we introduce dependencies to our application and show Otto manages those dependencies.
---

# Dependencies

A modern application usually has many dependencies: other applications they
depend on to function. These can be things such as databases or other applications
that they need to communicate with such as in a
[microservice architecture](http://martinfowler.com/articles/microservices.html).

Dependencies introduce many complexities to the development and deployment
process: they need to be installed, configured, and started in development.
This puts a burden on the developer. And in production, they need to be
separately deployed, scaled, etc.

Otto supports microservices and dependencies as first class features.

## Declaring Dependencies

Declaring dependencies is simple.

Modify the Appfile `application` section to look like the following:

```
application {
  name = "otto-getting-started"
  type = "ruby"

  dependency {
    source = "github.com/hashicorp/otto/examples/mongodb"
  }
}
```

The `dependency` block specifies a dependency of the application. Multiple
of these can be specified. Within the block, the `source` is a URL
pointer to the application. Otto expects an Appfile to exist there to describe
how to develop and deploy that dependency.

In the above example, we've added MongoDB as a dependency to our application.

An important feature of Otto is that you only need to specify the immediate
dependencies; dependencies of dependencies do not need to be specified.
Otto automatically fetches the dependency and inspects it for more
dependencies, and continues this process.

## Develop

Once you've declared a dependency, run `otto dev` again. If you have
a development environment already created, run `otto dev destroy` first.

Once this process completes, you should notice some output about
setting up "mongodb" in various places. Otto automatically installs
and configures dependencies within the development environment. The
dependencies are also automatically started.

Dependencies are exposed via DNS by their name. If you SSH into the
development environment and look up `mongodb.service.consul`, then
you'll find the IP address for MongoDB (which should be local). In production,
this will likely point to another machine.

Notice how with only a few lines, and zero instructions for how to setup
the dependency, Otto knows how to automatically install, configure, and
start it in a development environment. As the number of dependencies
increases, this becomes more and more valuable.

## Next

Otto makes developing and deploying applications with dependencies
incredibly simple, and a huge improvement over other tools.

Modern application development is moving further and further towards
microservices and many separate applications communicating together.
By representing dependencies as a first class citizen, Otto eliminates
much of the friction associated with this model.

Next, we'll learn how to
[teardown](/intro/getting-started/teardown.html)
all the resources we've created with Otto.

