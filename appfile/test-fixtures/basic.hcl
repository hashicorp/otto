application {
    name = "foo"

    dependency {
        source = "foo"
    }
}

project {
    name = "foo"
    infrastructure = "aws"

    stack "bar" {}
}

infrastructure "aws" {
    flavor = "foo"
}
