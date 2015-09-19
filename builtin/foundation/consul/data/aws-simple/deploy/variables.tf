#--------------------------------------------------------------------
# Access Info
#--------------------------------------------------------------------

variable "aws_access_key" {
    description = "Access key for AWS"
}

variable "aws_secret_key" {
    description = "Secret key for AWS"
}

variable "region" {
    description = "Region where we will operate."
}

#--------------------------------------------------------------
# General settings
#--------------------------------------------------------------

variable "ami" {
    description = "AMI to launch with Consul"
    default = "ami-7f6a1f1a"
}

variable "key_name" {
    description = "SSH key name"
}

variable "subnet_public" {
    description = "Public subnet"
}

variable "vpc_id" {
    description = "VPC ID"
}

variable "vpc_cidr" {
    description = "VPC CIDR"
}
