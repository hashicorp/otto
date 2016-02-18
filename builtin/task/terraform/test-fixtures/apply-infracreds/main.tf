variable "foo" {}

output "bar" {
    value = "${var.foo}!"
}
