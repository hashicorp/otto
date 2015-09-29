# Images

This folder contains the various configuration files to build the
pre-built images that are used as defaults throughout Otto.

Subfolders are organized by what the image being built is. Almost all
images will be built by Packer, but other image types may exist as well.

## Running Packer

To run the Packer templates, the templates should be executed from this
directory (containing the README). The templates themselves are setup to
expect the pwd to be set to this directory for paths.

Example:

```
$ packer build app-go-dev/template.json
...
```
