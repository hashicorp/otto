#--------------------------------------------------------------
# AWS settings
#--------------------------------------------------------------
variable "aws_region" {
    description = "AWS region"
    default = "us-east-1"
}

#--------------------------------------------------------------
# General settings
#--------------------------------------------------------------
variable "ssh-key" {
    description = "SSH key name"
}

variable "subnet-private-1" {
    description = "Private subnet #1"
}

variable "subnet-private-2" {
    description = "Private subnet #2"
}

variable "vpc-id" {
    description = "VPC ID"
}
