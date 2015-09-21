---
layout: "docs"
page_title: "Install Otto"
sidebar_current: "docs-install"
description: |-
  Learn how to install Otto.
---

# Install Otto

Installing Otto is simple.

Otto is packaged as a zip and pre-built for various platforms.
Go to the [download page](/downloads.html) and download the appropriate package
for your system.

Once the zip is downloaded, unzip it into any directory. The
`otto` binary inside is all that is necessary to run Otto (or
`otto.exe` for Windows). Any additional files, if any, aren't
required to run Otto.

The final step is to make sure that `otto` is available on the PATH.
See [this page](http://stackoverflow.com/questions/14637979/how-to-permanently-set-path-on-linux)
for instructions on setting the PATH on Linux and Mac.
[This page](http://stackoverflow.com/questions/1618280/where-can-i-set-path-to-make-exe-on-windows)
contains instructions for setting the PATH on Windows.

## Verifying the Installation

To verify Otto is properly installed, execute the `otto` binary on
your system. You should see help output.
If you get an error that Otto could not be found, then your PATH environment
variable was not setup properly. Please go back and ensure that your PATH
variable contains the directory where Otto was installed.

