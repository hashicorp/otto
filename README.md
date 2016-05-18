Otto
=========

- Website: https://www.ottoproject.io
- IRC: `#otto-tool` on Freenode
- Mailing list: [Google Groups](https://groups.google.com/group/otto-tool)

![Otto](https://cloud.githubusercontent.com/assets/37534/10147078/d400509e-65e0-11e5-9d66-c419914cbcf4.png)

Otto knows how to develop and deploy any application on any cloud platform,
all controlled with a single consistent workflow to maximize the productivity
of you and your team.

For more information, see the
[introduction section](https://www.ottoproject.io/intro)
of the Otto website.

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
  build development environments ([Vagrant](https://www.vagrantup.com)),
  launch servers ([Terraform](https://www.terraform.io)), configure
  services ([Consul](https://www.consul.io)), and more. Otto builds on
  tools that powers the world's largest websites.
  Otto automatically installs and manages all of this tooling, so you don't
  have to.

## Getting Started & Documentation

All documentation is available on the [Otto website](https://www.ottoproject.io).

## Developing Otto

If you wish to work on Otto itself or any of its built-in systems,
you'll first need [Go](https://www.golang.org) installed on your
machine (version 1.4+ is *required*).

For local dev first make sure Go is properly installed, including setting up a
[GOPATH](https://golang.org/doc/code.html#GOPATH).

Next, clone this repository into `$GOPATH/src/github.com/hashicorp/otto`.
Then use `make` to get the dependencies and run the tests.
If this exits with exit status 0,
then everything is working!

```sh
$ make updatedeps
...
$ make
...
```

To compile a development version of Otto, run `make dev`. This will put the
Otto binary in the `bin` and `$GOPATH/bin` folders:

```sh
$ make dev
...
$ bin/otto
...
```

If you're developing a specific package, you can run tests for just that
package by specifying the `TEST` variable. For example below, only
`otto` package tests will be run.

```sh
$ make test TEST=./otto
...
```

### Acceptance Tests

Otto also has a comprehensive
[acceptance test](https://en.wikipedia.org/wiki/Acceptance_testing)
suite covering most of the major features of the built-in app types.

If you're working on a feature of an app type and want to verify it is
functioning (as well as hasn't broken anything else), we recommend running
the acceptance tests. Note that we _do not require_ that you run or write
acceptance tests to have a PR merged. The acceptance tests are here for your
convenience, and we'll add them if they're missing.

**Note:** Acceptance tests are slow. They spin up real resources and take
time to run. We recommend filtering as documented below to speed up your
workflow.

To run the acceptance tests, invoke `make testacc`:

```sh
$ make testacc TEST=./builtin/app/go TESTARGS='-run=dev'
...
```

The `TEST` variable is required, and you should specify the folder where the
app type is. The `TESTARGS` variable is recommended to filter down to a specific
test to run, since running all of them at once can take a very long time.
