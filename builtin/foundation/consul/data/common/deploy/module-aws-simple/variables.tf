variable "ami" {
    description = "AMI to use for Consul"
}

variable "index" {
    description = "Index for the name"
}

variable "key-name" {
    description = "SSH key name"
}

variable "private-ip" {
    description = "IP to assign to the instance"
}

variable "subnet-id" {
    description = "Subnet ID"
}

variable "vpc-cidr" {
    description = "VPC CIDR"
}

variable "vpc-id" {
    description = "VPC ID"
}
