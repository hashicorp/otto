---
layout: "app_python"
page_title: "Build & Deploy - Python App Type"
sidebar_current: "docs-python-deploy"
description: |-
  Otto defaults to assuming your Python application is a public-facing web
  application, and deploys it with this assumption.
---

# Build & Deploy

Otto defaults to assuming your Python application is a public-facing web
application, and deploys it with this assumption.

This page documents
all the common deployment choices made for all infrastructures. The sidebar
on the left can be used to view infrastructure-specific choices that are
made for certain infrastructure targets.

## Common Points

Below is an unordered list of common points about the build and deploy
process. Please see the [customizations](/docs/apps/python/customization.html)
page for a list of behavior that can be changed.

  * The application is deployed behind [Gunicorn](http://gunicorn.org/)
    and [Nginx](http://nginx.org/).

  * `pip install -r requirements.txt` and `python setup.py install` is used
    during the build process to get all the dependencies for your applicaiton.
