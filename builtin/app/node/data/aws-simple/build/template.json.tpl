{
    "min_packer_version": "0.8.0",

    "variables": {
      "aws_access_key": null,
      "aws_secret_key": null,
      "aws_region": null,
      "aws_vpc_id": null,
      "aws_subnet_id": null,
      "slug_path": null
    },

    "provisioners": [
      {% for dir in foundation_dirs.build %}
      {
        "type": "shell",
        "inline": ["mkdir -p /tmp/otto/foundation-{{ forloop.Counter }}"]
      },
      {
        "type": "file",
        "source": "{{ dir }}/",
        "destination": "/tmp/otto/foundation-{{ forloop.Counter }}"
      },
      {
        "type": "shell",
        "inline": ["cd /tmp/otto/foundation-{{ forloop.Counter}} && bash ./main.sh"]
      },
      {% endfor %}
      {
        "type": "file",
        "source": "{% verbatim %}{{ user `slug_path` }}{% endverbatim %}",
        "destination": "/tmp/otto-app.tgz"
      },
      {
        "type": "shell",
        "script": "build-node.sh"
      }
    ],

    "builders": [{
      "name": "otto",
      "type": "amazon-ebs",
      "access_key": "{% verbatim %}{{ user `aws_access_key` }}{% endverbatim %}",
      "secret_key": "{% verbatim %}{{ user `aws_secret_key` }}{% endverbatim %}",
      "region": "{% verbatim %}{{ user `aws_region` }}{% endverbatim %}",
      "vpc_id": "{% verbatim %}{{ user `aws_vpc_id` }}{% endverbatim %}",
      "subnet_id": "{% verbatim %}{{ user `aws_subnet_id` }}{% endverbatim %}",
      "source_ami": "ami-21630d44",
      "instance_type": "c3.large",
      "ssh_username": "ubuntu",
      "ami_name": "{{name}} {% verbatim %}{{timestamp}}{% endverbatim %}"
    }]

}
