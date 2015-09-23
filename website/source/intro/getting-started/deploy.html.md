---
layout: "intro"
page_title: "Deploy"
sidebar_current: "gettingstarted-deploy"
description: |-
  Deploy an application with Otto in the Otto getting started guide.
---

# Deploy

We've now [launched infrastructure](/intro/getting-started/infra.html)
and [built the application](/intro/getting-started/build.html). It is time
to deploy it.

To deploy, run `otto deploy`.

Otto will now take the AMI built in the previous step and launch a
server in the infrastructure built previously. Otto will configure firewalls
properly to secure the server if necessary.

Once the deploy is done, you'll see the IP address of the application.
Open this address in your browser and you'll see that the application is
deployed!

TODO: SCREENSHOT

If you ever want to know more about a deploy, you can use `otto deploy info`.
This will include information such as the address of the deployed application.

For our simple application, Otto launched a single server. If the
infrastructure flavor (covered in the
[infrastructure page](/intro/getting-started/infra.html)) was
something other than "simple", Otto would launch multiple servers
behind a load balancer. This is another example of the power of
Otto: it has built-in knowledge of different ways to deploy an application
depending on your needs.

## Next

We've deployed the example application!

With a few simple steps, anybody could've deployed this application,
whether or not they know anything about operations, Ruby runtimes, etc.
Otto exposes unified workflow from development to deployment for any
application type.

As technology improves or best practices changes, these will be
encoded into Otto and Otto will roll out these change for you. For example,
if cheaper or more optimized hardware becomes available in AWS, then
Otto can begin to use those instead.

Up to this point, we've done everything with _zero configuration_
(other than the credentials to access AWS). Next, we'll introduce
a way to configure Otto using the [Appfile](/intro/getting-started/appfile.html).
