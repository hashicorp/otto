#--------------------------------------------------------------------
# Access Info
#--------------------------------------------------------------------

variable "aws_access_key" {
    description = "Access key for AWS"
}

variable "aws_secret_key" {
    description = "Secret key for AWS"
}

variable "aws_region" {
    description = "Region where we will operate."
}

#--------------------------------------------------------------
# General settings
#--------------------------------------------------------------
variable "ssh-key" {
    description = "SSH key name"
}

variable "subnet-private" {
    description = "Private subnet"
}

variable "vpc-id" {
    description = "VPC ID"
}
