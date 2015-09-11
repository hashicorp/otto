application {
    name = "foo"
}

infrastructure "aws" {
    flavor = "foo"

    foundation "consul" {
        foo = "bar"
    }
}
