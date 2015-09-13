variable "aws_access_key" {
    description = "Access key for AWS"
}

variable "aws_secret_key" {
    description = "Secret key for AWS"
}

variable "aws_region" {
    description = "Region where we will operate."
}

variable "ssh_public_key" {
    description = "Contents of an SSH public key to grant access to created instances"
}
