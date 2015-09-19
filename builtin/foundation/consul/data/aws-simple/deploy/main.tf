provider "aws" {
    access_key = "${var.aws_access_key}"
    secret_key = "${var.aws_secret_key}"
    region = "${var.region}"
}

module "consul-1" {
    source = "./module-aws-simple"

    index = "1"
    private-ip = "10.0.2.6"
    ami = "${var.ami}"
    key-name = "${var.key_name}"
    subnet-id = "${var.subnet_public}"
    vpc-id = "${var.vpc_id}"
    vpc-cidr = "${var.vpc_cidr}"
}

output "consul_address" {
    value = "${module.consul-1.address}"
}
