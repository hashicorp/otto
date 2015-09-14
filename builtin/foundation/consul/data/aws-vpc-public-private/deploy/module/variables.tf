variable "ami" {
    description = "AMI to use for Consul"
}

variable "index" {
    description = "Index for the name"
}

variable "key-name" {
    description = "SSH key name"
}

variable "bastion_host" {
    description = "SSH bastion host"
}

variable "bastion_user" {
    description = "SSH bastion user"
}

variable "private-ip" {
    description = "IP to assign to the instance"
}

variable "subnet-id" {
    description = "Subnet ID"
}

variable "vpc-id" {
    description = "VPC ID"
}

variable "join_addr" {
    description = "IP/Address to `consul join`"
}
