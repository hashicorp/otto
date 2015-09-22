---
layout: "docs"
page_title: "Dependency Sources - Appfile"
sidebar_current: "docs-appfile-depsources"
description: |-
  Application dependencies and Appfile imports both take a URL as a parameter
  to tell Otto where to load the external Appfile.
---

# Dependency and Import Sources

[Application dependencies](/docs/concepts/deps.html) and
[Appfile imports](/docs/appfile/import.html) both take a URL as a parameter
to mark the source of where to load the external Appfile data. Both of these URLs
are in the same format and this page documents the supported values.

Otto supports the following sources:

  * Local file paths

  * GitHub

  * BitBucket

  * Generic Git, Mercurial repositories

  * HTTP URLs

Each is documented further below. The examples will use either an
import statement or a dependency block. However, both imports and dependencies
support all the values below.

## Local File Paths

The easiest source is the local file path.

In Otto, local file paths are relatively rare, except for testing purposes.
If you are using local file paths, try to use relative paths into a sub-directory
for maximum portability.

An example is shown below:

```
import "./folder" {}
```

File paths are the only source URL that update automatically: Otto
creates a symbolic link to the specified directory. Therfore, any changes
are instantly available (versus at compile-time, like every other URL).

## GitHub

Otto will automatically recognize GitHub URLs and turn them into
the proper Git repository. The syntax is simple:

```
dependency {
	source = "github.com/hashicorp/otto/examples/mongodb"
}
```

GitHub source URLs will require that Git is installed on your system
and that you have the proper access to the repository.

You can use the same parameters to GitHub repositories as you can generic
Git repositories (such as tags or branches). See the documentation for generic
Git repositories for more information.

## BitBucket

Otto will automatically recognize BitBucket URLs and turn them into
the proper Git or Mercurial repository. An example:

```
dependency {
	source = "bitbucket.org/hashicorp/example"
}
```

BitBucket URLs will require that Git or Mercurial is installed on your
system, depending on the source URL.

## Generic Git Repository

Generic Git repositories are also supported. The value of `source` in this
case should be a complete Git-compatible URL. Using Git requires that
Git is installed on your system. Example:

```
dependency {
	source = "git://hashicorp.com/module.git"
}
```

You can also use protocols such as HTTP or SSH, but you'll have to hint
to Otto (using the forced source type syntax documented below) to use
Git:

```
// force https source
dependency {
	source = "git::https://hashicorp.com/repo.git"
}

// force ssh source
dependency {
	source = "git::ssh://git@github.com/owner/repo.git"
}
```

URLs for Git repositories (of any protocol) support the following query
parameters:

  * `ref` - The ref to checkout. This can be a branch, tag, commit, etc.

An example of using these parameters is shown below:

```
dependency {
	source = "git::https://hashicorp.com/repo.git?ref=master"
}
```

## Generic Mercurial Repository

Generic Mercurial repositories are supported. The value of `source` in this
case should be a complete Mercurial-compatible URL. Using Mercurial requires that
Mercurial is installed on your system. Example:

```
dependency {
	source = "hg::http://hashicorp.com/module.hg"
}
```

In the case of above, we used the forced source type syntax documented below.
Mercurial repositories require this.

URLs for Mercurial repositories (of any protocol) support the following query
parameters:

  * `rev` - The rev to checkout. This can be a branch, tag, commit, etc.

## HTTP URLs

Any HTTP endpoint can serve up Otto Appfiles. For HTTP URLs (SSL is
supported, as well), Otto will make a GET request to the given URL.
An additional GET parameter `otto-get=1` will be appended, allowing
you to optionally render the page differently when Otto is requesting it.

Otto then looks for the resulting module URL in the following order.

First, if a header `X-Otto-Get` is present, then it should contain
the source URL of the actual module. This will be used.

If the header isn't present, Otto will look for a `<meta>` tag
with the name of "otto-get". The value will be used as the source
URL.

## Forced Source Type

In a couple places above, we've referenced "forced source type." Forced
source type is a syntax added to URLs that allow you to force a specific
method for download/updating the module. It is used to disambiguate URLs.

For example, the source "http://hashicorp.com/foo.git" could just as
easily be a plain HTTP URL as it might be a Git repository speaking the
HTTP protocol. The forced source type syntax is used to force Otto
one way or the other.

Example:

```
module "consul" {
	source = "git::http://hashicorp.com/foo.git"
}
```

The above will force Otto to get the module using Git, despite it
being an HTTP URL.

If a forced source type isn't specified, Otto will match the exact
protocol if it supports it. It will not try multiple methods. In the case
above, it would've used the HTTP method.

## Double-Slash to Split the Root and Subdirectory

Some Appfiles reference files that are above the directory where the
Appfile is. For example:

```
application {
    dependency { source = "../parent" }
}
```

If you depended on an application with an Appfile that looks like the
above, it would need access to the parent folder. Let's pretend that this
Appfile is at the URL local path "/foo/bar".

If you specify the dependency like below, then Otto will be unable to get
the next level dependency "../parent" because it won't exist since Otto
only copies the root of the path given.

```
application {
    dependency { source = "/foo/bar" }
}
```

A double-slash can be used to separate the root of the dependency
with the directory where the Appfile is. If this is used, then the entire
root is copied, and the Appfile in the sub-directory is loaded.
The example below allows this dependency to work:

```
application {
    dependency { source = "/foo//bar" }
}
```
