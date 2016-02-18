resource "null_resource" "foo" {
    provisioner "local-exec" {
        command = "otto-idontexistprobably"
    }
}
