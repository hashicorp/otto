{
    "min_packer_version": "0.8.0",

    "variables": {
      "do_token": null,
      "do_region": "sfo1"
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
        "script": "build-ruby.sh"
      }
    ],

    "builders": [{
      "name": "otto",
      "type": "digitalocean",
      "api_token": "{% verbatim %}{{ user `do_token` }}{% endverbatim %}",
      "region": "{% verbatim %}{{ user `do_region` }}{% endverbatim %}",
      "image": "ubuntu-14-04-x64",
      "size": "512mb",
      "ssh_timeout": "5m",
      "snapshot_name": "{{name}} {% verbatim %}{{timestamp}}{% endverbatim %}"
    }]

}
