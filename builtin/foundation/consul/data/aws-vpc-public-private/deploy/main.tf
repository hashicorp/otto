provider "aws" {
    region = "${var.aws_region}"
}

module "consul-1" {
    source = "./module"

    index = "1"
    private-ip = "10.0.1.8"
    ami = "ami-de253bb6"
    key-name = "${var.ssh-key}"
    subnet-id = "${var.subnet-private-2}"
    vpc-id = "${var.vpc-id}"
}

module "consul-2" {
    source = "./module"

    index = "2"
    private-ip = "10.0.1.7"
    ami = "ami-de253bb6"
    key-name = "${var.ssh-key}"
    subnet-id = "${var.subnet-private-2}"
    vpc-id = "${var.vpc-id}"
}

module "consul-3" {
    source = "./module"

    index = "3"
    private-ip = "10.0.1.6"
    ami = "ami-de253bb6"
    key-name = "${var.ssh-key}"
    subnet-id = "${var.subnet-private-2}"
    vpc-id = "${var.vpc-id}"
}
