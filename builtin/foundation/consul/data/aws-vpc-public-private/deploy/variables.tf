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
variable "key_name" {
    description = "SSH key name"
}

variable "subnet-private" {
    description = "Private subnet"
}

variable "vpc_id" {
    description = "VPC ID"
}
