---
layout: "intro"
page_title: "Teardown"
sidebar_current: "gettingstarted-teardown"
description: |-
  In this step, we teardown all the resources we made with Otto.
---

# Teardown

During this getting started guide, we created a development environment,
launched infrastructure, built an AMI, and deployed an application.
This created a lot of real resources.

Otto makes it just as simple to teardown all these resources. In this
page, we'll do this as a final step to clean up everything we've done.

-> **Important!** If you don't teardown the cloud resources, they will
   eventually cost you real money. Please make sure you teardown these
   resources.

## Teardown

Run the following:

1. `otto deploy destroy`
1. `otto infra destroy`
1. `otto dev destroy`

The ordering of the deploy and infra steps is important. The dev
environment can be destroyed at any time.

At the end of all of these commands, the resources associated with
each will be fully destroyed. Cloud resources and the local development
environment are completely deleted.

## Next

We've now completed the getting started guide. See
[next steps](/intro/getting-started/next-steps.html) for what to do
next with Otto.
