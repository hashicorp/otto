resource "aws_instance" "consul" {
    ami = "${var.ami}"
    instance_type = "t2.micro"
    key_name = "${var.key-name}"
    subnet_id = "${var.subnet-id}"
    vpc_security_group_ids = ["${aws_security_group.consul.id}"]
    private_ip = "${var.private-ip}"

    tags {
        Name = "consul ${var.index}"
    }

    connection {
        user         = "ubuntu"
        host         = "${self.private_ip}"
        bastion_host = "${var.bastion_host}"
        bastion_user = "${var.bastion_user}"
    }

    provisioner "file" {
        source = "${path.module}/join.sh"
        destination = "/tmp/script.sh"
    }

    provisioner "remote-exec" {
        inline = [
            "chmod +x /tmp/script.sh",
            "/tmp/script.sh ${var.join_addr}",
        ]
    }
}

resource "aws_security_group" "consul" {
    name = "consul ${var.index}"
    description = "Security group for Consul ${var.index}"
    vpc_id = "${var.vpc-id}"

    ingress {
        from_port = 1
        to_port = 65535
        protocol = "udp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    ingress {
        from_port = 1
        to_port = 65535
        protocol = "tcp"
        cidr_blocks = ["0.0.0.0/0"]
    }

    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}
