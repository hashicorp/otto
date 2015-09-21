---
layout: "docs"
page_title: "Uninstalling Otto"
sidebar_current: "docs-install-uninstall"
description: |-
  How to uninstall Otto.
---

# Uninstalling Otto

To uninstall Otto, remove the `otto` binary from your system.

If you want to completely uninstall Otto, including any supporting files,
then also remove `~/.otto.d`. If you have any active infrastructures, deployed
applications, etc. then this step will cause Otto to "lose" access, so
destroy any applications prior to this.

The `.ottoid` and `.otto` folders within applications themselves can also
be removed if you and anyone using that application has no intention of
using Otto with it. If your application is being shared with other people
who may use Otto, it is very important to keep the `.ottoid` file.
