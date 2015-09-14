provider "aws" {
    access_key = "${var.aws_access_key}"
    secret_key = "${var.aws_secret_key}"
    region = "${var.region}"
}

module "consul-1" {
    source = "./module"

    index = "1"
    private-ip = "10.0.1.8"
    ami = "${var.ami}"
    key-name = "${var.key_name}"
    subnet-id = "${var.subnet-private}"
    vpc-id = "${var.vpc_id}"
    join_addr = "10.0.1.6"
    bastion_host = "${var.bastion_host}"
    bastion_user = "${var.bastion_user}"
}

module "consul-2" {
    source = "./module"

    index = "2"
    private-ip = "10.0.1.7"
    ami = "${var.ami}"
    key-name = "${var.key_name}"
    subnet-id = "${var.subnet-private}"
    vpc-id = "${var.vpc_id}"
    join_addr = "10.0.1.6"
    bastion_host = "${var.bastion_host}"
    bastion_user = "${var.bastion_user}"
}

module "consul-3" {
    source = "./module"

    index = "3"
    private-ip = "10.0.1.6"
    ami = "${var.ami}"
    key-name = "${var.key_name}"
    subnet-id = "${var.subnet-private}"
    vpc-id = "${var.vpc_id}"
    join_addr = "10.0.1.6"
    bastion_host = "${var.bastion_host}"
    bastion_user = "${var.bastion_user}"
}
