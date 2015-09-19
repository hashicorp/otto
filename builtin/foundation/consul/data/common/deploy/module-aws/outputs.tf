output "address" {
    value = "${aws_instance.consul.private_ip}"
}
