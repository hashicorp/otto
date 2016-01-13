---
layout: "docs"
page_title: "Appfile Compilation"
sidebar_current: "docs-concepts-compile"
description: |-
  The compilation step in Otto takes the [Appfile](/docs/concepts/appfile.html),
  validates it, fetches any dependencies, and generates dozens of output
  files that are used by subsequent Otto operations.
---

# Compilation

The first command run on any project with Otto is `otto compile`.
The idea of "compilation" is core to the function of Otto, and is
also the secret behind a lot of the magic of Otto.

The compilation step in Otto takes the [Appfile](/docs/concepts/appfile.html),
validates it, fetches any dependencies, and generates dozens of output
files that are used by subsequent Otto operations.

Because compilation is relatively rare, it gives Otto an opportunity to
perform more expensive operations, such as resolving dependencies,
fetching Appfile imports, semantic validation, verifying state,
generating upgrade paths for deploys, etc. Despite this, compilation in
Otto is still incredibly fast, taking milliseconds, except in the case
of network operations (dependencies, imports).

Once compilation is complete, the remaining Otto operations such as
`otto dev` use the output of compilation and can assume validation passed,
making their load times much faster.

## Output

The output of Otto compilation goes into the ".otto" directory relative
to the location of the Appfile. This directory is local to every Otto
compilation and _should not_ be committed to version control.

In addition to the output in ".otto", the first compilation will generate
a file called ".ottoid" alongside the Appfile. This file contains a unique
ID that is used to identify this application. It _should be_ committed to
version control. If you're starting a new project, you should run
`otto compile` just to generate and commit the ".ottoid" file early, even
if the Appfile will change dramatically.

More details on the ".ottoid" are in its own dedicated section below.

### .otto

Within the ".otto" directory, Otto outputs dozens of files that are used
for the remainder of Otto execution (until another compile). An example
of some of these files are shown below:

```
$ tree .otto
.otto
├── appfile
│   ├── Appfile.compiled
│   └── version
├── compiled
│   ├── app
│   │   ├── build
│   │   │   ├── build-ruby.sh
│   │   │   └── template.json
│   │   ├── deploy
│   │   │   └── main.tf
│   │   ├── dev
│   │   │   └── Vagrantfile
...

24 directories, 41 files
```

This folder contains [Vagrant](https://www.vagrantup.com) files,
[Packer](https://www.packer.io) templates, [Terraform](https://www.terraform.io)
configurations, and more. Otto uses these tools under the covers to
manage many aspects of development and deployment. If you're familiar with
this tooling, you can inspect the files to see exactly what Otto will do.

In addition to these files, you can see the "Appfile.compiled" file, which
is the Appfile structure used by the other `otto` commands.

The ".otto" folder _should not_ be committed to version control. It is
local to the system that ran `otto compile`.

### .ottoid

The ".ottoid" file is generated one time per application and contains
a unique identifier for that application. This UUID is used to track the
application deployment in the [directory](/docs/concepts/directory.html).
It is a critical piece of information that Otto uses to ensure an application
is only deployed once (unless explicitly requested otherwise), to maintain
state history, etc.

This file should be committed to version control.
