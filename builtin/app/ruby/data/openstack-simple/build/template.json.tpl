{
    "min_packer_version": "0.8.0",

    "variables": {
      "openstack_auth_url": null,
      "openstack_username": null,
      "openstack_tenant_name": null,
      "openstack_password": null,
      "openstack_region_name": null,
      "openstack_image_id": null,
      "openstack_flavor_id": null,
      "openstack_floating_ip_pool": null,
      "network_id": null,
      "security_group": null,
      "openstack_ssh_username": null,
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
        "script": "build-ruby.sh"
      }
    ],

    "builders": [{
      "name": "otto",
      "type": "openstack",
      "region": "{% verbatim %}{{ user `openstack_region_name` }}{% endverbatim %}",
      "source_image": "{% verbatim %}{{ user `openstack_image_id` }}{% endverbatim %}",
      "flavor": "{% verbatim %}{{ user `openstack_flavor_id` }}{% endverbatim %}",
      "ssh_username": "{% verbatim %}{{ user `openstack_ssh_username` }}{% endverbatim %}",
      "floating_ip_pool": "{% verbatim %}{{ user `openstack_floating_ip_pool` }}{% endverbatim %}",
      "use_floating_ip": true,
      "networks": ["{% verbatim %}{{ user `network_id` }}{% endverbatim %}"],
      "security_groups": ["{% verbatim %}{{ user `security_group` }}{% endverbatim %}"],
      "image_name": "{{name}} {% verbatim %}{{timestamp}}{% endverbatim %}"
    }]

}
