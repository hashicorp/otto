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
}

infrastructure "aws" {
    flavor = "foo"
}
