## 0.1.2 (unreleased)

IMPROVEMENTS:

  * core: IP addresses for dev environments are now in the RFC 6598 space [GH-113]
  * core: Added `otto dev halt` action to halt the Vagrant machine [GH-195]
  * app/rails: Support for Rails projects [GH-190]
  * app/ruby: Use --no-document when installing bundler in dev [GH-130]
  * app/ruby: Install apt deps based on detected gems [GH-137] [GH-250]
  * app/ruby: Bundle automatically [GH-156]
  * app/php: support customizing PHP version [GH-105]
  * app: support Vagrant parallels provider in dev [GH-85]

BUG FIXES:

  * appfile: some git dependencies that weren't working now do
  * appfile: application type is merged separately from other applications
      fields so it is optional [GH-192] [GH-212]
  * app/node: fix node download directory in build [GH-125]
  * app/ruby: allow `gem install` to work as `vagrant` user in dev [GH-129]
  * app/ruby, app/php: Fix `package.json` causing apps to be detected as Node.js [GH-149]
  * app: fix Vagrant warning about box name [GH-110]

PLUGIN DEV CHANGES:

  * Template `extends` and `include` support: you can now include/extend
      templates for better reusability.
  * Template shares: there are now shared templates to include/extend from
      that contain common behavior (such as Vagrantfiles).

## 0.1.1 (September 28, 2015)

BUG FIXES:

  * core: Fix marshalling format problem preventing successful parsing of an
      Appfile with multiple dependencies declared [GH-83]
  * app: Fix string escaping issues affecting `otto dev` on Windows hosts [GH-79]
  * app/node: Fix issue preventing `node` apps from interpolating configs properly [GH-73]

## 0.1.0 (September 28, 2015)

  * Initial release
