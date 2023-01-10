go_test(
    name = "engine_test",
    srcs = ["engine_test.go"],
    library = ":engine",
    deps = ["//third_party/golang/dns"],
)

go_binary(
    name = "engine",
    srcs = ["engine.go"],
    deps = ["//third_party/golang/dns"],
)
