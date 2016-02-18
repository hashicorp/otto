plan {
    task "store" {
        foo = "hello"
    }

    task "test" {
        result = "${foo}"
    }
}
