plan {
    task "test-1" {
        description = "desc foo"
    }

    task "test-1" {
        description = "desc foo"
        foo = "${result.Result}"
    }

    task "test-1" {
        description = "desc foo"
        foo = "${bar}"
    }
}
