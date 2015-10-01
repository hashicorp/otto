---
layout: "docs"
page_title: "Syntax - Appfile"
sidebar_current: "docs-appfile-syntax"
description: |-
  The syntax of Appfiles is HCL.
---

# Syntax

The syntax of Appfiles is [HCL](https://github.com/hashicorp/hcl).

## HCL

Here is an example of an Appfile:

```
# Application stanza
application {
    name = "otto"
    type = "go"

    dependency { source = "github.com/hashicorp/otto/examples/mongodb" }
}

/*
Let's disable this for now
customization "go" {
    go_version = "1.4.2"
}
*/
```

Basic bullet point reference:

  * Single line comments start with `#`

  * Multi-line comments are wrapped with `/*` and `*/`

  * Values are assigned with the syntax of `key = value` (whitespace
    doesn't matter). The value can be any primitive: a string,
    number, or boolean.

  * Strings are in double-quotes.

  * Numbers are assumed to be base 10. If you prefix a number with
    `0x`, it is treated as a hexadecimal number.

  * Boolean values: `true`, `false`.

  * Lists of primitive types can be made by wrapping it in `[]`.
    Example: `["foo", "bar", 42]`.

  * Maps can be made with the `{}` syntax:
	`{ foo = "bar" }`.

