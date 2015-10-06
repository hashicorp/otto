---
layout: "intro"
page_title: "Install Otto"
sidebar_current: "gettingstarted-install"
description: |-
  The first step to using Otto is to get it installed.
---

# Install Otto

To get started with Otto, Otto must first be installed on your machine.
Otto is distributed as a [binary package](/downloads.html) for all
supported platforms and architectures.

## Installing Otto

To install Otto, find the [appropriate package](/downloads.html) for
your system and download it. Otto is packaged as a zip archive.

After downloading Otto, unzip the package. Otto runs as a single binary
named `otto`. Any other files in the package can be safely removed and
Otto will still function.

The final step is to make sure that `otto` is available on the PATH.
See [this page](http://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux)
for instructions on setting the PATH on Linux and Mac.
[This page](http://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows)
contains instructions for setting the PATH on Windows. For each of those
pages, refer to the question and accepted answer for best practices. Create
any files that they refer to that don't exist.

Also, if you have not already done so, please install [VirtualBox](http://virtualbox.org).
The development environment for Otto will use this. A future version of
Otto will do this for you automatically.

## Verifying the Installation

After installing Otto, verify the installation worked by opening a new
terminal session and checking that `otto` is available. By executing
`otto`, you should see help output similar to the following:

```
$ otto
usage: otto [--version] [--help] <command> [<args>]

Available commands are:
    build      Build the deployable artifact for the app
    compile    Prepares your project for being run.
    deploy     Deploy the application
    dev        Start and manage a development environment
    infra      Builds the infrastructure for the Appfile
    status     Status of the stages of this application
    version    Prints the Otto version
```

If you get an error that Otto could not be found, then your PATH environment
variable was not setup properly. Please go back and ensure that your PATH
variable contains the directory where Otto was installed.

Otherwise, Otto is installed and ready to go!
