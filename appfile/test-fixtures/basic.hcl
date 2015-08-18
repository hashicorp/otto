application {
    name = "foo"

    dependency {
        source = "foo"
    }

    dependency {
        source = "bar"
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
