plan {
    description = "Create the base infrastructure: VPC, subnets, etc."

    task "terraform.apply" {
        description = "Create the AWS infrastructure"

        pwd = "${input.context.compile_dir}/v0"
        infra = "${input.context.infra.name}"
    }
}
