---
layout: "app_node"
page_title: "Development - Node.js App Type"
sidebar_current: "docs-node-dev"
description: |-
  The development environment built for Node.js applications is built for
  general Node.js development with a lean towards web development.
---

# Development

The development environment built for Node.js applications is built for
general Node.js development with a lean towards web development.

Please see the [customizations](/docs/apps/node/customization.html)
page for details on how to customize some of the behavior on this page.

## Pre-Installed Software

  * **Node.js**
  * **Git, Mercurial** - Useful for fetching dependencies
  * **GCC/G++ 0.8** - Required for newer versions of Node

## Usage

To work on your project, edit files locally on your own machine. The file changes
will be synced to the development environment.

When you're ready to build or test your project, run `otto dev ssh` to enter
the development environment. You'll be placed directly into the working
directory where you can run `npm install`, `npm run`, etc.

You can access the environment from your machine using development
environment's IP address. For example, if your app is running on port 5000,
then access it using the IP address from `otto dev address` plus that port.
