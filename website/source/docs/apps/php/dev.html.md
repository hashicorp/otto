---
layout: "app_php"
page_title: "Development - PHP App Type"
sidebar_current: "docs-php-dev"
description: |-
  The development environment built for PHP applications is built for
  general PHP development with a lean towards web development.
---

# Development

The development environment built for PHP applications is built for
general PHP development with a lean towards web development.

Please see the [customizations](/docs/apps/php/customization.html)
page for details on how to customize some of the behavior on this page.

## Pre-Installed Software

  * **PHP**
  * **Common PHP Modules** - mcrypt, mysql, fpm, gd, readline, pgsql
  * **Composer** - PHP dependency management
  * **Git, Mercurial, Bazaar** - Useful for pulling
    Composer dependencies

## Usage

You can access your environment via SSH to run your application.

For example, if you start up a development server like this...

```
$ otto dev ssh
vagrant@precise64:/vagrant$ php -S 0.0.0.0:5000
PHP 5.6.13-1+deb.sury.org~precise+3 Development Server started at Tue Sep 22 12:38:55 2015
Listening on http://0.0.0.0:5000
Document root is /vagrant
Press Ctrl-C to quit.
```

...then you should be able to access your app on port 5000 of the IP address reported by `otto dev address`.
