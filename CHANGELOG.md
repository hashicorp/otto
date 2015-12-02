## Next Version

FEATURES:

  * **Layered Dev Environments**: Dev environments are now layered. Each
    layer is cached. When bringing up a new development environment, cached
    layers are used to speed it up immensely.
  * **App Type Plugins**: You can now add custom app types (or even override
    Otto's built-in types) using app type plugins. Want to support a new
    language? A new framework? App type plugins are for you.

IMPROVEMENTS:

  * app/ruby: Automatically detect desired Ruby version and install it [GH-293]
  * app/ruby: Install PhantomJS when poltergeist is detected [GH-313]
  * core: Use releases.hashicorp.com to download HashiCorp binaries [GH-353]

BUG FIXES:

  * app/ruby: Fix libxml2 package name [GH-320]
  * command/compile: compilation works if Appfile is a directory (it
      ignores the directory and detects an Appfile) [GH-280]

## 0.1.2 (October 20, 2015)

IMPROVEMENTS:

  * core: IP addresses for dev environments are now in the RFC 6598 space [GH-113]
  * core: Added `otto dev halt` action to halt the Vagrant machine [GH-195]
  * core: Otto will error if it detects a compiled environment from a newer
      version of Otto [GH-254]
  * app/custom: Vagrantfile for dev is rendered as a template [GH-168]
  * app/rails: Support for Rails projects [GH-190]
  * app/ruby: Use --no-document when installing bundler in dev [GH-130]
  * app/ruby: Install apt deps based on detected gems [GH-137] [GH-250]
  * app/ruby: Bundle automatically [GH-156]
  * app/php: support customizing PHP version [GH-105]
  * app: support Vagrant parallels provider in dev [GH-85]
  * command/help: A "help" command was introduced which does nothing except
      guide people to the proper way to ask for help. [GH-74]

BUG FIXES:

  * core: Ctrl-C now works when asking for credential password [GH-252]
  * appfile: some git dependencies that weren't working now do
  * appfile: application type is merged separately from other applications
      fields so it is optional [GH-192] [GH-212]
  * appfile: dependencies don't need an Appfile (but they do need
      a .ottoid) [GH-237]
  * app: fix Vagrant warning about box name [GH-110]
  * app: support Vagrant dev versions
  * app: don't error if no internet is availabile
  * app: `VAGRANT_CWD` won't cause dev to break [GH-262]
  * app: Friendly error message if you attempt to SSH into a dev environment
      that hasn't been created yet. [GH-69]
  * app/node: fix node download directory in build [GH-125]
  * app/ruby: allow `gem install` to work as `vagrant` user in dev [GH-129]
  * app/ruby, app/php: Fix `package.json` causing apps to be detected as Node.js [GH-149]
  * command/build: show help if any args are given [GH-245]

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
