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

#--------------------------------------------------------------------
# Deploy Info
#--------------------------------------------------------------------

variable "ami" {
    description = "AMI to deploy"
}

variable "instance_type" {
    description = "Instance type"
    default = "t2.small"
}
