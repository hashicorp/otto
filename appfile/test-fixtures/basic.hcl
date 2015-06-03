application {
    name = "foo"
}

project {
    name = "foo"
    infrastructure = "aws"

    stack "bar" {}
}

infrastructure "aws" {
    flavor = "foo"
}
