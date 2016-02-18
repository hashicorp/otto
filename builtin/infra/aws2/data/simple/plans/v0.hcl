plan {
    description = "Create the base infrastructure: VPC, subnets, etc."

    task "otto.infra.deploy_version" {
        description = "Update our deployment version to v0"

        infra = "${input.context.infra.name}"
        deploy_version = "0.0.0"
    }

    task "otto.infra.creds" {
        description = "Load AWS credentials"
    }

    task "terraform.apply" {
        description = "Create the AWS infrastructure"

        pwd = "${input.context.compile_dir}/v0"
        infra = "${input.context.infra.name}"
        "var.aws_region" = "us-east-1"
    }
}
