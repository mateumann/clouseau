# A go_binary builds a self-contained Go executable.
# It is possible to put these in the same directory as the go_library rule,
# but keeping them separate plays a little nicer with normal Go tools.
go_binary(
    name = 'collector',
    srcs = ['cmd/collector/main.go',],
    out = './collector',
    visibility = ['PUBLIC'],
)

go_library(
    name = 'clouseau',
    srcs = ['pkg/'],
    #deps = [
    #    '//proto:kitten',
    #    '//third_party/go:grpc',
    #],
    visibility = ['//go/...'],
)

# A go_test is similar to a go_binary in that it builds a self-contained
# executable, but it's typically run through 'plz test' to run all its
# tests and report their results.
#go_test(
#    name = 'kitten_test',
#    srcs = ['kitten_test.go'],
#    deps = [
#        ':lib',
#        '//third_party/go:testify',
#    ],
#)
