resource "digitalocean_droplet" "app" {
  image              = "${var.image}"
  region             = "${var.region}"
  name               = "${var.name}"
  size               = "${var.size}"
  private_networking = true
  user_data          = "{role: consul-${var.index}}"
  ssh_keys = [
    "${var.key-name}"
  ]

    provisioner "file" {
        source = "${path.module}/setup.sh"
        destination = "/tmp/script.sh"
    }

    provisioner "remote-exec" {
        inline = [
            "chmod +x /tmp/script.sh",
            "/tmp/script.sh",
       ]
    }
}

