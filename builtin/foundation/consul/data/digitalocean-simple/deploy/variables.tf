#--------------------------------------------------------------------
# Access Info
#--------------------------------------------------------------------

variable "api_token" {
    description = "Access key for AWS"
}

variable "region" {
    description = "Region where we will operate."
}

#--------------------------------------------------------------
# General settings
#--------------------------------------------------------------

variable "image" {
    description = "AMI to launch with Consul"
    default = "ubuntu-14-04-x64"
}

variable "key_name" {
    description = "SSH key name"
}
