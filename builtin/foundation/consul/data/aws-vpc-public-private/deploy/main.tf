provider "aws" {
    access_key = "${var.aws_access_key}"
    secret_key = "${var.aws_secret_key}"
    region = "${var.region}"
}

module "consul-1" {
    source = "./module"

    index = "1"
    private-ip = "10.0.1.8"
    ami = "ami-de253bb6"
    key-name = "${var.key_name}"
    subnet-id = "${var.subnet-private}"
    vpc-id = "${var.vpc_id}"
}

module "consul-2" {
    source = "./module"

    index = "2"
    private-ip = "10.0.1.7"
    ami = "ami-de253bb6"
    key-name = "${var.key_name}"
    subnet-id = "${var.subnet-private}"
    vpc-id = "${var.vpc_id}"
}

module "consul-3" {
    source = "./module"

    index = "3"
    private-ip = "10.0.1.6"
    ami = "ami-de253bb6"
    key-name = "${var.key_name}"
    subnet-id = "${var.subnet-private}"
    vpc-id = "${var.vpc_id}"
}
