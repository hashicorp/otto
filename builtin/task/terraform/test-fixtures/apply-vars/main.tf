variable "foo" {}

output "output" {
    value = "${var.foo}!"
}
