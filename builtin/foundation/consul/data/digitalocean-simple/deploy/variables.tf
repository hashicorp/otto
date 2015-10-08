#--------------------------------------------------------------------
# Access Info
#--------------------------------------------------------------------

variable "do_token" {
	description = "API Token for DigitalOcean"
}

variable "region" {
    description = "Region where we will operate."
}

#--------------------------------------------------------------
# General settings
#--------------------------------------------------------------

variable "image" {
    description = "Droplet to launch with Consul"
    default = "ubuntu-14-04-x64"
}

variable "key_name" {
    description = "SSH key name"
}
