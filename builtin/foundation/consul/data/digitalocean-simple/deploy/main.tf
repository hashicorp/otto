provider "digitalocean" {
  token = "${var.do_token}"
}

module "consul-1" {
    source = "./module-digitalocean-simple"

    index = "1"
    image = "${var.image}"
    key-name = "${var.key_name}"
}

output "consul_address" {
    value = "${module.consul-1.address}"
}
