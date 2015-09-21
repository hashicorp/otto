{
    "min_packer_version": "0.8.0",

    "variables": {
      "aws_access_key": null,
      "aws_secret_key": null,
      "aws_region": null,
      "slug_path": null
    },

    "provisioners": [
      {
        "type": "file",
        "source": "{% verbatim %}{{ user `slug_path` }}{% endverbatim %}",
        "destination": "/tmp/otto-app.tgz"
      },
      {
        "type": "shell",
        "script": "build-php.sh"
      }
    ],

    "builders": [{
      "name": "otto",
      "type": "amazon-ebs",
      "access_key": "{% verbatim %}{{ user `aws_access_key` }}{% endverbatim %}",
      "secret_key": "{% verbatim %}{{ user `aws_secret_key` }}{% endverbatim %}",
      "region": "{% verbatim %}{{ user `aws_region` }}{% endverbatim %}",
      "source_ami": "ami-21630d44",
      "instance_type": "c3.large",
      "ssh_username": "ubuntu",
      "ami_name": "{{name}} {% verbatim %}{{timestamp}}{% endverbatim %}"
    }]

}
