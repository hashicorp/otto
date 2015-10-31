---
layout: "intro"
page_title: "Otto vs. Heroku"
sidebar_current: "vs-other-heroku"
description: |-
  Comparison between Otto and Heroku.
---

# Otto vs. Heroku

Heroku is a platform for deploying just about any kind of application.
Heroku manages the underlying infrastructure, and provides easy knobs
for scaling up and down. Heroku addons allow applications to easily
talk to external services such as databases, queues, etc.

Heroku provides no solutions for development environments. Heroku can only
deploy to their infrastructure (which is on AWS for now, but you don't have
direct access to the underlying infrastructure).

Otto aims to replicate the ease of Heroku deploys on any infrastructure
while also providing a complete development experience. Otto can integrate
more easily with custom services and custom infrastructure components that
may not fit easily within the Heroku model.

Additionally, the price of Heroku rises dramatically, very quickly becoming
unaffordable for hobbyists and even startups. Otto itself is of course free,
and it has multiple options for infrastructure flavors to balance cost.
It is certainly cheaper than Heroku.

Admittedly, the deployment experience with Otto is still not as clean
and simple as Heroku, but for an initial release it is close. As time goes
on, the Otto deployment experience will only get better. Otto's development
experience is already fantastic.
